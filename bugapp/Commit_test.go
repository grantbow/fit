package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

func runcommit(expected string, t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		Commit(ArgumentList{}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	if stdout != "" {
		t.Error("Unexpected output on STDOUT")
		fmt.Printf("Expected: %s\nGot %s\n", "", stdout)
	}
	stdoutnew, errnew := exec.Command("git", "log",).Output()
	if errnew != nil {
		log.Fatal(errnew)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(string(stdoutnew))
	if  ! matched {
		t.Error("Unexpected output on STDOUT for bugapp/Commit_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdoutnew)
	}
}
func TestCommit(t *testing.T) {
	dir, err := ioutil.TempDir("", "commit")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	os.Chdir(dir)
	os.MkdirAll("issues/Test", 0700)
	defer os.RemoveAll(dir)

	// On MacOS, /tmp is a symlink, which causes GetDirectory() to return
	// a different path than expected in these tests, so make the issues
	// directory explicit with an environment variable
	err = os.Setenv("PMIT", dir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}
	cmd := exec.Command("git", "init", "-q")
	err = cmd.Run() ; if err != nil {
		log.Fatal(err)
	}
	// create
	ioutil.WriteFile(dir+"/issues/Test/Description", []byte("TestBug\n"), 0600)
	expected := "bug. Create issue .Test."
	runcommit(expected, t)
	// update
	ioutil.WriteFile(dir+"/issues/Test/Description", []byte("TestBug-changed\n"), 0600)
	expected = "bug. Update issue .Test."
	runcommit(expected, t)
	// close
	os.RemoveAll(dir+"/issues/Test")
	expected = "bug. Close issue .Test."
	runcommit(expected, t)
}
