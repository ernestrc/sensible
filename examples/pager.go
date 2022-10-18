package main

import (
	"strings"

	"github.com/ernestrc/sensible/pager"
)

func main() {
	if err := pager.PageReader(strings.NewReader("a\nb\nc")); err != nil {
		panic(err)
	}
}
