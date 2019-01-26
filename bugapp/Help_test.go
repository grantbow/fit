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
	runhelp(t, "Usage:.*")
}
func TestHelpEmptyArg(t *testing.T) {
	runhelp(t, "Usage:.*", "")
}
func TestHelpAnyArg(t *testing.T) {
	runhelp(t, "Usage:.*", "any")
}
func TestHelpHelpArg(t *testing.T) {
	runhelp(t, "Usage:.*", "help")
}
func TestHelpValidArg(t *testing.T) {
	runhelp(t, "Usage:.*", "help", "create")
	runhelp(t, "Usage:.*", "help", "list")
	runhelp(t, "Usage:.*", "help", "edit")
	runhelp(t, "Usage:.*", "help", "status")
	runhelp(t, "Usage:.*", "help", "priority")
	runhelp(t, "Usage:.*", "help", "milestone")
	runhelp(t, "Usage:.*", "help", "retitle")
	runhelp(t, "Usage:.*", "help", "rm")
	runhelp(t, "Usage:.*", "help", "find")
	runhelp(t, "Usage:.*", "help", "purge")
	runhelp(t, "Usage:.*", "help", "commit")
	runhelp(t, "Usage:.*", "help", "env")
	runhelp(t, "Usage:.*", "help", "dir")
	runhelp(t, "Usage:.*", "help", "tag")
	runhelp(t, "Usage:.*", "help", "roadmap")
	runhelp(t, "Usage:.*", "help", "id")
	runhelp(t, "Usage:.*", "help", "about")
	runhelp(t, "Bugs can be.*", "help", "identifiers")
}
func TestHelpValidArgs(t *testing.T) {
	runhelp(t, "Usage:.*", "create", "list")
}
