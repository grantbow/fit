package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	//	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

func TestPurgeNoEditor(t *testing.T) {
	config := bugs.Config{}
	config.FitDirName = "fit"
	dir, err := ioutil.TempDir("", "purgetest")
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

    // Create
	stdout, stderr := captureOutput(func() {
		Create(argumentList{"-n", "Test", "bug"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	if stdout != "Created issue: Test bug\n" {
		fmt.Printf("Expected: %s\nGot %s\n", "", stdout)
		t.Error("Unexpected output on STDOUT")
	}
	issuesDir, err := ioutil.ReadDir(fmt.Sprintf("%s%s%s%s", dir, sops, config.FitDirName, sops))
	if err != nil {
		t.Error("Could not read " + config.FitDirName + " directory")
		return
	}
	if len(issuesDir) != 1 {
		t.Errorf("Expected 1 issue  : %v\n", dirDumpFI(issuesDir))
	}

    // Create does not automatically Commit
    //Commit(argumentList{}, config)

    // Purge
	stdout, stderr = captureOutput(func() {
		Purge(config)
	}, t)
	issuesDir, err = ioutil.ReadDir(fmt.Sprintf("%s%s%s%s", dir, sops, config.FitDirName, sops))
	if err != nil {
		t.Error("Could not purge " + config.FitDirName + " directory")
		return
	}
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	expected := "Removing " + config.FitDirName + ".Test-bug."
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		fmt.Printf("Expected: %s\nGot %s\n", expected, stdout)
		t.Error("Unexpected output on STDOUT")
	}
	if len(issuesDir) != 0 {
		t.Errorf("Expected 0 issues : %v\n", dirDumpFI(issuesDir))
	}
	os.Chdir(pwd)
}
