package main

import (
	"io/ioutil"
	"os"
	_ "os/exec"
	"regexp"
	_ "strings"
	"testing"
)

var bugargtests = []struct {
	input  string
	output string
}{
	{"", ``},
	//{"", `^Warn:`},
	{"--version", ``},
	{"pwd", ``},
	{"env", ``},
	{"find", ``},
	//{"find", `Usage:`},
	{"status", ``},
	//{"status", `Usage:`},
	{"list", ``},
	{"help", ``},
	//{"help", `^Warn:`},
	{"pwd --help aha yes", ``},
	//{"pwd --help aha yes", `^Warn:`},
}

func TestBugArgParser(t *testing.T) {
	//log.Print("PATH " + os.Getenv("PATH"))
	for _, tt := range bugargtests {
		// //runcmd := exec.Command("sh", "-c", tt.input) // input
		// //out, err := runcmd.CombinedOutput()
		//numArgs := 0
		//args := strings.Fields(tt.input)
		//for _, arg := range args {
		//	os.Args[numArgs] = arg
		//	numArgs += 1
		//}
		//for i := numArgs; i < 4; i++ {
		//	os.Args[i] = ""
		//}

		out, err := captureOutput(main, t) // TODO: needs tt.input for the main func via an env var for testing
		if err != "" {
			//t.Error("Could not exec command bug: " + err.Error())
			t.Error("Could not exec command bug: " + err)
		}
		found, ferr := regexp.Match(tt.output, []byte(out)) // output
		if ferr != nil {
			t.Error("Usage output: " + ferr.Error())
		} else if !found {
			t.Errorf("Unexpected usage, wanted to match %q, got %q", tt.output, tt.input)
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
