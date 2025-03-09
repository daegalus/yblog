package data

import (
	"fmt"
	"image"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
	"yblog/cache"
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

func encodeFilesInDir(file fs.DirEntry, imageStoragePath string, postImageDir fs.DirEntry, gen *Generator, wg *sync.WaitGroup) {
	defer wg.Done()
	filename := strings.Split(file.Name(), ".")[0]
	extension := strings.Split(file.Name(), ".")[1]
	if extension != "png" {
		return
	}
	postImageDirPath := filepath.Join(ContentPrefix, imageStoragePath, postImageDir.Name())
	postImagePath := filepath.Join(postImageDirPath, file.Name())
	avifPath := filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "avif"))
	jxlPath := filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "jxl"))
	webpPath := filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "webp"))
	pngPath := filepath.Join(imageStoragePath, postImageDir.Name(), file.Name())

	avifExists, _ := afero.Exists(gen.Output, avifPath)
	jxlExists, _ := afero.Exists(gen.Output, jxlPath)
	webpExists, _ := afero.Exists(gen.Output, webpPath)
	pngExists, _ := afero.Exists(gen.Output, pngPath)

	imageByteData, err := afero.ReadFile(gen.Input, postImagePath)
	if err != nil {
		log.WithField("error", err).Fatal("Error reading image")
	}

	if !pngExists {
		log.WithField("file", postImagePath).Info("PNG Doesn't Exist, writing data to output")
		err := afero.WriteFile(gen.Output, pngPath, imageByteData, fs.ModePerm)
		if err != nil {
			log.WithField("error", err).Fatal("Error writing png to output")
		}
	}

	savedImageHashes, _ := cache.LoadImageHashes()
	currentHashes, _ := cache.CalculateImageHashes()

	imageDifferent := currentHashes[postImagePath] != savedImageHashes[postImagePath]

	if !imageDifferent {
		if avifExists && jxlExists && webpExists && pngExists {
			log.WithField("file", postImagePath).Info("Skipping encoding")
			return // Skip encoding if the image is the same
		}
	}

	imageData := utils.ImageFromPNG(afero.FromIOFS{FS: Content}, postImagePath)

	var avifChan, jxlChan, webpChan chan []byte
	if !avifExists || imageDifferent {
		wg.Add(1)
		avifChan = make(chan []byte)
		go threadEncode(avifChan, imageData, wg, utils.ImageToAVIF)
	}
	if !jxlExists || imageDifferent {
		wg.Add(1)
		jxlChan = make(chan []byte)
		go threadEncode(jxlChan, imageData, wg, utils.ImageToJXL)
	}
	if !webpExists || imageDifferent {
		wg.Add(1)
		webpChan = make(chan []byte)
		go threadEncode(webpChan, imageData, wg, utils.ImageToWebP)
	}

	for range 3 {
		select {
		case avifBytes, _ := <-avifChan:
			log.WithField("file", filepath.Join("images", postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "avif"))).Info("Converting to AVIF")
			afero.WriteFile(gen.Output, filepath.Join("images", postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "avif")), avifBytes, fs.ModePerm)
		case jxlBytes, _ := <-jxlChan:
			log.WithField("file", filepath.Join("images", postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "jxl"))).Info("Converting to JXL")
			afero.WriteFile(gen.Output, filepath.Join("images", postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "jxl")), jxlBytes, fs.ModePerm)
		case webpBytes, _ := <-webpChan:
			log.WithField("file", filepath.Join("images", postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "webp"))).Info("Converting to WebP")
			afero.WriteFile(gen.Output, filepath.Join("images", postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "webp")), webpBytes, fs.ModePerm)
		}
	}
}

