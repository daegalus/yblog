package utils

import (
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
