package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"sort"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// find does the work of finding bugs.
func find(findType string, findValues []string, config bugs.Config) {
	issuesroot := bugs.IssuesDirer(config)
	//issues, _ := ioutil.ReadDir(string(issuesroot))
	issues := readIssues(string(issuesroot))
	sort.Sort(byDir(issues))
	for idx, issue := range issues {
		var dir bugs.Directory = issuesroot + dops + bugs.Directory(issue.Name())
		b := bugs.Bug{Dir: dir}
		name := bugNamer(b, idx)
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
