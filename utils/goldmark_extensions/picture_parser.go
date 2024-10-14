package goldmark_extensions

import (
	"fmt"
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
	line, _ := block.PeekLine()
	fmt.Printf("line: %s\n", string(line))
	block.Advance(1) // skip @
	fmt.Println(string(block.Peek()))

	linkParser := parser.NewLinkParser()
	linkNode := linkParser.Parse(parent, block, pc)
	linkNode2 := linkParser.Parse(parent, block, pc)

	fmt.Println(linkNode2)

	fmt.Printf("block: %s", string(block.Peek()))

	fmt.Printf("title: %s, destination: %s, block: %s", linkNode.(*ast.Image).Title, linkNode.(*ast.Image).Destination, string(block.Source()))

	if block.Peek() == '{' {
		link := s.parseSources(linkNode, block, pc)
		return NewPictureSet(linkNode.(*ast.Image), link)
	} else {
		return linkNode
	}
}

func (s *pictureSetParser) parseSources(parent ast.Node, block text.Reader, pc parser.Context) *[]string {
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
