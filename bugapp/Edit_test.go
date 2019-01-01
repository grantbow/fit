package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"regexp"
	"testing"
)

func TestEditInvalid(t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		Edit(ArgumentList{"Test"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error:")
		fmt.Printf("Expected: %s\nGot %s\n", "", stderr)
	}
	expected := "Invalid BugID Test\n"
	if stdout != expected {
		t.Error("Unexpected output on STDOUT")
		fmt.Printf("Expected: %s\nGot %s\n", expected, stdout)
	}

	stdout, stderr  = captureOutput(func() {
		Edit(ArgumentList{"a","b","c"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error:")
		fmt.Printf("Expected: %s\nGot %s\n", "", stderr)
	}
	expected  = "Usage.* edit \\[fieldname\\] BugID\n\nNo BugID specified\n"
	re := regexp.MustCompile(expected)
	matched := re.MatchString(string(stdout))
	if  ! matched {
		t.Error("Unexpected output on STDOUT for bugapp/Commit_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}

