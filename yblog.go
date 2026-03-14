package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func startServeWatcher(config data.Config, site *data.SiteState) {
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
				
				site.Lock()
				site.Generator = genNew
				site.Unlock()
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

	site := &data.SiteState{
		Generator: generator,
	}

	utils.CopyFiles(generator.Output, generator.Cache, ".", "", "./public")

	startServeWatcher(*config, site)

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

	handlrs := handlers.NewHandlers(config, site)

	e.GET("", handlrs.IndexHandler)
	e.GET("/", handlrs.IndexHandler)
	e.GET("/blog", handlrs.BlogIndexHandler)
	e.FileFS("/favicon.ico", "favicon.ico", afero.NewIOFS(generator.Output))
	// e.FileFS("/index.xml", "/feed.atom", outputIO)
	// e.FileFS("/feed.atom", "feed.atom", outputIO)
	// e.FileFS("/feed.rss", "feed.rss", outputIO)
	// e.FileFS("/feed.json", "feed.json", outputIO)
	e.GET("/feed.:format", handlrs.FeedHandler)
	e.GET("/feed_:type.:format", handlrs.SpecificFeedHandler)
	e.GET("/index.xml", handlrs.FeedHandler)
	e.GET("/posts/:post", handlrs.PostHandler)
	e.GET("/:post", handlrs.PostHandler)
	e.GET("/tags/:tag", handlrs.TagsHandler)
	e.GET("/kb", handlrs.KBIndexHandler)
	e.GET("/kb/*", handlrs.KBHandler)
	e.GET("/historical/*", func(c echo.Context) error {
		filePath := c.Param("*")
		if filePath == "" {
			return c.String(http.StatusNotFound, "Not Found")
		}
		return handlrs.ServeFile(c, fmt.Sprintf("historical/%s", filePath))
	})
	e.GET("/gallery", handlrs.GalleryIndexHandler)
	e.GET("/gallery/*", handlrs.GalleryHandler)
	e.GET("/images/*", func(c echo.Context) error {
		filePath := c.Param("*")
		if filePath == "" {
			return c.String(http.StatusNotFound, "Not Found")
		}
		return handlrs.ServeFile(c, fmt.Sprintf("images/%s", filePath))
	})
	e.GET("/favicon.ico", func(c echo.Context) error {
		return handlrs.ServeFile(c, "favicon.ico")
	})
	e.GET("/bookmarks", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "https://links.mistlyric.net")
	})
	if err := e.Start(":8080"); err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}

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

type New struct {
	Type  string `arg:"" help:"Content type (post, kb, til, page)"`
	Title string `arg:"" help:"Title of the content"`
}

func (n *New) Run(ctx *Context) error {
	slug := strings.ToLower(strings.ReplaceAll(n.Title, " ", "-"))
	baseDir := map[string]string{
		"post": "data/content/blog",
		"kb":   "data/content/kb",
		"til":  "data/content/til",
		"page": "data/content/pages",
	}

	targetDir, exists := baseDir[n.Type]
	if !exists {
		log.Fatalf("Invalid content type: %s. Must be post, kb, til, or page", n.Type)
	}

	dateStr := time.Now().Format("2006-01-02T15:04:05Z")
	content := fmt.Sprintf(`---
author: "Yulian Kuncheff"
date: %s
draft: false
slug: "%s"
title: "%s"
tags: []
type: "%s"
---

Write your content here...`, dateStr, slug, n.Title, n.Type)

	filePath := filepath.Join(targetDir, fmt.Sprintf("%s.md", slug))
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		log.WithError(err).Fatal("Failed to create directories")
	}

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.WithError(err).Fatal("Failed to write new file")
	}

	log.WithField("file", filePath).Info("Successfully created new content")
	return nil
}

type Context struct {
	Debug bool
}

var cli struct {
	Debug bool `help:"Enable debug mode"`

	Serve  Serve  `cmd:"" help:"Serve the generated files"`
	Deploy Deploy `cmd:"" help:"Deploy the generated files"`
	New    New    `cmd:"" help:"Create new content scaffolding"`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
