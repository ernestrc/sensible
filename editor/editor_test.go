package editor

import (
	"reflect"
	"testing"
)

func resetGlobal() {
	selectedExec = ""
	selectedEditor = nil
}

func assertEditor(t *testing.T, ed *Editor, path string, args []string) {
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
	editors := []string{"vi"}
	if ed, err := findEditor(editors); err != nil {
		t.Error(err)
	} else {
		assertEditor(t, ed, "/usr/bin/vi", nil)
	}
}
