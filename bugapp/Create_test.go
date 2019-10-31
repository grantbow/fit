package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	//	"io"
	"io/ioutil"
	"os"
	"testing"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

func runCreateOutput(args argumentList, expected string, t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		Create(args, config)
	}, t)
	if stdout != expected {
		t.Error("Unexpected output on STDOUT for bugapp/Create_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
	if stderr[:7] != "Usage: " {
		t.Error("Expected usage information with no arguments")
	}
}

// Captures stdout and stderr to ensure that
// a usage line gets printed to Stderr when
// no parameters are specified
func TestCreateHelpOutput(t *testing.T) {
	runCreateOutput(argumentList{}, "", t)
}

// Test "Create" without an issues directory
func TestCreateWithoutIssues(t *testing.T) {
	t.Skip("see bugapp/Create_test.go+41 and bugapp/utils.go+96")
	config := bugs.Config{}
	config.DescriptionFileName = "Description"
	config.IssuesDirName = "fit"
	dir, err := ioutil.TempDir("", "createtest")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	pwd, _ := os.Getwd()
	os.Chdir(dir)
	// this test should comment MkdirAll.
	// Oddly that causes a test halt with "exit status 1".
	// I tracked this down to bugapp/utils.go +96, os.Stdout = op
	// Capturing the output of the RUNNING process for testing
	// is a bit sneaky. I don't see another way to make it work.
	// Even though I can't run this test as a function it passes.
	// I added t.Skip above.
	os.MkdirAll(config.IssuesDirName, 0700) // the real test
	defer os.RemoveAll(dir)
	err = os.Setenv("FIT", dir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	//fmt.Print("1")
	//fmt.Print(err)
	stdout, stderr := captureOutput(func() {
		Create(argumentList{"-n", "Test", "bug"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected output on STDERR for Test-bug: " + stderr)
	}
	if stdout != "Created issue: Test bug\n" {
		t.Error("Unexpected output on STDOUT for Test-bug: " + stdout)
	}
	//fmt.Print("2")
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s%s%s%s", dir, sops, config.IssuesDirName, sops))
	//fmt.Print("3")
	if err != nil {
		t.Error("Could not read " + config.IssuesDirName + " directory")
		return
	}
	if len(issuesDir) != 1 {
		t.Error("Unexpected number of issues in " + config.IssuesDirName + " dir\n")
	}
	//fmt.Print("4")
	os.Chdir(pwd)
}

// Test a very basic invocation of "Create" with the -n
// argument. We can't yet try it without -n, since it means
// an editor will be spawned..
func TestCreateNoEditor(t *testing.T) {
	config := bugs.Config{}
	config.DescriptionFileName = "Description"
	config.IssuesDirName = "fit"
	dir, err := ioutil.TempDir("", "createtest")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	pwd, _ := os.Getwd()
	os.Chdir(dir)
	os.MkdirAll(config.IssuesDirName, 0700)
	defer os.RemoveAll(dir)
	// On MacOS, /tmp is a symlink, which causes GetDirectory() to return
	// a different path than expected in these tests, so make the issues
	// directory explicit with an environment variable
	err = os.Setenv("FIT", dir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	///// without an issue
	runCreateOutput(argumentList{"-n"}, "", t)

	///// first issue
	stdout, stderr := captureOutput(func() {
		Create(argumentList{"-n", "Test", "bug"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected output on STDERR for Test-bug")
	}
	if stdout != "Created issue: Test bug\n" {
		t.Error("Unexpected output on STDOUT for Test-bug")
	}
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s%s%s%s", dir, sops, config.IssuesDirName, sops))
	if err != nil {
		t.Error("Could not read " + config.IssuesDirName + " directory")
		return
	}
	if len(issuesDir) != 1 {
		t.Error("Unexpected number of issues in " + config.IssuesDirName + " dir\n")
	}

	bugDir, err := ioutil.ReadDir(fmt.Sprintf("%s%s%s%sTest-bug", dir, sops, config.IssuesDirName, sops))
	if len(bugDir) != 1 {
		t.Error("Unexpected number of files found in Test-bug dir\n")
	}
	if err != nil {
		t.Error("Could not read Test-bug directory")
		return
	}

	file, err := ioutil.ReadFile(fmt.Sprintf("%s%s%s%sTest-bug%sDescription", dir, sops, config.IssuesDirName, sops, sops))
	if err != nil {
		t.Error("Could not load description file for Test bug" + err.Error())
	}
	if len(file) != 0 {
		t.Error("Expected empty file for Test bug")
	}

	///// second issue
	config.DefaultDescriptionFile = dir + sops + "ddf" // put ABOVE issues so len(issuesDir) check later is unaltered
	ioutil.WriteFile(config.DefaultDescriptionFile,
		[]byte("text used in default description file (ddf) issue template"), 0755)

	stdout, stderr = captureOutput(func() {
		Create(argumentList{"-n", "--generate-id", "Test2", "bug"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected output on STDERR for Test2-bug")
	}
	if stdout != "Created issue: Test2 bug\n" {
		t.Error("Unexpected output on STDOUT for Test2-bug")
	}
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s%s%s%s", dir, sops, config.IssuesDirName, sops))
	if err != nil {
		t.Error("Could not read " + config.IssuesDirName + " directory")
		return
	}
	if len(issuesDir) != 2 {
		t.Error("Unexpected number of issues in " + config.IssuesDirName + " dir\n")
	}

	bugDir, err = ioutil.ReadDir(fmt.Sprintf("%s%s%s%sTest2-bug", dir, sops, config.IssuesDirName, sops))
	if len(bugDir) != 2 {
		t.Error("Unexpected number of files found in Test2-bug dir\n")
	}
	if err != nil {
		t.Error("Could not read Test2-bug directory")
		return
	}

	file, err = ioutil.ReadFile(fmt.Sprintf("%s%s%s%sTest2-bug%sDescription", dir, sops, config.IssuesDirName, sops, sops))
	if err != nil {
		t.Error("Could not load description file for Test2 bug" + err.Error())
	}
	if len(file) == 0 {
		t.Error("Unexpected empty file for Test2 bug")
	}
	os.Chdir(pwd)
}

/* currently hangs spawning editor

// this test will not spawn editor
func TestCreateNoIssuesDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "createtest")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	// This is what we are testing for
	//os.MkdirAll("issues", 0700)
	//defer os.RemoveAll(dir)
	err = os.Setenv("FIT", string(dir))
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	stdout, stderr := captureOutput(func() {
		Create(argumentList{"Test", "bug"}) // fire editor without -n
	}, t)
	_ = stderr
	_ = stdout
	// stderr is expected from log.Fatal(err)
	//if stdout != "" {
	//	t.Error("Unexpected output on STDOUT for Test-bug")
	//}
}
*/
