package app

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type ESLintServerInstaller struct {
	baseInstaller
}

var _ Installer = (*ESLintServerInstaller)(nil)

func NewESLintServerInstaller(baseDir string) *ESLintServerInstaller {
	var i ESLintServerInstaller
	i.baseInstaller = newBaseInstaller(filepath.Join(baseDir, i.Name()))
	return &i
}

func (i *ESLintServerInstaller) Name() string {
	return "eslint-server"
}

func (i *ESLintServerInstaller) BinName() string {
	return noExecutable
}

func (i *ESLintServerInstaller) Version() string {
	return "2.1.4-next.1"
}

func (i *ESLintServerInstaller) Requires() []string {
	return noRequires
}

func (i *ESLintServerInstaller) normalizedVersion() string {
	return strings.Split(i.Version(), "-")[0] // 2.1.4-next.1 -> 2.1.4
}

func (i *ESLintServerInstaller) Install(ctx context.Context) error {
	slash, err := url.PathUnescape("/")
	if err != nil {
		return err
	}
	u := fmt.Sprintf("https://github.com/microsoft/vscode-eslint/releases/download/release%s%s/vscode-eslint-%s.vsix", slash, i.Version(), i.normalizedVersion())
	const archive = "vscode-eslint.zip"
	return i.FetchWithExtract(ctx, u, archive)
}