func encodeFilesInDirSync(file fs.DirEntry, imageStoragePath string, postImageDir fs.DirEntry, gen *Generator) {
	filename := strings.Split(file.Name(), ".")[0]
	extension := strings.Split(file.Name(), ".")[1]
	if extension != "png" {
		return
	}
	postImageDirPath := filepath.Join(ContentPrefix, imageStoragePath, postImageDir.Name())
	postImagePath := filepath.Join(postImageDirPath, file.Name())
	avifPath := filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "avif"))
	jxlPath := filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "jxl"))
	webpPath := filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "webp"))
	pngPath := filepath.Join(imageStoragePath, postImageDir.Name(), file.Name())
	log.
		WithField("avif", avifPath).
		WithField("jxl", jxlPath).
		WithField("webp", webpPath).
		WithField("png", pngPath).
		Info("Images in Post")
	avifExists, _ := afero.Exists(gen.Output, avifPath)
	jxlExists, _ := afero.Exists(gen.Output, jxlPath)
	webpExists, _ := afero.Exists(gen.Output, webpPath)
	pngExists, _ := afero.Exists(gen.Output, pngPath)

	imageByteData, err := afero.ReadFile(gen.Input, postImagePath)
	if err != nil {
		log.WithField("error", err).Fatal("Error reading image")
	}

	if !pngExists {
		log.WithField("file", postImagePath).Info("PNG Doesn't Exist, writing data to output")
		err := afero.WriteFile(gen.Output, pngPath, imageByteData, fs.ModePerm)
		if err != nil {
			log.WithField("error", err).Fatal("Error writing png to output")
		}
	}

	savedImageHashes, _ := cache.LoadImageHashes()
	currentHashes, _ := cache.CalculateImageHashes()

	imageDifferent := currentHashes[postImagePath] != savedImageHashes[postImagePath]

	if !imageDifferent {
		if avifExists && jxlExists && webpExists && pngExists {
			log.WithField("file", postImagePath).Info("Skipping encoding")
			return // Skip encoding if the image is the same
		}
	}

	imageData := utils.ImageFromPNG(afero.FromIOFS{FS: Content}, postImagePath)

	if !avifExists || imageDifferent {
		avifData := utils.ImageToAVIF(imageData)
		log.WithField("file", filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "avif"))).Info("Converting to AVIF")
		afero.WriteFile(gen.Output, filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "avif")), avifData, fs.ModePerm)
	}

	if !jxlExists || imageDifferent {
		jxlData := utils.ImageToJXL(imageData)
		log.WithField("file", filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "jxl"))).Info("Converting to JXL")
		afero.WriteFile(gen.Output, filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "jxl")), jxlData, fs.ModePerm)
	}

	if !webpExists || imageDifferent {
		webpData := utils.ImageToWebP(imageData)
		log.WithField("file", filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "webp"))).Info("Converting to WebP")
		afero.WriteFile(gen.Output, filepath.Join(imageStoragePath, postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "webp")), webpData, fs.ModePerm)
	}
}

func encodeAllDirs(imageStoragePath string, imageStorageDirEntry fs.DirEntry, gen *Generator, wg *sync.WaitGroup) {
	defer wg.Done()
	log.WithField("dir", imageStorageDirEntry.Name()).Info("Converting in Post")
	if imageStorageDirEntry.IsDir() {
		postImageDir, _ := fs.ReadDir(Content, filepath.Join(ContentPrefix, imageStoragePath, imageStorageDirEntry.Name()))
		wg.Add(len(postImageDir))
		for _, postImageFile := range postImageDir {
			go encodeFilesInDir(postImageFile, imageStoragePath, imageStorageDirEntry, gen, wg)
		}
	}
}

func encodeAllDirsSync(imageStoragePath string, imageStorageDirEntry fs.DirEntry, gen *Generator) {
	log.WithField("dir", imageStorageDirEntry.Name()).Info("Converting in Post")
	if imageStorageDirEntry.IsDir() {
		postImageDir, _ := fs.ReadDir(Content, filepath.Join(ContentPrefix, imageStoragePath, imageStorageDirEntry.Name()))
		for _, postImageFile := range postImageDir {
			encodeFilesInDirSync(postImageFile, imageStoragePath, imageStorageDirEntry, gen)
		}
	}
}

func (gen *Generator) CompileMarkdown() {
	log.Info("Converting Images")
	imageStoragePath := "images"

	utils.CopyFiles(afero.FromIOFS{FS: Content}, gen.Output, filepath.Join(ContentPrefix, imageStoragePath), ContentPrefix, "")

	if ok, _ := afero.Exists(gen.Cache, CachePrefix); ok {
		if ok, _ := afero.Exists(gen.Cache, filepath.Join(CachePrefix, imageStoragePath)); ok {
			utils.CopyFiles(gen.Cache, gen.Output, filepath.Join(CachePrefix, imageStoragePath), CachePrefix, "")
		}
	}
	utils.CopyFiles(afero.FromIOFS{FS: Content}, gen.Cache, filepath.Join(ContentPrefix, imageStoragePath), ContentPrefix, CachePrefix)

	imageStorageDirEntries, _ := fs.ReadDir(Content, filepath.Join(ContentPrefix, imageStoragePath))
	wg := sync.WaitGroup{}
	wg.Add(len(imageStorageDirEntries))
	for _, imageStorageDirEntry := range imageStorageDirEntries {
		go encodeAllDirs(imageStoragePath, imageStorageDirEntry, gen, &wg)
		//encodeAllDirsSync(imageStoragePath, imageStorageDirEntry, gen)
	}
	wg.Wait()

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
	contentPathBlog := "content/blog"
	files, _ := fs.ReadDir(Content, contentPathBlog)
	for _, file := range files {
		gen.CompileMarkdownFile(file, contentPathBlog, "post", md)
	}

	html := gen.generateIndex()
	afero.WriteFile(gen.Output, "blog/index.html", []byte(html), fs.ModePerm)

	contentPathKB := "content/kb"
	files, _ = fs.ReadDir(Content, contentPathKB)
	for _, file := range files {
		gen.CompileMarkdownFile(file, contentPathKB, "kb", md)
	}

	atom, rss, json := gen.generateFeeds()
	afero.WriteFile(gen.Output, "feed.atom", []byte(atom), fs.ModePerm)
	afero.WriteFile(gen.Output, "feed.rss", []byte(rss), fs.ModePerm)
	afero.WriteFile(gen.Output, "feed.json", []byte(json), fs.ModePerm)
}

