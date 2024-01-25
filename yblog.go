package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/radovskyb/watcher"
	"github.com/spf13/afero"

	"yblog/data"
	"yblog/handlers"
	"yblog/wasm"
)

var Version string
var BuildDate string

func initialize() (*data.Generator, *data.Config) {
	wasm.NewImageProcessing()
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//}
	//os.Exit(0)
	log.WithField("version", Version).WithField("build-date", BuildDate).Info("yblog")
	memfsInput := afero.NewMemMapFs()
	memfsOutput := afero.NewMemMapFs()

	copyContentFiles(memfsInput)

	config := loadConfig()
	log.WithField("config", config).Info("Config loaded")

	generator := data.NewGenerator(&config, memfsInput, memfsOutput)

	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write, watcher.Create, watcher.Remove, watcher.Rename)

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
				config = loadConfig()
				reloadFilesFromFS(memfsInput)
				generator = data.NewGenerator(&config, memfsInput, memfsOutput)
				generator.CompileMarkdown()
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

	return generator, &config
}

func loadConfig() data.Config {
	var config data.Config
	_, err := toml.DecodeFile("yblog.toml", &config)
	if err != nil {
		log.WithError(err).Fatal("Error reading config file")
	}
	return config
}

func copyContentFiles(input afero.Fs) {
	err := fs.WalkDir(data.Content, ".", func(path string, info os.DirEntry, err error) error {
		if err != nil {
			log.WithError(err).Fatal("Error reading files")
		}
		if info.IsDir() {
			return nil
		}

		input.MkdirAll(filepath.Dir(path), 0755)
		in, err := data.Content.ReadFile(path)
		if err != nil {
			log.WithError(err).Fatal("Error reading files")
		}
		afero.WriteFile(input, path, in, 0644)

		return nil
	})
	if err != nil {
		log.WithError(err).Fatal("Error reading files")
	}
}

func reloadFilesFromFS(input afero.Fs) {
	err := filepath.WalkDir("./data", func(path string, info os.DirEntry, err error) error {
		if err != nil {
			log.WithError(err).Fatal("Error reading files")
		}
		if info.IsDir() {
			return nil
		}

		input.MkdirAll(filepath.Dir(path), 0755)
		in, err := os.ReadFile(path)
		if err != nil {
			log.WithError(err).Fatal("Error reading files")
		}
		afero.WriteFile(input, path, in, 0644)

		return nil
	})
	if err != nil {
		log.WithError(err).Fatal("Error reading files")
	}
}

func main() {
	generator, config := initialize()
	generator.CompileMarkdown()

	// file, _ := afero.ReadFile(generator.Output, "/feed.atom")
	// fmt.Println(string(file))

	// fmt.Println("Input:")
	// afero.Walk(generator.Input, ".", func(path string, info fs.FileInfo, err error) error {
	// 	if err != nil {
	// 		log.WithError(err).Fatal("Error reading files")
	// 	}
	// 	if info.IsDir() {
	// 		return nil
	// 	}

	// 	fmt.Println(path)
	// 	return nil
	// })

	// fmt.Println("Output:")
	// afero.Walk(generator.Output, ".", func(path string, info fs.FileInfo, err error) error {
	// 	if err != nil {
	// 		log.WithError(err).Fatal("Error reading files")
	// 	}
	// 	if info.IsDir() {
	// 		return nil
	// 	}

	// 	fmt.Println(path)
	// 	return nil
	// })

	// os.Exit(0)

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

	imageFS := echo.MustSubFS(afero.NewIOFS(generator.Input), "content/images")
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
}
