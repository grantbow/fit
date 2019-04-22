package bugs

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetRootDirWithGoodEnvironmentVariable(t *testing.T) {
	var gdir string
	gdir, err := ioutil.TempDir("", "rootdirbug")
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
	err = os.MkdirAll("abc/issues", 0755)
	if err != nil {
		t.Error("Failed creating abc/issues")
		return
	}
	//os.Mkdir("issues", 0755)
	expected := Directory(gdir + "/abc")
	os.Setenv("PMIT", string(expected))
	defer os.Unsetenv("PMIT")
	// PMIT exists and overrides wd
	config := Config{}
	dir := GetRootDir(config)
	if dir != expected {
		t.Errorf("Expected directory %s from environment variable, got %s", expected, string(dir))
	}
}

func TestMissingGetRootDirWithEnvironmentVariable(t *testing.T) {
	var gdir string
	config := Config{}
	gdir, err := ioutil.TempDir("", "rootdirbug")
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
	// PMIT/issues missing so doesn't override wd
	os.Mkdir("../pmit", 0755) // missing issues directory
	defer os.RemoveAll(gdir + "../pmit")
	//os.Mkdir("../pmit/issues", 0755)
	os.Setenv("PMIT", gdir+"../pmit")
	defer os.Unsetenv("PMIT")
	dir := GetRootDir(config)
	if dir != "" {
		t.Errorf("GetRootDir %s environment variable %s", dir, gdir+"../pmit")
	}
}

func TestGetRootDirFromDirectoryTree(t *testing.T) {
	var gdir string
	config := Config{}
	gdir, err := ioutil.TempDir("", "rootdirbug")
	if err == nil {
		os.Chdir(gdir)
		os.Unsetenv("PMIT")
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
	dir := GetRootDir(config)
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
	dir = GetRootDir(config)
	if dir != Directory(gdir) {
		t.Error("Did not get proper directory according to walking the tree: " + dir)
	}
}

func TestNoRoot(t *testing.T) {
	var gdir string
	config := Config{}
	gdir, err := ioutil.TempDir("", "rootdirbug")
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
	if dir := GetRootDir(config); dir != "" {
		t.Error("Found unexpected issues directory." + string(dir))
	}

}

// TestGetIssuesDir was deprecated

func TestGetNoIssuesDir(t *testing.T) {
	var gdir string
	config := Config{}
	gdir, err := ioutil.TempDir("", "rootdirbug")
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
	if dir := GetIssuesDir(config); dir != "" {
		t.Error("Found unexpected issues directory." + string(dir))
	}

}
func TestShortName(t *testing.T) {
	var dir Directory = "/hello/i/am/a/test"
	if short := dir.GetShortName(); short != Directory("test") {
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
