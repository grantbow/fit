package bugs

import (
	"os"
	"regexp"
	"strings"
	"time"
)

// Directory type is a string path name.
type Directory string

// GetRootDir returns the directory usually containing the issues subdirectory.
// TODO: should be GetIssuesDir ! no wonder !! fix here and everywhere called
func GetRootDir(config Config) Directory {
	dir := os.Getenv("PMIT")
	if dir != "" {
		// TODO: that PMIT dir exists is a bad assumption
		if dirinfo, err := os.Stat(string(dir)); err == nil && dirinfo.IsDir() {
			config.BugDir = dir
			os.Chdir(dir)
			return Directory(dir)
			// better to start looking rather than
			//} else {
			//	return ""
		}
	}

	wd, _ := os.Getwd()

	if dirinfo, err := os.Stat(wd + "/issues"); err == nil && dirinfo.IsDir() {
		config.BugDir = string(wd)
		//os.Chdir(dir) // already there
		return Directory(wd)
	}

	// There's no environment variable and no issues
	// directory, so walk up the tree until we find one
	pieces := strings.Split(wd, "/")

	for i := len(pieces); i > 0; i -= 1 {
		dir := strings.Join(pieces[0:i], "/")
		if dirinfo, err := os.Stat(dir + "/issues"); err == nil && dirinfo.IsDir() {
			config.BugDir = dir
			os.Chdir(dir)
			return Directory(dir)
		}
	}
	return "" // out of luck
}

// GetIssuesDir returns the directory containing the issues.
// The root directory contains the issues directory.
func GetIssuesDir(config Config) Directory {
	root := GetRootDir(config)
	if root == "" {
		return root
	}
	return Directory(root + "/issues/") // TODO: remove trailing /
	/* then edit these $ grep -ils getissuesdir ...
	bug-import/be.go
	bug-import/github.go
	bugapp/Commit.go
	bugapp/Create.go
	bugapp/Env.go
	bugapp/Find.go
	bugapp/List.go
	bugapp/Purge.go
	bugapp/Pwd.go
	bugapp/Retitle.go
	bugs/Bug.go
	bugs/Bug_test.go
	bugs/Directory.go
	bugs/Directory_test.go
	bugs/Find.go
	*/
}

// GetShortName returns the directory name of a bug
func (d Directory) GetShortName() Directory {
	pieces := strings.Split(string(d), "/")
	return Directory(pieces[len(pieces)-1])
}

// ToTitle decodes the human string from the filesystem directory name.
func (d Directory) ToTitle() string {
	multidash := regexp.MustCompile("([_]*)-([-_]*)")
	dashReplacement := strings.Replace(string(d), " ", "/", -1)
	return multidash.ReplaceAllStringFunc(dashReplacement, func(match string) string {
		if match == "-" {
			return " "
		}
		if strings.Count(match, "_") == 0 {
			return match[1:]
		}
		return strings.Replace(match, "_", " ", -1)
	})
}

// ModTime returns the last modified time from the file system.
func (d Directory) ModTime() time.Time {
	var t time.Time
	stat, err := os.Stat(string(d))
	if err != nil {
		panic("Directory " + string(d) + " is not a directory.")
	}

	if stat.IsDir() == false {
		return stat.ModTime()
	}

	dir, _ := os.Open(string(d))
	files, _ := dir.Readdir(-1)
	if len(files) == 0 {
		t = stat.ModTime()
	}
	for _, file := range files {
		if file.IsDir() {
			mtime := (d + "/" + Directory(file.Name())).ModTime()
			if mtime.After(t) {
				t = mtime
			}
		} else {
			mtime := file.ModTime()
			if mtime.After(t) {
				t = mtime
			}

		}
	}
	return t
}
