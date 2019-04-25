package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/driusan/bug/scm"
)

// Env is a subcommand to output detected editor, directory and scm type.
func Env(config bugs.Config) {
	vcs, scmdir, scmerr := scm.DetectSCM(make(map[string]bool), config)
	fmt.Printf("Settings used by this command:\n")
	fmt.Printf("\nEditor: %s", getEditor())
	fmt.Printf("\nIssues Directory: %s", bugs.GetIssuesDir(config))

	if scmerr == nil {
		t := vcs.GetSCMType()
		fmt.Printf("\n\nVCS Type:\t%s", t)
		fmt.Printf("\n%s Directory:\t%s", t, scmdir)
		fmt.Printf("\nNeed Staging:\t")
		if err := vcs.GetSCMIssueUpdates(); err == nil {
			fmt.Print("(No files in issues)")
		} else {
			fmt.Print("Issues")
		}
	} else {
		fmt.Printf("\n\nVCS Type: None (purge and commit commands unavailable)")
	}

	fmt.Printf("\n")
}
