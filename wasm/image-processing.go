package wasm

import (
	"fmt"
	_ "image/png"
	"runtime"

	"github.com/ebitengine/purego"
)

type ImageProcessing struct {
}

func NewImageProcessing() {
	libc, err := purego.Dlopen(getSystemLibrary(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	var puts func(string)
	purego.RegisterLibFunc(&puts, libc, "puts")
	puts("Calling C from Go without Cgo!")
}

func getSystemLibrary() string {
	switch runtime.GOOS {
	case "linux":
		return "libc.so.6"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}
