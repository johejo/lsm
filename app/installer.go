package app

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
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
	RequireHook(ctx context.Context) error
	Supports() []Support
	Version() string
	Install(ctx context.Context) error
	SetWriter(w io.Writer)
}

type Support struct {
	os, arch string
}

type baseInstaller struct {
	dir            string
	stdout, stderr io.Writer
}

func newBaseInstaller(dir string) baseInstaller {
	stdout := colorable.NewColorableStdout()
	stderr := colorable.NewColorableStderr()
	return baseInstaller{dir: dir, stdout: stdout, stderr: stderr}
}

func (i *baseInstaller) RequireHook(ctx context.Context) error {
	return nil
}

func (i *baseInstaller) Supports() []Support {
	return []Support{}
}

func (i *baseInstaller) SetWriter(w io.Writer) {
	i.stderr = w
}

func (i *baseInstaller) Dir() string {
	return i.dir
}

func (i *baseInstaller) Version() string {
	return ""
}

func (i *baseInstaller) Download(req *http.Request, dst string) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("invalid status code: %v, body=%v", resp.StatusCode, string(b))
	}

	f, err := os.Create(filepath.Clean(dst))
	if err != nil {
		return err
	}
	defer f.Close()

	bar := pb.Full.Start64(resp.ContentLength)
	defer bar.Finish()
	bar.SetWriter(i.stderr)
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
