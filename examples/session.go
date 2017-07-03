package main

import (
	"fmt"

	editor "github.com/ernestrc/sensible-editor"
)

func main() {
	var err error
	var out string

	in := "what do you want to print?"

	if out, err = editor.NewSession(in); err != nil {
		panic(err)
	}

	fmt.Println(out)
}
