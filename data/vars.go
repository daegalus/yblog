package data

import (
	"embed"
	"sync"
	"yblog/utils"
)

type SiteState struct {
	sync.RWMutex
	Generator *Generator
	Pages     map[string]*StandalonePage // Standalone root pages (e.g. /resume)
}

//go:embed content
//go:embed themes
var Content embed.FS

type Backlink struct {
	Slug    string
	Context string
}

type Post struct {
	FrontMatter    utils.FrontMatter
	Tagsline       string
	Markdown       []byte
	HTML           string
	Summary        string
	LegacyComments []*LegacyComment
	PrimaryDomain string
}

type KB struct {
	FrontMatter utils.FrontMatter
	Tagsline    string
	Path        string // relative path within kb (e.g., "devops/docker")
	Markdown    []byte
	HTML        string
	Backlinks     []Backlink // paths of KB pages that link to this one
	PrimaryDomain string
}

type Album struct {
	Slug        string
	Title       string   `toml:"title"`
	Description string   `toml:"description"`
	Cover       string   `toml:"cover"`
	Images      []string // filenames in images/ subfolder
}

type TIL struct {
	FrontMatter utils.FrontMatter
	Tagsline    string
	HTML        string
}

type Page struct {
	Header        string
	Nav           string
	Posts         []*Post
	SingleKB      *KB
	KBPages       []*KB
	GalleryAlbums []*Album
	SingleAlbum   *Album
	TILList       []*TIL
	CurrentTag    string
	AllTags       map[string]int // tag -> count across all content types
	Footer        string
	SinglePage    *StandalonePage
	Graph         string
}

type StandalonePage struct {
	FrontMatter utils.FrontMatter
	HTML        string
}



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
