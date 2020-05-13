package lsm

import (
	"bufio"
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestApp_List(t *testing.T) {
	a, err := NewApp("")
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
}
