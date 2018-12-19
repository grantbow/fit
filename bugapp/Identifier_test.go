package bugapp

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func rungenid(t *testing.T, expected string, input string) {
	out := generateID(input)
	re := regexp.MustCompile(expected)
	matched := re.MatchString(out)
	if ! matched {
		t.Error("Unexpected output bugapp/Identifier_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, out)
	}
}
func runid(t *testing.T, expected string, args ArgumentList) {
	stdout, stderr := captureOutput(func() {
		Identifier(args)
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
		t.Error("Unexpected output on STDOUT for bugapp/Identifier_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}
func TestIdGen(t *testing.T) {
	rungenid(t, "b6612", "test string")
}
func TestIdUsage(t *testing.T) {
	runid(t, "Usage: .* identifier BugID \\[value\\]\n", ArgumentList{})
}
func TestIdInvalid(t *testing.T) {
	runid(t, "Invalid BugID: Could not find bug test\n", ArgumentList{"test"})
}
func TestIdGenerate(t *testing.T) {
	var gdir string
	gdir, err := ioutil.TempDir("", "idgit")
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating temporary directory for detect")
		return
	}
	// Make an issues Directory
	os.Mkdir("issues", 0755)
	err = os.Setenv("PMIT", gdir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	// bug
	_, _ = captureOutput(func() {
		Create(ArgumentList{"-n", "no_id_bug"})
	}, t)
	runid(t, "Identifier not defined\n", ArgumentList{"1"})

	runid(t, "Generated id .* for bug\n", ArgumentList{"1", "--generate"})
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/issues/no_id_bug/Identifier", gdir))
	if err != nil {
		t.Error("Could not load description file for Test bug" + err.Error())
	}
	if len(file) == 0 {
		t.Error("Expected an Identifier file")
	}
}

