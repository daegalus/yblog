package goldmark_extensions

import (
	"strings"

	"github.com/yuin/goldmark/ast"
)

// An Image struct represents an image of the Markdown text.
type PictureSet struct {
	ast.Image
	Sources []string
}

// Dump implements Node.Dump.
func (n *PictureSet) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Destination"] = string(n.Destination)
	m["Title"] = string(n.Title)

	strSources := []string{}
	strSources = append(strSources, n.Sources...)
	m["Sources"] = strings.Join(strSources, "")
	ast.DumpHelper(n, source, level, m, nil)
}

// KindImage is a NodeKind of the Image node.
var KindPicture = ast.NewNodeKind("PictureSet")

// Kind implements Node.Kind.
func (n *PictureSet) Kind() ast.NodeKind {
	return KindPicture
}

// NewImage returns a new Image node.
func NewPictureSet(image *ast.Image, sources *[]string) *PictureSet {
	c := &PictureSet{
		Image:   *image,
		Sources: *sources,
	}
	c.Destination = image.Destination
	c.Title = image.Title
	for n := image.FirstChild(); n != nil; {
		next := n.NextSibling()
		image.RemoveChild(image, n)
		c.AppendChild(c, n)
		n = next
	}

	return c
}
