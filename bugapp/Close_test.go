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

// ensure that
// a usage line gets printed to Stderr when
// no parameters are specified
func TestCloseHelpOutput(t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		Close(argumentList{}, config)
	}, t)

	if stdout != "" {
		t.Error("Unexpected output on stdout.")
	}
	if stderr[:7] != "Usage: " {
		t.Error("Expected usage information with no arguments")
	}

}

// Test closing a bug given it's directory index
func TestCloseByIndex(t *testing.T) {
	config := bugs.Config{}
	dir, err := ioutil.TempDir("", "closetest")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	os.MkdirAll("issues"+sops+"Test", 0700)

	// On MacOS, /tmp is a symlink, which causes GetDirectory() to return
	// a different path than expected in these tests, so make the issues
	// directory explicit with an environment variable
	err = os.Setenv("FIT", dir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	ioutil.WriteFile(dir+sops+"issues"+sops+"Test"+sops+"Identifier", []byte("TestBug\n"), 0600)

	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s%sissues%s", dir, sops, sops))
	// Assert that there's 1 bug to start, otherwise what are we closing?
	if err != nil || len(issuesDir) != 1 {
		t.Error("Could not read issues directory")
		return
	}
	// check error
	stdout, stderr := captureOutput(func() {
		Close(argumentList{"FooBug"}, config)
	}, t)
	if stderr != "Could not close bug FooBug: Not found FooBug\n" {
		t.Error("Unexpected output on STDERR for Foo-bug")
		fmt.Printf("Got %s\nExpected %s\n", stderr, "Could not close bug FooBug: Not found FooBug")
	}
	if stdout != "" {
		t.Error("Unexpected output on STDOUT for Foo-bug")
	}
	// now success
	stdout, stderr = captureOutput(func() {
		Close(argumentList{"TestBug"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected output on STDERR for Test-bug")
	}
	if stdout != fmt.Sprintf("Removing %s%sissues%sTest\n", dir, sops, sops) {
		t.Error("Unexpected output on STDOUT for Test-bug")
		fmt.Printf("Got %s\nExpected%s", stdout, fmt.Sprintf("Removing %s%sissues%sTest\n", dir, sops, sops))
	}
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s%sissues%s", dir, sops, sops))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	// After closing, there should be 0 bugs.
	if len(issuesDir) != 0 {
		t.Error(fmt.Sprintf("Unexpected number %v not %v of issues in issues dir\n", len(issuesDir), 0))
	}
}

func TestCloseBugByIdentifier(t *testing.T) {
	config := bugs.Config{}
	dir, err := ioutil.TempDir("", "close")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	os.MkdirAll("issues"+sops+"Test", 0700)
	defer os.RemoveAll(dir)

	// On MacOS, /tmp is a symlink, which causes GetDirectory() to return
	// a different path than expected in these tests, so make the issues
	// directory explicit with an environment variable
	err = os.Setenv("FIT", dir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s%sissues%s", dir, sops, sops))
	// Assert that there's 1 bug to start, otherwise what are we closing?
	if err != nil || len(issuesDir) != 1 {
		t.Error("Could not read issues directory")
		return
	}
	stdout, stderr := captureOutput(func() {
		Close(argumentList{"1"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected output on STDERR for Test-bug")
	}
	if stdout != fmt.Sprintf("Removing %s%sissues%sTest\n", dir, sops, sops) {
		t.Error("Unexpected output on STDOUT for Test-bug")
		fmt.Printf("Got %s\nExpected: %s\n", stdout, dir)
	}
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s%sissues%s", dir, sops, sops))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	// After closing, there should be 0 bugs.
	if len(issuesDir) != 0 {
		t.Error(fmt.Sprintf("Unexpected number %v not %v of issues in issues dir\n", len(issuesDir), 0))
	}
}

func TestCloseMultipleIndexesWithLastIndex(t *testing.T) {
	config := bugs.Config{}
	dir, err := ioutil.TempDir("", "closetest")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	os.Setenv("FIT", dir)
	os.MkdirAll("issues"+sops+"Test", 0700)
	os.MkdirAll("issues"+sops+"Test2", 0700)
	os.MkdirAll("issues"+sops+"Test3", 0700)
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s%sissues%s", dir, sops, sops))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	if len(issuesDir) != 3 {
		t.Error(fmt.Sprintf("Unexpected number %v not %v of issues in issues dir\n", len(issuesDir), 3))
	}
	_, stderr := captureOutput(func() {
		Close(argumentList{"1", "3"}, config)
	}, t)
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s%sissues%s", dir, sops, sops))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	// After closing, there should be 1 bug. Otherwise, it probably
	// means that the last error was "invalid index" since indexes
	// were renumbered after closing the first bug.
	if len(issuesDir) != 1 {
		fmt.Printf("%s\n\n", stderr)
		t.Error(fmt.Sprintf("Unexpected number %v not %v of issues in issues dir\n", len(issuesDir), 1))
	}
}

func TestCloseMultipleIndexesAtOnce(t *testing.T) {
	config := bugs.Config{}
	dir, err := ioutil.TempDir("", "closetest")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	os.Setenv("FIT", dir)
	os.MkdirAll("issues"+sops+"Test", 0700)
	os.MkdirAll("issues"+sops+"Test2", 0700)
	os.MkdirAll("issues"+sops+"Test3", 0700)
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s%sissues%s", dir, sops, sops))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	if len(issuesDir) != 3 {
		t.Error(fmt.Sprintf("Unexpected number %v not %v of issues in issues dir\n", len(issuesDir), 3))
	}
	_, _ = captureOutput(func() {
		Close(argumentList{"1", "2"}, config)
	}, t)
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s%sissues%s", dir, sops, sops))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	if len(issuesDir) != 1 {
		t.Error(fmt.Sprintf("Unexpected number %v not %v of issues in issues dir\n", len(issuesDir), 1))
		return
	}

	// 1 and 2 should have closed. If 3 was renumbered after 1 was closed,
	// it would be closed instead.
	if issuesDir[0].Name() != "Test3" {
		t.Error("Closed incorrect issue when closing multiple issues.")
	}
}
