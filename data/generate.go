package data

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/caarlos0/log"

	"github.com/daegalus/feeds"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/spf13/afero"
)

type Generator struct {
	config         *Config
	Input          afero.Fs
	Output         afero.Fs
	Cache          afero.Fs
	Posts          map[string]*Post
	KBs            map[string]*KB
	Albums         map[string]*Album
	TILs           map[string]*TIL
	Pages          map[string]*StandalonePage
	TaggedPosts    map[string][]*Post
	TaggedKBs      map[string][]*KB
	TaggedTILs     map[string][]*TIL
	templates      map[string]*template.Template
	legacyComments map[string]*LegacyPost
}

func NewGenerator(config *Config, input afero.Fs, output afero.Fs, cache afero.Fs) *Generator {
	return &Generator{
		config:      config,
		Input:       input,
		Output:      output,
		Cache:       cache,
		Posts:       map[string]*Post{},
		KBs:         map[string]*KB{},
		Albums:      map[string]*Album{},
		TILs:        map[string]*TIL{},
		Pages:       map[string]*StandalonePage{},
		TaggedPosts: map[string][]*Post{},
		TaggedKBs:   map[string][]*KB{},
		TaggedTILs:  map[string][]*TIL{},
		templates:   map[string]*template.Template{},
	}
}

// themeTemplate returns the parsed template for the named file in the active theme,
// parsing and caching it on first use. Template read/parse failures panic; the panic
// is recovered in CompileMarkdown and surfaced as an error.
func (gen *Generator) themeTemplate(name string) *template.Template {
	if t, ok := gen.templates[name]; ok {
		return t
	}
	path := fmt.Sprintf("themes/%s/%s", gen.config.Site.Theme, name)
	raw, err := afero.ReadFile(gen.Input, path)
	if err != nil {
		panic(fmt.Errorf("reading template %s: %w", path, err))
	}
	t, err := template.New(name).Funcs(template.FuncMap{"now": time.Now}).Parse(string(raw))
	if err != nil {
		panic(fmt.Errorf("parsing template %s: %w", path, err))
	}
	gen.templates[name] = t
	return t
}

// renderTheme executes the named (cached) theme template with data.
func (gen *Generator) renderTheme(name string, data any) string {
	var buf strings.Builder
	if err := gen.themeTemplate(name).Execute(&buf, data); err != nil {
		panic(fmt.Errorf("executing template %s: %w", name, err))
	}
	return buf.String()
}

// renderInline parses and executes a one-off template string (e.g. post/KB content
// that may itself contain template directives). Not cached, since the string varies
// per item.
func (gen *Generator) renderInline(label, tmplStr string, data any) string {
	t, err := template.New(label).Funcs(template.FuncMap{"now": time.Now}).Parse(tmplStr)
	if err != nil {
		panic(fmt.Errorf("parsing %s template: %w", label, err))
	}
	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		panic(fmt.Errorf("executing %s template: %w", label, err))
	}
	return buf.String()
}

// chrome renders the header/nav/footer shared by every page, with the given data.
func (gen *Generator) chrome(data any) Page {
	return Page{
		Header: gen.renderTheme("header.html", data),
		Nav:    gen.renderTheme("nav.html", data),
		Footer: gen.renderTheme("footer.html", data),
	}
}

func (gen *Generator) sortedPosts() []*Post {
	posts := make([]*Post, 0, len(gen.Posts))
	for _, post := range gen.Posts {
		posts = append(posts, post)
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].FrontMatter.OrigDate.After(posts[j].FrontMatter.OrigDate)
	})
	return posts
}

func (gen *Generator) generatePost(post Post) string {
	page := gen.chrome(post)
	page.Posts = []*Post{&post}
	// The post body may itself contain template directives referencing the page.
	post.HTML = gen.renderInline("post content", post.HTML, page)
	return gen.renderTheme("post.html", page)
}

func (gen *Generator) generateIndex() string {
	page := gen.chrome(Post{})
	page.Posts = gen.sortedPosts()
	return gen.renderTheme("index.html", page)
}

func (gen *Generator) generateFront() string {
	page := gen.chrome(Post{})
	page.SingleKB = gen.KBs["front"]
	return gen.renderTheme("front.html", page)
}

