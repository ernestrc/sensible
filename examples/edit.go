package main

import (
	"fmt"

	"github.com/ernestrc/sensible/editor"
)

func main() {
	var out string
	var err error
	if out, err = editor.EditTmp("what do you want to print?"); err != nil {
		panic(err)
	}

	fmt.Println(out)
}
