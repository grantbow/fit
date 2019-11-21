package fitapp

import (
	//"fmt"
	//	"io"
	//"io/ioutil"
	//"os"
	"regexp"
	//"runtime"
	"testing"
)

// Captures stdout and stderr to ensure that
// a usage line gets printed to Stderr when
// no parameters are specified
func TestVersionOutput(t *testing.T) {

	stdout, stderr := captureOutput(func() {
		PrintVersion()
	}, t)

	if stdout == "" {
		t.Error("No output on stdout.")
	}
	if stderr != "" {
		t.Error("Unexpected output on stderr.")
	}

	expected := "version .* built using .*"
	// was == fmt.Sprintf("%s version %s built using %s\n", os.Args[0], ProgramVersion(), runtime.Version())
	re := regexp.MustCompile(expected)
	matched := re.MatchString(string(stdout))
	if !matched {
		t.Error("Unexpected output on stdout.")
	}
}
