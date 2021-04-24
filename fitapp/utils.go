package fitapp

import (
	_ "flag"
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var dops = bugs.Directory(os.PathSeparator)
var sops = string(os.PathSeparator)

type argumentList []string

// HasArgument checks pkg global argumentList for an argument parameter. Returns true or false.
func (args argumentList) HasArgument(arg string) bool {
	for _, argCandidate := range args {
		if arg == argCandidate {
			return true
		}
	}
	return false
}

// GetArgument gets an argument from the pkg global argumentList. Returns a string.
func (args argumentList) GetArgument(argname, defaultVal string) string {
	for idx, argCandidate := range args {
		if argname == argCandidate {
			// If it's the last argument, then return string
			// "true" because we can't return idx+1, but it
			// shouldn't be the default value when the argument
			// isn't provided either..
			if idx >= len(args)-1 {
				return "true"
			}
			return args[idx+1]
		}
	}
	return defaultVal
}

// GetAndRemoveArguments returns an argumentList and corresponding values as an argumentList and a slice of strings.
func (args argumentList) GetAndRemoveArguments(argnames []string) (argumentList, []string) {
	var nextArgumentType int = -1
	matches := make([]string, len(argnames))
	var retArgs []string
	for idx, argCandidate := range args {
		// The last token was in argnames, so this one
		// is the value. Set it in matches and reset
		// nextArgumentType and continue to the next
		// possible token
		if nextArgumentType != -1 {
			matches[nextArgumentType] = argCandidate
			nextArgumentType = -1
			continue
		}

		// Check if this is a argname we're looking for
		for argidx, argname := range argnames {
			if argname == argCandidate {
				if idx >= len(args)-1 {
					matches[argidx] = "true"
				}
				nextArgumentType = argidx
				break
			}
		}

		// It wasn't an argname, so add it to the return
		if nextArgumentType == -1 {
			retArgs = append(retArgs, argCandidate)
		}
	}
	return retArgs, matches
}

// check will panic with an error
func check(e error) {
	if e != nil {
		//	fmt.Fprintf(os.Stderr, "err: %s", err.Error())
		//	return NoConfigError
		panic(e)
	}
}

// captureOutput accepts a function for testing and returns stdout string and stderr string.
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

// fieldHandler is used for Priority, Milestone and Status, not Identifier
func fieldHandler(command string, args argumentList,
	setCallback func(bugs.Issue, string, bugs.Config) error, retrieveCallback func(bugs.Issue) string, config bugs.Config) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s %s IssueID [set %s]\n", os.Args[0], command, command)
		return
	}

	b, err := bugs.LoadIssueByHeuristic(args[0], config)
	if err != nil {
		fmt.Printf("Invalid IssueID: %s\n", err.Error())
		return
	}
	if len(args) > 1 {
		newValue := strings.Join(args[1:], " ")
		err := setCallback(*b, newValue, config)
		if err != nil {
			fmt.Printf("Error setting %s: %s", command, err.Error())
		}
	} else {
		val := retrieveCallback(*b)
		if val == "" {
			fmt.Printf("%s not defined\n", command)
		} else {
			fmt.Printf("%s\n", val)
		}
	}
}

// dirDump accepts a string directory and returns a string.
func dirDump(dir string) string {
	a := []string{}
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			a = append(a, fmt.Sprintf("%v", path))
			return nil
		})
	if err != nil {
		fmt.Printf("dirDump error %s", err.Error())
	}
	//for _, file := range files {
	//	a = append(a, fmt.Sprintf("%v", file))
	//}
	return strings.Join(a, ",\n")
}

// dirDumpFI accepts an array of os.FileInfo and returns a string
func dirDumpFI(files []os.FileInfo) string {
	a := []string{}
	for _, file := range files {
		a = append(a, fmt.Sprintf("%v", file))
	}
	return strings.Join(a, ",\n")
}

// SkipRootCheck is a helper function to avoid unnecessary filesystem checking.
func SkipRootCheck(args *[]string) bool {
	ret := false
	switch len(*args) {
	case 0, 1:
		ret = true
	case 2:
		if (*args)[1] == "help" ||
			(*args)[1] == "--help" ||
			(*args)[1] == "-h" ||
			(*args)[1] == "status" ||
			(*args)[1] == "version" ||
			(*args)[1] == "about" ||
			(*args)[1] == "--version" ||
			(*args)[1] == "-v" ||
			(*args)[1] == "environment" ||
			(*args)[1] == "config" ||
			(*args)[1] == "settings" ||
			(*args)[1] == "env" {
			ret = true
		}
	case 3:
		if (*args)[2] == "--help" ||
			(*args)[1] == "help" {
			ret = true
		}
	}
	return ret
}

// also in issues/utils.go
func removeFi(slice []os.FileInfo, i int) []os.FileInfo {
	//flag.Parse()
	//Debug("debug len " + string(len(slice)) + " i " + string(i) + " slice[0].Name() " + string(slice[0].Name()) + "\n")
	//fmt.Printf("%+v\n", flag.Args()) // didn't seem to help, needs more work to make it active
	//
	//fmt.Printf("debug ok 01 \n")
	//fmt.Printf("debug len " + string(len(slice)) + " i " + string(i) + " slice[0].Name() " + string(slice[0].Name()) + "\n")
	//fmt.Printf("debug removeFi args len " + string(len(slice)) + " i " + string(i) + "\n")
	if (len(slice) == 1) && (i == 0) {
		return []os.FileInfo{}
	} else if i < len(slice)-2 {
		copy(slice[i:], slice[i+1:])
	}
	return slice[:len(slice)-1]
}

// also in issues/utils.go
func readIssues(dirname string) []os.FileInfo {
	//var issueList []os.FileInfo
	fis, _ := ioutil.ReadDir(string(dirname))
	issueList := fis
	for idx, fi := range issueList {
		//Debug("debug fi " + string(fi.Name()) + "idx " + string(idx) + "\n")
		//fmt.Printf("debug readIssues loop fi " + string(fi.Name()) + "idx " + string(idx) + "\n")
		if fi.IsDir() != true {
			//fmt.Printf("debug before removeFi name " + fi.Name() + " idx " + string(idx) + "\n")
			issueList = removeFi(issueList, idx)
		}
	}
	return issueList
}

//byDir allows sort.Sort(byDir(issues))
// type and three functions are needed - also see Issue.go for type byIssue
// rather than a custom Len function for os.FileInfo, Len is calculated in Less
type byDir []os.FileInfo

func (t byDir) Len() int {
	return len(t)
}
func (t byDir) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t byDir) Less(i, j int) bool {
	return (t[i]).ModTime().Unix() < (t[j]).ModTime().Unix()
}
