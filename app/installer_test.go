package app

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGoInstaller(t *testing.T) {
	h := newInstallerTestHelper(t, "gopls")
	h.Run(context.Background())
}

func TestNpmInstaller(t *testing.T) {
	tests := []string{
		"typescript-language-server",
		"vim-language-server",
		"dockerfile-language-server-nodejs",
		"bash-language-server",
		"yaml-language-server",
		"vscode-json-languageserver",
		"vscode-css-languageserver",
		"vscode-html-languageserver",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			t.Parallel()
			h := newInstallerTestHelper(t, tt)
			h.Run(context.Background())
		})
	}
}

func TestMetalsInstaller(t *testing.T) {
	t.Parallel()
	h := newInstallerTestHelper(t, "metals")
	h.Run(context.Background())
}

func TestKotlinLSInstaller(t *testing.T) {
	t.Parallel()
	h := newInstallerTestHelper(t, "kotlin-language-server")
	h.Run(context.Background())
}

func TestRustAnalyzerInstaller(t *testing.T) {
	t.Parallel()
	h := newInstallerTestHelper(t, "rust-analyzer")
	h.Run(context.Background())
}

func TestEfmLSInstaller(t *testing.T) {
	t.Parallel()
	h := newInstallerTestHelper(t, "efm-langserver")
	h.Run(context.Background())
}

type installerTestHelper struct {
	t    *testing.T
	a    *App
	name string
}

func (h *installerTestHelper) Run(ctx context.Context) {
	t, a, ls := h.t, h.a, h.name
	t.Helper()
	i, err := a.getInstaller(ls)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Install(ctx, ls); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(i.Dir()); err != nil {
			t.Fatal(err)
		}
	})
	info, err := os.Stat(filepath.Join(i.Dir(), i.BinName()))
	if err != nil {
		t.Fatal(err)
	}
	if !isExecutable(info.Mode()) {
		t.Fatal(ls + " is not executable")
	}
	if err := a.Uninstall(ctx, ls); err != nil {
		t.Fatal(err)
	}
	files, err := ioutil.ReadDir(filepath.Join(a.baseDir, ls))
	if err != nil {
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			t.Fatal(pathErr)
		}
	}
	if len(files) != 0 {
		t.Fatal("files should not exists", files)
	}
}

func newInstallerTestHelper(t *testing.T, name string) *installerTestHelper {
	t.Helper()
	if testing.Short() {
		t.Skip()
	}
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tmp); err != nil {
			t.Fatal(err)
		}
	})
	a, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	a.baseDir = tmp
	return &installerTestHelper{t: t, a: a, name: name}
}
