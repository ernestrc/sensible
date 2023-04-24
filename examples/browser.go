package main

import (
	"log"
	"net/url"

	"github.com/ernestrc/sensible/browser"
)

func main() {
	url, err := url.Parse("https://unstable.build")
	if err != nil {
		panic(err)
	}

	if err := browser.Browse(url); err != nil {
		log.Fatal(err)
	}
}
