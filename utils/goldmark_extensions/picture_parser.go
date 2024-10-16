package goldmark_extensions

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type pictureSetParser struct {
}

var defaultPictureSetParser = &pictureSetParser{}

// NewPictureSetParser return a new InlineParser that parses links.
func NewPictureSetParser() parser.InlineParser {
	return defaultPictureSetParser
}

func (s *pictureSetParser) Trigger() []byte {
	return []byte{'@'}
}

var linkBottom = parser.NewContextKey()

func (s *pictureSetParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	block.Advance(1) // skip @
	if block.Peek() != '!' {
		return nil
	}
	block.Advance(1) // skip !
	if block.Peek() != '[' {
		return nil
	}
	block.Advance(1) // skip [
	var altText []byte
	for block.Peek() != ']' {
		altText = append(altText, block.Peek())
		block.Advance(1)
	}
	block.Advance(1) // skip ]
	if block.Peek() != '(' {
		return nil
	}
	block.Advance(1) // skip (
	var destination []byte
	for block.Peek() != '"' && block.Peek() != ')' {
		if block.Peek() == '\\' {
			block.Advance(1)
		}
		destination = append(destination, block.Peek())
		block.Advance(1)
	}
	var title []byte
	if block.Peek() == '"' {
		destination = destination[:len(destination)-1]
		block.Advance(1)
		for block.Peek() != '"' {
			title = append(title, block.Peek())
			block.Advance(1)
		}
		block.Advance(1)
	}
	block.Advance(1) // skip )
	if block.Peek() != '{' {
		return nil
	}

	sources := s.parseSources(block)
	if sources == nil {
		return nil
	}

	linkNode := ast.NewLink()
	linkNode.Destination = destination
	linkNode.Title = title

	imageNode := ast.NewImage(linkNode)
	imageNode.SetAttribute([]byte("alt"), altText)

	pictureNode := NewPictureSet(imageNode, sources)
	return pictureNode
}

func (s *pictureSetParser) parseSources(block text.Reader) *[]string {
	block.Advance(1) // skip '{'
	block.SkipSpaces()
	var sources []byte
	var ok bool
	if block.Peek() == '}' { // empty sources like '[link](){}'
		block.Advance(1)
	} else {
		sources, ok = parsePictureSources(block)
		if !ok {
			return nil
		}
		block.SkipSpaces()
		if block.Peek() == '}' {
			block.Advance(1)
		} else {
			return nil
		}
	}

	sourceList := strings.Split(string(sources), ",")

	return &sourceList
}

var pictureFindClosureOptions text.FindClosureOptions = text.FindClosureOptions{
	Nesting: false,
	Newline: true,
	Advance: true,
}

func parsePictureSources(block text.Reader) ([]byte, bool) {
	block.SkipSpaces()
	opener := block.Peek()
	if opener != '"' && opener != '\'' && opener != '(' {
		return nil, false
	}
	closer := opener
	if opener == '(' {
		closer = ')'
	}
	block.Advance(1)
	segments, found := block.FindClosure(opener, closer, pictureFindClosureOptions)
	if found {
		if segments.Len() == 1 {
			return block.Value(segments.At(0)), true
		}
		var sources []byte
		for i := 0; i < segments.Len(); i++ {
			s := segments.At(i)
			sources = append(sources, block.Value(s)...)
		}
		return sources, true
	}
	return nil, false
}
