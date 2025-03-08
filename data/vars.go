package data

import (
	"embed"
	"yblog/utils"
)

//go:embed content
//go:embed themes
var Content embed.FS

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

var ContentPrefix string = "content"
var ThemesPrefix string = "themes"
var CachePrefix string = "public"

type Config struct {
	Site struct {
		Theme  string `toml:"theme"`
		Output string `toml:"output"`
	} `toml:"site"`
}
