package app

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lookPython(t *testing.T) {
	py3, err := exec.Command("python3", "-V").Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(py3))
	py, err := exec.Command("python", "-V").Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(py))
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
