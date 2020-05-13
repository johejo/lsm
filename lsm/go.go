package lsm

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

type GoInstaller struct {
	baseInstaller

	goGetURL string
	binName  string
}

var _ Installer = (*GoInstaller)(nil)

func NewGoInstaller(baseDir, goGetURL, binName string) *GoInstaller {
	return &GoInstaller{
		baseInstaller: baseInstaller{dir: filepath.Join(baseDir, binName)},
		goGetURL:      goGetURL,
		binName:       binName,
	}
}

func (i *GoInstaller) Name() string {
	return i.binName
}

func (i *GoInstaller) BinName() string {
	return i.binName
}

func (i *GoInstaller) Version() string {
	return "latest"
}

func (i *GoInstaller) Install(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "go", "get", i.goGetURL)
	cmd.Env = append(os.Environ(), "GOPATH="+i.dir, "GO111MODULE=on", "GOBIN="+i.dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.CommandContext(ctx, "go", "clean", "-modcache")
	cmd.Env = append(os.Environ(), "GOPATH="+i.dir, "GO111MODULE=on")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := os.RemoveAll(filepath.Join(i.dir, "src")); err != nil {
		return err
	}
	return nil
}

func (i *GoInstaller) Requires() []string {
	return []string{"go"}
}
