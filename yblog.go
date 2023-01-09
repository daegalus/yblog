package main

import (
	"github.com/caarlos0/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/afero"

	"yblog/data"
	"yblog/handlers"
)

var Version string
var BuildDate string

func initialize() {
	log.WithFields(log.Fields{"version": Version, "build-date": BuildDate}).Info("yblog")
	memfs := afero.NewMemMapFs()
	data.Output = memfs
	log.Info("Initialized")
}

func main() {
	initialize()
	data.CompileMarkdown()

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
			log.WithFields(log.Fields{
				"method":    values.Method,
				"client_ip": values.RemoteIP,
				"latency":   values.Latency.Microseconds(),
				"time":      values.StartTime,
				"status":    values.Status,
				"route":     values.RoutePath,
				"length":    values.ResponseSize,
			}).Info(values.URI)

			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.HideBanner = true

	//outputIO := afero.NewIOFS(data.Output)
	imageFS := echo.MustSubFS(data.Content, "content/images")

	e.GET("", handlers.IndexHandler)
	e.StaticFS("/images/", imageFS)
	e.FileFS("/favicon.ico", "favicon.ico", imageFS)
	// e.FileFS("/index.xml", "/feed.atom", outputIO)
	// e.FileFS("/feed.atom", "feed.atom", outputIO)
	// e.FileFS("/feed.rss", "feed.rss", outputIO)
	// e.FileFS("/feed.json", "feed.json", outputIO)
	e.GET("/feed.:format", handlers.FeedHandler)
	e.GET("/index.xml", handlers.FeedHandler)
	e.GET("/posts/:post", handlers.PostHandler)
	e.GET("/:post", handlers.PostHandler)
	e.GET("/tags/:tag", handlers.TagsHandler)
	e.Start(":8080")
}
