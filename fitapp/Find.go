package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"os"
	"sort"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// find does the work of finding bugs.
func find(findType string, findValues []string, config bugs.Config) {
	fitdir := bugs.FitDirer(config)
	//issues, _ := ioutil.ReadDir(string(fitdir))
	issues := readIssues(string(fitdir))
	sort.Sort(byDir(issues))
	for idx, issue := range issues {
		var dir bugs.Directory = fitdir + dops + bugs.Directory(issue.Name())
		b := bugs.Issue{Dir: dir}
		name := issueNamer(b, idx)
		var values []string
		switch findType {
		case "tags":
			values = b.StringTags()
		case "status":
			values = []string{b.Status()}
		case "priority":
			values = []string{b.Priority()}
		case "milestone":
			values = []string{b.Milestone()}
		default:
			fmt.Printf("Unknown find type: %s\n", findType)
			return
		}
		printed := false
		for _, findValue := range findValues {
			for _, value := range values {
				if value == findValue {
					fmt.Printf("%s: %s\n", name, b.Title(findType))
					printed = true
				}
			}
			if printed {
				break
			}
		}
	}
}

// Find is a subcommand to find issues.
func Find(args argumentList, config bugs.Config) {
	if len(args) < 2 {
		fmt.Printf("Usage: %s find {tags, status, priority, milestone} value1 [value2 ...]\n", os.Args[0])
		return
	}
	switch args[0] {
	case "tags":
		fallthrough
	case "status":
		fallthrough
	case "priority":
		fallthrough
	case "milestone":
		find(args[0], args[1:], config)
	default:
		fmt.Printf("Unknown command: %v\n", args)
		return
	}
}
