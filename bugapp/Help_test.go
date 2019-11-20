package bugapp

import (
	"fmt"
	//"io/ioutil"
	//"os"
	"regexp"
	"testing"
)

func runhelp(t *testing.T, expected string, args ...string) {
	stdout, stderr := captureOutput(func() {
		Help(args...)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Errorf("Unexpected output on STDOUT for bugapp/Find_test %v", args)
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

func TestHelpNilArg(t *testing.T) {
	runhelp(t, "usage:.*")
}
func TestHelpEmptyArg(t *testing.T) {
	runhelp(t, "usage:.*", "")
}
func TestHelpAnyArg(t *testing.T) {
	runhelp(t, "usage:.*", "any")
}
func TestHelpHelpArg(t *testing.T) {
	runhelp(t, "usage:.*", "help")
}
func TestHelpValidArg(t *testing.T) {
	runhelp(t, "usage:.*", "help", "create")
	runhelp(t, "usage:.*", "help", "list")
	runhelp(t, "usage:.*", "help", "edit")
	runhelp(t, "usage:.*", "help", "status")
	runhelp(t, "usage:.*", "help", "priority")
	runhelp(t, "usage:.*", "help", "milestone")
	runhelp(t, "usage:.*", "help", "retitle")
	runhelp(t, "usage:.*", "help", "rm")
	runhelp(t, "usage:.*", "help", "find")
	runhelp(t, "usage:.*", "help", "purge")
	runhelp(t, "usage:.*", "help", "commit")
	runhelp(t, "usage:.*", "help", "env")
	runhelp(t, "usage:.*", "help", "dir")
	runhelp(t, "usage:.*", "help", "tag")
	runhelp(t, "usage:.*", "help", "roadmap")
	runhelp(t, "usage:.*", "help", "id")
	runhelp(t, "usage:.*", "help", "about")
	runhelp(t, "Issues can be.*", "help", "identifiers")
}
func TestHelpValidArgs(t *testing.T) {
	runhelp(t, "usage:.*", "create", "list")
}
