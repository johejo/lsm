package lsm

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
	i.baseInstaller = baseInstaller{dir: filepath.Join(baseDir, i.Name())}
	return &i
}

func (i *RustAnalyzerInstaller) Name() string {
	return "rust-analyzer"
}

func (i *RustAnalyzerInstaller) BinName() string {
	return i.Name()
}

func (i *RustAnalyzerInstaller) Requires() []string {
	return []string{}
}

func (i *RustAnalyzerInstaller) Version() string {
	return "2020-05-11"
}

func (i *RustAnalyzerInstaller) Install(ctx context.Context) error {
	var binSuffix string
	switch runtime.GOOS {
	case "darwin":
		binSuffix = "mac"
	case "linux":
		binSuffix = "linux"
	case "windows":
		binSuffix = "windows.exe"
	default:
		return errors.New(runtime.GOOS + " is not supported")
	}

	u := fmt.Sprintf("https://github.com/rust-analyzer/rust-analyzer/releases/download/%s/rust-analyzer-%s", i.Version(), binSuffix)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	bin := filepath.Join(i.Dir(), i.BinName())
	if err := i.download(req, bin); err != nil {
		return err
	}
	// TODO Is this work on Windows?
	if err := os.Chmod(bin, 0777); err != nil {
		return err
	}
	return nil
}
