package app

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/mattn/go-colorable"
	"github.com/mholt/archiver/v3"
)

type Installer interface {
	Name() string
	BinName() string
	Dir() string
	Requires() []string
	RequiresHook() Hook
	Version() string
	Install(ctx context.Context) error
	SetWriter(w io.Writer)
}

type baseInstaller struct {
	dir string
	out io.Writer
}

func (i *baseInstaller) RequiresHook() Hook {
	return nil
}

func (i *baseInstaller) SetWriter(w io.Writer) {
	i.out = w
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
	bar := pb.Full.Start64(resp.ContentLength)
	defer bar.Finish()
	if i.out == nil {
		bar.SetWriter(colorable.NewColorableStderr())
	} else {
		bar.SetWriter(i.out)
	}
	pr := bar.NewProxyReader(resp.Body)
	if _, err := io.Copy(f, pr); err != nil {
		return err
	}
	return nil
}

func (i *baseInstaller) CmdRun(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = i.Dir()
	cmd.Stdout = colorable.NewColorableStdout()
	cmd.Stderr = colorable.NewColorableStderr()
	return cmd.Run()
}

func (i *baseInstaller) Extract(ctx context.Context, path string) error {
	return archiver.Unarchive(path, i.Dir())
}

type Hook func(ctx context.Context) error
