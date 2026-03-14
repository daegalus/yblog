package goldmark_extensions

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type wikiLinkParser struct{}

var defaultWikiLinkParser = &wikiLinkParser{}

// NewWikiLinkParser returns a new InlineParser for [[wiki links]].
func NewWikiLinkParser() parser.InlineParser {
	return defaultWikiLinkParser
}

// Trigger returns the byte that triggers this parser.
// We trigger on '[' and then check for the second '['.
func (s *wikiLinkParser) Trigger() []byte {
	return []byte{'['}
}

// Parse parses [[target]] or [[target|display text]] syntax.
func (s *wikiLinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 2 || line[0] != '[' || line[1] != '[' {
		return nil
	}

	// Find the closing ]]
	closingIdx := strings.Index(string(line[2:]), "]]")
	if closingIdx < 0 {
		return nil
	}

	inner := string(line[2 : 2+closingIdx])
	if len(inner) == 0 {
		return nil
	}

	// Advance past [[ + inner + ]]
	block.Advance(2 + closingIdx + 2)

	// Parse target and optional display text
	target := inner
	display := ""
	if pipeIdx := strings.Index(inner, "|"); pipeIdx >= 0 {
		target = strings.TrimSpace(inner[:pipeIdx])
		display = strings.TrimSpace(inner[pipeIdx+1:])
	}

	if target == "" {
		return nil
	}

	return NewWikiLink(target, display)
}