func (gen *Generator) generatePage(standalone *StandalonePage) string {
	page := gen.chrome(Post{})
	page.SinglePage = standalone
	return gen.renderTheme("page.html", page)
}

func (gen *Generator) generateKB(kb KB) string {
	page := gen.chrome(kb)
	page.SingleKB = &kb
	// The KB body may itself contain template directives referencing the page.
	kb.HTML = gen.renderInline("kb content", kb.HTML, page)
	return gen.renderTheme("kb.html", page)
}

func (gen *Generator) generateKBIndex() string {
	kbPages := []*KB{}
	for _, kb := range gen.KBs {
		if kb.FrontMatter.Slug != "front" {
			kbPages = append(kbPages, kb)
		}
	}
	sort.Slice(kbPages, func(i, j int) bool {
		return kbPages[i].Path < kbPages[j].Path
	})

	page := gen.chrome(Post{})
	page.KBPages = kbPages
	page.Graph = gen.generateKBGraph()
	return gen.renderTheme("kb_index.html", page)
}

func (gen *Generator) generateGalleryIndex() string {
	albums := []*Album{}
	for _, album := range gen.Albums {
		albums = append(albums, album)
	}
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].Slug < albums[j].Slug
	})

	page := gen.chrome(Post{})
	page.GalleryAlbums = albums
	return gen.renderTheme("gallery_index.html", page)
}

func (gen *Generator) generateAlbumPage(album *Album) string {
	page := gen.chrome(Post{})
	page.SingleAlbum = album
	return gen.renderTheme("gallery_album.html", page)
}

func (gen *Generator) generateTILIndex() string {
	tils := []*TIL{}
	for _, til := range gen.TILs {
		tils = append(tils, til)
	}
	sort.Slice(tils, func(i, j int) bool {
		return tils[i].FrontMatter.OrigDate.After(tils[j].FrontMatter.OrigDate)
	})

	page := gen.chrome(Post{})
	page.TILList = tils
	return gen.renderTheme("til_index.html", page)
}

func (gen *Generator) generateTagsPage() string {
	allTags := map[string]int{}
	for tag, posts := range gen.TaggedPosts {
		allTags[tag] += len(posts)
	}
	for tag, kbs := range gen.TaggedKBs {
		allTags[tag] += len(kbs)
	}
	for tag, tils := range gen.TaggedTILs {
		allTags[tag] += len(tils)
	}

	page := gen.chrome(Post{})
	page.AllTags = allTags
	return gen.renderTheme("tags.html", page)
}

func (gen *Generator) generateTagList(tag string) string {
	page := gen.chrome(Post{})
	page.Posts = gen.TaggedPosts[tag]
	page.KBPages = gen.TaggedKBs[tag]
	page.TILList = gen.TaggedTILs[tag]
	page.CurrentTag = tag
	return gen.renderTheme("tag_list.html", page)
}

// Generate 5 line summary from the post html
func (gen *Generator) generateSummary(post *Post) {
	strippedHTML := strip.StripTags(post.HTML)
	splitHTML := strings.Split(strippedHTML, ".")
	if len(splitHTML) < 5 {
		post.Summary = strippedHTML
		return
	}
	post.Summary = strings.Join(splitHTML[:4], ".") + "..."
}

// baseURL is the site root used for absolute links in feeds.
func (gen *Generator) baseURL() string {
	domain := gen.config.Site.PrimaryDomain
	if domain == "" {
		domain = "localhost"
	}
	return "https://" + domain
}

func (gen *Generator) buildFeed(title, description string, items []*feeds.Item) *feeds.Feed {
	return &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: gen.baseURL()},
		Description: description,
		Author:      &feeds.Author{Name: "Yulian Kuncheff", Email: "yulian@kuncheff.com"},
		Copyright:   fmt.Sprintf("© <a href=\"https://github.com/daegalus\">Yulian Kuncheff</a> %d", time.Now().Year()),
		Created:     time.Now(),
		Language:    "en-us",
		Generator:   "yblog",
		Items:       items,
	}
}

