package editor

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
)

var basePath = []string{"/usr/local/bin", "/usr/bin", "/usr/sbin", "/bin"}
var editors = []string{"$EDITOR", "$VISUAL", "nvim", "vim", "emacs", "nano", "vi", "pico", "qe", "mg", "jed", "gedit", "mc-edit"}

var userPath []string
var selected string

func init() {
	editors[0] = os.Getenv("EDITOR")
	editors[1] = os.Getenv("VISUAL")

	pathEnv := os.Getenv("PATH")

	if pathEnv == "" {
		userPath = basePath
	} else {
		userPath = strings.Split(pathEnv, ":")
	}
}

func isExecutable(f os.FileInfo) bool {
	return f.Mode().Perm()|0111 != 0
}

func findExec(name string) (execPath string, err error) {
	var files []os.FileInfo

	for _, dir := range userPath {
		if files, err = ioutil.ReadDir(dir); err != nil {
			return
		}
		for _, file := range files {
			if file.Mode().IsRegular() &&
				isExecutable(file) &&
				file.Name() == name {
				execPath = path.Join(dir, name)
				return
			}
		}
	}
	return "", nil
}

func tmpFile() string {
	return fmt.Sprintf("/tmp/sedit_%d", rand.Int())
}

// FindEditor TODO ...
func FindEditor() (string, error) {
	var err error
	if selected != "" {
		return selected, nil
	}
	for _, editor := range editors {
		selected, err = findExec(editor)
		if err != nil {
			return "", err
		}
		if selected != "" {
			return selected, nil
		}
	}

	return "", fmt.Errorf("could not find an editor; please set $VISUAL or $EDITOR environment variables or install one of the preferred editors: %v", editors)
}

// NewSession TODO ...
func NewSession(in string) (out string, err error) {
	var path string
	var f *os.File
	var p *os.Process
	var s *os.ProcessState
	var outBytes []byte

	if path, err = FindEditor(); err != nil {
		return
	}

	if f, err = ioutil.TempFile("/tmp", "sedit_"); err != nil {
		return
	}

	if err = ioutil.WriteFile(f.Name(), []byte(in), 0600); err != nil {
		return
	}

	args := []string{"", f.Name()}

	var fds = []*os.File{os.Stdin, os.Stdout, os.Stderr, f}
	var procAttrs = os.ProcAttr{
		Dir:   "",
		Env:   nil,
		Files: fds,
		Sys:   nil,
	}

	if p, err = os.StartProcess(path, args, &procAttrs); err != nil {
		return
	}

	if s, err = p.Wait(); err != nil {
		return
	}

	if !s.Success() {
		err = fmt.Errorf("editor process exited with non 0 status: %s", s.String())
		return
	}

	if outBytes, err = ioutil.ReadFile(f.Name()); err != nil {
		return
	}

	out = string(outBytes)

	return
}
