package bugs

import (
	"io/ioutil"
	"os"
	"testing"
)

//var dops = Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

func TestRootDirerWithGoodEnvironmentVariable(t *testing.T) {
	var gdir string
	gdir, err := ioutil.TempDir("", "rootdirbug")
    pwd, _ := os.Getwd()
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating TempDir")
		return
	}
	err = os.MkdirAll("abc"+sops+"issues", 0755)
	if err != nil {
		t.Error("Failed creating abc/issues")
		return
	}
	//os.Mkdir("issues", 0755)
	expected := Directory(gdir + string(os.PathSeparator) + "abc")
	os.Setenv("FIT", string(expected))
	defer os.Unsetenv("FIT")
	// FIT exists and overrides wd
	config := Config{}
	dir := RootDirer(config)
	if dir != expected {
		t.Errorf("Expected directory %s from environment variable, got %s", expected, string(dir))
	}
    os.Chdir(pwd)
}

func TestMissingRootDirerWithEnvironmentVariable(t *testing.T) {
	var gdir string
	config := Config{}
	gdir, err := ioutil.TempDir("", "rootdirbug")
    pwd, _ := os.Getwd()
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating temporary directory")
		return
	}
	// FIT/issues missing so doesn't override wd
	fitdir := ".." + string(os.PathSeparator) + "fit"
	os.Mkdir(fitdir, 0755) // missing issues directory
	defer os.RemoveAll(gdir + fitdir)
	//os.Mkdir("../fit/issues", 0755)
	os.Setenv("FIT", gdir+fitdir)
	defer os.Unsetenv("FIT")
	dir := RootDirer(config)
	if dir != "" {
		t.Errorf("RootDirer %s environment variable %s", dir, gdir+fitdir)
	}
    os.Chdir(pwd)
}

func TestRootDirerFromDirectoryTree(t *testing.T) {
	var gdir string
	config := Config{}
	gdir, err := ioutil.TempDir("", "rootdirbug")
    pwd, _ := os.Getwd()
	if err == nil {
		os.Chdir(gdir)
		os.Unsetenv("FIT")
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating temporary directory")
		return
	}
	// Make sure we get the right directory from the top level
	os.Mkdir("issues", 0755)
	dir := RootDirer(config)
	if dir != Directory(gdir) {
		t.Error("Did not get proper directory according to walking the tree: " + dir)
	}
	// Now go deeper into the tree and try the same thing..
	err = os.MkdirAll("abc/123", 0755)
	if err != nil {
		t.Error("Could not make directory for testing")
	}
	err = os.Chdir("abc/123")
	if err != nil {
		t.Error("Could not change directory for testing")
	}
	dir = RootDirer(config)
	if dir != Directory(gdir) {
		t.Error("Did not get proper directory according to walking the tree: " + dir)
	}
    os.Chdir(pwd)
}

func TestNoRoot(t *testing.T) {
	var gdir string
	config := Config{}
	gdir, err := ioutil.TempDir("", "rootdirbug")
    pwd, _ := os.Getwd()
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating temporary directory")
		return
	}
	// Don't create an issues directory. Just try and get the directory
	if dir := RootDirer(config); dir != "" {
		t.Error("Found unexpected issues directory." + string(dir))
	}
    os.Chdir(pwd)
}

// TestIssuesDirer was deprecated

func TestNoIssuesDirer(t *testing.T) {
	var gdir string
	config := Config{}
	gdir, err := ioutil.TempDir("", "rootdirbug")
    pwd, _ := os.Getwd()
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating temporary directory")
		return
	}
	// Don't create an issues directory. Just try and get the directory.
	// empty is accepted! why?
	if dir := IssuesDirer(config); dir != "" {
		t.Error("Found unexpected issues directory." + string(dir))
	}
    os.Chdir(pwd)
}
func TestShortName(t *testing.T) {
	var dir Directory = dops + "hello"+dops+"i"+dops+"am"+dops+"a"+dops+"test"
	if short := dir.ShortNamer(); short != Directory("test") {
		t.Error("Unexpected short name: " + string(short))
	}
}
func TestDirectoryToTitle(t *testing.T) {
	var assertTitle = func(directory, title string) {
		dir := Directory(directory)
		if dir.ToTitle() != title {
			t.Error("Failed on " + directory + ": got " + dir.ToTitle() + " but expected " + title)
		}
	}
	assertTitle("Test", "Test")
	assertTitle("Test-Multiword", "Test Multiword")
	assertTitle("Test--Dash", "Test-Dash")
	assertTitle("Test---Dash", "Test--Dash")
	assertTitle("Test_--TripleDash", "Test --TripleDash")
	assertTitle("Test_-_What", "Test - What")
}
