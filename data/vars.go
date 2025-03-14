package data

import (
	"embed"
	"yblog/utils"
)

//go:embed content
//go:embed themes
var Content embed.FS

type Post struct {
	FrontMatter    utils.FrontMatter
	Tagsline       string
	Markdown       []byte
	HTML           string
	Summary        string
	LegacyComments []*LegacyComment
}

type KB struct {
	FrontMatter utils.FrontMatter
	Tagsline    string
	Markdown    []byte
	HTML        string
	PrimaryDomain string
}

type Page struct {
	Header   string
	Nav      string
	Posts    []*Post
	SingleKB *KB
	Footer   string
	PrimaryDomain string
}

var Posts map[string]*Post = map[string]*Post{}
var KBs map[string]*KB = map[string]*KB{}
var TaggedPosts map[string][]*Post = map[string][]*Post{}
var TaggedKBs map[string][]*KB = map[string][]*KB{}

var ContentPrefix string = "content"
var ThemesPrefix string = "themes"
var CachePrefix string = "public"

type Config struct {
	Site struct {
		PrimaryDomain string `toml:"primary-domain"`
		Theme  string `toml:"theme"`
		Output string `toml:"output"`
	} `toml:"site"`
}
