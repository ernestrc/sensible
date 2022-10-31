package find

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	multierr "github.com/ernestrc/go-multierror"
)

var (
	userPath []string
	basePath = []string{"/usr/local/bin", "/usr/bin", "/usr/sbin", "/bin"}
)

func init() {
	pathEnv := os.Getenv("PATH")

	if pathEnv == "" {
		userPath = basePath
	} else {
		userPath = strings.Split(pathEnv, ":")
	}
}

func isExecutable(f os.FileMode) bool {
	return f.Perm()|0111 != 0
}

func isRegularOrSymlink(mode os.FileMode) bool {
	return mode.IsRegular() || mode&os.ModeSymlink != 0
}

// Executable finds the given executable in the host. It returns the executable's
// path or an error if the given executable could not be found on the host along
// with any errors incurred while reading the directories in PATH.
func Executable(name string) (execPath string, ret error) {
	if name == "" {
		ret = errors.New("invalid argument: empty name")
		return
	}

	for _, dir := range userPath {
		files, err := os.ReadDir(dir)
		if err != nil {
			ret == multierr.Append(ret, err)
		}
		for _, entry := range files {
			if isRegularOrSymlink(entry.Type()) &&
				isExecutable(entry.Type()) &&
				entry.Name() == name {
				execPath = path.Join(dir, name)
				ret = nil
				return
			}
		}
	}
	ret = multierr.Append(ret, fmt.Errorf("could not find %q in path", name))
	return
}