func (gen *Generator) postItems() []*feeds.Item {
	var items []*feeds.Item
	for _, post := range gen.sortedPosts() {
		link := fmt.Sprintf("%s/blog/%s/", gen.baseURL(), post.FrontMatter.Slug)
		items = append(items, &feeds.Item{
			Id:          link,
			Title:       post.FrontMatter.Title,
			Link:        &feeds.Link{Href: link},
			Created:     post.FrontMatter.OrigDate,
			Description: post.Summary,
			Author:      &feeds.Author{Name: "Yulian Kuncheff", Email: "yulian@kuncheff.com"},
			Content:     post.HTML,
		})
	}
	return items
}

func (gen *Generator) tilItems() []*feeds.Item {
	tils := []*TIL{}
	for _, til := range gen.TILs {
		tils = append(tils, til)
	}
	sort.Slice(tils, func(i, j int) bool {
		return tils[i].FrontMatter.OrigDate.After(tils[j].FrontMatter.OrigDate)
	})

	link := fmt.Sprintf("%s/til/", gen.baseURL())
	var items []*feeds.Item
	for _, til := range tils {
		items = append(items, &feeds.Item{
			Id:      fmt.Sprintf("%s#%s", link, til.FrontMatter.Slug),
			Title:   til.FrontMatter.Title,
			Link:    &feeds.Link{Href: link},
			Created: til.FrontMatter.OrigDate,
			Author:  &feeds.Author{Name: "Yulian Kuncheff", Email: "yulian@kuncheff.com"},
			Content: til.HTML,
		})
	}
	return items
}

func (gen *Generator) kbItems() []*feeds.Item {
	kbs := []*KB{}
	for _, kb := range gen.KBs {
		if kb.FrontMatter.Slug != "front" {
			kbs = append(kbs, kb)
		}
	}
	sort.Slice(kbs, func(i, j int) bool {
		return kbs[i].FrontMatter.OrigDate.After(kbs[j].FrontMatter.OrigDate)
	})

	var items []*feeds.Item
	for _, kb := range kbs {
		link := fmt.Sprintf("%s/kb/%s/", gen.baseURL(), kb.Path)
		items = append(items, &feeds.Item{
			Id:      link,
			Title:   kb.FrontMatter.Title,
			Link:    &feeds.Link{Href: link},
			Created: kb.FrontMatter.OrigDate,
			Author:  &feeds.Author{Name: "Yulian Kuncheff", Email: "yulian@kuncheff.com"},
			Content: kb.HTML,
		})
	}
	return items
}

// writeFeedSet writes atom, rss and json files for the feed under the given base name
// (e.g. "feed" -> feed.atom/feed.rss/feed.json) in the output FS.
func (gen *Generator) writeFeedSet(name string, feed *feeds.Feed) {
	atom, err := feed.ToAtom()
	if err != nil {
		panic(fmt.Errorf("generating %s atom feed: %w", name, err))
	}
	rss, err := feed.ToRss()
	if err != nil {
		panic(fmt.Errorf("generating %s rss feed: %w", name, err))
	}
	json, err := (&feeds.JSON{Feed: feed}).ToJSON()
	if err != nil {
		panic(fmt.Errorf("generating %s json feed: %w", name, err))
	}
	afero.WriteFile(gen.Output, name+".atom", []byte(atom), fs.ModePerm)
	afero.WriteFile(gen.Output, name+".rss", []byte(rss), fs.ModePerm)
	afero.WriteFile(gen.Output, name+".json", []byte(json), fs.ModePerm)
}

// generateFeeds writes the combined blog feed (feed.*) plus per-type feeds
// (feed_blog.*, feed_til.*, feed_kb.*) to the output.
func (gen *Generator) generateFeeds() {
	blog := gen.buildFeed("Yulian Kuncheff — Blog", "my blog of thoughts and experiments", gen.postItems())
	gen.writeFeedSet("feed", blog)
	gen.writeFeedSet("feed_blog", blog)
	gen.writeFeedSet("feed_til", gen.buildFeed("Yulian Kuncheff — TIL", "Today I Learned", gen.tilItems()))
	gen.writeFeedSet("feed_kb", gen.buildFeed("Yulian Kuncheff — Knowledge Base", "Knowledge base / digital garden", gen.kbItems()))
	log.Info("Generated feeds")
}
