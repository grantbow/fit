package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"strings"
)

func getAllTags(config bugs.Config) []string {
	bugs := bugs.GetAllBugs(config)

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
	return keys
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
			b.RemoveTag(bugs.Tag(tag))
		} else {
			b.TagBug(bugs.Tag(tag))
		}
	}

}
