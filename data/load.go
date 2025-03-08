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
	"github.com/zeebo/blake3"
)

func threadEncode(returnChannel chan []byte, imageData image.Image, wg *sync.WaitGroup, encodeFunc func(image.Image) []byte) {
	defer wg.Done()
	encodedImage := encodeFunc(imageData)
	returnChannel <- encodedImage
}

func encodeFilesInDir(file fs.DirEntry, imagePath string, imageDirEntry fs.DirEntry, gen *Generator, wg *sync.WaitGroup) {
	defer wg.Done()
	postImagePath := filepath.Join(imagePath, imageDirEntry.Name(), file.Name())
	filename := strings.Split(file.Name(), ".")[0]
	imageByteData, err := afero.ReadFile(gen.Input, postImagePath)
	if err != nil {
		log.WithField("error", err).Fatal("Error reading image")
	}
	cachedImageByteData, err := afero.ReadFile(gen.Cache, postImagePath)
	if err != nil {
		log.WithField("error", err).Warn("Error reading cached image")
	}
	bi := blake3.Sum512(imageByteData)
	bci := blake3.Sum512(cachedImageByteData)
	if bi == bci {
		log.WithField("file", postImagePath).Info("Skipping encoding")
		return // Skip encoding if the image is the same
	}

	imageData := utils.ImageFromPNG(afero.FromIOFS{FS: Content}, postImagePath)

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
	cachedImageByteData, err := afero.ReadFile(gen.Cache, filepath.Join(CachePrefix, imageStoragePath, postImageDir.Name(), file.Name()))
	if err != nil {
		log.WithField("error", err).Warn("Error reading cached image")
	}

	bi := blake3.Sum512(imageByteData)
	bci := blake3.Sum512(cachedImageByteData)
	imageDifferent := bi != bci
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

func encodeAllDirs(imagePath string, imageDirEntry fs.DirEntry, gen *Generator, wg *sync.WaitGroup) {
	defer wg.Done()
	log.WithField("dir", imageDirEntry.Name()).Info("Converting in Post")
	if imageDirEntry.IsDir() {
		postImageDir, _ := fs.ReadDir(Content, filepath.Join(ContentPrefix, imagePath, imageDirEntry.Name()))
		wg.Add(len(postImageDir))
		for _, file := range postImageDir {
			go encodeFilesInDir(file, imagePath, imageDirEntry, gen, wg)
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
	utils.CopyFiles(gen.Cache, gen.Output, filepath.Join(CachePrefix, imageStoragePath), CachePrefix, "")
	utils.CopyFiles(afero.FromIOFS{FS: Content}, gen.Cache, filepath.Join(ContentPrefix, imageStoragePath), ContentPrefix, CachePrefix)

	// afero.Walk(gen.Output, ".", func(path string, info fs.FileInfo, err error) error {
	// 	log.WithField("path", path).Info("Walking")
	// 	return err
	// })
	// os.Exit(0)

	imageStorageDirEntries, _ := fs.ReadDir(Content, filepath.Join(ContentPrefix, imageStoragePath))
	//var wg sync.WaitGroup
	//wg.Add(len(imageDirEntries))
	for _, imageStorageDirEntry := range imageStorageDirEntries {
		//go encodeAllDirs(imagePath, imageDirEntry, gen, &wg)
		encodeAllDirsSync(imageStoragePath, imageStorageDirEntry, gen)
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
