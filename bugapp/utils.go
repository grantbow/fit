package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func check(e error) {
	if e != nil {
		//	fmt.Fprintf(os.Stderr, "err: %s", err.Error())
		//	return NoConfigError
		panic(e)
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

// fieldHandler is used for Priority, Milestone and Status
func fieldHandler(command string, args argumentList,
	setCallback func(bugs.Bug, string) error, retrieveCallback func(bugs.Bug) string, config bugs.Config) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s %s BugID [set %s]\n", os.Args[0], command, command)
		return
	}

	b, err := bugs.LoadBugByHeuristic(args[0], config)
	if err != nil {
		fmt.Printf("Invalid BugID: %s\n", err.Error())
		return
	}
	if len(args) > 1 {
		newValue := strings.Join(args[1:], " ")
		err := setCallback(*b, newValue)
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

// dirDump accepts a string directory and returns details.
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

// dirDumpFI accepts an array of os.FileInfo and returns details.
func dirDumpFI(files []os.FileInfo) string {
	a := []string{}
	for _, file := range files {
		a = append(a, fmt.Sprintf("%v", file))
	}
	return strings.Join(a, ",\n")
}

// SkipRootCheck is a helper function to avoid unnecessary filesystem checking..
func SkipRootCheck(args *[]string) bool {
	ret := false
	switch len(*args) {
	case 0, 1:
		ret = true
	case 2:
		if (*args)[1] == "help" {
			ret = true
		}
	case 3:
		if (*args)[2] == "--help" {
			ret = true
		}
	}
	return ret
}
