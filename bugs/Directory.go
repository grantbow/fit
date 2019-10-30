package bugs

import (
	"os"
	"regexp"
	"strings"
	"time"
)

// Directory type is a string path name.
type Directory string

func findIssuesDir(dir string, config *Config) Directory {
	if dirinfo, err := os.Stat(dir); err == nil && dirinfo.IsDir() {
		if dirinfo, err = os.Stat(dir + sops + "fit"); err == nil && dirinfo.IsDir() {
			// has a fit dir
			config.BugDir = dir
			config.IssuesDirName = "fit"
			os.Chdir(dir)
			return Directory(dir)
		} else if dirinfo, err = os.Stat(dir + sops + "issues"); err == nil && dirinfo.IsDir() {
			// has an issues dir
			config.BugDir = dir
			config.IssuesDirName = "issues"
			os.Chdir(dir)
			return Directory(dir)
		}
		// better to fall through and start looking rather than
		//} else {
		//	return ""
	}
	return ""
}

// RootDirer returns the directory usually containing the issues subdirectory.
func RootDirer(config *Config) Directory {
	dir := os.Getenv("FIT") // new first
	if dir != "" {
		if x := findIssuesDir(dir, config); x != "" {
			return x
		}
	} else {
		dir = os.Getenv("PMIT") // for backwards compatibility
		if dir != "" {
			if x := findIssuesDir(dir, config); x != "" {
				return x
			}
		}
	}

	wd, _ := os.Getwd()

	if x := findIssuesDir(wd, config); x != "" {
		return x
	}

	// There's no environment variable and no issues
	// directory, so walk up the tree until we find one
	pieces := strings.Split(wd, sops) // sops is string(os.PathSeparator)

	for i := len(pieces); i > 0; i -= 1 {
		dir := strings.Join(pieces[0:i], sops)
		if x := findIssuesDir(dir, config); x != "" {
			return x
			//if dirinfo, err := os.Stat(dir + sops + "issues"); err == nil && dirinfo.IsDir() {
			//	config.BugDir = dir
			//	os.Chdir(dir)
			//	return Directory(dir)
		}
	}
	return "" // out of luck
}

// IssuesDirer returns the directory containing the issues.
// The root directory contains the issues directory.
func IssuesDirer(config Config) Directory {
	root := RootDirer(&config)
	if root == "" {
		return root
	}
	return Directory(root + dops + "issues") // dops is Directory(string(os.PathSeparator))
	/* edited the following
	   when changed from /issues/ to /issues
	   $ grep -ils issuesdirer ...
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

// ShortNamer returns the directory name of a bug
func (d Directory) ShortNamer() Directory {
	pieces := strings.Split(string(d), sops)
	return Directory(pieces[len(pieces)-1])
}

// ToTitle decodes the human string from the filesystem directory name.
func (d Directory) ToTitle() string {
	multidash := regexp.MustCompile("([_]*)-([-_]*)")
	dashReplacement := strings.Replace(string(d), " ", sops, -1)
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
		panic("Directory " + string(d) + " stat error : " + err.Error())
	}

	if stat.IsDir() == false {
		return stat.ModTime()
	}

	dir, _ := os.Open(string(d))
	defer dir.Close() // discards error for now
	files, _ := dir.Readdir(-1)
	if len(files) == 0 {
		t = stat.ModTime()
	}
	for _, file := range files {
		if file.IsDir() {
			mtime := (d + dops + Directory(file.Name())).ModTime()
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
