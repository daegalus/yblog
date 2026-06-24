package utils

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/caarlos0/log"
	"github.com/spf13/afero"
)

// CopyFiles walks rootPath in input and copies every file to output. When stripPrefix
// is set the leading path segment is dropped and, if replacementPrefix is also set,
// replaced with it. Individual file errors are logged and skipped rather than aborting
// the whole copy.
func CopyFiles(input afero.Fs, output afero.Fs, rootPath string, stripPrefix string, replacementPrefix string) {
	err := afero.Walk(input, rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		targetPath := path
		if stripPrefix != "" {
			pathList := strings.Split(path, string(filepath.Separator))
			if len(pathList) > 1 {
				targetPath = filepath.Join(slices.Delete(pathList, 0, 1)...)
			}
			if replacementPrefix != "" {
				targetPath = filepath.Join(replacementPrefix, targetPath)
			}
		}

		in, err := afero.ReadFile(input, path)
		if err != nil {
			log.WithError(err).WithField("file", path).Warn("Skipping file during copy: cannot read")
			return nil
		}
		if err := output.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			log.WithError(err).WithField("dir", filepath.Dir(targetPath)).Warn("Skipping file during copy: cannot create directory")
			return nil
		}
		if err := afero.WriteFile(output, targetPath, in, 0644); err != nil {
			log.WithError(err).WithField("file", targetPath).Warn("Failed to write file during copy")
		}
		return nil
	})
	if err != nil {
		log.WithError(err).Error("Error walking files during copyFiles")
	}
}
