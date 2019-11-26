package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"os"
	"regexp"
	"testing"
)

func runretitle(label string, args argumentList, config bugs.Config, expected string, t *testing.T) {
	stdout, _ := captureOutput(func() {
		Relabel(args, config) // name might change
	}, t)
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Errorf("Unexpected output on STDOUT for fitapp/Retitle_test %s.", label)
		fmt.Printf("Expected to match: %s\nGot: %s\n", expected, stdout)
	}
}

func TestRetitle(t *testing.T) {
	config := bugs.Config{}
	args := argumentList{"1"} // < 2
	test := tester{}          // from Pwd_test.go ; originally from Bug_test.go
	test.Setup()
	defer test.Teardown()
	//bugDir := bugs.FitDirer(config) + bugs.TitleToDir(args[0])
	fitDir := bugs.FitDirer(config)
	//bugDir := fitDir + bugs.TitleToDir("Test Bug")
	//fmt.Print("bugDir ", bugDir, "\n")

	expected := "Usage: .*"
	runretitle("usage", args, config, expected, t)

	args = argumentList{"bad", "bar"} // bad
	expected = "Could not load issue: .*"
	runretitle("bad", args, config, expected, t)

	/*
		        this test fails on windows
		        removing write on directory doesn't cause the same error
		        TODO: finish making tests on Windows pass then redo this test

			args = argumentList{"1", "Error Bug"} // rename err
			// before chmod
			//fi, _ := os.Open(string(fitDir))
			//stat, _ := fi.Stat()
			//fmt.Println(dirDumpFI([]os.FileInfo{stat}))
			//fmt.Println(dirDump(string(fitDir)))
			//fmt.Printf("mode %v\n", stat.Mode())

			// chmod 500 temp parent directory, read and execute
			err := os.Chmod(string(fitDir), 0500) // remove write permission
			check(err)
			// after chmod
			//fi, _ = os.Open(string(fitDir))
			//stat, _ = fi.Stat()
			//fmt.Println(dirDumpFI([]os.FileInfo{stat}))
			//fmt.Println(dirDump(string(fitDir)))
			//fmt.Printf("mode %v\n", stat.Mode())
			expected = "Moving .*\\nError moving directory\\n"
			runretitle("rename err", args, config, expected, t)
	*/

	args = argumentList{"1", "Success Bug"} // good
	// chmod 700 temp parent directory
	err := os.Chmod(string(fitDir), 0700) // change
	check(err)
	expected = "Moving .*"
	runretitle("good", args, config, expected, t)
}
