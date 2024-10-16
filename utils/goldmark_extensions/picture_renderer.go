package goldmark_extensions

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

const optPictureParentPath renderer.OptionName = "base64ParentPath"

// PictureConfig embeds html.Config to refer to some fields like unsafe and xhtml.
type PictureConfig struct {
	html.Config
	ParentPath   string
	ImageFormats []string
}

// SetOption implements renderer.NodeRenderer.SetOption
func (c *PictureConfig) SetOption(name renderer.OptionName, value any) {
	c.Config.SetOption(name, value)

	switch name {
	case optPictureParentPath:
		c.ParentPath = value.(string)
	}
}

type PictureOption interface {
	renderer.Option
	SetPictureOption(*PictureConfig)
}

func WithParentPath(path string) interface {
	renderer.Option
	PictureOption
} {
	return &withParentPath{path}
}

type withParentPath struct {
	path string
}

func (o *withParentPath) SetConfig(c *renderer.Config) {
	c.Options[optPictureParentPath] = o.path
}

func (o *withParentPath) SetPictureOption(c *PictureConfig) {
	c.ParentPath = o.path
}

type PictureRenderer struct {
	PictureConfig
}

func NewPictureRenderer(opts ...PictureOption) renderer.NodeRenderer {
	r := &PictureRenderer{
		PictureConfig: PictureConfig{
			Config: html.NewConfig(),
		},
	}
	for _, o := range opts {
		o.SetPictureOption(&r.PictureConfig)
	}
	return r
}

func (r *PictureRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindPicture, r.renderImage)
}

// see https://developer.mozilla.org/ja/docs/Web/Media/Formats/Image_types
var commonWebImages = func() map[string]struct{} {
	types := []string{
		"image/apng",
		"image/avif",
		"image/gif",
		"image/jpeg",
		"image/jxl",
		"image/png",
		"image/svg+xml",
		"image/webp",
	}
	m := map[string]struct{}{}
	for _, t := range types {
		m[t] = struct{}{}
	}
	return m
}()

// renderImage adds image embedding function to github.com/yuin/goldmark/renderer/html (MIT).
func (r *PictureRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*PictureSet)
	_, _ = w.WriteString("<picture>")

	urlSplit := strings.Split(string(n.Destination), ".")
	url := strings.Join(urlSplit[:len(urlSplit)-1], ".")

	for _, source := range n.Sources {
		_, _ = w.WriteString(` <source srcset="`)
		if r.Unsafe || !html.IsDangerousURL([]byte(fmt.Sprintf("%s.%s", url, source))) {
			_, _ = w.Write([]byte(fmt.Sprintf("%s.%s", url, source)))
		}
		_ = w.WriteByte('"')
		if r.XHTML {
			_, _ = w.WriteString(" />")
		} else {
			_, _ = w.WriteString(">")
		}
	}
	_, _ = w.WriteString(` <img src="`)
	if r.Unsafe || !html.IsDangerousURL(n.Destination) {
		_, _ = w.Write(n.Destination)
	}
	_, _ = w.WriteString(`" alt="`)
	altAttr, _ := n.AttributeString("alt")
	_, _ = w.WriteString(string(altAttr.([]byte)))
	_ = w.WriteByte('"')
	if n.Title != nil {
		_, _ = w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		_ = w.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	if r.XHTML {
		_, _ = w.WriteString(" />")
	} else {
		_, _ = w.WriteString(">")
	}
	return ast.WalkSkipChildren, nil
}

// Picture implements goldmark.Extender
type Picture struct {
	options []PictureOption
}

// Picture is an implementation of goldmark.Extender
var picture = &Picture{}

// NewPicture initializes Picture: goldmark's extension with its options.
// Using default Picture with goldmark.WithRendereOptions(opts) give the same result.
func NewPicture(opts ...PictureOption) goldmark.Extender {
	return &Picture{
		options: opts,
	}
}

func (e *Picture) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewPictureSetParser(), 100),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewPictureRenderer(e.options...), 100),
	))
}
