package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"regexp"
	"testing"
)

func runrelabel(label string, args argumentList, config bugs.Config, expected string, t *testing.T) {
	stdout, _ := captureOutput(func() {
		Relabel(args, config)
	}, t)
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Errorf("Unexpected output on STDOUT for bugapp/Pwd_test %s.", label)
		fmt.Printf("Expected to match: %s\nGot: %s\n", expected, stdout)
	}
}

func TestRelabel(t *testing.T) {
	config := bugs.Config{}
	args := argumentList{"1"} // < 2
	test := tester{}          // from Bug_test.go
	test.Setup()
	defer test.Teardown()
	//config.BugDir = bugs.GetRootDir(config)
	//bugDir := bugs.GetIssuesDir(config) + bugs.TitleToDir(args[0])
	rootDir := bugs.GetIssuesDir(config)
	//bugDir := rootDir + bugs.TitleToDir("Test Bug")
	//fmt.Print("bugDir ", bugDir, "\n")

	expected := "Usage: .*"
	runrelabel("usage", args, config, expected, t)

	args = argumentList{"bad", "bar"} // bad
	expected = "Could not load bug: .*"
	runrelabel("bad", args, config, expected, t)

	args = argumentList{"1", "Error Bug"} // rename err
	// before chmod
	//fi, _ := os.Open(string(rootDir))
	//stat, _ := fi.Stat()
	//fmt.Println(dirDumpFI([]os.FileInfo{stat}))
	//fmt.Println(dirDump(string(rootDir)))
	//fmt.Printf("mode %v\n", stat.Mode())
	// chmod 500 temp parent directory, read and execute
	err := os.Chmod(string(rootDir), 0500) // change
	check(err)
	// after chmod
	//fi, _ = os.Open(string(rootDir))
	//stat, _ = fi.Stat()
	//fmt.Println(dirDumpFI([]os.FileInfo{stat}))
	//fmt.Println(dirDump(string(rootDir)))
	//fmt.Printf("mode %v\n", stat.Mode())
	expected = "Moving .*\\nError moving directory\\n"
	runrelabel("rename err", args, config, expected, t)

	args = argumentList{"1", "Success Bug"} // good
	// chmod 700 temp parent directory
	err = os.Chmod(string(rootDir), 0700) // change
	check(err)
	expected = "Moving .*"
	runrelabel("good", args, config, expected, t)
}
