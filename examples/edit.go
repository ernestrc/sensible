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

	var out string
	if out, err = editor.EditTmp("what do you want to print?"); err != nil {
		panic(err)
	}

	fmt.Println(out)
}
