# sensible-editor
Utilities to find and spawn a sensible editor. 

EDITOR and VISUAL environment variables are used first to find the user's preferred editor. If these are not set, a list of commonly installed editors is used.

# Usage
Find a sensible editor
```go
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
```


Edit a file

```go
package main

import (
	"io/ioutil"

	editor "github.com/ernestrc/sensible-editor"
)

func main() {
	f1, _ := ioutil.TempFile("/tmp", "")
	if err := editor.Edit(f1); err != nil {
		panic(err)
	}
}
```

`go run` some of the [examples](examples) to see more advanced use-cases.
