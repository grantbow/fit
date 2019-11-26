package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"log"
	"os"
	"os/exec"
	"strings"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// Edit is a subcommand to modify an issue.
func Edit(args argumentList, config bugs.Config) {

	var file, bugID string
	switch len(args) {
	case 1:
		// If there's only 1 argument, it's an issue
		// identifier and it's editing the Description.
		// So set the variables and fallthrough to the
		// 2 argument (editing a specific fieldname)
		// case
		bugID = args[0]
		file = config.DescriptionFileName
		fallthrough
	case 2:
		// If there's exactly 2 arguments, idx and
		// file haven't been set by the first case
		// statement, so set them, but everything else
		// is the same
		if len(args) == 2 {
			bugID = args[0]
			file = args[1]
		}

		b, err := bugs.LoadIssueByHeuristic(bugID, config)
		if err != nil {
			fmt.Printf("Invalid IssueID %s\n", bugID)
			return
		}

		dir := b.Direr()

		switch title := strings.Title(file); title {
		case "Description", "Milestone", "Status", "Priority", "Identifier":
			// enforces Title case
			file = title
		} // else falls through
		fmt.Printf("Editing %s%s%s\n", dir, sops, file)
		cmd := exec.Command(getEditor(), string(dir)+sops+file)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Printf("Usage: %s edit [fieldname] IssueID\n", os.Args[0])
		fmt.Printf("\nNo IssueID specified\n")
	}
}
