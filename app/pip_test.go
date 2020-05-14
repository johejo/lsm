package app

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lookPython(t *testing.T) {
	b, err := exec.Command("python", "-V").Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
	if !isWindows {
		t.Skip()
	}
	baseDir := filepath.Clean("./testdata")
	i := NewPipInstaller(baseDir, "python-language-server", "pyls")
	assert.Equal(t, i.python, "python")
}

func Test_lookPython3(t *testing.T) {
	if isWindows {
		t.Skip()
	}
	baseDir := filepath.Clean("./testdata")
	i := NewPipInstaller(baseDir, "python-language-server", "pyls")
	assert.Equal(t, i.python, "python3")
}
