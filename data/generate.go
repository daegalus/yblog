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
	config *Config
	Input  afero.Fs
	Output afero.Fs
	Cache  afero.Fs
}

func NewGenerator(config *Config, input afero.Fs, output afero.Fs, cache afero.Fs) *Generator {
	return &Generator{
		config: config,
		Input:  input,
		Output: output,
		Cache:  cache,
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
	for _, post := range Posts {
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
		SingleKB: KBs["front"],
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
	for _, post := range Posts {
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
