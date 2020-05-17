package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type KotlinLSInstaller struct {
	baseInstaller
}

var _ Installer = (*KotlinLSInstaller)(nil)

func NewKotlinLSInstaller(baseDir string) *KotlinLSInstaller {
	var i KotlinLSInstaller
	i.baseInstaller = newBaseInstaller(filepath.Join(baseDir, i.Name()))
	return &i
}

func (i *KotlinLSInstaller) Name() string {
	return "kotlin-language-server"
}

func (i *KotlinLSInstaller) BinName() string {
	if isWindows {
		return i.Name() + ".bat"
	}
	return i.Name()
}

func (i *KotlinLSInstaller) Requires() []string {
	return noRequires
}

func (i *KotlinLSInstaller) Version() string {
	return "0.5.2"
}

func (i *KotlinLSInstaller) Install(ctx context.Context) error {
	u := fmt.Sprintf("https://github.com/fwcd/kotlin-language-server/releases/download/%s/server.zip", i.Version())
	archive := filepath.Join(i.Dir(), "server.zip")
	if err := i.FetchWithExtract(ctx, u, archive); err != nil {
		return err
	}
	src := filepath.Join("server", "bin", i.BinName())
	dst := filepath.Join(i.Dir(), i.BinName())
	if err := os.Symlink(src, dst); err != nil {
		return err
	}
	return nil
}
