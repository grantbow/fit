package bugapp

import (
	"fmt"
	//	"io"
	//"io/ioutil"
	"os"
	"runtime"
	"testing"
)

// Captures stdout and stderr to ensure that
// a usage line gets printed to Stderr when
// no parameters are specified
func TestVersionOutput(t *testing.T) {

	stdout, stderr := captureOutput(func() {
		Version()
	}, t)

	if stdout == "" {
		t.Error("No output on stdout.")
	}
	if stderr != "" {
		t.Error("Unexpected output on stderr.")
	}
	if stdout != fmt.Sprintf("%s version 0.4 built using %s\n", os.Args[0], runtime.Version()) {
		t.Error("Unexpected output on stdout.")
	}
}
