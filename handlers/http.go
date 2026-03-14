package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"yblog/data"

	"github.com/caarlos0/log"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	config *data.Config
	state  *data.SiteState
}

func NewHandlers(config *data.Config, state *data.SiteState) *Handler {
	return &Handler{
		config: config,
		state:  state,
	}
}

func (h *Handler) ServeFile(c echo.Context, filePath string) error {
	h.state.RLock()
	outFS := h.state.Generator.Output
	h.state.RUnlock()

	file, err := outFS.Open(filePath)
	if err != nil {
		return c.String(http.StatusNotFound, "Not Found")
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	http.ServeContent(c.Response(), c.Request(), stat.Name(), stat.ModTime(), file)
	return nil
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
	return h.ServeFile(c, feedFilename)
}

func (h *Handler) SpecificFeedHandler(c echo.Context) error {
	feedType := c.Param("type")
	format := c.Param("format")

	if format != "atom" && format != "rss" && format != "json" && format != "xml" {
		log.WithField("format", format).Fatal("Invalid feed format")
		c.String(http.StatusNotFound, "")
		return nil
	}

	if format == "xml" {
		format = "atom" // Default xml explicitly to atom standard.
	}

	feedFilename := fmt.Sprintf("feed_%s.%s", feedType, format)
	log.WithField("file", feedFilename).Info("Serving explicit file")
	return h.ServeFile(c, feedFilename)
}

func (h *Handler) BlogIndexHandler(c echo.Context) error {
	log.WithField("file", "blog/index.html").Info("Serving file")
	return h.ServeFile(c, "blog/index.html")
}

func (h *Handler) IndexHandler(c echo.Context) error {
	log.WithField("file", "index.html").Info("Serving file")
	return h.ServeFile(c, "index.html")
}

func (h *Handler) PostHandler(c echo.Context) error {
	prefix := "/blog"
	slug := c.Param("post")

	// Check if this is a Standalone Page first
	h.state.RLock()
	_, isPage := h.state.Generator.Pages[slug]
	h.state.RUnlock()

	if isPage {
		log.WithField("page", slug).Info("Serving standalone page")
		return h.ServeFile(c, fmt.Sprintf("%s/index.html", slug))
	}

	if strings.HasPrefix(c.Request().URL.Path, fmt.Sprintf("%s/", prefix)) {
		log.WithField("file", fmt.Sprintf("%s.html", c.Request().URL.Path)).Info("Serving post")
		filePath := fmt.Sprintf("%s/index.html", c.Request().URL.Path[1:])
		return h.ServeFile(c, filePath)
	} else {
		log.WithField("file", fmt.Sprintf("%s%s", prefix, c.Request().URL.Path)).Info("Redirecting post to new URL")
		return c.Redirect(http.StatusFound, fmt.Sprintf("%s%s", prefix, c.Request().URL.Path))
	}
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

func (h *Handler) KBIndexHandler(c echo.Context) error {
	log.WithField("file", "kb/index.html").Info("Serving KB index")
	return h.ServeFile(c, "kb/index.html")
}

func (h *Handler) KBHandler(c echo.Context) error {
	// Wildcard param captures the full nested path, e.g., "devops/docker"
	pagePath := c.Param("*")
	if pagePath == "" {
		return h.KBIndexHandler(c)
	}

	// Strip trailing slash if present
	pagePath = strings.TrimSuffix(pagePath, "/")

	filePath := fmt.Sprintf("kb/%s/index.html", pagePath)
	log.WithField("file", filePath).Info("Serving KB page")
	return h.ServeFile(c, filePath)
}

func (h *Handler) GalleryIndexHandler(c echo.Context) error {
	log.WithField("file", "gallery/index.html").Info("Serving gallery index")
	return h.ServeFile(c, "gallery/index.html")
}

func (h *Handler) GalleryHandler(c echo.Context) error {
	pagePath := c.Param("*")
	if pagePath == "" {
		return h.GalleryIndexHandler(c)
	}

	pagePath = strings.TrimSuffix(pagePath, "/")

	filePath := fmt.Sprintf("gallery/%s/index.html", pagePath)
	log.WithField("file", filePath).Info("Serving gallery page")
	return h.ServeFile(c, filePath)
}

func (h *Handler) TILHandler(c echo.Context) error {
	log.WithField("file", "til/index.html").Info("Serving TIL index")
	return h.ServeFile(c, "til/index.html")
}

func (h *Handler) TagsHandler(c echo.Context) error {
	tag := c.Param("tag")
	if tag == "" {
		// Serve tags index page
		log.WithField("file", "tags/index.html").Info("Serving tags index")
		return h.ServeFile(c, "tags/index.html")
	}

	log.WithField("tag", tag).Info("Serving tag")
	return h.ServeFile(c, fmt.Sprintf("tags/%s/index.html", tag))
}
