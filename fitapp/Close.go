package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"os"
)

// Close is a subcommand to close issues.
func Close(args argumentList, config bugs.Config) {
	// No parameters, print a list of all bugs
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s close <IssueID>\n\nMust provide an ID to close as parameter\n", os.Args[0])
		return
	}

	// There were parameters, so show the full description of each
	// of those issues
	var bugsToClose []string
	for _, bugID := range args {
		if bug, err := bugs.LoadIssueByHeuristic(bugID, config); err == nil {
			dir := bug.Direr()
			bugsToClose = append(bugsToClose, string(dir))
			if config.CloseStatusTag {
				fmt.Printf("Tag status closed %s\n", dir)
				err = bug.SetField("Status", "closed", config)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error setting %s %s : %s\n", "Status", "closed", err.Error())
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Could not close issue %s: %s\n", bugID, err.Error())
		}
	}
	for _, dir := range bugsToClose {
		if config.CloseMove {
			fmt.Printf("Moving %s to %s\n", dir, config.ClosedDirName)
			err := os.Rename(dir, config.ClosedDirName + sops + dir )
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error renaming %s : %s\n", dir, err.Error())
			}
		} else if !config.ClosePreventDelete {
			fmt.Printf("Removing %s\n", dir)
			err := os.RemoveAll(dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error removing %s : %s\n", dir, err.Error())
			}
		}
	}
}
