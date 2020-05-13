package lsm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type KotlinLSInstaller struct {
	baseInstaller
}

var _ Installer = (*KotlinLSInstaller)(nil)

func NewKotlinLSInstaller(baseDir string) *KotlinLSInstaller {
	var i KotlinLSInstaller
	i.baseInstaller = baseInstaller{dir: filepath.Join(baseDir, i.Name())}
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
	return []string{"java"}
}

func (i *KotlinLSInstaller) Version() string {
	return "0.5.2"
}

func (i *KotlinLSInstaller) Install(ctx context.Context) error {
	u := fmt.Sprintf("https://github.com/fwcd/kotlin-language-server/releases/download/%s/server.zip", i.Version())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	archive := filepath.Join(i.Dir(), "server.zip")
	if err := i.Download(req, archive); err != nil {
		return err
	}
	if err := i.Extract(ctx, archive); err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(archive); err != nil {
			log.Println(err)
		}
	}()
	src := filepath.Join("server", "bin", i.BinName())
	dst := filepath.Join(i.Dir(), i.BinName())
	if err := os.Symlink(src, dst); err != nil {
		return err
	}
	return nil
}
