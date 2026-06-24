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
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// imageHashes holds the previously-saved and freshly-computed content hashes of the
// source images, used to decide which images need re-encoding.
type imageHashes struct {
	saved   map[string]string
	current map[string]string
}

func (h imageHashes) changed(srcPath string) bool {
	return h.current[srcPath] != h.saved[srcPath]
}

// encodeImage converts a single source PNG into avif/jxl/webp under images/<dir>/ in
// the output FS. Formats that already exist are skipped when the source is unchanged
// (the original PNG itself is copied to the output by convertImages). Each needed
// format is encoded concurrently, then written sequentially.
func (gen *Generator) encodeImage(imageStoragePath, dirName, fileName string, hashes imageHashes) {
	if strings.ToLower(filepath.Ext(fileName)) != ".png" {
		return
	}
	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	srcPath := filepath.Join(ContentPrefix, imageStoragePath, dirName, fileName)
	outDir := filepath.Join(imageStoragePath, dirName)
	changed := hashes.changed(srcPath)

	type encodeJob struct {
		ext    string
		encode func(image.Image) []byte
	}
	var jobs []encodeJob
	for _, job := range []encodeJob{
		{"avif", utils.ImageToAVIF},
		{"jxl", utils.ImageToJXL},
		{"webp", utils.ImageToWebP},
	} {
		exists, _ := afero.Exists(gen.Output, filepath.Join(outDir, base+"."+job.ext))
		if changed || !exists {
			jobs = append(jobs, job)
		}
	}
	if len(jobs) == 0 {
		return
	}

	img := utils.ImageFromPNG(afero.FromIOFS{FS: Content}, srcPath)
	if img == nil {
		return
	}

	results := make([][]byte, len(jobs))
	var wg sync.WaitGroup
	for i, job := range jobs {
		wg.Add(1)
		go func(i int, job encodeJob) {
			defer wg.Done()
			results[i] = job.encode(img)
		}(i, job)
	}
	wg.Wait()

	for i, job := range jobs {
		if results[i] == nil {
			continue
		}
		outPath := filepath.Join(outDir, base+"."+job.ext)
		log.WithField("file", outPath).Infof("Converting to %s", strings.ToUpper(job.ext))
		if err := afero.WriteFile(gen.Output, outPath, results[i], fs.ModePerm); err != nil {
			log.WithError(err).WithField("file", outPath).Warn("Failed to write encoded image")
		}
	}
}

// convertImages copies source images to the output, restores any cached encodes from
// a previous build, and (re)encodes images whose source content changed. The per-image
// content hashes are loaded and computed once, then persisted for incremental builds.
func (gen *Generator) convertImages() {
	imageStoragePath := "images"

	utils.CopyFiles(afero.FromIOFS{FS: Content}, gen.Output, filepath.Join(ContentPrefix, imageStoragePath), ContentPrefix, "")

	if ok, _ := afero.Exists(gen.Cache, CachePrefix); ok {
		if ok, _ := afero.Exists(gen.Cache, filepath.Join(CachePrefix, imageStoragePath)); ok {
			utils.CopyFiles(gen.Cache, gen.Output, filepath.Join(CachePrefix, imageStoragePath), CachePrefix, "")
		}
	}
	utils.CopyFiles(afero.FromIOFS{FS: Content}, gen.Cache, filepath.Join(ContentPrefix, imageStoragePath), ContentPrefix, CachePrefix)

	saved, _ := cache.LoadImageHashes()
	current, _ := cache.CalculateImageHashes()
	hashes := imageHashes{saved: saved, current: current}

	dirs, _ := fs.ReadDir(Content, filepath.Join(ContentPrefix, imageStoragePath))
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU())
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		files, _ := fs.ReadDir(Content, filepath.Join(ContentPrefix, imageStoragePath, dir.Name()))
		for _, file := range files {
			wg.Add(1)
			go func(dirName, fileName string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				gen.encodeImage(imageStoragePath, dirName, fileName, hashes)
			}(dir.Name(), file.Name())
		}
	}
	wg.Wait()

	if err := cache.SaveImageHashes(current); err != nil {
		log.WithError(err).Warn("Failed to save image hashes")
	}
}

func (gen *Generator) CompileMarkdown() (err error) {
	// Template and feed generation panic on failure (e.g. a broken theme file);
	// recover here so a bad build returns an error instead of taking down the
	// process — important for the serve watcher, which keeps the last good build.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("compile failed: %v", r)
		}
	}()

	log.Info("Converting Images")
	gen.convertImages()

	if comments, cerr := GetLegacyComments(); cerr != nil {
		log.WithError(cerr).Warn("Could not load legacy comments; continuing without them")
	} else {
		gen.legacyComments = comments
	}

	log.Info("Compiling markdown")
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
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

	// Copy static passthrough files verbatim to the output root, preserving hierarchy.
	// e.g. content/static/.well-known/oauth-authorization-server -> /.well-known/oauth-authorization-server
	gen.copyStaticPassthrough()

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
	return nil
}

