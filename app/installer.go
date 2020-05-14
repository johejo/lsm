package app

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mholt/archiver/v3"
	"github.com/schollz/progressbar/v3"
)

type Installer interface {
	Name() string
	BinName() string
	Dir() string
	Requires() []string
	Version() string
	Install(ctx context.Context) error
}

type baseInstaller struct {
	dir string
}

func (i *baseInstaller) Dir() string {
	return i.dir
}

func (i *baseInstaller) Download(req *http.Request, dst string) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(filepath.Clean(dst))
	if err != nil {
		return err
	}
	defer f.Close()
	bar := progressbar.DefaultBytes(resp.ContentLength, "downloading")
	if _, err := io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
		return err
	}
	return nil
}

func (i *baseInstaller) CmdRun(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = i.Dir()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (i *baseInstaller) Extract(ctx context.Context, path string) error {
	if err := archiver.Unarchive(path, i.Dir()); err != nil {
		return err
	}
	return nil
}