func sortComments(legacyComments []*LegacyComment) []*LegacyComment {
	slices.SortFunc(legacyComments, func(i, j *LegacyComment) int {
		iCA, _ := time.Parse(time.RFC3339, i.CreatedAt+"Z")
		jCA, _ := time.Parse(time.RFC3339, j.CreatedAt+"Z")
		return int(iCA.UnixMilli() - jCA.UnixMilli())
	})
	return legacyComments
}

func sortChildrenComments(legacyComments []*LegacyComment) []*LegacyComment {
	legacyComments = sortComments(legacyComments)
	for _, comment := range legacyComments {
		comment.Children = sortChildrenComments(comment.Children)
	}
	return legacyComments
}

func (gen *Generator) CompileMarkdownFile(file fs.DirEntry, contentPath string, contentType string, md goldmark.Markdown) {
	if !strings.Contains(file.Name(), ".md") {
		return
	}
	data, err := Content.ReadFile(filepath.Join(contentPath, file.Name()))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading file")
		panic(err)
	}

	content, _ := utils.StripFrontMatter(data)

	context := parser.NewContext()
	var out strings.Builder

	err = md.Convert(data, &out, parser.WithContext(context))
	if err != nil {
		log.WithField("error", err).Fatal("Error converting markdown")
		panic(err)
	}

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

	switch contentType {
	case "post":
		gen.postSpecific(frontmatter, content, out.String())
	case "kb":
		gen.kbSpecific(frontmatter, content, out.String())
	}

}

func (gen *Generator) postSpecific(frontmatter utils.FrontMatter, content []byte, html string) {
	legacyPostComments := GetLegacyComments()

	legacyComments := []*LegacyComment{}
	if _, ok := (*legacyPostComments)[frontmatter.Slug]; ok {
		for _, comment := range *(*legacyPostComments)[frontmatter.Slug].Comments {
			legacyComments = append(legacyComments, comment)
		}
	}

	legacyComments = sortChildrenComments(legacyComments)

	post := Post{
		FrontMatter:    frontmatter,
		Markdown:       content,
		HTML:           html,
		LegacyComments: legacyComments,
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

	htmlOut := gen.generatePost(post)

	gen.Output.MkdirAll("posts", fs.ModeDir)
	afero.WriteFile(gen.Output, fmt.Sprintf("blog/posts/%s.html", frontmatter.Slug), []byte(htmlOut), fs.ModePerm)
	gen.Output.MkdirAll(fmt.Sprintf("$s", frontmatter.Slug), fs.ModeDir)
	afero.WriteFile(gen.Output, fmt.Sprintf("blog/%s/index.html", frontmatter.Slug), []byte(htmlOut), fs.ModePerm)
}

func (gen *Generator) kbSpecific(frontmatter utils.FrontMatter, content []byte, html string) {

	kb := KB{
		FrontMatter: frontmatter,
		Markdown:    content,
		HTML:        html,
	}

	KBs[frontmatter.Slug] = &kb

	for _, tag := range kb.FrontMatter.Tags {
		TaggedKBs[tag] = append(TaggedKBs[tag], &kb)

	}
	hashtags := []string{}
	for _, tag := range kb.FrontMatter.Tags {
		hashtags = append(hashtags, fmt.Sprintf("#%s", tag))
	}
	kb.Tagsline = strings.Join(hashtags, " | ")

	var htmlOut string
	if kb.FrontMatter.Slug == "front" {
		htmlOut = gen.generateFront()
		afero.WriteFile(gen.Output, "index.html", []byte(htmlOut), fs.ModePerm)
	} else {
		htmlOut = gen.generateKB(kb)
		gen.Output.MkdirAll("kb", fs.ModeDir)
		afero.WriteFile(gen.Output, fmt.Sprintf("kb/pages/%s.html", frontmatter.Slug), []byte(htmlOut), fs.ModePerm)
		gen.Output.MkdirAll(fmt.Sprintf("$s", frontmatter.Slug), fs.ModeDir)
		afero.WriteFile(gen.Output, fmt.Sprintf("kb/%s/index.html", frontmatter.Slug), []byte(htmlOut), fs.ModePerm)
	}

}
