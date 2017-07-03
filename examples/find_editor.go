package main

import (
	"fmt"

	editor "github.com/ernestrc/sensible-editor"
)

func main() {
	e, err := editor.FindEditor()

	if err != nil {
		panic(err)
	}

	fmt.Println(e.GetPath())
}
