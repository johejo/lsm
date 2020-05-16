package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew_unix_home(t *testing.T) {
	if isWindows {
		t.Skip()
	}
	p, err := filepath.Abs("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("HOME", p); err != nil {
		t.Fatal(err)
	}
	a, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, filepath.Join(p, ".local", "share", "lsm", "servers"), a.baseDir)
}

func TestNew_xdgDataHome(t *testing.T) {
	if isWindows {
		t.Skip()
	}
	p, err := filepath.Abs("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("XDG_DATA_HOME", p); err != nil {
		t.Fatal(err)
	}
	a, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, filepath.Join(p, "lsm", "servers"), a.baseDir)
}

func TestNew_windows(t *testing.T) {
	if !isWindows {
		t.Skip()
	}
	a, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	local := os.Getenv("LOCALAPPDATA")
	assert.Equal(t, filepath.Join(local, "lsm", "servers"), a.baseDir)
}

func TestApp_List(t *testing.T) {
	baseDir := filepath.Clean("./testdata/lsm/servers")

	t.Run("not installed any language servers", func(t *testing.T) {
		_ = os.RemoveAll(baseDir)
		a, err := New(baseDir)
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		a.out = &buf
		if err := a.List(context.Background(), ListStyleJSON); err != nil {
			t.Fatal(err)
		}
		var list []languageServer
		if err := json.NewDecoder(&buf).Decode(&list); err != nil {
			t.Fatal(err)
		}
		for _, ls := range list {
			assert.False(t, ls.Installed)
		}
	})

	t.Run("installed only one", func(t *testing.T) {
		const efmls = "efm-langserver"
		_ = os.RemoveAll(baseDir)
		a, err := New(baseDir)
		if err != nil {
			t.Fatal(err)
		}
		if err := a.Uninstall(context.Background(), efmls); err != nil {
			t.Fatal(err)
		}
		if err := a.Install(context.Background(), efmls); err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		a.out = &buf
		if err := a.List(context.Background(), ListStyleJSON); err != nil {
			t.Fatal(err)
		}
		var list []languageServer
		if err := json.NewDecoder(&buf).Decode(&list); err != nil {
			t.Fatal(err)
		}
		for _, ls := range list {
			if ls.Name == efmls {
				assert.True(t, ls.Installed)
			} else {
				assert.False(t, ls.Installed)
			}
		}
	})

	t.Run("table", func(t *testing.T) {
		_ = os.RemoveAll(baseDir)
		a, err := New(baseDir)
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		a.out = &buf
		if err := a.List(context.Background(), ListStyleTable); err != nil {
			t.Fatal(err)
		}
	})
}
