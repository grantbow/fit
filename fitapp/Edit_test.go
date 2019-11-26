package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"regexp"
	"testing"
)

func TestEditInvalid(t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		Edit(argumentList{"Test"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error:")
		fmt.Printf("Expected: %s\nGot %s\n", "", stderr)
	}
	expected := "Invalid IssueID Test\n"
	if stdout != expected {
		t.Error("Unexpected output on STDOUT")
		fmt.Printf("Expected: %s\nGot %s\n", expected, stdout)
	}

	stdout, stderr = captureOutput(func() {
		Edit(argumentList{"a", "b", "c"}, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error:")
		fmt.Printf("Expected: %s\nGot %s\n", "", stderr)
	}
	expected = "Usage.* edit \\[fieldname\\] IssueID\n\nNo IssueID specified\n"
	re := regexp.MustCompile(expected)
	matched := re.MatchString(string(stdout))
	if !matched {
		t.Error("Unexpected output on STDOUT for fitapp/Commit_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}
