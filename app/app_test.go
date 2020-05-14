package app

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, `%LOCALAPPDATA%\lsm\servers`, a.baseDir)
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
		if err := a.List(context.Background()); err != nil {
			t.Fatal(err)
		}
		s := bufio.NewScanner(&buf)
		for s.Scan() {
			found := false
			line := s.Text()
			assert.Contains(t, line, "not")
			for _, i := range a.installers {
				if strings.Contains(line, i.Name()) {
					found = true
					break
				}
			}
			if !found {
				t.Error("line should contain language server name", line)
			}
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
		if err := a.List(context.Background()); err != nil {
			t.Fatal(err)
		}
		s := bufio.NewScanner(&buf)
		for s.Scan() {
			line := s.Text()
			if strings.Contains(line, efmls) {
				assert.NotContains(t, line, "not")
			} else {
				assert.Contains(t, line, "not")
			}
		}
	})
}
