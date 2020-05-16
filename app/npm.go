package app

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type NpmInstaller struct {
	baseInstaller

	moduleName string
	binName    string
}

var _ Installer = (*NpmInstaller)(nil)

func NewNpmInstaller(baseDir, moduleName, binName string) *NpmInstaller {
	return &NpmInstaller{
		baseInstaller: baseInstaller{dir: filepath.Join(baseDir, moduleName)},
		moduleName:    moduleName,
		binName:       binName,
	}
}

func (i *NpmInstaller) Name() string {
	return i.moduleName
}

func (i *NpmInstaller) BinName() string {
	return i.binName
}

func (i *NpmInstaller) Requires() []string {
	return []string{"node", "npm"}
}

func (i *NpmInstaller) Version() string {
	return "latest"
}

func (i *NpmInstaller) Install(ctx context.Context) error {
	f, err := os.Create(filepath.Join(i.Dir(), "package.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.WriteString(f, `{"name": ""}`); err != nil {
		return err
	}

	if err := i.CmdRun(ctx, "npm", "install", i.Name()); err != nil {
		return err
	}

	src := filepath.Join("node_modules", ".bin", i.BinName())
	dst := filepath.Join(i.Dir(), i.BinName())
	if err := os.Symlink(src, dst); err != nil {
		return err
	}
	return nil
}
