package app

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GoInstaller struct {
	baseInstaller

	goPath, binName string
	cgo             bool
}

var _ Installer = (*GoInstaller)(nil)

func NewGoInstaller(baseDir, goPath, binName string, cgo bool) *GoInstaller {
	i := &GoInstaller{
		goPath:  goPath,
		binName: binName,
		cgo:     cgo,
	}
	i.baseInstaller = newBaseInstaller(filepath.Join(baseDir, i.Name()))
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

func (i *GoInstaller) cmdRun(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = i.Dir()
	cmd.Env = append(os.Environ(), "GOPATH="+i.Dir(), "GOBIN="+i.Dir(), "GO111MODULE=on")
	cmd.Stdout = i.stdout
	cmd.Stderr = i.stderr
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

func (i *GoInstaller) RequireHook(ctx context.Context) error {
	if !i.cgo {
		return nil
	}
	out, err := exec.CommandContext(ctx, "go", "env", "CC").Output()
	if err != nil {
		return err
	}
	cc := strings.TrimSpace(string(out))
	if _, err := exec.LookPath(cc); err != nil {
		return err
	}
	return nil
}
