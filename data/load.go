package data

import (
	"fmt"
	"image"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
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

func threadEncode(returnChannel chan []byte, imageData image.Image, wg *sync.WaitGroup, encodeFunc func(image.Image) []byte) {
	defer wg.Done()
	encodedImage := encodeFunc(imageData)
	returnChannel <- encodedImage
}

func encodeFilesInDir(file fs.DirEntry, imagePath string, imageDirEntry fs.DirEntry, gen *Generator, wg *sync.WaitGroup) {
	defer wg.Done()
	filename := strings.Split(file.Name(), ".")[0]
	imageData := utils.ImageFromPNG(Content, filepath.Join(imagePath, imageDirEntry.Name(), file.Name()))

	wg.Add(2)
	avifChan := make(chan []byte)
	jxlChan := make(chan []byte)
	webpChan := make(chan []byte)

	go threadEncode(avifChan, imageData, wg, utils.ImageToAVIF)
	go threadEncode(jxlChan, imageData, wg, utils.ImageToJXL)
	go threadEncode(webpChan, imageData, wg, utils.ImageToWebP)

	for i := 0; i < 2; i++ {
		select {
		case avifBytes := <-avifChan:
			log.WithField("file", filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "avif"))).Info("Converting to AVIF")
			afero.WriteFile(gen.Output, filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "avif")), avifBytes, fs.ModePerm)
		case jxlBytes := <-jxlChan:
			log.WithField("file", filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "jxl"))).Info("Converting to JXL")
			afero.WriteFile(gen.Output, filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "jxl")), jxlBytes, fs.ModePerm)
		case webpBytes := <-webpChan:
			log.WithField("file", filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "webp"))).Info("Converting to WebP")
			afero.WriteFile(gen.Output, filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "webp")), webpBytes, fs.ModePerm)
		}
	}
}

func encodeFilesInDirSync(file fs.DirEntry, imagePath string, imageDirEntry fs.DirEntry, gen *Generator) {
	filename := strings.Split(file.Name(), ".")[0]
	imageData := utils.ImageFromPNG(Content, filepath.Join(imagePath, imageDirEntry.Name(), file.Name()))

	avifData := utils.ImageToAVIF(imageData)
	jxlData := utils.ImageToJXL(imageData)
	webpData := utils.ImageToWebP(imageData)

	log.WithField("file", filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "avif"))).Info("Converting to AVIF")
	afero.WriteFile(gen.Output, filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "avif")), avifData, fs.ModePerm)

	log.WithField("file", filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "jxl"))).Info("Converting to JXL")
	afero.WriteFile(gen.Output, filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "jxl")), jxlData, fs.ModePerm)

	log.WithField("file", filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "webp"))).Info("Converting to WebP")
	afero.WriteFile(gen.Output, filepath.Join("images", imageDirEntry.Name(), fmt.Sprintf("%s.%s", filename, "webp")), webpData, fs.ModePerm)
}

func encodeAllDirs(imagePath string, imageDirEntry fs.DirEntry, gen *Generator, wg *sync.WaitGroup) {
	defer wg.Done()
	log.WithField("dir", imageDirEntry.Name()).Info("Converting in Post")
	if imageDirEntry.IsDir() {
		postImageDir, _ := fs.ReadDir(Content, filepath.Join(imagePath, imageDirEntry.Name()))
		wg.Add(len(postImageDir))
		for _, file := range postImageDir {
			go encodeFilesInDir(file, imagePath, imageDirEntry, gen, wg)
		}
	}
}

func encodeAllDirsSync(imagePath string, imageDirEntry fs.DirEntry, gen *Generator) {
	log.WithField("dir", imageDirEntry.Name()).Info("Converting in Post")
	if imageDirEntry.IsDir() {
		postImageDir, _ := fs.ReadDir(Content, filepath.Join(imagePath, imageDirEntry.Name()))
		for _, file := range postImageDir {
			encodeFilesInDirSync(file, imagePath, imageDirEntry, gen)
		}
	}
}

func (gen *Generator) CompileMarkdown() {
	log.Info("Converting Images")
	imagePath := "content/images"
	imageDirEntries, _ := fs.ReadDir(Content, imagePath)
	//var wg sync.WaitGroup
	//wg.Add(len(imageDirEntries))
	for _, imageDirEntry := range imageDirEntries {
		//go encodeAllDirs(imagePath, imageDirEntry, gen, &wg)
		encodeAllDirsSync(imagePath, imageDirEntry, gen)
	}

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
