package main

import (
	"io/ioutil"

	"github.com/ernestrc/sensible/editor"
)

func main() {
	f1, _ := ioutil.TempFile("/tmp", "")
	if err := editor.Edit(f1); err != nil {
		panic(err)
	}
}
