package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/driusan/bug/scm"
)

// Env is a subcommand to output detected editor, directory and scm type.
func Env(config bugs.Config) {
	scm, scmdir, scmerr := scm.DetectSCM(make(map[string]bool), config)
	fmt.Printf("Settings used by this command:\n")
	fmt.Printf("\nEditor: %s", getEditor())
	fmt.Printf("\nIssues Directory: %s", bugs.GetIssuesDir(config))

	if scmerr == nil {
		fmt.Printf("\n\nVCS Type:\t%s", scm.GetSCMType())
		fmt.Printf("\n%s Directory:\t%s", scm.GetSCMType(), scmdir)
	} else {
		fmt.Printf("\n\nVCS Type: None (purge and commit commands unavailable)")
	}

	fmt.Printf("\n")
}
