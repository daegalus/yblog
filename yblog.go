package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/kong"
	"github.com/caarlos0/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/radovskyb/watcher"
	"github.com/spf13/afero"

	"yblog/data"
	"yblog/handlers"
	"yblog/utils"
)

var Version string
var BuildDate string

func initialize() (*data.Generator, *data.Config) {
	log.WithField("version", Version).WithField("build-date", BuildDate).Info("yblog")

	config := loadConfig()

	memfsInput := afero.NewMemMapFs()
	memfsOutput := afero.NewMemMapFs()
	memfsCache := afero.NewMemMapFs()

	copyFiles(afero.FromIOFS{FS: data.Content}, memfsInput, ".")
	if _, err := os.Stat(config.Site.Output); err == nil {
		copyFiles(afero.NewOsFs(), memfsCache, config.Site.Output)
	}

	log.WithField("config", config).Info("Config loaded")

	generator := data.NewGenerator(&config, memfsInput, memfsOutput, memfsCache)

	return generator, &config
}

func loadConfig() data.Config {
	var config data.Config
	_, err := toml.DecodeFile("yblog.toml", &config)
	if err != nil {
		log.WithError(err).Fatal("Error reading config file")
	}

	if config.Site.Theme == "" {
		config.Site.Theme = "simple"
	}

	if config.Site.Output == "" {
		config.Site.Output = "./public"
	}
	return config
}

func copyFiles(input afero.Fs, output afero.Fs, rootPath string) {
	err := afero.Walk(input, rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.WithError(err).Fatal("Error reading files during walk")
		}
		if info.IsDir() {
			return nil
		}

		input.MkdirAll(filepath.Dir(path), 0755)

		in, err := afero.ReadFile(input, path)
		if err != nil {
			log.WithError(err).Fatal("Error reading file during walk")
		}
		afero.WriteFile(output, path, in, 0644)

		return nil
	})
	if err != nil {
		log.WithError(err).Fatal("Error walking files")
	}
}

func startServeWatcher(config data.Config, generator *data.Generator, memfsInput afero.Fs, memfsOutput afero.Fs, memfsCache afero.Fs) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write, watcher.Create, watcher.Remove, watcher.Rename)

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
				genNew, _ := initialize()
				genNew.CompileMarkdown()
				generator.Input = genNew.Input
				generator.Output = genNew.Output
				generator.Cache = genNew.Cache
			case err := <-w.Error:
				log.WithError(err).Fatal("Error watching files")
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive("./data/content"); err != nil {
		log.WithError(err).Fatal("Error watching contents folder")
	}

	if err := w.AddRecursive("./data/themes"); err != nil {
		log.WithError(err).Fatal("Error watching contents folder")
	}

	if err := w.Add("./yblog.toml"); err != nil {
		log.WithError(err).Fatal("Error watching contents folder")
	}

	// if err := w.Start(time.Second * 1); err != nil {
	// 	log.WithError(err).Fatal("Error starting watcher")
	// }

	go w.Start(time.Second * 1)
}

type Serve struct{}

func (s *Serve) Run(ctx *Context) error {
	generator, config := initialize()
	generator.CompileMarkdown()

	utils.CopyFiles(generator.Output, generator.Cache, ".", "", "./public")

	startServeWatcher(*config, generator, generator.Input, generator.Output, generator.Cache)

	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:      true,
		LogMethod:       true,
		LogURI:          true,
		LogStatus:       true,
		LogRemoteIP:     true,
		LogRoutePath:    true,
		LogResponseSize: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithField("method", values.Method).
				WithField("client_ip", values.RemoteIP).
				WithField("latency", values.Latency.Microseconds()).
				WithField("time", values.StartTime).
				WithField("status", values.Status).
				WithField("route", values.RoutePath).
				WithField("length", values.ResponseSize).
				Info(values.URI)

			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.HideBanner = true

	imageFS := echo.MustSubFS(afero.NewIOFS(generator.Output), "images")
	handlrs := handlers.NewHandlers(config, generator.Output)

	e.GET("", handlrs.IndexHandler)
	e.StaticFS("/images/", imageFS)
	e.FileFS("/favicon.ico", "favicon.ico", afero.NewIOFS(generator.Output))
	// e.FileFS("/index.xml", "/feed.atom", outputIO)
	// e.FileFS("/feed.atom", "feed.atom", outputIO)
	// e.FileFS("/feed.rss", "feed.rss", outputIO)
	// e.FileFS("/feed.json", "feed.json", outputIO)
	e.GET("/feed.:format", handlrs.FeedHandler)
	e.GET("/index.xml", handlrs.FeedHandler)
	e.GET("/posts/:post", handlrs.PostHandler)
	e.GET("/:post", handlrs.PostHandler)
	e.GET("/tags/:tag", handlrs.TagsHandler)
	e.Start(":8080")

	return nil
}

type Deploy struct {
	Watch      bool   `help:"Watch for changes, and rebuild"`
	OutputPath string `help:"Output path of generated site" default:"./public"`
}

func (d *Deploy) Run(ctx *Context) error {
	generator, _ := initialize()
	generator.CompileMarkdown()

	utils.CopyFiles(generator.Output, generator.Cache, ".", "", d.OutputPath)

	osfs := afero.NewOsFs()
	afero.Walk(generator.Output, ".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.WithError(err).Fatal("Error reading files during deploy walk")
			return err
		}
		if info.IsDir() {
			return nil
		} else {
			osfs.MkdirAll(filepath.Join(d.OutputPath, filepath.Dir(path)), 0755)
		}

		in, err := generator.Output.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			log.WithError(err).Fatal("Error opening file for output")
			return err
		}
		content := make([]byte, info.Size())
		in.Read(content)

		oFile, err := osfs.Create(filepath.Join(d.OutputPath, path))
		if err != nil {
			log.WithError(err).Fatal("Error creating files for output")
			return err
		}
		oFile.Write(content)

		return nil
	})
	return nil
}

type Context struct {
	Debug bool
}

var cli struct {
	Debug bool `help:"Enable debug mode"`

	Serve  Serve  `cmd:"" help:"Serve the generated files"`
	Deploy Deploy `cmd:"" help:"Deploy the generated files"`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
