package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func runlist(args argumentList, expected string, t *testing.T) {
	config := bugs.Config{}
	config.DescriptionFileName = "Description"
	stdout, stderr := captureOutput(func() {
		List(args, config, true)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for bugapp/List_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

func TestList(t *testing.T) {
	config := bugs.Config{}
	config.DescriptionFileName = "Description"
	var gdir string
	gdir, err := ioutil.TempDir("", "listgit")
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
	os.Mkdir("issues", 0755)

	err = os.Setenv("FIT", gdir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}
	// before
	//fmt.Println("before")
	runlist(argumentList{""}, "", t)
	// bug
	_, _ = captureOutput(func() {
		Create(argumentList{"-n", "no_list_bug"}, config)
	}, t)
	// after
	//fmt.Println("a")
	runlist(argumentList{""}, "no_list_bug", t)

	//fmt.Println("b")
	runlist(argumentList{"1"}, "no_list_bug", t)

	//fmt.Println("c")
	runlist(argumentList{"-t"}, "no_list_bug", t)

	//fmt.Println("d")
	runlist(argumentList{"-m", "list"}, "no_list_bug", t)

	//fmt.Println("e")
	//runlist(argumentList{"-r"}, "no_list_bug", t)

	// after
	//fmt.Println("after")
	//file, err := ioutil.ReadFile(fmt.Sprintf("%s%sissues%sno_list_bug%sMilestone", gdir, sops, sops, sops))
	//if err != nil {
	//	t.Error("Could not load milestone file" + err.Error())
	//}
	//if len(file) == 0 {
	//	t.Error("Expected a Milestone file")
	//}
	//runfind(argumentList{"milestone", "foo"}, "Issue 1: no_list_bug\n", t)
    os.Chdir(pwd)
}
