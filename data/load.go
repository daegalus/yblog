package data

import (
	"fmt"
	"image"
	"io/fs"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
	"yblog/cache"
	"yblog/utils"
	"yblog/utils/goldmark_extensions"

	"github.com/BurntSushi/toml"

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

	osFs := afero.NewOsFs()
	localImagesDir := filepath.Join("data", ContentPrefix, imageStoragePath, postImageDir.Name())

	for range 3 {
		select {
		case avifBytes, _ := <-avifChan:
			outPath := filepath.Join("images", postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "avif"))
			log.WithField("file", outPath).Info("Converting to AVIF")
			afero.WriteFile(gen.Output, outPath, avifBytes, fs.ModePerm)
			if ok, _ := afero.Exists(osFs, localImagesDir); ok {
				localPath := filepath.Join(localImagesDir, fmt.Sprintf("%s.%s", filename, "avif"))
				if exists, _ := afero.Exists(osFs, localPath); !exists {
					afero.WriteFile(osFs, localPath, avifBytes, fs.ModePerm)
				}
			}
		case jxlBytes, _ := <-jxlChan:
			outPath := filepath.Join("images", postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "jxl"))
			log.WithField("file", outPath).Info("Converting to JXL")
			afero.WriteFile(gen.Output, outPath, jxlBytes, fs.ModePerm)
			if ok, _ := afero.Exists(osFs, localImagesDir); ok {
				localPath := filepath.Join(localImagesDir, fmt.Sprintf("%s.%s", filename, "jxl"))
				if exists, _ := afero.Exists(osFs, localPath); !exists {
					afero.WriteFile(osFs, localPath, jxlBytes, fs.ModePerm)
				}
			}
		case webpBytes, _ := <-webpChan:
			outPath := filepath.Join("images", postImageDir.Name(), fmt.Sprintf("%s.%s", filename, "webp"))
			log.WithField("file", outPath).Info("Converting to WebP")
			afero.WriteFile(gen.Output, outPath, webpBytes, fs.ModePerm)
			if ok, _ := afero.Exists(osFs, localImagesDir); ok {
				localPath := filepath.Join(localImagesDir, fmt.Sprintf("%s.%s", filename, "webp"))
				if exists, _ := afero.Exists(osFs, localPath); !exists {
					afero.WriteFile(osFs, localPath, webpBytes, fs.ModePerm)
				}
			}
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

func encodeAllDirs(imageStoragePath string, imageStorageDirEntry fs.DirEntry, gen *Generator, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	log.WithField("dir", imageStorageDirEntry.Name()).Info("Converting in Post")
	if imageStorageDirEntry.IsDir() {
		postImageDir, _ := fs.ReadDir(Content, filepath.Join(ContentPrefix, imageStoragePath, imageStorageDirEntry.Name()))
		for _, postImageFile := range postImageDir {
			wg.Add(1)
			go func(file fs.DirEntry) {
				sem <- struct{}{}
				defer func() { <-sem }()
				encodeFilesInDir(file, imageStoragePath, imageStorageDirEntry, gen, wg)
			}(postImageFile)
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
	sem := make(chan struct{}, runtime.NumCPU())
	wg.Add(len(imageStorageDirEntries))
	for _, imageStorageDirEntry := range imageStorageDirEntries {
		go encodeAllDirs(imageStoragePath, imageStorageDirEntry, gen, &wg, sem)
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
			&goldmark_extensions.WikiLinkExtension{},
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

	contentPathPages := "content/pages"
	pageFiles, _ := fs.ReadDir(Content, contentPathPages)
	if pageFiles != nil {
		for _, file := range pageFiles {
			gen.CompileMarkdownFilePage(file, contentPathPages, md)
		}
	}

	html := gen.generateIndex()
	afero.WriteFile(gen.Output, "blog/index.html", []byte(html), fs.ModePerm)

	contentPathTIL := "content/til"
	tilFiles, _ := fs.ReadDir(Content, contentPathTIL)
	if tilFiles != nil {
		for _, file := range tilFiles {
			gen.CompileMarkdownFile(file, contentPathTIL, "til", md)
		}
		tilIndexHTML := gen.generateTILIndex()
		gen.Output.MkdirAll("til", fs.ModeDir)
		afero.WriteFile(gen.Output, "til/index.html", []byte(tilIndexHTML), fs.ModePerm)
	}

	contentPathKB := "content/kb"
	fs.WalkDir(Content, contentPathKB, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		// Compute relative path within kb, e.g., "devops/docker.md" -> "devops/docker"
		relPath, _ := filepath.Rel(contentPathKB, path)
		kbPath := strings.TrimSuffix(relPath, ".md")
		gen.CompileMarkdownFileKB(path, kbPath, md)
		return nil
	})

	kbIndexHTML := gen.generateKBIndex()
	afero.WriteFile(gen.Output, "kb/index.html", []byte(kbIndexHTML), fs.ModePerm)

	// Copy historical HTML files directly to output
	contentPathHistorical := filepath.Join(ContentPrefix, "historical")
	if ok, _ := afero.Exists(afero.FromIOFS{FS: Content}, contentPathHistorical); ok {
		utils.CopyFiles(afero.FromIOFS{FS: Content}, gen.Output, contentPathHistorical, ContentPrefix, "")
	}

	// Load gallery albums
	contentPathGallery := "content/gallery"
	albumDirs, err := fs.ReadDir(Content, contentPathGallery)
	if err == nil {
		for _, albumDir := range albumDirs {
			if !albumDir.IsDir() {
				continue
			}
			albumSlug := albumDir.Name()
			albumTomlPath := filepath.Join(contentPathGallery, albumSlug, "album.toml")

			albumData, err := Content.ReadFile(albumTomlPath)
			if err != nil {
				log.WithField("album", albumSlug).Warn("Skipping album: no album.toml found")
				continue
			}

			var album Album
			if _, err := toml.Decode(string(albumData), &album); err != nil {
				log.WithField("album", albumSlug).WithField("error", err).Warn("Skipping album: invalid album.toml")
				continue
			}
			album.Slug = albumSlug

			// Discover images in images/ subfolder
			imagesDir := filepath.Join(contentPathGallery, albumSlug, "images")
			imageFiles, err := fs.ReadDir(Content, imagesDir)
			if err == nil {
				for _, imgFile := range imageFiles {
					if !imgFile.IsDir() {
						album.Images = append(album.Images, imgFile.Name())
					}
				}
			}

			gen.Albums[albumSlug] = &album

			// Copy album images to output
			utils.CopyFiles(afero.FromIOFS{FS: Content}, gen.Output, imagesDir, ContentPrefix, "")
		}

		// Generate gallery index and album pages
		galleryIndexHTML := gen.generateGalleryIndex()
		gen.Output.MkdirAll("gallery", fs.ModeDir)
		afero.WriteFile(gen.Output, "gallery/index.html", []byte(galleryIndexHTML), fs.ModePerm)

		for _, album := range gen.Albums {
			albumHTML := gen.generateAlbumPage(album)
			gen.Output.MkdirAll(filepath.Join("gallery", album.Slug), fs.ModeDir)
			afero.WriteFile(gen.Output, fmt.Sprintf("gallery/%s/index.html", album.Slug), []byte(albumHTML), fs.ModePerm)
		}
	}

	// Generate tags page
	tagsPageHTML := gen.generateTagsPage()
	gen.Output.MkdirAll("tags", fs.ModeDir)
	afero.WriteFile(gen.Output, "tags/index.html", []byte(tagsPageHTML), fs.ModePerm)

	// Generate individual tag pages
	allTagsSet := map[string]bool{}
	for tag := range gen.TaggedPosts {
		allTagsSet[tag] = true
	}
	for tag := range gen.TaggedKBs {
		allTagsSet[tag] = true
	}
	for tag := range gen.TaggedTILs {
		allTagsSet[tag] = true
	}
	for tag := range allTagsSet {
		tagListHTML := gen.generateTagList(tag)
		gen.Output.MkdirAll(fmt.Sprintf("tags/%s", tag), fs.ModeDir)
		afero.WriteFile(gen.Output, fmt.Sprintf("tags/%s/index.html", tag), []byte(tagListHTML), fs.ModePerm)
	}

	// Compute KB backlinks
	wikiLinkRegex := regexp.MustCompile(`(?i)([^.\n]*?)\[\[([^[\]|]+)(?:\|[^\]]+)?\]\]([^.\n]*\.?)`)
	for sourcePath, sourceKB := range gen.KBs {
		matches := wikiLinkRegex.FindAllSubmatch(sourceKB.Markdown, -1)
		for _, match := range matches {
			if len(match) >= 4 {
				targetPath := strings.ToLower(strings.TrimSpace(string(match[2])))
				if _, exists := gen.KBs[targetPath]; exists && sourcePath != targetPath {
					context := strings.TrimSpace(string(match[1]) + " " + string(match[2]) + " " + string(match[3]))
					gen.KBs[targetPath].Backlinks = append(gen.KBs[targetPath].Backlinks, Backlink{
						Slug:    sourcePath,
						Context: context,
					})
				}
			}
		}
	}

	gen.generateFeeds()
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

	content, _, err := utils.StripFrontMatter(data)
	if err != nil {
		log.WithField("error", err).WithField("file", file.Name()).Error("Error stripping frontmatter, skipping file")
		return
	}

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

	origDate, err := utils.ParseDate(metaData["date"].(string))
	if err != nil {
		log.WithField("error", err).WithField("file", file.Name()).Error("Error parsing date, skipping file")
		return
	}
	modDate, err := utils.ModifyDate(metaData["date"].(string))
	if err != nil {
		log.WithField("error", err).WithField("file", file.Name()).Error("Error modifying date, skipping file")
		return
	}

	frontmatter := utils.FrontMatter{
		Author:   metaData["author"].(string),
		OrigDate: origDate,
		Date:     modDate,
		Draft:    metaData["draft"].(bool),
		Slug:     metaData["slug"].(string),
		Title:    metaData["title"].(string),
		Type:     metaData["type"].(string),
		Tags:     tagsStringList,
	}

	switch contentType {
	case "post":
		gen.postSpecific(frontmatter, content, out.String())
	case "til":
		gen.tilSpecific(frontmatter, out.String())
	}

}

func (gen *Generator) tilSpecific(frontmatter utils.FrontMatter, html string) {
	til := TIL{
		FrontMatter: frontmatter,
		HTML:        html,
	}

	gen.TILs[frontmatter.Slug] = &til

	for _, tag := range til.FrontMatter.Tags {
		gen.TaggedTILs[tag] = append(gen.TaggedTILs[tag], &til)
	}
	hashtags := []string{}
	for _, tag := range til.FrontMatter.Tags {
		hashtags = append(hashtags, fmt.Sprintf("#%s", tag))
	}
	til.Tagsline = strings.Join(hashtags, " | ")
}

// CompileMarkdownFileKB handles KB markdown files with a path for nested hierarchy.
func (gen *Generator) CompileMarkdownFileKB(filePath string, kbPath string, md goldmark.Markdown) {
	data, err := Content.ReadFile(filePath)
	if err != nil {
		log.WithField("error", err).Fatal("Error reading file")
		panic(err)
	}

	content, _, err := utils.StripFrontMatter(data)
	if err != nil {
		log.WithField("error", err).WithField("file", filePath).Error("Error stripping frontmatter, skipping file")
		return
	}

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

	origDate, err := utils.ParseDate(metaData["date"].(string))
	if err != nil {
		log.WithField("error", err).WithField("file", filePath).Error("Error parsing date, skipping file")
		return
	}
	modDate, err := utils.ModifyDate(metaData["date"].(string))
	if err != nil {
		log.WithField("error", err).WithField("file", filePath).Error("Error modifying date, skipping file")
		return
	}

	frontmatter := utils.FrontMatter{
		Author:   metaData["author"].(string),
		OrigDate: origDate,
		Date:     modDate,
		Draft:    metaData["draft"].(bool),
		Slug:     metaData["slug"].(string),
		Title:    metaData["title"].(string),
		Type:     metaData["type"].(string),
		Tags:     tagsStringList,
	}

	gen.kbSpecific(frontmatter, content, out.String(), kbPath)
}

func (gen *Generator) CompileMarkdownFilePage(file fs.DirEntry, contentPath string, md goldmark.Markdown) {
	if !strings.Contains(file.Name(), ".md") {
		return
	}
	data, err := Content.ReadFile(filepath.Join(contentPath, file.Name()))
	if err != nil {
		log.WithField("error", err).Fatal("Error reading page file")
		panic(err)
	}

	_, _, err = utils.StripFrontMatter(data)
	if err != nil {
		log.WithField("error", err).WithField("file", file.Name()).Error("Error stripping frontmatter, skipping page file")
		return
	}

	context := parser.NewContext()
	var out strings.Builder

	err = md.Convert(data, &out, parser.WithContext(context))
	if err != nil {
		log.WithField("error", err).Fatal("Error converting page markdown")
		panic(err)
	}

	metaData := meta.Get(context)
	
	frontmatter := utils.FrontMatter{
		Author:   "",
		Draft:    false,
		Slug:     metaData["slug"].(string),
		Title:    metaData["title"].(string),
		Type:     "page",
	}

	page := StandalonePage{
		FrontMatter: frontmatter,
		HTML:        out.String(),
	}

	gen.Pages[frontmatter.Slug] = &page
	
	htmlOut := gen.generatePage(&page)
	gen.Output.MkdirAll(fmt.Sprintf("%s", frontmatter.Slug), fs.ModeDir)
	afero.WriteFile(gen.Output, fmt.Sprintf("%s/index.html", frontmatter.Slug), []byte(htmlOut), fs.ModePerm)
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
		PrimaryDomain:  gen.config.Site.PrimaryDomain,
	}

	gen.generateSummary(&post)

	gen.Posts[frontmatter.Slug] = &post

	for _, tag := range post.FrontMatter.Tags {
		gen.TaggedPosts[tag] = append(gen.TaggedPosts[tag], &post)
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

func (gen *Generator) kbSpecific(frontmatter utils.FrontMatter, content []byte, html string, kbPath string) {

	kb := KB{
		FrontMatter:   frontmatter,
		Markdown:      content,
		HTML:          html,
		Path:          kbPath,
		PrimaryDomain: gen.config.Site.PrimaryDomain,
	}

	gen.KBs[kbPath] = &kb

	for _, tag := range kb.FrontMatter.Tags {
		gen.TaggedKBs[tag] = append(gen.TaggedKBs[tag], &kb)
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
		gen.Output.MkdirAll(filepath.Join("kb", kbPath), fs.ModeDir)
		afero.WriteFile(gen.Output, fmt.Sprintf("kb/%s/index.html", kbPath), []byte(htmlOut), fs.ModePerm)
	}

}
