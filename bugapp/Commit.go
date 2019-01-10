package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/driusan/bug/scm"
)

// Commit is a subcommand to save issues to the git or mercurial (hg) SCMs.
func Commit(args ArgumentList, config bugs.Config) {
	options := make(map[string]bool)
	if !args.HasArgument("--no-autoclose") {
		options["autoclose"] = true
	} else {
		options["autoclose"] = false
	}
	options["use_bug_prefix"] = true // SCM will ignore this option if it doesn't know it

	scm, _, err := scm.DetectSCM(options)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	err = scm.Commit(bugs.GetIssuesDir(config), "Added or removed issues with the tool \"bug\"")

	if err != nil {
		fmt.Printf("Could not commit: %s\n", err.Error())
		return
	}
}
