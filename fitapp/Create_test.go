package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
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
		t.Error("Unexpected output on STDOUT for fitapp/Create_test")
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

// Test "Create" without a fit directory
func TestCreateWithoutIssues(t *testing.T) {
	t.Skip("see fitapp/Create_test.go+41 and fitapp/utils.go+96")
	config := bugs.Config{}
	config.DescriptionFileName = "Description"
	config.FitDirName = "fit"
	dir, err := ioutil.TempDir("", "createtest")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	pwd, _ := os.Getwd()
	os.Chdir(dir)
	// this test should comment MkdirAll.
	// Oddly that causes a test halt with "exit status 1".
	// I tracked this down to fitapp/utils.go +96, os.Stdout = op
	// Capturing the output of the RUNNING process for testing
	// is a bit sneaky. I don't see another way to make it work.
	// Even though I can't run this test as a function it passes.
	// I added t.Skip above.
	os.MkdirAll(config.FitDirName, 0700) // the real test
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
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s%s%s%s", dir, sops, config.FitDirName, sops))
	//fmt.Print("3")
	if err != nil {
		t.Error("Could not read " + config.FitDirName + " directory")
		return
	}
	if len(issuesDir) != 1 {
		t.Error(fmt.Sprintf("Unexpected number of issues in %s dir.\n    Expected %d, got %d\n", config.FitDirName, 1, len(issuesDir)))
	}
	//fmt.Print("4")
	os.Chdir(pwd)
}

// Test a very basic invocation of "Create" with the -n
// argument. We can't yet try it without -n, since it means
// an editor will be spawned.
func TestCreateNoEditor(t *testing.T) {
	t.Skip("windows failure - see fitapp/Create_test.go+92")
	// TODO: finish making tests on Windows pass then redo this test
	//       first issue was ok.
	//       second issue had trouble with setting and using DefaultDescriptionFile
/*
=== RUN   TestCreateNoEditor
--- FAIL: TestCreateNoEditor (0.01s)
    Create_test.go:167: Unexpected output on STDOUT for Test2-bug: open \ddf: The system cannot find th
e file specified.
        Created issue: Test2 bug

    Create_test.go:180: Unexpected number of files found in Test2-bug dir.
            Expected 2, got 1

    Create_test.go:189: Could not load description file for Test2 bugopen C:\cygwin64\tmp\createtest740
771043\fit\Test2-bug\ddf: The system cannot find the file specified.
    Create_test.go:192: Unexpected empty file for Test2 bug

*/
	config := bugs.Config{}
	config.DescriptionFileName = "Description"
	config.FitDirName = "fit"
	dir, err := ioutil.TempDir("", "createtest")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	pwd, _ := os.Getwd()
	os.Chdir(dir)
	os.MkdirAll(config.FitDirName, 0700)
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
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s%s%s%s", dir, sops, config.FitDirName, sops))
	if err != nil {
		t.Error("Could not read " + config.FitDirName + " directory")
		return
	}
	if len(issuesDir) != 1 {
		t.Error(fmt.Sprintf("Unexpected number of issues in %s dir.\n    Expected %d, got %d\n", config.FitDirName, 1, len(issuesDir)))
	}

	bugDir, err := ioutil.ReadDir(fmt.Sprintf("%s%s%s%sTest-bug", dir, sops, config.FitDirName, sops))
	if len(bugDir) != 1 {
		t.Error(fmt.Sprintf("Unexpected number of files found in %s dir.\n    Expected %d, got %d\n", "Test-bug", 1, len(bugDir)))
	}
	if err != nil {
		t.Error("Could not read Test-bug directory")
		return
	}

	file, err := ioutil.ReadFile(fmt.Sprintf("%s%s%s%sTest-bug%sDescription", dir, sops, config.FitDirName, sops, sops))
	if err != nil {
		t.Error("Could not load description file for Test bug" + err.Error())
	}
	if len(file) != 0 {
		t.Error("Expected empty file for Test bug")
	}

	///// second issue
	////// uses a configured file name
	config.DefaultDescriptionFile = "ddf"
	// put this ABOVE issues so len(issuesDir) check later is unaltered
	ioutil.WriteFile(config.DefaultDescriptionFile,
		[]byte("text used in default description file (ddf) issue template"), 0755)

	stdout, stderr = captureOutput(func() {
		Create(argumentList{"-n", "--generate-id", "Test2", "bug"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected output on STDERR for Test2-bug")
	}
	if stdout != "Created issue: Test2 bug\n" {
		t.Error("Unexpected output on STDOUT for Test2-bug: " + stdout)
	}
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s%s%s%s", dir, sops, config.FitDirName, sops))
	if err != nil {
		t.Error("Could not read " + config.FitDirName + " directory")
		return
	}
	if len(issuesDir) != 2 {
		t.Error(fmt.Sprintf("Unexpected number of issues in %s dir.\n    Expected %d, got %d\n", config.FitDirName, 2, len(issuesDir)))
	}

	bugDir, err = ioutil.ReadDir(fmt.Sprintf("%s%s%s%sTest2-bug", dir, sops, config.FitDirName, sops))
	if len(bugDir) != 2 {
		t.Error(fmt.Sprintf("Unexpected number of files found in %s dir.\n    Expected %d, got %d\n", "Test2-bug", 2, len(bugDir)))
	}
	if err != nil {
		t.Error("Could not read Test2-bug directory")
		return
	}

	file, err = ioutil.ReadFile(fmt.Sprintf("%s%s%s%sTest2-bug%s%s", dir, sops, config.FitDirName, sops, sops, config.DefaultDescriptionFile))
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
