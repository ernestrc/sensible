package find

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
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

func isExecutable(f os.FileInfo) bool {
	return f.Mode().Perm()|0111 != 0
}

func isRegularOrSymlink(finfo os.FileInfo) bool {
	mode := finfo.Mode()
	return mode.IsRegular() || mode&os.ModeSymlink != 0
}

// Executable finds the given executable in the host. It returns the executable's
// path and default arguments or an error if the given executable could not
// be found on the host.
func Executable(name string) (execPath string, err error) {
	if name == "" {
		err = errors.New("invalid argument: empty name")
		return
	}

	var files []os.FileInfo
	for _, dir := range userPath {
		if files, err = ioutil.ReadDir(dir); err != nil {
			return
		}
		for _, finfo := range files {
			if isRegularOrSymlink(finfo) &&
				isExecutable(finfo) &&
				finfo.Name() == name {
				execPath = path.Join(dir, name)
				return
			}
		}
	}
	return "", nil
}
