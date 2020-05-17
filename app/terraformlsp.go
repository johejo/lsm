package app

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
)

type TerraformLSPInstaller struct {
	baseInstaller
}

var _ Installer = (*TerraformLSPInstaller)(nil)

func NewTerraformLSPInstaller(baseDir string) *TerraformLSPInstaller {
	var i TerraformLSPInstaller
	i.baseInstaller = newBaseInstaller(filepath.Join(baseDir, i.Name()))
	return &i
}

func (i *TerraformLSPInstaller) BinName() string {
	if isWindows {
		return i.Name() + ".exe"
	}
	return i.Name()
}

func (i *TerraformLSPInstaller) Requires() []string {
	return []string{}
}

func (i *TerraformLSPInstaller) Name() string {
	return "terraform-lsp"
}

func (i *TerraformLSPInstaller) Version() string {
	return "0.0.11-beta1"
}

func (i *TerraformLSPInstaller) Supports() []Support {
	return []Support{
		{os: darwin, arch: amd64},
		{os: linux, arch: amd64},
		{os: windows, arch: amd64},
	}
}

func (i *TerraformLSPInstaller) Install(ctx context.Context) error {
	u := fmt.Sprintf("https://github.com/juliosueiras/terraform-lsp/releases/download/v%[1]s/terraform-lsp_%[1]s_%s_amd64.tar.gz", i.Version(), runtime.GOOS)
	return i.FetchWithExtract(ctx, u, filepath.Join(i.Dir(), i.Name()+".tar.gz"))
}
