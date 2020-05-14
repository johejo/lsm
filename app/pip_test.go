package app

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lookPython(t *testing.T) {
	t.Run("look python3", func(t *testing.T) {
		if isWindows {
			t.Skip()
		}
		baseDir := filepath.Clean("./testdata")
		i := NewPipInstaller(baseDir, "python-language-server", "pyls")
		assert.Equal(t, i.python, "python3")
	})
}
