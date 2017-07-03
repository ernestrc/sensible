package main

import (
	"io/ioutil"
	"os"

	editor "github.com/ernestrc/sensible-editor"
)

func main() {
	var err error
	var ed *editor.Editor

	if ed, err = editor.FindEditor(); err != nil {
		panic(err)
	}

	var f1, f2, f3 *os.File

	f1, _ = ioutil.TempFile("/tmp", "")
	f2, _ = ioutil.TempFile("/tmp", "")
	f3, _ = ioutil.TempFile("/tmp", "")

	if err = ed.Edit(f1, f2, f3); err != nil {
		panic(err)
	}
}
