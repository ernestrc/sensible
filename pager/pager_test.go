package pager

import (
	"reflect"
	"testing"
)

func resetGlobal() {
	selectedExec = ""
	selectedPager = nil
}

func assertPager(t *testing.T, ed *Pager, path string, args []string) {
	if ed.path != path {
		t.Error("unexpected path", path, ed.path)
	}
	if len(ed.Args) == 0 && len(args) == 0 {
		return
	}
	if !reflect.DeepEqual(ed.Args, args) {
		t.Error("unexpected args:", ed.Args, args)
	}
}

// this tests will only work if vi is in PATH
func TestEnvPath(t *testing.T) {
	resetGlobal()
	pagers := []string{"less"}
	if ed, err := findPager(pagers); err != nil {
		t.Error(err)
	} else {
		assertPager(t, ed, "/usr/bin/less", nil)
	}
}
