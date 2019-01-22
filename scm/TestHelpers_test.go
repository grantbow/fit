package scm

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

type Commit interface {
	CommitID() string
	LogMsg() string
	Diff() (string, error)
}

type ManagerTester interface {
	GetLogs() ([]Commit, error)
	Setup() error
	GetWorkDir() string
	TearDown()
	StageFile(string) error
	AssertCleanTree(t *testing.T)
	AssertStagingIndex(*testing.T, []FileStatus)
	GetManager() SCMHandler
}

type argumentList []string

func runCmd(cmd string, options ...string) (string, error) {
	runcmd := exec.Command(cmd, options...)
	out, err := runcmd.CombinedOutput()

	return string(out), err
}

func assertLogs(tester ManagerTester, t *testing.T, titles []map[string]bool, diffs []string) {
	logs, err := tester.GetLogs()
	if err != nil {
		t.Error("Could not get scm logs" + err.Error())
		return
	}

	if len(diffs) != len(titles) {
		t.Errorf("Different number of diffs (%v) from titles(%v)", len(diffs), len(titles))
		return
	}
	if len(logs) != len(titles) || len(logs) != len(diffs) {
		t.Errorf("Unexpected, len(logs) %v != len(titles) %v || len(diffs) %v", len(logs), len(titles), len(diffs))
		return
	}

	for i := range titles {
		if _, ok := titles[i][logs[i].LogMsg()]; !ok {
			t.Error("Unexpected commit message:" + logs[i].LogMsg())
		}

		if diff, err := logs[i].Diff(); err != nil {
			t.Error("Could not get diff of commit")
		} else {
			if diff != diffs[i] {
				// get shortest commit msg to keep errors simple
				var s string
				for k := range titles[i] {
					if len(s) == 0 || len(k) < len(s) {
						s = k
					}
				}
				t.Error(fmt.Sprintf("Incorrect diff for i=%d, title=%s", i, s))
				fmt.Fprintf(os.Stderr, "Got %s, expected %s", diff, diffs[i])
			}
		}
	}
}

func runtestRenameCommitsHelper(tester ManagerTester, t *testing.T, expectedDiffs []string) {
	err := tester.Setup()
	defer tester.TearDown()
	if err != nil {
		t.Error("Could not initialize repo")
		return
	}

	m := tester.GetManager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	os.MkdirAll("issues/Test-bug", 0755)
	ioutil.WriteFile("issues/Test-bug/Description", []byte(""), 0644)
	m.Commit(bugs.Directory(tester.GetWorkDir()), "Initial commit")
	//runCmd("bug", "relabel", "1", "Renamed", "bug")
	config := bugs.Config{}
	args := argumentList{"1", "Renamed bug"}
	expected := "Moving .*"
	runrelabel("scm/TestHelpers_test runtestRenameCommitsHelper", args, config, expected, t)
	m.Commit(bugs.Directory(tester.GetWorkDir()), "This is a test rename")

	tester.AssertCleanTree(t)

	assertLogs(tester, t, []map[string]bool{{
		"Initial commit":          true, // simple format
		`Create issue "Test-bug"`: true, // rich format
	}, {
		"This is a test rename":                    true, // simple format
		`Update issues: "Test-bug", "Renamed-bug"`: true, // rich format
		`Update issues: "Renamed-bug", "Test-bug"`: true, // has two alternatives equally good
	}}, expectedDiffs)

}

func runtestCommitDirtyTree(tester ManagerTester, t *testing.T) {
	err := tester.Setup()
	if err != nil {
		panic("Something went wrong trying to initialize git:" + err.Error())
	}
	defer tester.TearDown()
	m := tester.GetManager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	os.MkdirAll("issues/Test-bug", 0755)
	ioutil.WriteFile("issues/Test-bug/Description", []byte(""), 0644)
	if err = ioutil.WriteFile("donotcommit.txt", []byte(""), 0644); err != nil {
		t.Error("Could not write file")
		return
	}
	tester.AssertStagingIndex(t, []FileStatus{
		FileStatus{"donotcommit.txt", "?", "?"},
	})

	//fmt.Print("pre  1 runtestCommitDirtyTree\n")
	m.Commit(bugs.Directory(tester.GetWorkDir()+"/issues"), "Initial commit")
	//fmt.Print("post 1 runtestCommitDirtyTree\n")
	tester.AssertStagingIndex(t, []FileStatus{
		FileStatus{"donotcommit.txt", "?", "?"},
	})
	tester.StageFile("donotcommit.txt")
	tester.AssertStagingIndex(t, []FileStatus{
		FileStatus{"donotcommit.txt", "A", " "},
	})
	//fmt.Print("pre  2 runtestCommitDirtyTree\n")
	m.Commit(bugs.Directory(tester.GetWorkDir()+"/issues"), "Initial commit")
	//errCommit := m.Commit(bugs.Directory(tester.GetWorkDir()+"/issues"), "Initial commit")
	//fmt.Printf("post 2 runtestCommitDirtyTree error %v\n", errCommit) // nil here
	//    running test shows output here. actually HgManager.go Commit() returns *expected* error not fully handled.
	//    stdout not captured this time.
	//fmt.Print("post 2 runtestCommitDirtyTree\n")
	tester.AssertStagingIndex(t, []FileStatus{
		FileStatus{"donotcommit.txt", "A", " "},
	})
}

