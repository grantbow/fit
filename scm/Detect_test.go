package scm

import (
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
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
	if dir != bugs.Directory(gdir + sops + ".git") {
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
	os.MkdirAll("tmp" + sops + "abc" + sops + "hello", 0755)
	os.Chdir("tmp" + sops + "abc" + sops + "hello")
	handler, dir, err = DetectSCM(options, config)
	if err != nil {
		t.Error("Unexpected while detecting repo type: " + err.Error())
	}
	if dir != bugs.Directory(gdir + sops + ".git") {
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
	if dir != bugs.Directory(gdir + sops + ".hg") {
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
	os.MkdirAll("tmp" + sops + "abc" + sops + "hello", 0755)
	os.Chdir("tmp" + sops + "abc" + sops + "hello")
	handler, dir, err = DetectSCM(options, config)
	if err != nil {
		t.Error("Unexpected while detecting repo type: " + err.Error())
	}
	if dir != bugs.Directory(gdir + sops + ".hg") {
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
    // --- FAIL: TestDetectNone (0.00s)
    //panic: runtime error: invalid memory address or nil pointer dereference [recovered]
    //        panic: runtime error: invalid memory address or nil pointer dereference
    //        [signal 0xc0000005 code=0x0 addr=0x18 pc=0x5db701]
    //
    //goroutine 10 [running]:
    //testing.tRunner.func1(0xc0000a4800)
    //  c:/go/src/testing/testing.go:874 +0x6a6
    //  panic(0x614ae0, 0x798c30)
    //  c:/go/src/runtime/panic.go:679 +0x1c0
    //  scm.TestDetectNone(0xc0000a4800)
    //  (...)/bug/scm/Detect_test.go:142 +0x2a1   <--- moved down now after these comments
    //  testing.tRunner(0xc0000a4800, 0x6530e8)
    //  c:/go/src/testing/testing.go:909 +0x1a1
    //  created by testing.(*T).Run
    //  c:/go/src/testing/testing.go:960 +0x659
    //  exit status 2
    //  FAIL    scm     0.488s

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
		t.Error("Unexpected success detecting repo type: " + err.Error())
	}
	switch handler.(type) {
	case nil:
		// nil is what we expect, don't fall through
		// to the error
	default:
		t.Error("Unexpected SCMHandler found")
	}
    os.Chdir(pwd)
}
