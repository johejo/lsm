package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type RustAnalyzerInstaller struct {
	baseInstaller
}

var _ Installer = (*RustAnalyzerInstaller)(nil)

func NewRustAnalyzerInstaller(baseDir string) *RustAnalyzerInstaller {
	var i RustAnalyzerInstaller
	i.baseInstaller = newBaseInstaller(filepath.Join(baseDir, i.Name()))
	return &i
}

func (i *RustAnalyzerInstaller) Name() string {
	return "rust-analyzer"
}

func (i *RustAnalyzerInstaller) BinName() string {
	if isWindows {
		return i.Name() + ".exe"
	}
	return i.Name()
}

func (i *RustAnalyzerInstaller) Requires() []string {
	return []string{}
}

func (i *RustAnalyzerInstaller) Version() string {
	return "2020-05-11"
}

func (i *RustAnalyzerInstaller) Supports() []Support {
	return []Support{
		{os: darwin, arch: amd64},
		{os: linux, arch: amd64},
		{os: windows, arch: amd64},
	}
}

func (i *RustAnalyzerInstaller) Install(ctx context.Context) error {
	var suffix string
	switch runtime.GOOS {
	case darwin:
		suffix = "mac"
	case linux:
		suffix = linux
	case windows:
		suffix = "windows.exe"
	default:
		return errors.New(runtime.GOOS + " is not supported")
	}

	u := fmt.Sprintf("https://github.com/rust-analyzer/rust-analyzer/releases/download/%s/rust-analyzer-%s", i.Version(), suffix)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	bin := filepath.Join(i.Dir(), i.BinName())
	if err := i.Download(req, bin); err != nil {
		return err
	}
	if err := os.Chmod(bin, 0777); err != nil {
		return err
	}
	return nil
}
