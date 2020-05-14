package app

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPipInstaller_lookPython(t *testing.T) {
	var want string
	if isWindows {
		want = "python"
	} else {
		want = "python3"
	}
	baseDir := filepath.Clean("./testdata")
	i := NewPipInstaller(baseDir, "python-language-server", "pyls")
	assert.Equal(t, want, i.python)
}
