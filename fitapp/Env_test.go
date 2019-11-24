package fitapp

import (
	"fmt"
	bugs "github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func TestEnvGit(t *testing.T) {
	config := bugs.Config{}
	config.DescriptionFileName = "Description"
	config.FitDirName = "fit"
	var gdir string
	gdir, err := ioutil.TempDir("", "envgit")
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
	// Fake a Fit Directory
	os.Mkdir("Fit", 0755)

	stdout, stderr := captureOutput(func() {
		Env(config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	expected := "Settings:\n\nEditor:.*\nRoot Directory:.*\nFit Directory:.*\nSettings file:.*\n\nVCS Type:.*\ngit Directory:.*\nNeed Committing or Staging:.*"
	// TODO: fix Need Staging output and test
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for fitapp/Env_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
	os.Chdir(pwd)
}
