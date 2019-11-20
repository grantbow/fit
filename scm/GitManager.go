package scm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// GitManager type has fields Autoclose and UseBugPrefix.
type GitManager struct {
	Autoclose    bool
	UseBugPrefix bool
}

// Purge runs git clean -fd on the directory containing the issues directory.
func (mgr GitManager) Purge(dir bugs.Directory) error {
	cmd := exec.Command("git", "clean", "-fd", string(dir)+sops)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type issueStatus struct {
	a, d, m bool // Added, Deleted, Modified
}

type issuesStatus map[string]issueStatus

// Get list of created, updated, closed and closed-on-github issues.
//
// In general following rules to categorize issues are applied:
// * closed if Description file is deleted (D);
// * created if Description file is created (A) (TODO: handle issue renamings);
// * closed issue will also close issue on GH when Autoclose is true (see Identifier example);
// * updated if Description file is modified (M);
// * updated if Description is unchanged but any other files are touched. (' '+x)
//
// eg output from `from git status --porcelain`, appendix mine
// note that `git add -A issues` was invoked before
//
// D  issues/First-GH-issue/Description		issue closed (GH issues are also here)
// D  issues/First-GH-issue/Identifier		maybe it is GH issue, maybe not
// M  issues/issue--2/Description		desc updated
// A  issues/issue--2/Status			new field added (status); considered as update unless Description is also created
// D  issues/issue1/Description			issue closed
// A  issues/issue3/Description			new issue, description field is mandatory for rich format

func (mgr GitManager) currentStatus(dir bugs.Directory, config bugs.Config) (closedOnGitHub []string, _ issuesStatus) {
	ghRegex := regexp.MustCompile("(?im)^-Github:(.*)$")
	closesGH := func(file string) (issue string, ok bool) {
		if !mgr.Autoclose {
			return "", false
		}
		if !strings.HasSuffix(file, "Identifier") {
			return "", false
		}
		diff := exec.Command("git", "diff", "--staged", "--", file)
		diffout, _ := diff.CombinedOutput()
		matches := ghRegex.FindStringSubmatch(string(diffout))
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1]), true
		}
		return "", false
	}
	short := func(path string) string {
		beg := strings.Index(path, sops)
		end := strings.LastIndex(path, sops)
		if beg+1 >= end {
			return "???"
		}
		return path[beg+1 : end]
	}

	cmd := exec.Command("git", "status", "-z", "--porcelain", string(dir))
	out, _ := cmd.CombinedOutput()
	files := strings.Split(string(out), "\000")

	issues := issuesStatus{}
	var ghClosed []string
	const minLineLen = 3 /*for path*/ + 2 /*for issues dir with path sep*/ + 3 /*for issue name, path sep and any file under issue dir*/
	for _, file := range files {
		if len(file) < minLineLen {
			continue
		}

		path := file[3:]
		op := file[0]
		desc := strings.HasSuffix(path, sops+config.DescriptionFileName)
		name := short(path)
		issue := issues[name]

		switch {
		case desc && op == 'D':
			issue.d = true
		case desc && op == 'A':
			issue.a = true
		default:
			issue.m = true
			if op == 'D' {
				if ghIssue, ok := closesGH(path); ok {
					ghClosed = append(ghClosed, ghIssue)
					issue.d = true // to be sure
				}
			}
		}

		issues[name] = issue
	}
	return ghClosed, issues
}

