package utils

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/caarlos0/log"
	"github.com/spf13/afero"
)

func ReadFile(fileSys afero.Fs, filename string) ([]byte, error) {
	out, err := afero.ReadFile(fileSys, filename)
	return out, err
}

func ReadFileToString(fileSys afero.Fs, filename string) (string, error) {
	out, err := ReadFile(fileSys, filename)
	return string(out), err
}

func CopyFiles(input afero.Fs, output afero.Fs, rootPath string, stripPrefix string, replacementPrefix string) {
	err := afero.Walk(input, rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.WithError(err).Fatal("Error reading files")
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

		input.MkdirAll(filepath.Dir(path), 0755)

		in, err := afero.ReadFile(input, path)
		if err != nil {
			log.WithError(err).Fatal("Error reading files")
		}
		afero.WriteFile(output, targetPath, in, 0644)

		return nil
	})
	if err != nil {
		log.WithError(err).Fatal("Error reading files")
	}
}
