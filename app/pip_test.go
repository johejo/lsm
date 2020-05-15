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

func Test_isSupportedPython(t *testing.T) {
	tests := []struct {
		python string
		want   bool
	}{
		{"3.4.10", false},
		{"3.5.8", true},
		{"3.6", true},
		{"3.6.0", true},
		{"3.6.9", true},
		{"3.7.7", true},
		{"3.8.3", true},
	}
	for _, tt := range tests {
		t.Run(tt.python, func(t *testing.T) {
			ok, err := isSupportedPython(tt.python)
			if tt.want != ok {
				t.Fatalf("want=%v, got=%v, err=%v", tt.want, ok, err)
			}
		})
	}
}
