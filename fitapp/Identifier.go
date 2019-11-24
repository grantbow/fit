package fitapp

import (
	"crypto/sha1"
	"fmt"
	bugs "github.com/driusan/bug/bugs"
	"os"
	"sort"
	"strings"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// getAllIds returns all the ids
func getAllIds(config bugs.Config) []string {
	bugs := bugs.GetAllIssues(config)
	idMap := make(map[string]int, 0)
	for _, bug := range bugs {
		idMap[strings.ToLower(string(bug.Identifier()))] += 1
	}
	var ids []string
	for k := range idMap {
		ids = append(ids, k)
	}
	sort.Strings(ids)
	return ids
}

// IdsNone is a subcommand to print issues with no assigned tags.
func IdsNone(config bugs.Config) {
	fitdir := bugs.FitDirer(config)
	issues := readIssues(string(fitdir))
	sort.Sort(byDir(issues))
	var wantTags bool = false

	allbugs := bugs.GetAllIssues(config)
	idMap := make(map[string]int, 0)
	for _, bug := range allbugs {
		if bug.Identifier() == "" {
			title := bug.Dir.ShortNamer()
			idMap[string(title)] += 1
		}
	}
	fmt.Printf("No ids assigned:\n")
	for idx, issue := range issues {
		//fmt.Printf("%v\n", issue)
		for k, _ := range idMap {
			if issue.Name() == k {
				//fmt.Printf("1in: %v\n2tm: %v\n", issue.Name(), k)
				var dir bugs.Directory = fitdir + dops + bugs.Directory(issue.Name())
				//fmt.Printf("dir %v\n", dir)
				b := bugs.Issue{Dir: dir, DescriptionFileName: config.DescriptionFileName}
				name := issueNamer(b, idx) // Issue x:
				//fmt.Printf("name %v\n", name)
				if wantTags == false { // always
					fmt.Printf("%s: %s\n", name, b.Title(""))
					//keys = append(keys, fmt.Sprintf("%s: %s\n", name, b.Title("")))
				} else {
					fmt.Printf("%s: %s\n", name, b.Title("tags"))
					//keys = append(keys, fmt.Sprintf("%s: %s\n", name, b.Title("tags")))
				}
			}
		}
	}
	//return keys
	return
}

// IdsAssigned is a subcommand to print the assigned ids.
func IdsAssigned(config bugs.Config) {
	//fmt.Printf("here\n")
	get := getAllIds(config)
	fmt.Printf("Ids used in current tree: ")
	if len(get) > 0 {
		fmt.Printf("%s\n", strings.Join(get, ", "))
	} else {
		fmt.Print("<none assigned yet>\n")
	}
}

func generateID(val string) string {
	hash := sha1.Sum([]byte(val))
	return fmt.Sprintf("b%x", hash)[0:5]
}

// Identifier is a subcommand to assign tags to issues.
func Identifier(args argumentList, config bugs.Config) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s id <IssueID> [value]\n", os.Args[0])
		return
	}

	b, err := bugs.LoadIssueByHeuristic(args[0], config)
	if err != nil {
		fmt.Printf("Invalid IssueID: %s\n", err.Error())
		return
	}
	if len(args) > 1 {
		var newValue string
		if args.HasArgument("--generate-id") {
			newValue = generateID(b.Title(""))
			fmt.Printf("Generated id %s for issue\n", newValue)
		} else {
			newValue = strings.Join(args[1:], " ")
		}
		err := b.SetIdentifier(newValue, config)
		if err != nil {
			fmt.Printf("Error setting id: %s", err.Error())
		}
	} else {
		val := b.Identifier()
		if val == "" {
			fmt.Printf("Id not defined\n")
		} else {
			fmt.Printf("%s\n", val)
		}
	}
}
