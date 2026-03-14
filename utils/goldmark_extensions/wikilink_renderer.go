package goldmark_extensions

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// WikiLinkRenderer renders WikiLink nodes to HTML.
type WikiLinkRenderer struct{}

// NewWikiLinkRenderer returns a new WikiLinkRenderer.
func NewWikiLinkRenderer() renderer.NodeRenderer {
	return &WikiLinkRenderer{}
}

// RegisterFuncs registers the render function for WikiLink nodes.
func (r *WikiLinkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindWikiLink, r.renderWikiLink)
}

func (r *WikiLinkRenderer) renderWikiLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*WikiLink)
	display := n.Display
	if display == "" {
		display = n.Target
	}

	_, _ = w.WriteString(fmt.Sprintf(`<a href="/kb/%s" class="wikilink">%s</a>`, n.Target, display))
	return ast.WalkSkipChildren, nil
}

// WikiLinkExtension implements goldmark.Extender.
type WikiLinkExtension struct{}

// Extend adds the wiki-link parser and renderer to Goldmark.
func (e *WikiLinkExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		// Higher priority (lower number) so we match [[ before the standard link parser
		util.Prioritized(NewWikiLinkParser(), 50),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewWikiLinkRenderer(), 100),
	))
}
