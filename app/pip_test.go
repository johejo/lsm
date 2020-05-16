package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lookPython(t *testing.T) {
	var want string
	if isWindows {
		want = "python"
	} else {
		want = "python3"
	}
	got, err := _lookPython()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, got, want)
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
