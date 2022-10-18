// Package pager is a collection of utilities to find and spawn a sensible pagination client (i.e. less)
package pager

import (
	"fmt"
	"io"
	"os"

	"github.com/ernestrc/sensible/find"
)

var pagers = []string{ /* "$PAGER", "$GIT_PAGER", */ "less", "more"}

var selectedExec string
var selectedPager *Pager

func init() {
	if pager := os.Getenv("PAGER"); pager != "" {
		pagers = append([]string{pager}, pagers...)
	}
	if pager := os.Getenv("GIT_PAGER"); pager != "" {
		pagers = append([]string{pager}, pagers...)
	}
}

func (e *Pager) clean() {
	e.proc = nil
	e.procState = nil
}

func findPager(pagers []string) (pager *Pager, err error) {
	// cached
	if selectedExec != "" {
		return NewPager(selectedExec), nil
	}
	for _, pager := range pagers {
		selectedExec, err = find.Executable(pager)
		if err != nil {
			return nil, err
		}
		if selectedExec != "" {
			return NewPager(selectedExec), nil
		}
	}

	return nil, fmt.Errorf("FindPager: could not find an pager; please set $PAGER environment variable or install one of the following pagers: %v", pagers)
}

// NewPager will create a new Pager struct with the given executable path
func NewPager(abspath string, args ...string) *Pager {
	return &Pager{path: abspath, Args: args}
}

// FindPager will attempt to find the user's preferred pager
// by scanning the PATH in search of PAGER and GIT_PAGER env variables
// or will default to one of the commonly installed pagers.
// Failure to find a suitable pager will result in an error
func FindPager() (pager *Pager, err error) {
	return findPager(pagers)
}

// PageReader will attempt to view the given reader with the user's preferred pager.
// Check the documentation of Pager.PageReader and FindPager for more information.
func PageReader(r io.Reader) error {
	var err error
	if selectedPager == nil {
		if selectedPager, err = FindPager(); err != nil {
			return err
		}
	}

	return selectedPager.PageReader(r)
}

// Pager stores the information about a pager and its processes
type Pager struct {
	path      string
	proc      *os.Process
	procState *os.ProcessState
	// extra arguments to be passed to the pager process before filename(s)
	Args []string
	// extra process attributes to be used when spawning pager process
	ProcAttrs *os.ProcAttr
}

// GetPath returns the pagers executable path
func (e *Pager) GetPath() string {
	return e.path
}

// PageReader will start a new pager process render the given data,
// and wait for the process to exit. If process exists with non 0 status,
// this will be reported as an error
func (e *Pager) PageReader(r io.Reader) error {
	var err error

	if err = e.Start(r); err != nil {
		return err
	}

	if err = e.Wait(); err != nil {
		return err
	}

	return nil
}

// Start will start a new process and pass the given io.Reader
// to the Pager's standard input for it to render it.
func (e *Pager) Start(r io.Reader) error {
	if e.proc != nil {
		return fmt.Errorf("Pager.Start: there is already an ongoing session")
	}

	f, err := os.CreateTemp("", "sensible_pager_")
	if err != nil {
		return err
	}

	defer f.Close()
	defer os.Remove(f.Name())

	if _, err = io.Copy(f, r); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	args := []string{""}
	var fds = []*os.File{f, os.Stdout, os.Stderr}

	for _, arg := range e.Args {
		args = append(args, arg)
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

	if e.proc, err = os.StartProcess(e.path, args, procAttrs); err != nil {
		return err
	}

	return nil
}

// Wait waits for the current pager process to exit and returns
// an error if pager exited with non 0 status.
func (e *Pager) Wait() error {
	var err error

	if e.proc == nil {
		return fmt.Errorf("Pager.Wait: no process is currently running")
	}

	if e.procState, err = e.proc.Wait(); err != nil {
		return err
	}

	if !e.procState.Success() {
		return fmt.Errorf("Pager.Wait: pager process exited with non 0 status: %s", e.procState.String())
	}

	e.clean()

	return nil
}
