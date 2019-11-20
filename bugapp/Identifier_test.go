package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

func rungenid(input string, expected string, t *testing.T) {
	out := generateID(input)
	re := regexp.MustCompile(expected)
	matched := re.MatchString(out)
	if !matched {
		t.Error("Unexpected output bugapp/Identifier_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, out)
	}
}
func runid(t *testing.T, expected string, args argumentList) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		Identifier(args, config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for bugapp/Identifier_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}
func runidsassigned(args argumentList, expected string, t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		IdsAssigned(config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for bugapp/IdsAssigned_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}
func runidsnone(args argumentList, expected string, t *testing.T) {
	config := bugs.Config{}
	stdout, stderr := captureOutput(func() {
		IdsNone(config)
	}, t)
	if stderr != "" {
		t.Error("Unexpected error: " + stderr)
	}
	re := regexp.MustCompile(expected)
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for bugapp/IdsNone_test")
		fmt.Printf("Expected: %s\nGot: %s\n", expected, stdout)
	}
}
func TestIdGen(t *testing.T) {
	rungenid("test string", "b6612", t)
}
func TestIdUsage(t *testing.T) {
	runid(t, "Usage: .* id <IssueID> \\[value\\]\n", argumentList{})
}
func TestIdInvalid(t *testing.T) {
	runid(t, "Invalid IssueID: Not found test\n", argumentList{"test"})
}
func TestIdGenerate(t *testing.T) {
	config := bugs.Config{}
	config.IssuesDirName = "fit"
	var gdir string
	gdir, err := ioutil.TempDir("", "idgit")
	pwd, _ := os.Getwd()
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
	os.Mkdir(config.IssuesDirName, 0755)
	err = os.Setenv("FIT", gdir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	// bug
	_, _ = captureOutput(func() {
		Create(argumentList{"-n", "no_id_bug"}, config)
	}, t)
	runid(t, "Id not defined\n", argumentList{"1"})

	runid(t, "Generated id .* for issue\n", argumentList{"1", "--generate-id"})
	file, err := ioutil.ReadFile(fmt.Sprintf("%s%s%s%sno_id_bug%sIdentifier", gdir, sops, config.IssuesDirName, sops, sops))
	if err != nil {
		t.Error("Could not load an Identifier file for Test bug" + err.Error())
	}
	if len(file) == 0 {
		t.Error("Expected an Identifier file")
	}
	os.Chdir(pwd)
}
func TestIdsAssigned(t *testing.T) {
	runidsassigned(argumentList{""}, "Ids used in current tree", t)
}
func TestIdNone(t *testing.T) {
	runidsnone(argumentList{""}, "No ids assigned", t)
}
