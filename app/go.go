package app

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mattn/go-colorable"
)

type GoInstaller struct {
	baseInstaller

	goPath  string
	binName string
}

var _ Installer = (*GoInstaller)(nil)

func NewGoInstaller(baseDir, goPath, binName string) *GoInstaller {
	i := &GoInstaller{
		goPath:  goPath,
		binName: binName,
	}
	i.baseInstaller = baseInstaller{dir: filepath.Join(baseDir, i.Name())}
	return i
}

func (i *GoInstaller) Name() string {
	return i.binName
}

func (i *GoInstaller) BinName() string {
	if isWindows {
		return i.binName + ".exe"
	}
	return i.binName
}

func (i *GoInstaller) Version() string {
	return "latest"
}

func (i *GoInstaller) cmdRun(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = i.Dir()
	cmd.Env = append(os.Environ(), "GOPATH="+i.Dir(), "GOBIN="+i.Dir())
	cmd.Stdout = colorable.NewColorableStdout()
	cmd.Stderr = colorable.NewColorableStderr()
	return cmd.Run()
}

func (i *GoInstaller) Install(ctx context.Context) error {
	if err := i.cmdRun(ctx, "go", "get", i.goPath); err != nil {
		return err
	}
	if err := i.cmdRun(ctx, "go", "clean", "-modcache"); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(i.Dir(), "src")); err != nil {
		return err
	}
	return nil
}

func (i *GoInstaller) Requires() []string {
	return []string{"go"}
}
