package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

// getAllTags reads the tags subdir only
func getAllTags(config bugs.Config) []string {
	bugs := bugs.GetAllBugs(config)
	//fmt.Printf("%+v\n", bugs)

	// Put all the tags in a map, then iterate over
	// the keys so that only unique tags are included
	tagMap := make(map[string]int, 0)
	for _, bug := range bugs {
		for _, tag := range bug.Tags() {
			tagMap[string(tag)] += 1
		}
	}

	keys := make([]string, 0, len(tagMap))
	for k := range tagMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// TagsNone is a subcommand to return print ready issues with no assigned tags.
func TagsNone(config bugs.Config) {
	issuesroot := bugs.GetIssuesDir(config)
	issues, _ := ioutil.ReadDir(string(issuesroot)) // TODO: should be a method elsewhere
	var wantTags bool = false

	allbugs := bugs.GetAllBugs(config)
	tagMap := make(map[string]int, 0)
	for _, bug := range allbugs {
		if len(bug.Tags()) == 0 {
			title := bug.Dir.GetShortName()
			tagMap[string(title)] += 1
		}
	}

	//keys := make([]string, 0, len(tagMap))
	/*for k, _ := range tagMap {
		//fmt.Printf("%v\n", k)
		name := getBugName(b, idx) // Issue x:
		fmt.Printf("%v\n", k)
		//keys = append(keys, k) // TODO: should just append not tagmap intermediary
	} */

	fmt.Printf("No tags assigned:\n")
	//fmt.Printf("%v\n", len(issues))
	for idx, issue := range issues {
		//fmt.Printf("%v\n", issue)
		for k, _ := range tagMap {
			if issue.Name() == k {
				//fmt.Printf("1in: %v\n2tm: %v\n", issue.Name(), k)
				var dir bugs.Directory = issuesroot + bugs.Directory(issue.Name()) //issuesroot + issue.Dir
				//fmt.Printf("dir %v\n", dir)
				b := bugs.Bug{Dir: dir, DescriptionFileName: config.DescriptionFileName}
				name := getBugName(b, idx) // Issue x:
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

// TagsAssigned is a subcommand to print the assigned tags.
func TagsAssigned(config bugs.Config) {
	//fmt.Printf("here\n")
	get := getAllTags(config)
	fmt.Printf("Tags used in current tree: ")
	if len(get) > 0 {
		fmt.Printf("%s\n", strings.Join(get, ", "))
	} else {
		fmt.Print("<none assigned yet>\n")
	}
}

// Tag is a subcommand to assign a tag to an issue.
func Tag(Args argumentList, config bugs.Config) {
	if len(Args) < 2 {
		fmt.Printf("Usage: %s tag [--rm] BugID tagname [more tagnames]\n", os.Args[0])
		fmt.Printf("\nBoth issue number and tagname to set are required.\n")
		fmt.Printf("\nCurrently used tags in entire tree: %s\n", strings.Join(getAllTags(config), ", "))
		return
	}
	var removeTags bool = false
	if Args[0] == "--rm" {
		removeTags = true
		Args = Args[1:]
	}

	b, err := bugs.LoadBugByHeuristic(Args[0], config)

	if err != nil {
		fmt.Printf("Could not load bug: %s\n", err.Error())
		return
	}
	for _, tag := range Args[1:] {
		if removeTags {
			b.RemoveTag(bugs.TagBoolTrue(tag), config)
		} else {
			b.TagBug(bugs.TagBoolTrue(tag), config)
		}
	}

}
