package bugapp

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func TestEnvGit(t *testing.T) {
	var gdir string
	gdir, err := ioutil.TempDir("", "envgit")
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
	// Fake an Issues Directory
	os.Mkdir(".git", 0755)

	stdout, stderr := captureOutput(func() {
		Env()
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	expected := fmt.Sprintf("Settings used by this command:\n\nEditor:.*\nIssues Directory:.*\n\nSCM Type:.*\ngit Directory:.*\n")
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if ! matched {
		t.Error("Unexpected output on STDOUT for bugapp/Env_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

