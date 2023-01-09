package data

import (
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/caarlos0/log"

	"github.com/daegalus/feeds"
	strip "github.com/grokify/html-strip-tags-go"
)

func generateTemplateHTML(file string, post Post) string {
	header, err := Content.ReadFile(file)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "file": file}).Fatal("Error reading template")
		panic(err)
	}

	tmpl, err := template.New("header").Parse(string(header))
	if err != nil {
		log.WithFields(log.Fields{"error": err, "file": file}).Fatal("Error parsing header template")
		panic(err)
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, post)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "file": file}).Fatal("Error executing header template")
		panic(err)
	}
	return buf.String()
}

func GeneratePage(post Post) string {
	postHTML, err := Content.ReadFile("content/templates/post.html")
	if err != nil {
		log.WithField("error", err).Fatal("Error reading post template")
		panic(err)
	}

	page := Page{
		Header: generateTemplateHTML("content/templates/header.html", post),
		Nav:    generateTemplateHTML("content/templates/nav.html", post),
		Posts:  []*Post{&post},
		Footer: generateTemplateHTML("content/templates/footer.html", post),
	}

	tmpl, err := template.New("page").Parse(string(postHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing post template")
		panic(err)
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, page)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing post template")
		panic(err)
	}
	return buf.String()
}

func GenerateIndex() string {
	indexHTML, err := Content.ReadFile("content/templates/index.html")
	if err != nil {
		log.WithField("error", err).Fatal("Error reading index template")
		panic(err)
	}

	local_posts := []*Post{}
	for _, post := range Posts {
		local_posts = append(local_posts, post)
	}

	sort.Slice(local_posts, func(i, j int) bool {
		return local_posts[i].FrontMatter.OrigDate.After(local_posts[j].FrontMatter.OrigDate)
	})

	index := Page{
		Header: generateTemplateHTML("content/templates/header.html", Post{}),
		Nav:    generateTemplateHTML("content/templates/nav.html", Post{}),
		Posts:  local_posts,
		Footer: generateTemplateHTML("content/templates/footer.html", Post{}),
	}

	tmpl, err := template.New("index").Parse(string(indexHTML))
	if err != nil {
		log.WithField("error", err).Fatal("Error parsing index template")
		panic(err)
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, index)
	if err != nil {
		log.WithField("error", err).Fatal("Error executing index template")
		panic(err)
	}
	return buf.String()
}

// Generate 5 line summary from the post html
func generateSummary(post *Post) {
	strippedHTML := strip.StripTags(post.HTML)
	splitHTML := strings.Split(strippedHTML, ".")
	if len(splitHTML) < 5 {
		post.Summary = strippedHTML
		return
	}
	post.Summary = strings.Join(splitHTML[:4], ".") + "..."
}

func generateFeeds() (string, string, string) {
	feed := &feeds.Feed{
		Title:       "my blog of thoughts and experiments",
		Link:        &feeds.Link{Href: "https://yulian.kuncheff.com"},
		Description: "Yulian Kuncheff's Blog",
		Author:      &feeds.Author{Name: "Yulian Kuncheff", Email: "yulian@kuncheff.com"},
		Copyright:   "© <a href=\"https://github.com/daegalus\">Yulian Kuncheff</a> 2022",
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
			Id:          fmt.Sprintf("https://yulian.kuncheff.com/%s/", post.FrontMatter.Slug),
			Title:       post.FrontMatter.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://yulian.kuncheff.com/%s/", post.FrontMatter.Slug)},
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

	// jsonObj := &feeds.JSON{feed}
	// jsonFeedObj := jsonObj.JSONFeed()
	// jsonFeedObj.Extensions = append(jsonFeedObj.Extensions, &feeds.JSONExtensions{
	// 	Key:   "Generator",
	// 	Value: "yblog",
	// })
	json, err := feed.ToJSON()
	if err != nil {
		log.WithField("error", err).Fatal("Error generating json feed")
	}

	fmt.Println(atom[0:100])
	return atom, rss, json
}
