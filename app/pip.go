package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type PipInstaller struct {
	baseInstaller

	python, moduleName, binName string
}

var _ Installer = (*PipInstaller)(nil)

func isSupportedPython(v string) (bool, error) {
	min, err := semver.NewVersion("3.5")
	if err != nil {
		return false, err
	}
	version, err := semver.NewVersion(v)
	if err != nil {
		return false, err
	}
	return version.GreaterThan(min), nil
}

func lookPython3() (string, error) {
	if _, err := exec.LookPath("python3"); err != nil {
		return "", nil
	}
	return "python3", nil
}

func lookPython() (string, error) {
	if _, err := exec.LookPath("python"); err != nil {
		return "", errors.New("no python found")
	}
	_out, err := exec.Command("python", "--version").Output()
	if err != nil {
		return "", err
	}
	out := strings.TrimSuffix(string(_out), "\n")
	v := strings.Split(out, " ") // ["Python", "3.x.y"]
	if len(v) != 2 {
		return "", fmt.Errorf("invalid python version output %s", string(_out))
	}
	ok, err := isSupportedPython(v[1])
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("unsupported python version: %v", v)
	}
	return "python", nil
}

func NewPipInstaller(baseDir, moduleName, binName string) *PipInstaller {
	i := &PipInstaller{moduleName: moduleName, binName: binName}
	i.baseInstaller = baseInstaller{dir: filepath.Join(baseDir, i.Name())}
	if !isWindows {
		python, err := lookPython3()
		if err == nil {
			i.python = python
			return i
		}
	}
	python, err := lookPython()
	if err != nil {
		panic(err) //FIXME
	}
	i.python = python
	return i
}

func (i *PipInstaller) Name() string {
	return i.moduleName
}

func (i *PipInstaller) BinName() string {
	return i.binName
}

func (i *PipInstaller) Version() string {
	return "latest"
}

func (i *PipInstaller) Requires() []string {
	return []string{i.python}
}

func (i *PipInstaller) Install(ctx context.Context) error {
	python := i.python
	venv := filepath.Join(i.Dir(), "venv")
	if err := i.CmdRun(ctx, python, "-m", "venv", venv); err != nil {
		return err
	}
	var bin string
	if isWindows {
		bin = "Scripts"
	} else {
		bin = "bin"
	}
	vpython := filepath.Join(venv, bin, python)
	if err := i.CmdRun(ctx, vpython, "-m", "pip", "install", "--upgrade", "pip"); err != nil {
		return err
	}
	if err := i.CmdRun(ctx, vpython, "-m", "pip", "install", i.Name()); err != nil {
		return err
	}
	src := filepath.Join("venv", bin, i.BinName())
	dst := filepath.Join(i.Dir(), i.BinName())
	if err := os.Symlink(src, dst); err != nil {
		return err
	}
	return nil
}
