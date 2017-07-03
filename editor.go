// Package editor is a collection of utilities to find and spawn a sensible editor
package editor

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const bugMessage = "This is a bug in sensible-editor; please file at https://github.com/ernestrc/sensible-editor/issues"

// inspired by i3-sensible-editor
// The order has been altered to make the world a better place
var editors = []string{"$EDITOR", "$VISUAL", "vim", "nvim", "vi", "emacs", "nano", "pico", "qe", "mg", "jed", "gedit", "mc-edit"}
var basePath = []string{"/usr/local/bin", "/usr/bin", "/usr/sbin", "/bin"}

var userPath []string
var selectedExec string
var selectedArgs []string
var selectedEditor *Editor

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

func getFileName(f os.FileInfo) string {
	_, fileName := filepath.Split(f.Name())
	return fileName
}

func isRegularOrSymlink(finfo os.FileInfo) bool {
	mode := finfo.Mode()
	return mode.IsRegular() || mode&os.ModeSymlink != 0
}

func parseAlias(alias string) (name string, args []string) {
	split := strings.Split(alias, " ")
	if len(split) == 0 {
		return "", nil

	}
	_, name = filepath.Split(split[0])
	return name, split[1:]
}

func findExec(alias string) (execPath string, execArgs []string, err error) {
	var files []os.FileInfo
	name, args := parseAlias(alias)

	for _, dir := range userPath {
		if files, err = ioutil.ReadDir(dir); err != nil {
			return
		}
		for _, finfo := range files {
			if isRegularOrSymlink(finfo) &&
				isExecutable(finfo) &&
				getFileName(finfo) == name {
				execPath = path.Join(dir, name)
				execArgs = args
				return
			}
		}
	}
	return "", nil, nil
}

func (e *Editor) clean() {
	e.proc = nil
	e.procState = nil
}

func findEditor(editors []string) (editor *Editor, err error) {
	// cached
	if selectedExec != "" {
		if selectedArgs == nil {
			panic(fmt.Sprintf("parsed args is empty but selected has been cached. %s", bugMessage))
		}
		return NewEditor(selectedExec, selectedArgs...), nil
	}
	for _, editor := range editors {
		selectedExec, selectedArgs, err = findExec(editor)
		if err != nil {
			return nil, err
		}
		if selectedExec != "" {
			return NewEditor(selectedExec, selectedArgs...), nil
		}
	}

	return nil, fmt.Errorf("FindEditor: could not find an editor; please set $VISUAL or $EDITOR environment variables or install one of the following editors: %v", editors)
}

// NewEditor will create a new Editor struct with the given executable path
func NewEditor(abspath string, args ...string) *Editor {
	return &Editor{path: abspath, Args: args}
}

// FindEditor will attempt to find the user's preferred editor
// by scanning the PATH in search of EDITOR and VISUAL env variables
// or will default to one of the commonly installed editors.
// Failure to find a suitable editor will result in an error
func FindEditor() (editor *Editor, err error) {
	return findEditor(editors)
}

// Edit will attempt to edit the passed files with the user's preferred editor.
// Check the documentation of Editor.Edit and FindEditor for more information.
func Edit(files ...*os.File) error {
	var err error
	if selectedEditor == nil {
		if selectedEditor, err = FindEditor(); err != nil {
			return err
		}
	}

	return selectedEditor.Edit(files...)
}

// Editor stores the information about an editor and its processes
type Editor struct {
	path      string
	proc      *os.Process
	procState *os.ProcessState
	// extra arguments to be passed to the editor process before filename(s)
	Args []string
	// extra process attributes to be used when spawning editor process
	ProcAttrs *os.ProcAttr
}

// GetPath returns the editors executable path
func (e *Editor) GetPath() string {
	return e.path
}

// Edit will start a new process and wait for the process to exit.
// If process exists with non 0 status, this will be reported as an error
func (e *Editor) Edit(files ...*os.File) error {
	var err error

	if err = e.Start(files...); err != nil {
		return err
	}

	if err = e.Wait(); err != nil {
		return err
	}

	return nil
}

// Start will start a new process and pass the list of files as arguments
func (e *Editor) Start(f ...*os.File) error {
	if e.proc != nil {
		return fmt.Errorf("Editor.Start: there is already an ongoing session")
	}

	args := []string{""}
	var fds = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	if e.Args != nil {
		for _, arg := range e.Args {
			args = append(args, arg)
		}
	}

	for _, file := range f {
		args = append(args, file.Name())
		fds = append(fds, file)
	}

	var procAttrs *os.ProcAttr
	if e.ProcAttrs == nil {
		procAttrs = &os.ProcAttr{
			Dir:   "",
			Env:   nil,
			Files: fds,
			Sys:   nil,
		}
	} else {
		procAttrs = e.ProcAttrs
	}

	var err error
	if e.proc, err = os.StartProcess(e.path, args, procAttrs); err != nil {
		return err
	}

	return nil
}

// Wait waits for the current editor process to exit and returns
// an error if editor exited with non 0 status
func (e *Editor) Wait() error {
	var err error

	if e.proc == nil {
		return fmt.Errorf("Editor.Wait: no process is currently running")
	}

	if e.procState, err = e.proc.Wait(); err != nil {
		return err
	}

	if !e.procState.Success() {
		return fmt.Errorf("Editor.Wait: editor process exited with non 0 status: %s", e.procState.String())
	}

	e.clean()

	return nil
}

// EditTmp will place the contents of in in a temp file,
// start a editor process to edit the tmp file, and return
// the contents of the tmp file after the process exits, or an error
// if editor exited with non 0 status
func (e *Editor) EditTmp(in string) (out string, err error) {
	var f *os.File
	var outBytes []byte

	if f, err = ioutil.TempFile("/tmp", "sedit_"); err != nil {
		return
	}

	if err = ioutil.WriteFile(f.Name(), []byte(in), 0600); err != nil {
		return
	}

	if err = e.Edit(f); err != nil {
		return
	}

	if outBytes, err = ioutil.ReadFile(f.Name()); err != nil {
		return
	}

	out = string(outBytes)

	return
}
