package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"strings"
)

// Relabel is a subcommand to change an issue title.
func Relabel(Args argumentList, config bugs.Config) {
	if len(Args) < 2 {
		fmt.Printf("Usage: %s relabel BugID New Title\n", os.Args[0])
		return
	}

	b, err := bugs.LoadBugByHeuristic(Args[0], config)

	if err != nil {
		fmt.Printf("Could not load bug: %s\n", err.Error())
		return
	}

	currentDir := b.GetDirectory()
	newDir := bugs.GetIssuesDir(config) + bugs.TitleToDir(strings.Join(Args[1:], " "))
	fmt.Printf("Moving %s to %s\n", currentDir, newDir)
	err = os.Rename(string(currentDir), string(newDir))
	if err != nil {
		fmt.Printf("Error moving directory\n")
	}
}
