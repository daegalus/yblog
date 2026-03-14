package data

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/caarlos0/log"

	"github.com/daegalus/feeds"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/spf13/afero"
)

type Generator struct {
	config      *Config
	Input       afero.Fs
	Output      afero.Fs
	Cache       afero.Fs
	Posts       map[string]*Post
	KBs         map[string]*KB
	Albums      map[string]*Album
	TILs        map[string]*TIL
	Pages       map[string]*StandalonePage
	TaggedPosts map[string][]*Post
	TaggedKBs   map[string][]*KB
	TaggedTILs  map[string][]*TIL
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
	}
}

func generateTemplateHTML[T Post | KB](fs afero.Fs, file string, in T) string {
	templateString, err := afero.ReadFile(fs, file)
	if err != nil {
		log.WithField("error", err).WithField("file", file).Fatal("Error read template file")
	}

	tmpl, err := template.New(file).Funcs(template.FuncMap{"now": time.Now}).Parse(string(templateString))
	if err != nil {
		log.WithField("error", err).WithField("file", file).Fatal("Error parsing template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, in)
	if err != nil {
		log.WithField("error", err).WithField("file", file).Fatal("Error executing template")
	}
	return buf.String()
}

func (gen *Generator) generatePost(post Post) string {
	postHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/post.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error read post template")
	}

	page := Page{
		Header: generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), post),
		Nav:    generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), post),
		Posts:  []*Post{&post},
		Footer: generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), post),
	}

	// if strings.Contains(post.FrontMatter.Slug, "bloat") {
	// 	log.WithField("post", post.HTML).Info("Generated post")
	// }

	uttmpl, _ := template.New("utPage").Funcs(template.FuncMap{"now": time.Now}).Parse(post.HTML)
	var templatedOut strings.Builder
	err = uttmpl.Execute(&templatedOut, page)

	post.HTML = templatedOut.String()

	tmpl, err := template.New("page").Funcs(template.FuncMap{"now": time.Now}).Parse(string(postHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing post template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, page)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing post template")
	}
	return buf.String()
}

func (gen *Generator) generateIndex() string {
	indexHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/index.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading index template")
	}

	local_posts := []*Post{}
	for _, post := range gen.Posts {
		local_posts = append(local_posts, post)
	}

	sort.Slice(local_posts, func(i, j int) bool {
		return local_posts[i].FrontMatter.OrigDate.After(local_posts[j].FrontMatter.OrigDate)
	})

	index := Page{
		Header: generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), Post{}),
		Nav:    generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), Post{}),
		Posts:  local_posts,
		Footer: generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), Post{}),
	}

	tmpl, err := template.New("index").Funcs(template.FuncMap{"now": time.Now}).Parse(string(indexHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing index template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, index)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing index template")
	}
	return buf.String()
}

func (gen *Generator) generateFront() string {
	frontHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/front.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading index template")
	}

	page := Page{
		Header:   generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), Post{}),
		Nav:      generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), Post{}),
		SingleKB: gen.KBs["front"],
		Footer:   generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), Post{}),
	}

	tmpl, err := template.New("front").Funcs(template.FuncMap{"now": time.Now}).Parse(string(frontHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing index template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, page)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing index template")
	}
	return buf.String()
}

func (gen *Generator) generatePage(page *StandalonePage) string {
	pageHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/page.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading page template")
	}

	p := Page{
		Header:     generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), Post{}),
		Nav:        generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), Post{}),
		SinglePage: page,
		Footer:     generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), Post{}),
	}

	tmpl, err := template.New("page").Funcs(template.FuncMap{"now": time.Now}).Parse(string(pageHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing page template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, p)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing page template")
	}
	return buf.String()
}

func (gen *Generator) generateKB(kb KB) string {
	kbHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/kb.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error read post template")
	}

	page := Page{
		Header:   generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), kb),
		Nav:      generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), kb),
		SingleKB: &kb,
		Footer:   generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), kb),
	}

	// if strings.Contains(post.FrontMatter.Slug, "bloat") {
	// 	log.WithField("post", post.HTML).Info("Generated post")
	// }

	uttmpl, _ := template.New("utPage").Funcs(template.FuncMap{"now": time.Now}).Parse(kb.HTML)
	var templatedOut strings.Builder
	err = uttmpl.Execute(&templatedOut, page)

	kb.HTML = templatedOut.String()

	tmpl, err := template.New("page").Funcs(template.FuncMap{"now": time.Now}).Parse(string(kbHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing post template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, page)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing post template")
	}
	return buf.String()
}

func (gen *Generator) generateKBIndex() string {
	kbIndexHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/kb_index.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading kb_index template")
	}

	kbPages := []*KB{}
	for _, kb := range gen.KBs {
		if kb.FrontMatter.Slug != "front" {
			kbPages = append(kbPages, kb)
		}
	}

	sort.Slice(kbPages, func(i, j int) bool {
		return kbPages[i].Path < kbPages[j].Path
	})

	page := Page{
		Header:  generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), Post{}),
		Nav:     generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), Post{}),
		KBPages: kbPages,
		Graph:   gen.generateKBGraph(),
		Footer:  generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), Post{}),
	}

	tmpl, err := template.New("kb_index").Funcs(template.FuncMap{"now": time.Now}).Parse(string(kbIndexHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing kb_index template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, page)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing kb_index template")
	}
	return buf.String()
}

