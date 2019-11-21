package scm

import (
	"flag"
	"fmt"
	bugs "github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

type GitCommit struct {
	commit string
	log    string
}

func (c GitCommit) CommitID() string {
	return c.commit
}
func (c GitCommit) LogMsg() string {
	return c.log
}
func (c GitCommit) Diff() (string, error) {
	return runCmd("git", "show", "--pretty=format:%b", c.CommitID())
}

func (c GitCommit) CommitMessage() (string, error) {
	return runCmd("git", "show", "--pretty=format:%B", "--quiet", c.CommitID())
}

type GitTester struct {
	handler       SCMHandler
	workdir       string
	pwd           string
	issuesdirname string
}

func (g GitTester) Loggers() ([]Commit, error) {
	logs, err := runCmd("git", "log", "--oneline", "--reverse", "-z")
	if err != nil {
		wd, _ := os.Getwd()
		fmt.Fprintf(os.Stderr, "Error retrieving git logs: %s in directory %s\n", logs, wd)
		return nil, err
	}
	logMsgs := strings.Split(logs, "\000")
	// the last line is empty, so don't allocate 1 for
	// it
	commits := make([]Commit, len(logMsgs)-1)
	for idx, commitText := range logMsgs {
		if commitText == "" {
			continue
		}
		spaceIdx := strings.Index(commitText, " ")
		if spaceIdx >= 0 {
			commits[idx] = GitCommit{commitText[0:spaceIdx], commitText[spaceIdx+1:]}
		}
	}
	return commits, nil
}

func (g GitTester) AssertStagingIndex(t *testing.T, f []FileStatus) {
	for _, file := range f {
		out, err := runCmd("git", "status", "--porcelain", file.Filename)
		if err != nil {
			t.Error("Could not run git status")
		}
		expected := file.IndexStatus + file.WorkingStatus + " " + file.Filename + "\n"
		if out != expected {
			t.Error("Incorrect file status")
			t.Error("Got" + out + " not " + expected)
		}
	}
	// extra FileStatus entries not asserted yet
}

func (g GitTester) StageFile(file string) error {
	_, err := runCmd("git", "add", file)
	return err
}
func (g *GitTester) Setup() error {
	pwd, _ := os.Getwd()
	g.pwd = pwd
	if gdir, err := ioutil.TempDir("", "gitmanager"); err == nil {
		g.workdir = gdir
		os.Chdir(g.workdir)
		os.Unsetenv("FIT")
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		//gdir, _ = os.Getwd() // gdir not used later
	} else {
		panic("Failed creating temporary directory")
	}
	// Make sure we get the right directory from the top level
	os.Mkdir(g.issuesdirname, 0755)

	out, err := runCmd("git", "init")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing git: %s", out)
		return err
	}

	return nil
}

var git bool

func init() {
	flag.BoolVar(&git, "git", true, "git presence")
	//flag.Parse()
	_, err := runCmd("git", "help")
	if err != nil {
		git = false
	}
}

func (g GitTester) TearDown() {
	os.Chdir(g.pwd)
	os.RemoveAll(g.workdir)
}
func (g GitTester) WorkDir() string {
	return g.workdir
}

func (g GitTester) AssertCleanTree(t *testing.T) {
	out, err := runCmd("git", "status", "--porcelain")
	if err != nil {
		t.Error("Error running git status")
	}
	if out != "" {
		t.Error("Unexpected Output from git status (expected nothing):\n" + out)
	}
}

func (g GitTester) Manager() SCMHandler {
	return g.handler
}

func TestGitBugRenameCommits(t *testing.T) {
	t.Skip("windows failure - see scm/GitManager_test.go+143")
	// TODO: finish making tests on Windows pass then redo this test
	// This test fakes output of the main bug command then tries to rename
	// what looks like with os.rename and not hg rename. Maybe scrap the test
	// and start over. The simulations of simulations feel unnecessary.
	if git == false {
		t.Skip("WARN git executable not found")
	}
	g := GitTester{}
	g.handler = GitManager{}

	expectedDiffs := []string{
		`
diff --git a/fit/Test-bug/Description b/fit/Test-bug/Description
new file mode 100644
index 0000000..e69de29
`, `
diff --git a/fit/Test-bug/Description b/fit/Renamed-bug/Description
similarity index 100%
rename from fit/Test-bug/Description
rename to fit/Renamed-bug/Description
`}

	runtestRenameCommitsHelper(&g, t, expectedDiffs)
}

