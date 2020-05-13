package lsm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type EfmLSInstaller struct {
	baseInstaller
}

var _ Installer = (*EfmLSInstaller)(nil)

func NewEfmLSInstaller(baseDir string) *EfmLSInstaller {
	var i EfmLSInstaller
	i.baseInstaller = baseInstaller{dir: filepath.Join(baseDir, i.Name())}
	return &i
}

func (i *EfmLSInstaller) Name() string {
	return "efm-langserver"
}

func (i *EfmLSInstaller) BinName() string {
	if isWindows {
		return i.Name() + ".exe"
	}
	return i.Name()
}

func (i *EfmLSInstaller) Requires() []string {
	return []string{}
}

func (i *EfmLSInstaller) Version() string {
	return "0.0.14"
}

func (i *EfmLSInstaller) Install(ctx context.Context) error {
	var ext string
	switch runtime.GOOS {
	case "linux":
		ext = "tar.gz"
	case "darwin", "windows":
		ext = "zip"
	default:
		return errors.New(runtime.GOOS + " is not supported")
	}
	target := fmt.Sprintf("efm-langserver_v%s_%s_amd64", i.Version(), runtime.GOOS)
	archive := fmt.Sprintf("%s.%s", target, ext)
	u := fmt.Sprintf("https://github.com/mattn/efm-langserver/releases/download/v%s/%s", i.Version(), archive)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	tmp := filepath.Join(i.Dir(), archive)
	if err := i.Download(req, tmp); err != nil {
		return err
	}
	if err := i.Extract(ctx, tmp); err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(tmp); err != nil {
			log.Println(err)
		}
	}()
	src := filepath.Join(target, i.BinName())
	dst := filepath.Join(i.Dir(), i.BinName())
	if err := os.Symlink(src, dst); err != nil {
		return err
	}
	return nil
}
