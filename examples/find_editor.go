package main

import (
	"fmt"

	editor "github.com/ernestrc/sensible-editor"
)

func main() {
	editor, err := editor.FindEditor()

	if err != nil {
		panic(err)
	}

	fmt.Println(editor.GetPath())
}