func runtestPurgeFiles(tester ManagerTester, t *testing.T) {
	err := tester.Setup()
	if err != nil {
		panic("Something went wrong trying to initialize: " + err.Error())
	}
	defer tester.TearDown()
	m := tester.GetManager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	// Commit a bug which should stay around after the purge
	//runCmd("bug", "create", "-n", "Test", "bug")
	os.MkdirAll("issues/Test-bug", 0755)
	ioutil.WriteFile("issues/Test-bug/Description", []byte(""), 0644)
	m.Commit(bugs.Directory(tester.GetWorkDir()+"/issues"), "Initial commit")

	// Create another bug to elimate with purge
	//runCmd("bug", "create", "-n", "Test", "Purge", "bug")
	os.MkdirAll("issues/Test-Purge-Bug", 0755)
	err = m.Purge(bugs.Directory(tester.GetWorkDir() + "/issues"))
	if err != nil {
		t.Error("Error purging bug directory: " + err.Error())
	}
	issuesDir, err := ioutil.ReadDir("issues") //fmt.Sprintf("debug: %s/issues/", tester.GetWorkDir()))
	if err != nil {
		t.Error("Error reading issues directory: " + err.Error())
	}
	if len(issuesDir) != 1 {
		t.Errorf("Unexpected number of directories (%v, expected 1) in issues/ after purge.", len(issuesDir))
	}
	if len(issuesDir) > 0 && issuesDir[0].Name() != "Test-bug" {
		t.Error("Expected Test-bug to remain.")
	}
}

func runrelabel(label string, args argumentList, config bugs.Config, expected string, t *testing.T) {
	stdout, _ := captureOutput(func() {
		scmRelabel(args, config)
	}, t)
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Errorf("Unexpected output on STDOUT for %s.", label)
		fmt.Printf("Expected to match: %s\nGot: %s\n", expected, stdout)
	}
}

func captureOutput(f func(), t *testing.T) (string, string) {
	// Capture STDOUT with a pipe
	stdout := os.Stdout
	stderr := os.Stderr
	so, op, _ := os.Pipe() //outpipe
	oe, ep, _ := os.Pipe() //errpipe
	defer func(stdout, stderr *os.File) {
		os.Stdout = stdout
		os.Stderr = stderr
	}(stdout, stderr)

	os.Stdout = op
	os.Stderr = ep

	f()

	os.Stdout = stdout
	os.Stderr = stderr

	op.Close()
	ep.Close()

	errOutput, err := ioutil.ReadAll(oe)
	if err != nil {
		t.Error("Could not get output from stderr")
	}
	stdOutput, err := ioutil.ReadAll(so)
	if err != nil {
		t.Error("Could not get output from stdout")
	}
	return string(stdOutput), string(errOutput)
}

func scmRelabel(Args argumentList, config bugs.Config) {
	if len(Args) < 2 {
		fmt.Printf("Usage: %s relabel BugID New Title\n", os.Args[0])
		return
	}

	b, err := bugs.LoadBugByHeuristic(Args[0], config)

	if err != nil {
		fmt.Printf("Could not load bug: %s\n", err.Error())
		return
	}

	currentDir := b.GetDirectory()
	newDir := bugs.GetIssuesDir(config) + bugs.TitleToDir(strings.Join(Args[1:], " "))
	fmt.Printf("Moving %s to %s\n", currentDir, newDir)
	err = os.Rename(string(currentDir), string(newDir))
	if err != nil {
		fmt.Printf("Error moving directory\n")
	}
}
