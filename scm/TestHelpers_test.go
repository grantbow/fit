package scm

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

type Commit interface {
	CommitID() string
	LogMsg() string
	Diff() (string, error)
}

type ManagerTester interface {
	Loggers() ([]Commit, error)
	Setup() error
	WorkDir() string
	TearDown()
	StageFile(string) error
	AssertCleanTree(t *testing.T)
	AssertStagingIndex(*testing.T, []FileStatus)
	Manager() SCMHandler
}

type argumentList []string

func runCmd(cmd string, options ...string) (string, error) {
	runcmd := exec.Command(cmd, options...)
	out, err := runcmd.CombinedOutput()

	return string(out), err
}

func assertLogs(tester ManagerTester, t *testing.T, titles []map[string]bool, diffs []string) {
	logs, err := tester.Loggers()
	if err != nil {
		t.Error("scm Loggers error : " + err.Error())
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
	var config bugs.Config
	config.DescriptionFileName = "Description"
	config.FitDirName = "fit"
	err := tester.Setup()
	defer tester.TearDown()
	if err != nil {
		t.Error("Could not ihnitialize repo")
		return
	}

	m := tester.Manager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	os.MkdirAll(config.FitDirName+sops+"Test-bug", 0755)
	ioutil.WriteFile(config.FitDirName+sops+"Test-bug"+sops+"Description", []byte(""), 0644)
	m.Commit(bugs.Directory(tester.WorkDir()), "Initial commit", config)
	//runCmd("bug", "retitle", "1", "Renamed", "bug")
	args := argumentList{"1", "Renamed bug"}
	expected := "Moving .*"
	runretitle("scm"+sops+"TestHelpers_test runtestRenameCommitsHelper", args, config, expected, t)
	m.Commit(bugs.Directory(tester.WorkDir()), "This is a test rename", config)

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
	var config bugs.Config
	config.DescriptionFileName = "Description"
	config.FitDirName = "fit"
	err := tester.Setup()
	if err != nil {
		panic("Something went wrong trying to initialize the scm : " + err.Error())
	}
	defer tester.TearDown()
	m := tester.Manager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	//
	//
	os.MkdirAll(config.FitDirName+sops+"Test-bug", 0755)
	ioutil.WriteFile(config.FitDirName+sops+"Test-bug"+sops+"Description", []byte(""), 0644)
	if err = ioutil.WriteFile("donotcommit.txt", []byte(""), 0644); err != nil {
		t.Error("Could not write file")
		return
	}
	tester.AssertStagingIndex(t, []FileStatus{
		{"donotcommit.txt", "?", "?"},
	})

	//fmt.Print("pre  1 runtestCommitDirtyTree\n")
	m.Commit(bugs.Directory(tester.WorkDir()+sops+config.FitDirName), "Initial commit", config)
	//fmt.Print("post 1 runtestCommitDirtyTree\n")
	tester.AssertStagingIndex(t, []FileStatus{
		{"donotcommit.txt", "?", "?"},
	})
	tester.StageFile("donotcommit.txt")
	tester.AssertStagingIndex(t, []FileStatus{
		{"donotcommit.txt", "A", " "},
	})
	//fmt.Print("pre  2 runtestCommitDirtyTree\n")
	m.Commit(bugs.Directory(tester.WorkDir()+sops+config.FitDirName), "Initial commit", config)
	//errCommit := m.Commit(bugs.Directory(tester.WorkDir()+sops+config.FitDirName), "Initial commit", config)
	//fmt.Printf("post 2 runtestCommitDirtyTree error %v\n", errCommit) // nil here
	//    running test shows output here. actually HgManager.go Commit() returns *expected* error
	//        not fully handled by Hg though
	//        scm/GitManager.go#func SCMIssuesCacher() was implemented later and
	//        used with scm/GitManager.go#func SCMIssuesUpdaters in cmd/bug/main.go.
	//    stdout not captured this time.
	//fmt.Print("post 2 runtestCommitDirtyTree\n")
	tester.AssertStagingIndex(t, []FileStatus{
		{"donotcommit.txt", "A", " "},
	})
	//
	os.MkdirAll(config.FitDirName+sops+"Fresh-bug", 0755)
	ioutil.WriteFile(config.FitDirName+sops+"Fresh-bug"+sops+"Description", []byte(""), 0644)
	tester.AssertStagingIndex(t, []FileStatus{
		{"donotcommit.txt", "A", " "},
		{config.FitDirName + sops + "Fresh-bug" + sops + "Description", "?", "?"},
	})
	//errCommit := m.Commit(bugs.Directory(tester.WorkDir()+sops+config.FitDirName), "Initial commit", config)
	//fmt.Printf("post 2 runtestCommitDirtyTree error %v\n", errCommit) // shouldn't be nil
	scmoptions := make(map[string]bool)
	handler, _, _ := DetectSCM(scmoptions, config)
	if b, err := handler.SCMIssuesUpdaters(config); err != nil {
		for _, bline := range strings.Split(string(b), "\n") {
			fmt.Printf("1 Warn Updaters: %v\n", string(bline))
		}
		//if _, ErrCa := handler.SCMIssuesCacher(config); ErrCa != nil {
		//	fmt.Printf("1 Warn cacher: %s\n", ErrCa)
		//}
		if c, ErrCa := handler.SCMIssuesCacher(config); ErrCa != nil {
			for _, bline := range strings.Split(string(c), "\n") {
				fmt.Printf("2 Warn cacher: %v\n", string(bline))
			}
		}
	}
	tester.StageFile(config.FitDirName + sops + "Fresh-bug" + sops + "Description")
	handler, _, _ = DetectSCM(scmoptions, config)
	if b, err := handler.SCMIssuesUpdaters(config); err != nil {
		for _, bline := range strings.Split(string(b), "\n") {
			fmt.Printf("2 Warn Updaters: %v\n", string(bline))
		}
		if c, ErrCa := handler.SCMIssuesCacher(config); ErrCa != nil {
			for _, bline := range strings.Split(string(c), "\n") {
				fmt.Printf("2 Warn cacher: %v\n", string(bline))
			}
		}
	}
	/*
		if b, err := handler.SCMIssuesUpdaters(config); err != nil {
			fmt.Printf("Files in issues/ need committing, see $ git status --porcelain -u -- :/issues\nand for files already in index see $ git diff --name-status --cached HEAD -- :/issues\n")
			if _, ErrCach := handler.SCMIssuesCacher(config); ErrCach != nil {
				for _, bline := range strings.Split(string(b), "\n") {
					//if bline in c {
					//} else {
					fmt.Printf("%v\n", string(bline))
					//}
				}
			} else {
				fmt.Printf("%v\n", string(b))
			}
		} else {
			fmt.Printf("No files in issues/ need committing, see $ git status --porcelain -u issues \":top\"\n")
		}

		if ErrH != nil {
			fmt.Printf("Warn: to commit your issues first use {git|hg} init\n")
			//fmt.Printf("Warn: %s\n", ErrH) // No SCM found
			//a, b := handler.SCMIssuesUpdaters(config)
			//fmt.Printf("%+v %+v\n", a, b)
			if handler != nil {
				if _, ErrU := handler.SCMIssuesUpdaters(config); ErrU != nil {
					if _, ErrCa := handler.SCMIssuesCacher(config); ErrCa != nil {
						fmt.Printf("Warn: %s\n", ErrCa)
					} else {
						fmt.Printf("Warn: %s\n", ErrU)
					}
				}
			}
		}
	*/
}

func runtestPurgeFiles(tester ManagerTester, t *testing.T) {
	var config bugs.Config
	config.DescriptionFileName = "Description"
	config.FitDirName = "fit"
	err := tester.Setup()
	if err != nil {
		panic("Something went wrong trying to initialize: " + err.Error())
	}
	defer tester.TearDown()
	m := tester.Manager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	// Commit a bug which should stay around after the purge
	//runCmd("bug", "create", "-n", "Test", "bug")
	os.MkdirAll(config.FitDirName+sops+"Test-bug", 0755)
	ioutil.WriteFile(config.FitDirName+sops+"Test-bug"+sops+"Description", []byte(""), 0644)
	m.Commit(bugs.Directory(tester.WorkDir()+sops+config.FitDirName), "Initial commit", config)

	// Create another bug to elimate with purge
	//runCmd("bug", "create", "-n", "Test", "Purge", "bug")
	os.MkdirAll(config.FitDirName+sops+"Test-Purge-Bug", 0755)
	err = m.Purge(bugs.Directory(tester.WorkDir() + "" + sops + config.FitDirName))
	if err != nil {
		t.Error("Error purging directory: " + err.Error())
	}
	FitDir, err := ioutil.ReadDir(config.FitDirName) //fmt.Sprintf("debug: %s/issues/", tester.WorkDir()))
	if err != nil {
		t.Error("Error reading issues directory: " + err.Error())
	}
	if len(FitDir) != 1 {
		t.Errorf("Unexpected number of directories (%v, expected 1) in issues/ after purge.", len(FitDir))
	}
	if len(FitDir) > 0 && FitDir[0].Name() != "Test-bug" {
		t.Error("Expected Test-bug to remain.")
	}
	err = m.Commit(bugs.Directory(tester.WorkDir()+sops+config.FitDirName), "", config) // no changes
	if err != nil {                                                                     // should NOT be a nil
		t.Error("Expected no commit error : " + err.Error())
	}
}

func runretitle(label string, args argumentList, config bugs.Config, expected string, t *testing.T) {
	stdout, _ := captureOutput(func() {
		scmRetitle(args, config)
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

func scmRetitle(Args argumentList, config bugs.Config) {
	if len(Args) < 2 {
		fmt.Printf("Usage: %s retitle <IssueID> New Title\n", os.Args[0])
		return
	}

	b, err := bugs.LoadIssueByHeuristic(Args[0], config)

	if err != nil {
		fmt.Printf("Could not load issue: %s\n", err.Error())
		return
	}

	currentDir := b.Direr()
	newDir := bugs.FitDirer(config) + dops + bugs.TitleToDir(strings.Join(Args[1:], " "))
	fmt.Printf("Moving %s to %s\n", currentDir, newDir)
	err = os.Rename(string(currentDir), string(newDir))
	// uses os.rename
	// not git-mv and/or hg rename
	if err != nil {
		fmt.Printf("Error moving directory\n")
	}
}
