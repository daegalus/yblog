package main

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/log"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const Version = "0.0.1"

type Build mg.Namespace

func (b Build) All() {
	mg.SerialDeps(b.Deps, b.Linux, b.Windows, b.MacOS)
}

func (Build) Linux() {
	log.Info("Building Linux binaries")
	sh.RunWithV(buildEnv("linux", "amd64"), "go", buildCommand("yblog-linux-amd64")...)
	sh.RunWithV(buildEnv("linux", "arm64"), "go", buildCommand("yblog-linux-arm64")...)
	sh.RunWithV(buildEnv("linux", "riscv64"), "go", buildCommand("yblog-linux-riscv64")...)
}

func (Build) Windows() {
	log.Info("Building Windows binaries")
	sh.RunWithV(buildEnv("windows", "amd64"), "go", buildCommand("yblog-windows-amd64.exe")...)
}

func (Build) MacOS() {
	log.Info("Building MacOS binaries")
	sh.RunWithV(buildEnv("darwin", "amd64"), "go", buildCommand("yblog-macos-amd64")...)
	sh.RunWithV(buildEnv("darwin", "arm64"), "go", buildCommand("yblog-macos-arm64")...)
}

func (Build) Wasm() {
	log.Info("Building WASM binaries")
	sh.RunWithV(buildEnv("js", "wasm"), "go", buildCommand("yblog-js-wasm")...)
}

func (Build) Deps() {
	log.Info("Installing dependencies")
	os.RemoveAll("build")
	os.Mkdir("build", 0755)
	sh.Exec(map[string]string{}, os.Stdout, os.Stderr, "go", "mod", "download")
}

func buildEnv(goos string, goarch string) map[string]string {
	return map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        goos,
		"GOARCH":      goarch,
	}
}

func buildCommand(filename string) []string {
	return []string{
		"build",
		"--ldflags", fmt.Sprintf("-s -w -extldflags=-static-pie -X 'main.Version=%s' -X 'main.BuildDate=%s'", Version, time.Now().Format(time.RFC3339)),
		"-o", fmt.Sprintf("build/%s", filename),
		"yblog.go",
	}
}
