package data

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"yblog/utils"
	"yblog/utils/goldmark_extensions"

	gmEmbed "github.com/13rac1/goldmark-embed"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/caarlos0/log"
	"github.com/spf13/afero"
	fences "github.com/stefanfritsch/goldmark-fences"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func (gen *Generator) CompileMarkdown() {
	log.Info("Compiling markdown")
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithGuessLanguage(true),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			&goldmark_extensions.Picture{},
			gmEmbed.New(),
			&fences.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	contentPath := "content/blog"
	files, _ := fs.ReadDir(Content, contentPath)
	for _, file := range files {
		data, err := Content.ReadFile(filepath.Join(contentPath, file.Name()))
		if err != nil {
			log.WithField("error", err).Fatal("Error reading file")
			panic(err)
		}

		content, _ := utils.StripFrontMatter(data)
		//filename := strings.Split(file.Name(), ".")[0]

		context := parser.NewContext()
		var out strings.Builder

		err = md.Convert(data, &out, parser.WithContext(context))
		if err != nil {
			log.WithField("error", err).Fatal("Error converting markdown")
			panic(err)
		}

		// if strings.Contains(file.Name(), "bloat") {
		// 	fmt.Println(out.String())
		// }

		metaData := meta.Get(context)

		tagsInterface := metaData["tags"].([]interface{})
		tagsStringList := make([]string, len(tagsInterface))
		for i, v := range tagsInterface {
			tagsStringList[i] = v.(string)
		}

		frontmatter := utils.FrontMatter{
			Author:   metaData["author"].(string),
			OrigDate: utils.ParseDate(metaData["date"].(string)),
			Date:     utils.ModifyDate(metaData["date"].(string)),
			Draft:    metaData["draft"].(bool),
			Slug:     metaData["slug"].(string),
			Title:    metaData["title"].(string),
			Type:     metaData["type"].(string),
			Tags:     tagsStringList,
		}

		post := Post{
			FrontMatter: frontmatter,
			Markdown:    content,
			HTML:        out.String(),
		}

		gen.generateSummary(&post)

		Posts[frontmatter.Slug] = &post

		for _, tag := range post.FrontMatter.Tags {
			TaggedPosts[tag] = append(TaggedPosts[tag], &post)

		}
		hashtags := []string{}
		for _, tag := range post.FrontMatter.Tags {
			hashtags = append(hashtags, fmt.Sprintf("#%s", tag))
		}
		post.Tagsline = strings.Join(hashtags, " | ")
		html := gen.generatePost(post)
		gen.Output.MkdirAll("posts", fs.ModeDir)
		afero.WriteFile(gen.Output, fmt.Sprintf("posts/%s.html", frontmatter.Slug), []byte(html), fs.ModePerm)
	}

	html := gen.generateIndex()
	afero.WriteFile(gen.Output, "index.html", []byte(html), fs.ModePerm)

	atom, rss, json := gen.generateFeeds()
	afero.WriteFile(gen.Output, "feed.atom", []byte(atom), fs.ModePerm)
	afero.WriteFile(gen.Output, "feed.rss", []byte(rss), fs.ModePerm)
	afero.WriteFile(gen.Output, "feed.json", []byte(json), fs.ModePerm)
}
