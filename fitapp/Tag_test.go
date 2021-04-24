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

//func getAllTags() []string {
//func Tag(Args argumentList) {

func runtagsassigned(args argumentList, expected string, t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		TagsAssigned(args, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for fitapp/Tag_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

func runtagsnone(args argumentList, expected string, t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		TagsNone(config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for fitapp/Tag_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

func runtag(args argumentList, expected string, t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		Tag(args, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for fitapp/Tag_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

func TestTag(t *testing.T) {
	config := bugs.Config{}
	config.FitDirName = "fit"
	var gdir string
	gdir, err := ioutil.TempDir("", "taggit")
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
		Create(argumentList{"-n", "no_tag_bug"}, config)
	}, t)
	// before
	runfind(argumentList{"tags", "foo"}, "", t) // find uses tags but tag uses tag
	runtagsnone(argumentList{""}, "No tags assigned:", t)
	// add with too few args
	runtag(argumentList{}, "", t) // no cmd as argument
	// add
	runtag(argumentList{"1", "foo"}, "", t) // no cmd as argument
	// after
	runfind(argumentList{"tags", "foo"}, "Issue 1: no_tag_bug \\(foo\\)\n", t)                                              // boolean flags not tags
	_, err = ioutil.ReadFile(fmt.Sprintf("%s%s%s%sno_tag_bug%stags%sfoo", gdir, sops, config.FitDirName, sops, sops, sops)) // file is empty
	if err != nil {
		t.Error("Could not load tags/foo file" + err.Error())
	}
	// tags can have more than one
	runtag(argumentList{"1", "bar"}, "", t)                                                                                 // no cmd as argument
	runfind(argumentList{"tags", "foo"}, "Issue 1: no_tag_bug \\(bar, foo\\)\n", t)                                         // boolean flags not tags
	_, err = ioutil.ReadFile(fmt.Sprintf("%s%s%s%sno_tag_bug%stags%sbar", gdir, sops, config.FitDirName, sops, sops, sops)) // file is empty
	if err != nil {
		t.Error("Could not load tags/bar file" + err.Error())
	}
	// non existent bug
	runtag(argumentList{"3", "baz"}, "", t) // no cmd as argument
	// --rm a tag
	runtag(argumentList{"--rm", "1", "bar"}, "", t) // no cmd as argument
	os.Chdir(pwd)
}
func TestTagsAssigned(t *testing.T) {
	runtagsassigned(argumentList{""}, "<none assigned yet>", t)
}
func TestTagsNone(t *testing.T) {
	runtagsnone(argumentList{""}, "No tags assigned:", t)
}
