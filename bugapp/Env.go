package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/driusan/bug/scm"
	"strings"
)

// Env is a subcommand to output detected editor, directory and scm type.
func Env(config bugs.Config) {
	vcs, scmdir, scmerr := scm.DetectSCM(make(map[string]bool), config)
	fmt.Printf("Settings:\n\nEditor: %s\nRoot Directory: %s\nIssues Directory: %s\nSettings file: %s\n\n",
		getEditor(), config.BugDir, bugs.IssuesDirer(config), config.BugYml)

	if scmerr != nil {
		fmt.Printf("VCS Type: <missing> (purge and commit commands unavailable)\n\n")
	} else {
		t := vcs.SCMTyper()
		fmt.Printf("VCS Type:    %s\n", t)
		fmt.Printf("%s Directory:    %s\n", t, scmdir)
		//
		fmt.Printf("Need Committing or Staging:    ")
		if b, err := vcs.SCMIssuesUpdaters(); err == nil {
			fmt.Printf("(nothing)\n\n")
		} else {
			fmt.Printf("%v\n\n", string(b))
		}
	}
	fmt.Printf("Config:\n    " +
		strings.Replace(
			strings.TrimLeft(
				strings.Replace(
					fmt.Sprintf("%#v\n", config),
					", ", "\n    ", -1), // Replace
				"bugs.Config"), // TrimLeft
			":", " : ", -1), // Replace
	) // Printf
	fmt.Printf("\n")
}
