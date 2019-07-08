package scm

import (
	"flag"
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

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
	handler SCMHandler
	workdir string
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
	if gdir, err := ioutil.TempDir("", "gitbug"); err == nil {
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
	os.Mkdir("issues", 0755)

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
	flag.Parse()
	_, err := runCmd("git", "help")
	if err != nil {
		git = false
	}
}

func (g GitTester) TearDown() {
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
	if git == false {
		t.Skip("git executable not found")
	}
	g := GitTester{}
	g.handler = GitManager{}

	expectedDiffs := []string{
		`
diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
new file mode 100644
index 0000000..e69de29
`, `
diff --git a/issues/Test-bug/Description b/issues/Renamed-bug/Description
similarity index 100%
rename from issues/Test-bug/Description
rename to issues/Renamed-bug/Description
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
		t.Skip("git executable not found")
	}
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
		t.Skip("git executable not found")
	}
	g := GitTester{}
	g.handler = GitManager{}
	runtestPurgeFiles(&g, t)
}

func TestGitManagerAutoclosingGitHub(t *testing.T) {
	var config bugs.Config
	config.DescriptionFileName = "Description"
	// This test is specific to gitmanager, since GitHub
	// only supports git..
	if git == false {
		t.Skip("git executable not found")
	}
	tester := GitTester{}
	tester.handler = GitManager{Autoclose: true}

	err := tester.Setup()
	if err != nil {
		panic("Something went wrong trying to initialize git:" + err.Error())
	}
	defer tester.TearDown()
	m := tester.Manager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	//runCmd("bug", "create", "-n", "Test", "bug")
	os.MkdirAll("issues/Test-bug", 0755)
	ioutil.WriteFile("issues/Test-bug/Description", []byte(""), 0644)
	if err = ioutil.WriteFile("issues/Test-bug/Identifier", []byte("\n\nGitHub:#TestBug"), 0644); err != nil {
		t.Error("Could not write Test-bug/Identifier file")
		return
	}

	//runCmd("bug", "create", "-n", "Test", "Another", "bug")
	os.MkdirAll("issues/Test-Another-bug", 0755)
	ioutil.WriteFile("issues/Test-Another-bug/Description", []byte(""), 0644)
	if err = ioutil.WriteFile("issues/Test-Another-bug/Identifier", []byte("\n\nGITHuB:  #Whitespace   "), 0644); err != nil {
		t.Error("Could not write Test-Another-bug/Identifier file")
		return
	}

	// Commit the file, so that we can close it..
	m.Commit(bugs.Directory(tester.WorkDir()+"/issues"), "Adding commit", config)
	// Delete the bug
	os.RemoveAll(tester.WorkDir() + "/issues/Test-bug")
	os.RemoveAll(tester.WorkDir() + "/issues/Test-Another-bug")
	m.Commit(bugs.Directory(tester.WorkDir()+"/issues"), "Removal commit", config)

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
