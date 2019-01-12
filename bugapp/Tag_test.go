package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

//func getAllTags() []string {
//func Tag(Args argumentList) {

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
		t.Error("Unexpected output on STDOUT for bugapp/Tag_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

func TestTag(t *testing.T) {
	config := bugs.Config{}
	var gdir string
	gdir, err := ioutil.TempDir("", "taggit")
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

	err = os.Setenv("PMIT", gdir)
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
	// add
	runtag(argumentList{"1", "foo"}, "", t) // no cmd as argument
	// force it to test when runmiles doesn't work
	//val := []byte("foo\n")
	//fmt.Println(ioutil.WriteFile(string(gdir)+"/issues/no_tag_bug/Tag", []byte(val), 0644))
	// check
	//bugDir, _ := ioutil.ReadDir(fmt.Sprintf("%s/issues/no_tag_bug/tags", gdir))
	//fmt.Printf("readdir len %#v\n", len(bugDir))
	//fmt.Printf("readdir %#v\n", bugDir[0])
	//fmt.Printf("readdir %#v\n", bugDir[1])
	// after
	runfind(argumentList{"tags", "foo"}, "Issue 1: no_tag_bug \\(foo\\)\n", t)   // boolean flags not tags
	_, err = ioutil.ReadFile(fmt.Sprintf("%s/issues/no_tag_bug/tags/foo", gdir)) // file is empty
	if err != nil {
		t.Error("Could not load tags/foo file" + err.Error())
	}
	// tags can have more than one
	runtag(argumentList{"1", "bar"}, "", t)                                         // no cmd as argument
	runfind(argumentList{"tags", "foo"}, "Issue 1: no_tag_bug \\(bar, foo\\)\n", t) // boolean flags not tags
	_, err = ioutil.ReadFile(fmt.Sprintf("%s/issues/no_tag_bug/tags/bar", gdir))    // file is empty
	if err != nil {
		t.Error("Could not load tags/bar file" + err.Error())
	}
}
