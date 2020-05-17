package app

import (
	"context"
	"fmt"
	"path/filepath"
)

type EclipseJDTLSInstaller struct {
	baseInstaller
}

var _ Installer = (*EclipseJDTLSInstaller)(nil)

func NewEclipseJDTLSInstaller(baseDir string) *EclipseJDTLSInstaller {
	var i EclipseJDTLSInstaller
	i.baseInstaller = newBaseInstaller(filepath.Join(baseDir, i.Name()))
	return &i
}

func (i *EclipseJDTLSInstaller) Name() string {
	return "eclipse.jdt.ls"
}

func (i *EclipseJDTLSInstaller) Version() string {
	return "latest"
}

func (i *EclipseJDTLSInstaller) BinName() string {
	return ""
}

func (i *EclipseJDTLSInstaller) Requires() []string {
	return []string{}
}

func (i *EclipseJDTLSInstaller) Install(ctx context.Context) error {
	archive := fmt.Sprintf("jdt-language-server-%s.tar.gz", i.Version())
	u := fmt.Sprintf("https://download.eclipse.org/jdtls/snapshots/%s", archive)
	return i.FetchWithExtract(ctx, u, archive)
}
