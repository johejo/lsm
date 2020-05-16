package app

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ci() bool {
	ci, _ := strconv.ParseBool(os.Getenv("CI"))
	return ci
}

func skipCI(t *testing.T, skip bool) {
	t.Helper()
	if skip && ci() {
		t.Skip()
	}
}

func TestApp_InstallAll(t *testing.T) {
	skipCI(t, false)
	h := newInstallerTestHelper(t)
	for k := range h.a.installers {
		k, h := k, h
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			h.Run(context.Background(), k)
		})
	}
}

func TestGoInstaller(t *testing.T) {
	skipCI(t, true)
	tests := []string{
		"gopls",
		"sqls",
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt, func(t *testing.T) {
			t.Parallel()
			h := newInstallerTestHelper(t)
			h.Run(context.Background(), tt)
		})
	}
}

func TestNpmInstaller(t *testing.T) {
	skipCI(t, true)
	tests := []string{
		"typescript-language-server",
		"vim-language-server",
		"dockerfile-language-server-nodejs",
		"bash-language-server",
		"yaml-language-server",
		"vscode-json-languageserver",
		"vscode-css-languageserver",
		"vscode-html-languageserver",
		"vls",
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt, func(t *testing.T) {
			t.Parallel()
			h := newInstallerTestHelper(t)
			h.Run(context.Background(), tt)
		})
	}
}

func TestMetalsInstaller(t *testing.T) {
	skipCI(t, true)
	h := newInstallerTestHelper(t)
	h.Run(context.Background(), "metals")
}

func TestKotlinLSInstaller(t *testing.T) {
	skipCI(t, true)
	h := newInstallerTestHelper(t)
	h.Run(context.Background(), "kotlin-language-server")
}

func TestRustAnalyzerInstaller(t *testing.T) {
	skipCI(t, true)
	h := newInstallerTestHelper(t)
	h.Run(context.Background(), "rust-analyzer")
}

func TestEfmLSInstaller(t *testing.T) {
	skipCI(t, true)
	h := newInstallerTestHelper(t)
	h.Run(context.Background(), "efm-langserver")
}

func TestPipInstaller(t *testing.T) {
	skipCI(t, true)
	tests := []string{
		"python-language-server",
		"fortran-language-server",
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt, func(t *testing.T) {
			t.Parallel()
			h := newInstallerTestHelper(t)
			h.Run(context.Background(), tt)
		})
	}
}

func TestTerraformLSPInstaller(t *testing.T) {
	skipCI(t, true)
	t.Parallel()
	h := newInstallerTestHelper(t)
	h.Run(context.Background(), "terraform-lsp")
}

func TestTerraformLSInstaller(t *testing.T) {
	skipCI(t, true)
	t.Parallel()
	h := newInstallerTestHelper(t)
	h.Run(context.Background(), "terraform-ls")
}

type installerTestHelper struct {
	t *testing.T
	a *App
}

func (h *installerTestHelper) Run(ctx context.Context, name string) {
	t, a := h.t, h.a
	i, err := a.getInstaller(name)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, isSupported(i))
	i.SetWriter(ioutil.Discard)
	if err := a.Install(ctx, name); err != nil {
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
		t.Fatal(i.BinName() + " is not executable")
	}
	if err := a.Uninstall(ctx, name); err != nil {
		t.Fatal(err)
	}
	files, err := ioutil.ReadDir(i.Dir())
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

func newInstallerTestHelper(t *testing.T) *installerTestHelper {
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
	return &installerTestHelper{t: t, a: a}
}
