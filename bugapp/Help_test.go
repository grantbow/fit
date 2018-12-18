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

	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if ! matched {
		t.Errorf("Unexpected output on STDOUT for bugapp/Find_test %v", args)
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}


func TestHelpNilArg(t *testing.T) {
	runhelp(t, "Usage:.*", )
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
	runhelp(t, "Usage:.*", "create")
	runhelp(t, "Usage:.*", "list")
	runhelp(t, "Usage:.*", "edit")
	runhelp(t, "Usage:.*", "status")
	runhelp(t, "Usage:.*", "priorit)")
	runhelp(t, "Usage:.*", "milestone")
	runhelp(t, "Usage:.*", "retitle")
	runhelp(t, "Usage:.*", "rm")
	runhelp(t, "Usage:.*", "find")
	runhelp(t, "Usage:.*", "purge")
	runhelp(t, "Usage:.*", "commit")
	runhelp(t, "Usage:.*", "env")
	runhelp(t, "Usage:.*", "dir")
	runhelp(t, "Usage:.*", "tag")
	runhelp(t, "Usage:.*", "roadmap")
	runhelp(t, "Usage:.*", "id")
	runhelp(t, "Usage:.*", "about")
	runhelp(t, "Usage:.*", "identifiers")
}
func TestHelpValidArgs(t *testing.T) {
	runhelp(t, "Usage:.*", "create", "list")
}

