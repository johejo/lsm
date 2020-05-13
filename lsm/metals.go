package lsm

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type MetalsInstaller struct {
	baseInstaller
}

var _ Installer = (*MetalsInstaller)(nil)

func NewMetalsInstaller(baseDir string) *MetalsInstaller {
	return &MetalsInstaller{baseInstaller{dir: filepath.Join(baseDir, "metals")}}
}

func (i *MetalsInstaller) Name() string {
	return "metals"
}

func (i *MetalsInstaller) BinName() string {
	if isWindows {
		return i.Name() + ".bat"
	}
	return i.Name()
}

func (i *MetalsInstaller) Requires() []string {
	return []string{"java"}
}

func (i *MetalsInstaller) Version() string {
	return "0.9.0"
}

func (i *MetalsInstaller) Install(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://git.io/coursier-cli", nil)
	if err != nil {
		return err
	}
	coursier := filepath.Join(i.Dir(), "coursier")
	if err := i.Download(req, coursier); err != nil {
		return err
	}
	if isWindows {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://git.io/coursier-bat", nil)
		if err != nil {
			return err
		}
		coursierBat := filepath.Join(i.Dir(), "coursier.bat")
		if err := i.Download(req, coursierBat); err != nil {
			return err
		}
	}
	cmd := exec.CommandContext(ctx,
		"java", "-jar", "coursier", "bootstrap",
		"--ttl", "Inf", "org.scalameta:metals_2.12:"+i.Version(), "-r", "bintray:scalacenter/releases", "-r", "sonatype:public",
		"-o", filepath.Join(i.Dir(), i.Name()),
	)
	cmd.Dir = i.Dir()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
