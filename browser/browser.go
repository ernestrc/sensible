package browser

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ernestrc/sensible/find"
)

var browsers = []string{ /* "$BROWSER", */ "open", "google-chrome-stable", "firefox", "chromium"}

var selectedExec string
var selectedBrowser *Browser

func init() {
	if browser := os.Getenv("BROWSER"); browser != "" {
		browsers = append([]string{browser}, browsers...)
	}
}

func findBrowser(browsers []string) (browser *Browser, err error) {
	// cached
	if selectedExec != "" {
		return NewBrowser(selectedExec), nil
	}
	for _, browser := range browsers {
		selectedExec, err = find.Executable(browser)
		if err != nil {
			return nil, err
		}
		if selectedExec != "" {
			return NewBrowser(selectedExec), nil
		}
	}

	return nil, fmt.Errorf("FindBrowser: could not find an browser; please set $BROWSER environment variable or install one of the following browsers: %v", browsers)
}

// NewBrowser will create a new Browser struct with the given executable path
func NewBrowser(abspath string, args ...string) *Browser {
	return &Browser{path: abspath, Args: args}
}

// FindBrowser will attempt to find the user's preferred browser
// by scanning the PATH in search of BROWSER env variables
// or will default to one of the commonly installed browsers.
//
// Failure to find a suitable browser will result in an error
func FindBrowser() (browser *Browser, err error) {
	return findBrowser(browsers)
}

// Browse will attempt to open the passed urls with the user's preferred browser.
// Check the documentation of FindBrowser for more information.
func Browse(urls ...*url.URL) error {
	var err error
	if selectedBrowser == nil {
		if selectedBrowser, err = FindBrowser(); err != nil {
			return err
		}
	}

	return selectedBrowser.Browse(urls...)
}

// Browser stores the information about an browser and its processes
type Browser struct {
	path      string
	proc      *os.Process
	procState *os.ProcessState
	// extra arguments to be passed to the browser process before filename(s)
	Args []string
	// extra process attributes to be used when spawning browser process
	ProcAttrs *os.ProcAttr
}

// GetPath returns the browsers executable path
func (e *Browser) GetPath() string {
	return e.path
}

// Browse will start a new process and wait for the process to exit.
// If process exists with non 0 status, this will be reported as an error
func (e *Browser) Browse(urls ...*url.URL) error {
	var err error

	if err = e.Start(urls...); err != nil {
		return err
	}

	if err = e.Wait(); err != nil {
		return err
	}

	return nil
}

// Start will start a new process and pass the list of urls as arguments
func (e *Browser) Start(u ...*url.URL) error {
	if e.proc != nil {
		return fmt.Errorf("Browser.Start: there is already an ongoing session")
	}

	args := []string{""}
	for _, arg := range e.Args {
		args = append(args, arg)
	}

	var fds = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	for _, url := range u {
		args = append(args, url.String())
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
		return fmt.Errorf("start browser process: %v. "+
			"Make sure that $BROWSER environment variable is set correctly", err)
	}

	return nil
}

// Wait waits for the current browser process to exit and returns
// an error if browser exited with non 0 status
func (e *Browser) Wait() error {
	var err error

	if e.proc == nil {
		return fmt.Errorf("Browser.Wait: no process is currently running")
	}

	if e.procState, err = e.proc.Wait(); err != nil {
		return err
	}

	if !e.procState.Success() {
		return fmt.Errorf("Browser.Wait: browser process exited with non 0 status: %s", e.procState.String())
	}

	e.clean()

	return nil
}

func (e *Browser) clean() {
	e.proc = nil
	e.procState = nil
}
