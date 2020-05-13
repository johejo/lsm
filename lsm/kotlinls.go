package lsm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/koron-go/zipx"
)

type KotlinLSInstaller struct {
	baseInstaller
}

var _ Installer = (*KotlinLSInstaller)(nil)

func NewKotlinLSInstaller(dir string) *KotlinLSInstaller {
	var i KotlinLSInstaller
	i.baseInstaller = baseInstaller{dir: filepath.Join(dir, i.Name())}
	return &i
}

func (i *KotlinLSInstaller) Name() string {
	return "kotlin-language-server"
}

func (i *KotlinLSInstaller) BinName() string {
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
	zippped := filepath.Join(i.Dir(), "server.zip")
	if err := i.download(req, zippped); err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(zippped); err != nil {
			log.Println(err)
		}
	}()
	if err := zipx.New().ExtractFile(ctx, zippped, zipx.Dir(i.Dir())); err != nil {
		return err
	}
	src := filepath.Join("server", "bin", i.BinName())
	dst := filepath.Join(i.Dir(), i.BinName())
	if err := os.Symlink(src, dst); err != nil {
		return err
	}
	return nil
}
