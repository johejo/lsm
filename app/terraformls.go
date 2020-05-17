package app

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
)

type TerraformLSInstaller struct {
	baseInstaller
}

var _ Installer = (*TerraformLSInstaller)(nil)

func NewTerraformLSInstaller(baseDir string) *TerraformLSInstaller {
	var i TerraformLSInstaller
	i.baseInstaller = newBaseInstaller(filepath.Join(baseDir, i.Name()))
	return &i
}

func (i *TerraformLSInstaller) BinName() string {
	if isWindows {
		return i.Name() + ".exe"
	}
	return i.Name()
}

func (i *TerraformLSInstaller) Requires() []string {
	return []string{}
}

func (i *TerraformLSInstaller) Name() string {
	return "terraform-ls"
}

func (i *TerraformLSInstaller) Version() string {
	return "0.2.0"
}

func (i *TerraformLSInstaller) Supports() []Support {
	return []Support{
		{os: darwin, arch: amd64},

		{os: freebsd, arch: _386},
		{os: freebsd, arch: amd64},
		{os: freebsd, arch: arm},

		{os: linux, arch: _386},
		{os: linux, arch: amd64},
		{os: linux, arch: arm},

		{os: openbsd, arch: _386},
		{os: openbsd, arch: amd64},

		{os: solaris, arch: amd64},

		{os: windows, arch: _386},
		{os: windows, arch: amd64},
	}
}

func (i *TerraformLSInstaller) Install(ctx context.Context) error {
	u := fmt.Sprintf("https://github.com/hashicorp/terraform-ls/releases/download/v%[1]s/terraform-ls_%[1]s_%s_%s.zip", i.Version(), runtime.GOOS, runtime.GOARCH)
	return i.FetchWithExtract(ctx, u, filepath.Join(i.Dir(), i.Name()+".zip"))
}
