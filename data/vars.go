package data

import (
	"embed"
	"yblog/utils"

	"github.com/spf13/afero"
)

//go:embed content
var Content embed.FS
var Output afero.Fs

type Post struct {
	FrontMatter utils.FrontMatter
	Tagsline    string
	Markdown    []byte
	HTML        string
	Summary     string
}

type Page struct {
	Header string
	Nav    string
	Posts  []*Post
	Footer string
}

var Posts map[string]*Post = map[string]*Post{}
var TaggedPosts map[string][]*Post = map[string][]*Post{}
