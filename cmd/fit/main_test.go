package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	_ "strings"
	"testing"
)

type Config struct {
	FitDir                    string `json:"FitDir"`     // runtime only
	FitDirName                string `json:"FitDirName"` // runtime only
	ScmDir                    string `json:"ScmDir"`     // runtime only
	ScmType                   string `json:"ScmType"`    // runtime only
	FitYmlDir                 string `json:"FitYmlDir"`  // runtime only
	FitYml                    string `json:"FitYml"`     // runtime only
	DefaultDescriptionFile    string `json:"DefaultDescriptionFile"`
	ImportXmlDump             bool   `json:"ImportXmlDump"`
	ImportCommentsTogether    bool   `json:"ImportCommentsTogether"`
	ProgramVersion            string `json:"ProgramVersion"`
	DescriptionFileName       string `json:"DescriptionFileName"`
	TagKeyValue               bool   `json:"TagKeyValue"`
	NewFieldAsTag             bool   `json:"NewFieldAsTag"`
	NewFieldLowerCase         bool   `json:"NewFieldLowerCase"`
	GithubPersonalAccessToken string `json:"GithubPersonalAccessToken"`
	TwilioAccountSid          string `json:"TwilioAccountSid"`
	TwilioAuthToken           string `json:"TwilioAuthToken"`
	TwilioPhoneNumberFrom     string `json:"TwilioPhoneNumberFrom"`
	FitSite                   string `json:"FitSite"`
	MultipleFitDirs           bool   `json:"MultipleFitDirs"`
	CloseStatusTag            bool   `json:"CloseStatusTag"`
	CloseMove                 bool   `json:"CloseMove"`
	FitClosedDirName          string `json:"FitClosedDirName"`
	ClosePreventDelete        bool   `json:"ClosePreventDelete"`
	IdAbbreviate              bool   `json:"IdAbbreviate"`
	IdAutomatic               bool   `json:"IdAutomatic"`
}

var firstbugargtests = []struct {
	input  string
	output string
}{
	{"", ``},
}
var setupbugargtests = []struct {
	input  string
	output string
}{
	{"", `usage:`},
	{"--version", `version`},
	{"pwd", `fit`},
	{"env", `Editor`},
	{"find", `Usage:`},
	//{"status", `Usage:`},
	{"list", `list`},
	{"help", `usage:`},
	{"pwd --help aha yes", `usage:`},
}

var binaryname = "fit"
var binarypath string

func TestBugArgParser(t *testing.T) {
	for _, tt := range firstbugargtests {
		out, err := captureOutput(main, t) // TODO: needs tt.input for the main func via an env var for testing
		if err != "" {
			//t.Error("Could not exec command bug: " + err.Error())
			t.Error("Could not exec command bug: " + err)
		}
		found, ferr := regexp.Match(``, []byte(out)) // tt.output
		if ferr != nil {
			t.Error("Usage output: " + ferr.Error())
		} else if !found {
			t.Errorf("Unexpected usage, wanted to match %q, got %q", ``, tt.input) // tt.output
		}
	}

	// setup
	config := Config{}
	config.DescriptionFileName = "Description"
	config.FitDirName = "fit"
	var gdir string
	pwd, _ := os.Getwd()
	gdir, err := ioutil.TempDir("", "main")
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail.
		gdir, _ = os.Getwd()
	} else {
		t.Error("Failed creating temporary directory for detect")
		return
	}
	// Fake a git repo
	os.Mkdir(".git", 0755)
	// Make an issues Directory
	os.Mkdir(config.FitDirName, 0755)

	err = os.Setenv("FIT", gdir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	//log.Print("PATH " + os.Getenv("PATH"))
	for _, tt := range setupbugargtests {
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
		found, ferr := regexp.Match(``, []byte(out)) // tt.output
		if ferr != nil {
			t.Error("Usage output: " + ferr.Error())
		} else if !found {
			t.Errorf("Unexpected usage, wanted to match %q, got %q", ``, tt.input) // tt.output
		}
	}
	// cleanup
	os.Chdir(pwd)
	err = os.RemoveAll(gdir)
	if err != nil {
		t.Error("Could not RemoveAll(" + string(gdir) + ") : " + err.Error())
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

func TestMain(m *testing.M) {
	//err := os.Chdir("..")
	goPath, _ := exec.LookPath("go")
	build := exec.Command(goPath, "build")
	err := build.Run() // removed after TestCliArgs
	if err != nil {
		fmt.Printf("go build error %s: %v", binaryname, err)
		os.Exit(1)
	}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("go pwd of %s: %v", binaryname, err)
		os.Exit(1)
	}
	binarypath = path.Join(dir, binaryname)
	//fmt.Printf("binarypath %s\n", binarypath)
	os.Exit(m.Run())
}

func TestCliArgs(t *testing.T) {
	var gdir string
	config := Config{}
	//config.DescriptionFileName = "Description"
	config.FitDirName = "fit"
	pwd, _ := os.Getwd()
	gdir, err := ioutil.TempDir("", "main")
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
	} else {
		t.Error("Failed creating temporary directory for detect")
		return
	}
	// Fake a git repo
	os.Mkdir(".git", 0755)
	// Make an issues Directory
	os.Mkdir(config.FitDirName, 0755)

	err = os.Setenv("FIT", gdir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}

	for _, tt := range setupbugargtests {
		//dir, err := os.Getwd()
		//if err != nil {
		//	t.Fatal(err)
		//}
		t.Run(tt.input, func(t *testing.T) {
			cmd := exec.Command(binarypath, tt.input)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}
			//if string(output) != tt.output {
			//	t.Errorf(" expected output: %s\n actual output %s\n", tt.output, string(output))
			//}
			found, ferr := regexp.Match(tt.output, []byte(output)) // output
			if ferr != nil {
				t.Error("Usage output: " + ferr.Error())
			} else if !found {
				t.Errorf("Unexpected usage, wanted to match %q, got %q", tt.output, tt.input)
			}
		})
	}
	// cleanup
	os.Remove(binarypath) // created in TestMain
	os.Chdir(pwd)
	err = os.RemoveAll(gdir)
	if err != nil {
		t.Error("Could not RemoveAll(" + string(gdir) + ") : " + err.Error())
	}
}
