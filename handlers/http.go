package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"yblog/data"
	"yblog/utils"

	"github.com/caarlos0/log"
	"github.com/labstack/echo/v4"
)

func FeedHandler(c echo.Context) error {
	format := c.Param("format")

	if format != "atom" && format != "rss" && format != "json" && format != "" {
		log.WithField("format", format).Fatal("Invalid feed format")
		c.String(http.StatusNotFound, "")
		return nil
	}

	feedFilename := fmt.Sprintf("/feed.%s", format)
	if format == "" {
		feedFilename = "/feed.rss"
	}
	log.WithField("file", feedFilename).Info("Serving file")

	out, err := utils.ReadFileToString(data.Output, feedFilename)
	if err != nil {
		log.WithField("error", err).Fatal("Error reading file")
		c.Error(err)
		return err
	}

	c.HTML(http.StatusOK, string(out))
	return nil
}

func IndexHandler(c echo.Context) error {
	log.WithField("file", "index.html").Info("Serving file")

	out, err := utils.ReadFileToString(data.Output, "/index.html")
	if err != nil {
		log.WithField("error", err).Fatal("Error reading file")
		c.Error(err)
		return err
	}

	c.HTML(http.StatusOK, string(out))
	return nil
}

func PostHandler(c echo.Context) error {
	if strings.HasPrefix(c.Request().URL.Path, "/posts/") {
		log.WithField("file", fmt.Sprintf("%s.html", c.Request().URL.Path)).Info("Serving file")

		out, err := utils.ReadFileToString(data.Output, fmt.Sprintf("%s.html", c.Request().URL.Path))
		if err != nil {
			log.WithField("error", err).Fatal("Error reading file")
			c.Error(err)
			return err
		}

		c.HTML(http.StatusOK, out)
	} else {
		log.WithField("file", fmt.Sprintf("/posts/%s.html", c.Request().URL.Path)).Info("Serving file")

		out, err := utils.ReadFileToString(data.Output, fmt.Sprintf("/posts/%s.html", c.Request().URL.Path))
		if err != nil {
			log.WithField("error", err).Fatal("Error reading file")
			c.Error(err)
			return err
		}

		c.HTML(http.StatusOK, out)
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

func TagsHandler(c echo.Context) error {
	tag := c.Param("tag")
	log.WithField("tag", tag).Info("Serving tag")

	//taggedPosts[tag]
	c.String(http.StatusOK, "")
	return nil
}
