package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

func runmiles(args argumentList, expected string, t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		Milestone(args, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected runmiles error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for fitapp/Milestone_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

func TestMilestone(t *testing.T) {
	config := bugs.Config{}
	config.FitDirName = "fit"
	var gdir string
	gdir, err := ioutil.TempDir("", "milestonegit")
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
	os.Mkdir(config.FitDirName, 0755)

	err = os.Setenv("FIT", gdir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}
	// bug
	_, _ = captureOutput(func() {
		Create(argumentList{"-n", "no_miles_bug"}, config)
	}, t)
	// before
	runfind(argumentList{"milestone", "foo"}, "", t)
	// add
	runmiles(argumentList{"1", "foo"}, "", t) // no cmd as argument
	// force it to test when runmiles doesn't work
	//val := []byte("foo\n")
	//fmt.Println(ioutil.WriteFile(string(gdir)+sops+config.FitDirName+sops+"no_miles_bug"+sops+"Milestone", []byte(val), 0644))
	// check
	//bugDir, _ := ioutil.ReadDir(fmt.Sprintf("%s%s%s%sno_miles_bug", gdir, sops, config.FitDirName, sops))
	//fmt.Printf("readdir len %#v\n", len(bugDir))
	//fmt.Printf("readdir %#v\n", bugDir[0])
	//fmt.Printf("readdir %#v\n", bugDir[1])
	// after
	file, err := ioutil.ReadFile(fmt.Sprintf("%s%s%s%sno_miles_bug%sMilestone", gdir, sops, config.FitDirName, sops, sops))
	if err != nil {
		t.Error("Could not load Milestone file" + err.Error())
	}
	if len(file) == 0 {
		t.Error("Expected a Milestone file")
	}
	runfind(argumentList{"milestone", "foo"}, "Issue 1: no_miles_bug\n", t)
	os.Chdir(pwd)
}
