package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

//func Testfind(t *testing.T) {
//	// find(string, []string)
//}
func runfind(args argumentList, expected string, t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		Find(args, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for bugapp/Find_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

func TestFindUsage(t *testing.T) {
	args := argumentList{"any"} // < 2
	expected := "Usage: .* find \\{tags, status, priority, milestone\\} value1 \\[value2 ...\\]\n"
	runfind(args, expected, t)
}
func TestFindSubcommandUnknown(t *testing.T) {
	runfind(argumentList{"unk_sub", "matchstring"}, "Unknown command:.*\n", t)
}
func TestFindSubcommandUnknownGTOne(t *testing.T) {
	runfind(argumentList{"unk_sub", "not_found", "more"}, "Unknown command: .*\n", t)
}
func TestFindSubcommands(t *testing.T) {
	config := bugs.Config{}
	var gdir string
	gdir, err := ioutil.TempDir("", "findgit")
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

	runfind(argumentList{"tags", "matchstring"}, "", t)
	runfind(argumentList{"status", "matchstring"}, "", t)
	runfind(argumentList{"priority", "matchstring"}, "", t)
	runfind(argumentList{"milestone", "matchstring"}, "", t)

	// bug "id bug"
	_, _ = captureOutput(func() {
		Create(argumentList{"-n", "no_id_bug", "--tag", "foo"}, config)
	}, t)
	runfind(argumentList{"tags", "foo"}, "Issue 1: no_id_bug \\(foo\\)\n", t)
	runfind(argumentList{"tags", "matchstring"}, "", t) // still not found
}
