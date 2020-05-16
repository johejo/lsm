package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
)

type TerraformLSInstaller struct {
	baseInstaller
}

var _ Installer = (*TerraformLSInstaller)(nil)

func NewTerraformLSInstaller(baseDir string) *TerraformLSInstaller {
	var i TerraformLSInstaller
	return &TerraformLSInstaller{baseInstaller: newBaseInstaller(filepath.Join(baseDir, i.Name()))}
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
	switch runtime.GOOS {
	case windows, linux, darwin:
	default:
		return errors.New(runtime.GOOS + " is not supported")
	}
	u := fmt.Sprintf("https://github.com/hashicorp/terraform-ls/releases/download/v%[1]s/terraform-ls_%[1]s_%s_%s.zip", i.Version(), runtime.GOOS, runtime.GOARCH)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	archive := filepath.Join(i.Dir(), i.Name()+".zip")
	if err := i.Download(req, archive); err != nil {
		return err
	}
	if err := i.Extract(ctx, archive); err != nil {
		return err
	}
	return nil
}
