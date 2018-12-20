package main

import (
	"os/exec"
	"regexp"
	"testing"
)

var bugargtests = []struct {
	input string
	output string
}{
	{"bug", `^Usage:`},
	{"bug --help", `^Usage:`},
}
func TestBugArgParser(t *testing.T) {
	//log.Print("PATH " + os.Getenv("PATH"))
	for _, tt := range bugargtests {
		runcmd := exec.Command("sh", "-c", tt.input) // input
		out, err := runcmd.CombinedOutput()
		if err != nil {
			t.Error("Could not exec command bug: " + err.Error())
		}
		found, ferr := regexp.Match(tt.output, out) // output
		if ferr != nil {
			t.Error("Usage output: " + ferr.Error())
		} else if !found {
			t.Errorf("Unexpected usage, got %q, want to match %q", tt.input, tt.output )
		}
	}
}