// copyStaticPassthrough copies everything under content/static to the output root,
// stripping the content/static prefix and preserving the remaining folder hierarchy.
// This is for files that must be served verbatim at the site root (e.g. /.well-known/*,
// robots.txt, IndieAuth/OAuth metadata). No templating or processing is applied.
func (gen *Generator) copyStaticPassthrough() {
	staticRoot := filepath.Join(ContentPrefix, "static")
	if ok, _ := afero.Exists(afero.FromIOFS{FS: Content}, staticRoot); !ok {
		return
	}

	err := fs.WalkDir(Content, staticRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(staticRoot, path)
		if err != nil {
			return err
		}

		fileData, err := Content.ReadFile(path)
		if err != nil {
			return err
		}

		gen.Output.MkdirAll(filepath.Dir(relPath), fs.ModeDir)
		return afero.WriteFile(gen.Output, relPath, fileData, fs.ModePerm)
	})
	if err != nil {
		log.WithError(err).Error("Error copying static passthrough files")
	}
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
	if !strings.HasSuffix(file.Name(), ".md") {
		return
	}
	data, err := Content.ReadFile(filepath.Join(contentPath, file.Name()))
	if err != nil {
		log.WithError(err).WithField("file", file.Name()).Warn("Skipping file: cannot read")
		return
	}

	body, frontmatter, err := utils.StripFrontMatter(data)
	if err != nil {
		log.WithError(err).WithField("file", file.Name()).Warn("Skipping file: bad frontmatter")
		return
	}

	var out strings.Builder
	if err := md.Convert(body, &out); err != nil {
		log.WithError(err).WithField("file", file.Name()).Warn("Skipping file: markdown conversion failed")
		return
	}

	switch contentType {
	case "post":
		gen.postSpecific(frontmatter, body, out.String())
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
		log.WithError(err).WithField("file", filePath).Warn("Skipping KB: cannot read")
		return
	}

	body, frontmatter, err := utils.StripFrontMatter(data)
	if err != nil {
		log.WithError(err).WithField("file", filePath).Warn("Skipping KB: bad frontmatter")
		return
	}

	var out strings.Builder
	if err := md.Convert(body, &out); err != nil {
		log.WithError(err).WithField("file", filePath).Warn("Skipping KB: markdown conversion failed")
		return
	}

	gen.kbSpecific(frontmatter, body, out.String(), kbPath)
}

func (gen *Generator) CompileMarkdownFilePage(file fs.DirEntry, contentPath string, md goldmark.Markdown) {
	if !strings.HasSuffix(file.Name(), ".md") {
		return
	}
	data, err := Content.ReadFile(filepath.Join(contentPath, file.Name()))
	if err != nil {
		log.WithError(err).WithField("file", file.Name()).Warn("Skipping page: cannot read")
		return
	}

	body, frontmatter, err := utils.StripFrontMatter(data)
	if err != nil {
		log.WithError(err).WithField("file", file.Name()).Warn("Skipping page: bad frontmatter")
		return
	}
	if frontmatter.Slug == "" {
		log.WithField("file", file.Name()).Warn("Skipping page: missing slug")
		return
	}

	var out strings.Builder
	if err := md.Convert(body, &out); err != nil {
		log.WithError(err).WithField("file", file.Name()).Warn("Skipping page: markdown conversion failed")
		return
	}

	frontmatter.Type = "page"
	page := StandalonePage{
		FrontMatter: frontmatter,
		HTML:        out.String(),
	}

	gen.Pages[frontmatter.Slug] = &page

	htmlOut := gen.generatePage(&page)
	gen.Output.MkdirAll(frontmatter.Slug, fs.ModeDir)
	afero.WriteFile(gen.Output, fmt.Sprintf("%s/index.html", frontmatter.Slug), []byte(htmlOut), fs.ModePerm)
}

func (gen *Generator) postSpecific(frontmatter utils.FrontMatter, content []byte, html string) {
	legacyComments := []*LegacyComment{}
	if post, ok := gen.legacyComments[frontmatter.Slug]; ok && post.Comments != nil {
		for _, comment := range *post.Comments {
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

	// afero's MemMapFs creates parent directories on write, so no MkdirAll is needed.
	afero.WriteFile(gen.Output, fmt.Sprintf("blog/posts/%s.html", frontmatter.Slug), []byte(htmlOut), fs.ModePerm)
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
