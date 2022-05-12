package scm

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

type HgCommit struct {
	commit string
	log    string
}

func (c HgCommit) CommitID() string {
	return c.commit
}
func (c HgCommit) LogMsg() string {
	return c.log
}
func (c HgCommit) Diff() (string, error) {
	return runCmd("hg", "log", "-p", "-g", "-r", c.commit, "--template={changelog}")
}

type HgTester struct {
	handler SCMHandler
	workdir string
	pwd     string
}

func (h HgTester) Loggers() ([]Commit, error) {
	logs, err := runCmd("hg", "log", "-r", ":", "--template", "{node} {desc}\\n")
	if err != nil {
		return nil, err
	}
	logMsgs := strings.Split(logs, "\n")
	// the last line is empty, so don't allocate 1 for
	// it
	commits := make([]Commit, len(logMsgs)-1)
	for idx, commitText := range logMsgs {
		if commitText == "" {
			continue
		}
		spaceIdx := strings.Index(commitText, " ")
		if spaceIdx >= 0 {
			commits[idx] = HgCommit{commitText[0:spaceIdx], commitText[spaceIdx+1:]}
		}
	}
	return commits, nil
}

func (h HgTester) AssertStagingIndex(t *testing.T, f []FileStatus) {
	for _, file := range f {
		out, err := runCmd("hg", "status", file.Filename)
		if err != nil {
			t.Error("Could not get status of " + file.Filename)
		}

		// hg status doesn't include the working directory status
		expected := file.IndexStatus + " " + file.Filename + "\n"
		if out != expected {
			t.Error("Unexpected status. Got " + out + " not " + expected)
		}

	}
}

func (h HgTester) StageFile(file string) error {
	_, err := runCmd("hg", "add", file)
	return err
}
func (h *HgTester) Setup() error {
	pwd, _ := os.Getwd()
	h.pwd = pwd
	if dir, err := ioutil.TempDir("", "hgmanager"); err == nil {
		h.workdir = dir
		os.Chdir(h.workdir)
	} else {
		return err
	}

	_, err := runCmd("hg", "init")
	if err != nil {
		//h.Fail()
		log.Fatal(err)
	}
	h.handler = HgManager{}
	return nil
}

var hg bool

func init() {
	flag.BoolVar(&hg, "hg", true, "Mercurial (hg) presence")
	//flag.Parse()
	_, err := runCmd("hg")
	if err != nil {
		hg = false
	}
}

func (h HgTester) TearDown() {
	os.Chdir(h.pwd)
	os.RemoveAll(h.workdir)
}

func (h HgTester) WorkDir() string {
	return h.workdir
}
func (h HgTester) AssertCleanTree(t *testing.T) {
	out, err := runCmd("hg", "status")
	if err != nil {
		t.Error("Error running hg status")
	}
	if out != "" {
		fmt.Printf("\"%s\"\n", out)
		t.Error("Unexpected Output from hg status")
	}
}

func (h HgTester) Manager() SCMHandler {
	return h.handler
}

func TestHgBugRenameCommits(t *testing.T) {
	if hg == false {
		t.Skip("WARN hg executable not found")
	}
	t.Skip("windows failure - see scm/HgManager_test.go+132")
	// TODO: finish making tests on Windows pass then redo this test
	// This test fakes output of the main bug command then tries to rename
	// what looks like with os.rename and not hg rename. Maybe scrap the test
	// and start over. The simulations of simulations feel unnecessary.
	h := HgTester{}

	//t.Skip("TODO: fix HgBugRenameCommits changed output in some (hg version?) conditions")
	expectedDiffs := []string{
		`diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
new file mode 100644

`, `diff --git a/issues/Renamed-bug/Description b/issues/Renamed-bug/Description
new file mode 100644
diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
deleted file mode 100644

`}
	runtestRenameCommitsHelper(&h, t, expectedDiffs)
}
func TestHgFilesOutsideOfBugNotCommited(t *testing.T) {
	if hg == false {
		t.Skip("WARN hg executable not found")
	}
	t.Skip("windows failure - see scm/HgManager_test.go+156")
	// TODO: finish making tests on Windows pass then redo this test
	// the error codes need handling
	h := HgTester{}
	h.handler = HgManager{}
	runtestCommitDirtyTree(&h, t)
}

func TestHgTyper(t *testing.T) {
	manager := HgManager{}

	if getType := manager.SCMTyper(); getType != "hg" {
		t.Error("Incorrect SCM Type for HgManager. Got " + getType)
	}
}

func TestHgPurge(t *testing.T) {
	// This should eventually be replaced by something more
	// like:
	//      h := HgTester{}
	//      h.handler = HgManager{}
	//      runtestPurgeFiles(&h, t)
	// but since the current behaviour is to return a not
	// supported error, that would evidently fail..
	m := HgManager{}
	err := m.Purge(dops + "tmp" + dops + "imaginaryHgRepo")

	switch err.(type) {
	case UnsupportedType:
		// This is valid, do nothing.
	default:
		t.Error("Unexpected return value for Hg purge function.")
	}

}