func TestGitIssueStatus(t *testing.T) {
	i := issueStatus{true, true, true}
	if i.a != true || i.d != true || i.m != true {
		t.Error("issueStatus are not all true, a " + strconv.FormatBool(i.a) +
			", d " + strconv.FormatBool(i.d) +
			", m " + strconv.FormatBool(i.m))
	}
}

func TestGitFilesOutsideOfBugNotCommited(t *testing.T) {
	if git == false {
		t.Skip("WARN git executable not found")
	}
	t.Skip("windows failure - see scm/GitManager_test.go+182")
	// TODO: finish making tests on Windows pass then redo this test
	// the error codes need handling
	g := GitTester{}
	g.handler = GitManager{}
	runtestCommitDirtyTree(&g, t)
}

func TestGitManagerTyper(t *testing.T) {
	manager := GitManager{}

	if getType := manager.SCMTyper(); getType != "git" {
		t.Error("Incorrect SCM Type for GitManager. Got " + getType)
	}
}

func TestGitManagerPurge(t *testing.T) {
	if git == false {
		t.Skip("WARN git executable not found")
	}
	t.Skip("windows failure - see scm/GitManager_test.go+202")
	// TODO: finish making tests on Windows pass then redo this test
	// the error codes need handling
	g := GitTester{}
	g.handler = GitManager{}
	runtestPurgeFiles(&g, t)
}

func TestGitManagerAutoclosingGitHub(t *testing.T) {
	var config bugs.Config
	config.DescriptionFileName = "Description"
	config.IssuesDirName = "fit"
	// This test is specific to gitmanager, since GitHub
	// only supports git
	if git == false {
		t.Skip("WARN git executable not found")
	}
	t.Skip("windows failure - see scm/GitManager_test.go+218")
	// TODO: finish making tests on Windows pass then redo this test
	// the error codes need handling
	tester := GitTester{}
	tester.handler = GitManager{Autoclose: true}

	err := tester.Setup()
	if err != nil {
		panic("Something went wrong trying to initialize git : " + err.Error())
	}
	defer tester.TearDown()
	m := tester.Manager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	//runCmd("bug", "create", "-n", "Test", "bug")
	os.MkdirAll(config.IssuesDirName+sops+"Test-bug", 0755)
	ioutil.WriteFile(config.IssuesDirName+sops+"Test-bug"+sops+"Description", []byte("desc1"), 0644)
	if err = ioutil.WriteFile(config.IssuesDirName+sops+"Test-bug"+sops+"Identifier", []byte("\n\nGitHub:#TestBug"), 0644); err != nil {
		t.Error("Could not write Test-bug" + sops + "Identifier file")
		return
	}

	//runCmd("bug", "create", "-n", "Test", "Another", "bug")
	os.MkdirAll(config.IssuesDirName+sops+"Test-Another-bug", 0755)
	ioutil.WriteFile(config.IssuesDirName+sops+"Test-Another-bug"+sops+"Description", []byte("desc2"), 0644)
	if err = ioutil.WriteFile(config.IssuesDirName+sops+"Test-Another-bug"+sops+"Identifier", []byte("\n\nGITHuB:  #Whitespace   "), 0644); err != nil {
		t.Error("Could not write Test-Another-bug" + sops + "Identifier file")
		return
	}

	// Add and commit the file, so that we can close it..
	m.Commit(bugs.Directory(tester.WorkDir()+sops+config.IssuesDirName), "Adding commit", config)
	// Delete the bugs
	os.RemoveAll(tester.WorkDir() + sops + config.IssuesDirName + sops + "Test-bug")
	os.RemoveAll(tester.WorkDir() + sops + config.IssuesDirName + sops + "Test-Another-bug")
	m.Commit(bugs.Directory(tester.WorkDir()+sops+config.IssuesDirName), "Removal commit", config)

	commits, err := tester.Loggers()
	if len(commits) != 2 || err != nil {
		t.Error("Error getting git logs while attempting to test GitHub autoclosing")
		return
	}
	if msg, err := commits[1].(GitCommit).CommitMessage(); err != nil {
		t.Error("Error getting git logs while attempting to test GitHub autoclosing")
	} else {
		closing := func(issue string) bool {
			return strings.Contains(msg, "Closes #"+issue) ||
				strings.Contains(msg, ", closes #"+issue)
		}
		if !closing("Whitespace") || !closing("TestBug") {
			fmt.Printf("%s\n", msg)
			t.Error("GitManager did not autoclose Github issues")
		}
	}
}
