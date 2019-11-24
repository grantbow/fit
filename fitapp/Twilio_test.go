package fitapp

import (
	"fmt"
	bugs "github.com/driusan/bug/bugs"
	//	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

func TestTwilio(t *testing.T) {
	config := bugs.Config{}
	config.FitDirName = "fit"
	dir, err := ioutil.TempDir("", "twiliotest")
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
	cmd := exec.Command("git", "init", "-q")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	stdout, stderr := captureOutput(func() {
		Create(argumentList{"-n", "Test", "bug"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	if stdout != "Created issue: Test bug\n" {
		t.Error("Unexpected output on STDOUT")
		fmt.Printf("Expected: %s\nGot %s\n", "", stdout)
	}
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s%s%s%s", dir, sops, config.FitDirName, sops))
	if err != nil {
		t.Error("Could not read issues directory")
		return
	}
	if len(issuesDir) != 1 {
		t.Errorf("Expected 1 issue  : %v\n", dirDumpFI(issuesDir))
	}
	// TODO: test

	os.Chdir(pwd)
}