// Create commit message by iterating over issues in order:
// closed issues are most important (something is DONE, ok? ;), those issues will also become hidden)
// new issues are next, with just updates at the end
// TODO: do something if this message will be too long
func (mgr GitManager) commitMsg(dir bugs.Directory, config bugs.Config) []byte {
	ghClosed, issues := mgr.currentStatus(dir, config)

	done, add, update, together := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	var cntd, cnta, cntu int

	for issue, state := range issues {
		if state.d {
			fmt.Fprintf(done, ", %q", issue)
			cntd++
		} else if state.a {
			fmt.Fprintf(add, ", %q", issue)
			cnta++
		} else if state.m {
			fmt.Fprintf(update, ", %q", issue)
			cntu++
		}
	}

	outf := func(buf *bytes.Buffer, what string, many bool) {
		if buf.Len() == 0 {
			return
		}
		var plural string
		if many {
			plural = "s:"
		}
		item := buf.Bytes()[2:]
		fmt.Fprintf(together, "%s issue%s %s; ", what, plural, item)
	}
	outf(done, "Close", cntd > 1)
	outf(add, "Create", cnta > 1)
	outf(update, "Update", cntu > 1)
	if l := together.Len(); l > 0 {
		together.Truncate(l - 2) // "; " from last applied outf()
	}

	if len(ghClosed) > 0 {
		fmt.Fprintf(together, "\n\nCloses %s\n", strings.Join(ghClosed, ", closes "))
	}
	return together.Bytes()
}

// Commit saves files to the SCM. It runs git add -A.
func (mgr GitManager) Commit(dir bugs.Directory, backupCommitMsg string, config bugs.Config) error {
	cmd := exec.Command("git", "add", "-A", string(dir))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Could not add issues to be committed: %s?\n", err.Error())
		return err
	}

	msg := mgr.commitMsg(dir, config)

	file, err := ioutil.TempFile("", "bugCommit")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create temporary file.\nNothing committed.\n")
		return err
	}
	defer os.Remove(file.Name())

	if len(msg) == 0 {
		fmt.Fprintf(file, "%s\n", backupCommitMsg)
	} else {
		var pref string
		if mgr.UseBugPrefix {
			pref = "issue: "
		}
		fmt.Fprintf(file, "%s%s\n", pref, msg)
	}
	//fmt.Print("debug commit : git", "commit", "-o", string(dir), "-F", file.Name(), "-q\n")
	cmd = exec.Command("git", "commit", "-o", string(dir), "-F", file.Name(), "-q")
	if err := cmd.Run(); err != nil {
		// If nothing was added commit will have an error.
		// in some cases we didn't care, it just meant there's nothing to commit.
		// the stdout to test could be captured
		//fmt.Printf("No new issues committed.\n") // assumed this error incorrectly, same for HgManager
		fmt.Printf("git commit error %v\n", err.Error()) // $?
		return err
	} else {
		return nil
	}
}

// SCMTyper returns "git".
func (mgr GitManager) SCMTyper() string {
	return "git"
}

// SCMIssuesUpdaters returns []byte of uncommitted files staged AND working directory
func (mgr GitManager) SCMIssuesUpdaters() ([]byte, error) { // config bugs.Config
	cmd := exec.Command("git", "status", "--porcelain", "-u", "--", ":"+sops+"issues")
	// --porcelain output format
	// -u shows all unstaged files, not just directories
	// issues is the directory off of the git repo to show
	// the ":(top)" shows full paths when not at the git root directory
	// the shorthand is ":"+sops+"issues"
	co, _ := cmd.CombinedOutput()
	if string(co) == "" {
		return []byte(""), nil
	} else {
		return co, errors.New("Files In issues/ Need Committing")
	}
}

// SCMIssuesCacher returns []byte of uncommitted files staged NOT working directory
func (mgr GitManager) SCMIssuesCacher() ([]byte, error) { // config bugs.Config
	cmd := exec.Command("git", "diff", "--name-status", "--cached", "HEAD", "--", ":"+sops+"issues")
	// only whitespace differs from output of git status
	co, _ := cmd.CombinedOutput()
	if string(co) == "" {
		return []byte(""), nil
	} else {
		return co, errors.New("Files In issues/ Staged and Need Committing")
	}
}

//func (mgr GitManager) SCMChangedIssues() ([]byte, error) {
//output from SCMIssuesCacher(), accept unique first directory level of the file list
//then check if updates should be sent for these issues
//    then send the updates
//}
