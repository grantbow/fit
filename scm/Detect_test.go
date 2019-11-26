package scm

import (
	bugs "github.com/grantbow/fit/issues"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestDetectGit(t *testing.T) {
	var config bugs.Config
	config.DescriptionFileName = "Description"
	var gdir string
	gdir, err := ioutil.TempDir("", "gitdetect")
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

	options := make(map[string]bool)
	handler, dir, err := DetectSCM(options, config)
	if err != nil {
		t.Error("Unexpected while detecting repo type: " + err.Error())
	}
	if dir != bugs.Directory(gdir+sops+".git") {
		t.Error("Unexpected directory found when trying to detect git repo" + dir)
	}
	switch handler.(type) {
	case GitManager:
		// GitManager is what we expect, don't fall through
		// to the error
	default:
		t.Error("Unexpected SCMHandler found for Git")
	}

	// Go somewhere higher in the tree and do it again
	os.MkdirAll("tmp"+sops+"abc"+sops+"hello", 0755)
	os.Chdir("tmp" + sops + "abc" + sops + "hello")
	handler, dir, err = DetectSCM(options, config)
	if err != nil {
		t.Error("Unexpected while detecting repo type: " + err.Error())
	}
	if dir != bugs.Directory(gdir+sops+".git") {
		t.Error("Unexpected directory found when trying to detect git repo" + dir)
	}
	switch handler.(type) {
	case GitManager:
		// GitManager is what we expect, don't fall through
		// to the error
	default:
		t.Error("Unexpected SCMHandler found for Git")
	}
	os.Chdir(pwd)
}

func TestDetectHg(t *testing.T) {
	var config bugs.Config
	config.DescriptionFileName = "Description"
	var gdir string
	gdir, err := ioutil.TempDir("", "hgdetect")
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
	// Fake an hg repo
	os.Mkdir(".hg", 0755)

	options := make(map[string]bool)
	handler, dir, err := DetectSCM(options, config)
	if err != nil {
		t.Error("Unexpected while detecting repo type: " + err.Error())
	}
	if dir != bugs.Directory(gdir+sops+".hg") {
		t.Error("Unexpected directory found when trying to detect hg repo" + dir)
	}
	switch handler.(type) {
	case HgManager:
		// HgManager is what we expect, don't fall through
		// to the error
	default:
		t.Error("Unexpected SCMHandler found for Mercurial")
	}

	// Go somewhere higher in the tree and do it again
	os.MkdirAll("tmp"+sops+"abc"+sops+"hello", 0755)
	os.Chdir("tmp" + sops + "abc" + sops + "hello")
	handler, dir, err = DetectSCM(options, config)
	if err != nil {
		t.Error("Unexpected while detecting repo type: " + err.Error())
	}
	if dir != bugs.Directory(gdir+sops+".hg") {
		t.Error("Unexpected directory found when trying to detect hg repo" + dir)
	}
	switch handler.(type) {
	case HgManager:
		// HgManager is what we expect, don't fall through
		// to the error
	default:
		t.Error("Unexpected SCMHandler found for Mercurial")
	}
	os.Chdir(pwd)
}

func TestDetectNone(t *testing.T) {
	t.Skip("windows failure - see scm/Detect_test.go+121")
	// TODO: finish making tests on Windows pass then redo this test
	// seems to be run below a directory with a .git

	var config bugs.Config
	config.DescriptionFileName = "Description"
	var gdir string
	gdir, err := ioutil.TempDir("", "nonedetect") // almost same
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
	// Fake an hg repo
	//os.Mkdir(".hg", 0755)

	options := make(map[string]bool)
	handler, _, err := DetectSCM(options, config)
	if err == nil {
		t.Error("Unexpected success detecting repo type, no error.")
	}
	switch handler.(type) {
	case nil:
		// nil is what we expect, don't fall through
		// to the error
	default:
		t.Error("Unexpected SCMHandler found, type : " + reflect.TypeOf(handler).String())
	}
	os.Chdir(pwd)
}
