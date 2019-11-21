package fitapp

import (
	"fmt"
	bugs "github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// Priority and Status find printed differently

func runstatus(args argumentList, expected string, t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		Status(args, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for fitapp/Status_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

func TestStatus(t *testing.T) {
	config := bugs.Config{}
	config.IssuesDirName = "fit"
	var gdir string
	gdir, err := ioutil.TempDir("", "statusgit")
	pwd, _ := os.Getwd()
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating temporary directory for detect")
		return
	}
	// Fake a git repo
	os.Mkdir(".git", 0755)
	// Make an issues Directory
	os.Mkdir(config.IssuesDirName, 0755)

	err = os.Setenv("FIT", gdir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}
	// bug
	_, _ = captureOutput(func() {
		Create(argumentList{"-n", "no_status_bug"}, config)
	}, t)
	// before
	runfind(argumentList{"status", "foo"}, "", t)
	// add
	runstatus(argumentList{"1", "foo"}, "", t) // no cmd as argument
	// force it to test when runmiles doesn't work
	//val := []byte("foo\n")
	//fmt.Println(ioutil.WriteFile(string(gdir)+sops+config.IssuesDirName+sops+"no_status_bug"+sops+"Status", []byte(val), 0644))
	// check
	//bugDir, _ := ioutil.ReadDir(fmt.Sprintf("%s%s%s%sno_status_bug", gdir, sops, config.IssuesDirName, sops))
	//fmt.Printf("readdir len %#v\n", len(bugDir))
	//fmt.Printf("readdir %#v\n", bugDir[0])
	//fmt.Printf("readdir %#v\n", bugDir[1])
	// after
	runfind(argumentList{"status", "foo"}, "Issue 1: no_status_bug \\(Status: foo\\)\n", t)
	file, err := ioutil.ReadFile(fmt.Sprintf("%s%s%s%sno_status_bug%sStatus", gdir, sops, config.IssuesDirName, sops, sops))
	if err != nil {
		t.Error("Could not load Status file" + err.Error())
	}
	if len(file) == 0 {
		t.Error("Expected a Status file")
	}
	os.Chdir(pwd)
}
