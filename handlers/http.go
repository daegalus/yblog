package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"yblog/data"
	"yblog/utils"

	"github.com/caarlos0/log"
	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"
)

type Handler struct {
	config *data.Config
	output afero.Fs
}

func NewHandlers(config *data.Config, output afero.Fs) *Handler {
	return &Handler{
		config: config,
		output: output,
	}
}

func (h *Handler) FeedHandler(c echo.Context) error {
	format := c.Param("format")

	if format != "atom" && format != "rss" && format != "json" && format != "" {
		log.WithField("format", format).Fatal("Invalid feed format")
		c.String(http.StatusNotFound, "")
		return nil
	}

	feedFilename := fmt.Sprintf("feed.%s", format)
	if format == "" {
		feedFilename = "feed.rss"
	}
	log.WithField("file", feedFilename).Info("Serving file")

	out, err := utils.ReadFileToString(h.output, feedFilename)
	if err != nil {
		log.WithField("error", err).Fatal("Error reading file")
		c.Error(err)
		return err
	}

	c.HTML(http.StatusOK, string(out))
	return nil
}

func (h *Handler) IndexHandler(c echo.Context) error {
	log.WithField("file", "index.html").Info("Serving file")

	out, err := utils.ReadFileToString(h.output, "index.html")
	if err != nil {
		log.WithField("error", err).Fatal("Error reading file")
		c.Error(err)
		return err
	}

	c.HTML(http.StatusOK, string(out))
	return nil
}

func (h *Handler) PostHandler(c echo.Context) error {
	prefix := "/posts"
	if strings.HasPrefix(c.Request().URL.Path, fmt.Sprintf("%s/", prefix)) {
		log.WithField("file", fmt.Sprintf("%s.html", c.Request().URL.Path)).Info("Serving file")

		out, err := utils.ReadFileToString(h.output, fmt.Sprintf("%s.html", c.Request().URL.Path[1:]))
		if err != nil {
			log.WithField("error", err).Error("Error reading file")
			c.Error(err)
			return err
		}

		c.HTML(http.StatusOK, out)
	} else {
		log.WithField("file", fmt.Sprintf("%s%s", prefix, c.Request().URL.Path)).Info("Redirecting to new URL")
		c.Redirect(http.StatusFound, fmt.Sprintf("%s%s", prefix, c.Request().URL.Path))
	}
	return nil
}

// func ImagesHandler(c echo.Context) error {
// 	log.WithField("file", c.Request().URL.Path).Info("Serving file")

// 	out, err := utils.ReadFile(data.Content, fmt.Sprintf("content/images/%s", strings.Replace(c.Request().URL.Path, "/images/", "", 1)))
// 	if err != nil {
// 		log.WithField("error", err).Fatal("Error reading file")
// 		c.Error(err)
// 		return err
// 	}

// 	fmt.Printf("out %s", string(out[:100]))

// 	image.RegisterFormat("ico", "ico", ico.Decode, ico.DecodeConfig)
// 	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
// 	_, format, err := image.Decode(bytes.NewReader(out))
// 	if err != nil {
// 		log.WithField("error", err).Fatal("Error decoding image")
// 		c.Error(err)
// 		return err
// 	}

// 	c.Blob(http.StatusOK, format, out)
// 	return nil
// }

func (h *Handler) TagsHandler(c echo.Context) error {
	tag := c.Param("tag")
	log.WithField("tag", tag).Info("Serving tag")

	//taggedPosts[tag]
	c.String(http.StatusOK, "")
	return nil
}
