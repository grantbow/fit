package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

func runcommit(expected string, t *testing.T) {
	config := bugs.Config{}
	config.DescriptionFileName = "Description"
	stdout, stderr := captureOutput(func() {
		Commit(argumentList{}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	if stdout != "" {
		t.Error("Unexpected output on STDOUT")
		fmt.Printf("Expected: %s\nGot %s\n", "", stdout)
	}
	stdoutnew, errnew := exec.Command("git", "log").Output()
	if errnew != nil {
		log.Fatal(errnew)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(string(stdoutnew))
	if !matched {
		t.Error("Unexpected output on STDOUT for fitapp/Commit_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdoutnew)
	}
}
func TestCommit(t *testing.T) {
	t.Skip("windows failure - see fitapp/Commit_test.go+41")
	// TODO: finish making tests on Windows pass then redo this test
	dir, err := ioutil.TempDir("", "committest")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	pwd, _ := os.Getwd()
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
	cmd := exec.Command("git", "init", "-q")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	// create
	ioutil.WriteFile(dir+sops+"issues"+sops+"Test"+sops+"Description", []byte("TestBug\n"), 0600)
	expected := "issue. Create issue .Test."
	runcommit(expected, t)
	// update
	ioutil.WriteFile(dir+sops+"issues"+sops+"Test"+sops+"Description", []byte("TestBug-changed\n"), 0600)
	expected = "issue. Update issue .Test."
	runcommit(expected, t)
	// close
	os.RemoveAll(dir + sops + "issues" + sops + "Test")
	expected = "issue. Close issue .Test."
	runcommit(expected, t)
	os.Chdir(pwd)
}