func (gen *Generator) generateGalleryIndex() string {
	galleryHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/gallery_index.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading gallery_index template")
	}

	albums := []*Album{}
	for _, album := range gen.Albums {
		albums = append(albums, album)
	}
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].Slug < albums[j].Slug
	})

	page := Page{
		Header:        generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), Post{}),
		Nav:           generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), Post{}),
		GalleryAlbums: albums,
		Footer:        generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), Post{}),
	}

	tmpl, err := template.New("gallery_index").Funcs(template.FuncMap{"now": time.Now}).Parse(string(galleryHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing gallery_index template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, page)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing gallery_index template")
	}
	return buf.String()
}

func (gen *Generator) generateAlbumPage(album *Album) string {
	albumHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/gallery_album.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading gallery_album template")
	}

	page := Page{
		Header:      generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), Post{}),
		Nav:         generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), Post{}),
		SingleAlbum: album,
		Footer:      generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), Post{}),
	}

	tmpl, err := template.New("gallery_album").Funcs(template.FuncMap{"now": time.Now}).Parse(string(albumHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing gallery_album template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, page)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing gallery_album template")
	}
	return buf.String()
}

func (gen *Generator) generateTILIndex() string {
	tilHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/til_index.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading til_index template")
	}

	tils := []*TIL{}
	for _, til := range gen.TILs {
		tils = append(tils, til)
	}
	sort.Slice(tils, func(i, j int) bool {
		return tils[i].FrontMatter.OrigDate.After(tils[j].FrontMatter.OrigDate)
	})

	page := Page{
		Header:  generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), Post{}),
		Nav:     generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), Post{}),
		TILList: tils,
		Footer:  generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), Post{}),
	}

	tmpl, err := template.New("til_index").Funcs(template.FuncMap{"now": time.Now}).Parse(string(tilHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing til_index template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, page)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing til_index template")
	}
	return buf.String()
}

func (gen *Generator) generateTagsPage() string {
	tagsHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/tags.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading tags template")
	}

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

	page := Page{
		Header:  generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), Post{}),
		Nav:     generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), Post{}),
		AllTags: allTags,
		Footer:  generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), Post{}),
	}

	tmpl, err := template.New("tags").Funcs(template.FuncMap{"now": time.Now}).Parse(string(tagsHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing tags template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, page)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing tags template")
	}
	return buf.String()
}

func (gen *Generator) generateTagList(tag string) string {
	tagListHTML, err := afero.ReadFile(gen.Input, fmt.Sprintf("themes/%s/tag_list.html", gen.config.Site.Theme))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading tag_list template")
	}

	page := Page{
		Header:     generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/header.html", gen.config.Site.Theme), Post{}),
		Nav:        generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/nav.html", gen.config.Site.Theme), Post{}),
		Posts:      gen.TaggedPosts[tag],
		KBPages:    gen.TaggedKBs[tag],
		TILList:    gen.TaggedTILs[tag],
		CurrentTag: tag,
		Footer:     generateTemplateHTML(gen.Input, fmt.Sprintf("themes/%s/footer.html", gen.config.Site.Theme), Post{}),
	}

	tmpl, err := template.New("tag_list").Funcs(template.FuncMap{"now": time.Now}).Parse(string(tagListHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing tag_list template")
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, page)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing tag_list template")
	}
	return buf.String()
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

func (gen *Generator) generateFeeds() (string, string, string) {
	feed := &feeds.Feed{
		Title:       "my blog of thoughts and experiments",
		Link:        &feeds.Link{Href: "https://yulian.kuncheff.com"},
		Description: "Yulian Kuncheff's Blog",
		Author:      &feeds.Author{Name: "Yulian Kuncheff", Email: "yulian@kuncheff.com"},
		Copyright:   fmt.Sprintf("© <a href=\"https://github.com/daegalus\">Yulian Kuncheff</a> %s", strconv.Itoa(time.Now().Year())),
		Created:     time.Now(),
		Language:    "en-us",
		Generator:   "yblog",
	}

	var items []*feeds.Item

	var posts []*Post
	for _, post := range gen.Posts {
		posts = append(posts, post)
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].FrontMatter.OrigDate.After(posts[j].FrontMatter.OrigDate)
	})

	for _, post := range posts {
		item := &feeds.Item{
			Id:          fmt.Sprintf("https://yulian.kuncheff.com/blog/%s/", post.FrontMatter.Slug),
			Title:       post.FrontMatter.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://yulian.kuncheff.com/blog/%s/", post.FrontMatter.Slug)},
			Created:     post.FrontMatter.OrigDate,
			Description: post.Summary,
			Author:      &feeds.Author{Name: "Yulian Kuncheff", Email: "yulian@kuncheff.com"},
			Content:     post.HTML,
		}
		items = append(items, item)
	}

	feed.Items = items

	atom, err := feed.ToAtom()
	if err != nil {
		log.WithField("error", err).Fatal("Error generating atom feed")
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.WithField("error", err).Fatal("Error generating rss feed")
	}

	jsonObj := &feeds.JSON{Feed: feed}
	jsonFeedObj := jsonObj.JSONFeed()
	jsonFeedObj.Extensions = append(jsonFeedObj.Extensions, &feeds.JSONExtensions{
		Key:   "Generator",
		Value: "yblog",
	})
	json, err := jsonObj.ToJSON()
	if err != nil {
		log.WithField("error", err).Fatal("Error generating json feed")
	}

	return atom, rss, json
}
