package scm

import (
	"errors"
	"github.com/driusan/bug/bugs"
	"os"
	"strings"
)

type SCMNotFound struct {
	s string
}

func (e *SCMNotFound) Error() string {
	return e.s
}

type SCMDirty struct {
	s string
}

func (e *SCMDirty) Error() string {
	return e.s
}

func walkAndSearch(startpath string, dirnames []string) (fullpath, scmtype string) {
	for _, scmtype := range dirnames {
		if dirinfo, err := os.Stat(startpath + "/" + scmtype); err == nil && dirinfo.IsDir() {
			return startpath + "/" + scmtype, scmtype
		}
	}

	pieces := strings.Split(startpath, "/")

	for i := len(pieces); i > 0; i -= 1 {
		dir := strings.Join(pieces[0:i], "/")
		for _, scmtype := range dirnames {
			if dirinfo, err := os.Stat(dir + "/" + scmtype); err == nil && dirinfo.IsDir() {
				return dir + "/" + scmtype, scmtype
			}
		}
	}
	return "", ""
}

// DetectSCM takes options and returns an SCMHandler and directory.
func DetectSCM(options map[string]bool, config bugs.Config) (SCMHandler, bugs.Directory, error) {
	// First look for an SCM directory
	wd, _ := os.Getwd()

	dirFound, scmtype := walkAndSearch(wd, []string{".git", ".hg"})
	if dirFound != "" && scmtype == ".git" {
		var gm GitManager
		if val, ok := options["autoclose"]; ok {
			gm.Autoclose = val
		}
		if val, ok := options["use_bug_prefix"]; ok {
			gm.UseBugPrefix = val
		}
		return gm, bugs.Directory(dirFound), nil
	}
	if dirFound != "" && scmtype == ".hg" {
		return HgManager{}, bugs.Directory(dirFound), nil
	}

	return nil, "", errors.New("No SCM found")
}
