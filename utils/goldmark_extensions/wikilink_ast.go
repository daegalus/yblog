package goldmark_extensions

import (
	"github.com/yuin/goldmark/ast"
)

// WikiLink represents a wiki-style link [[target]] or [[target|display text]].
type WikiLink struct {
	ast.BaseInline
	Target  string // The target page slug/path (e.g., "devops/docker")
	Display string // Optional display text; if empty, Target is used
}

// Dump implements Node.Dump.
func (n *WikiLink) Dump(source []byte, level int) {
	m := map[string]string{
		"Target":  n.Target,
		"Display": n.Display,
	}
	ast.DumpHelper(n, source, level, m, nil)
}

// KindWikiLink is a NodeKind for WikiLink nodes.
var KindWikiLink = ast.NewNodeKind("WikiLink")

// Kind implements Node.Kind.
func (n *WikiLink) Kind() ast.NodeKind {
	return KindWikiLink
}

// NewWikiLink creates a new WikiLink node.
func NewWikiLink(target, display string) *WikiLink {
	return &WikiLink{
		Target:  target,
		Display: display,
	}
}
