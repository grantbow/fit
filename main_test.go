package main

import (
	"io/ioutil"
	"os"
	_ "os/exec"
	"regexp"
	"testing"
)

var bugargtests = []struct {
	input  string
	output string
}{
	{"./main", `^Usage:`},
	//{"./main --help", `^Usage:`},
}

func TestBugArgParser(t *testing.T) {
	//log.Print("PATH " + os.Getenv("PATH"))
	for _, tt := range bugargtests {
		//runcmd := exec.Command("sh", "-c", tt.input) // input
		//out, err := runcmd.CombinedOutput()
		out, err := captureOutput(main, t)
		if err != "" {
			//t.Error("Could not exec command bug: " + err.Error())
			t.Error("Could not exec command bug: " + err)
		}
		found, ferr := regexp.Match(tt.output, []byte(out)) // output
		if ferr != nil {
			t.Error("Usage output: " + ferr.Error())
		} else if !found {
			t.Errorf("Unexpected usage, got %q, want to match %q", tt.input, tt.output)
		}
	}
}

func captureOutput(f func(), t *testing.T) (string, string) {
	// Capture STDOUT with a pipe
	stdout := os.Stdout
	stderr := os.Stderr
	so, op, _ := os.Pipe() //outpipe
	oe, ep, _ := os.Pipe() //errpipe
	defer func(stdout, stderr *os.File) {
		os.Stdout = stdout
		os.Stderr = stderr
	}(stdout, stderr)

	os.Stdout = op
	os.Stderr = ep

	f()

	os.Stdout = stdout
	os.Stderr = stderr

	op.Close()
	ep.Close()

	errOutput, err := ioutil.ReadAll(oe)
	if err != nil {
		t.Error("Could not get output from stderr")
	}
	stdOutput, err := ioutil.ReadAll(so)
	if err != nil {
		t.Error("Could not get output from stdout")
	}
	return string(stdOutput), string(errOutput)
}
