package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"github.com/grantbow/fit/scm"
	"strings"
)

// Env is a subcommand to output detected editor, directory and scm type.
func Env(config bugs.Config) {
	vcs, scmdir, scmerr := scm.DetectSCM(make(map[string]bool), config)
	settingsOut := ""
	if config.FitYml == "" {
		settingsOut = "<missing>"
	} else {
		settingsOut = config.FitYml
	}
	fmt.Printf("Settings:\n\nEditor: %s\nRoot Directory: %s\nFit Directory: %s\nSettings file: %s\n\n",
		getEditor(), config.FitDir, bugs.FitDirer(config), settingsOut)

	if scmerr != nil {
		fmt.Printf("VCS Type: <missing> (purge and commit commands unavailable)\n\n")
	} else {
		t := vcs.SCMTyper()
		fmt.Printf("VCS Type:    %s\n", t)
		fmt.Printf("%s Directory:    %s\n", t, scmdir)
		//
		fmt.Printf("Need Committing or Staging:    ")
		if b, err := vcs.SCMIssuesUpdaters(config); err == nil {
			fmt.Printf("(nothing)\n\n")
		} else {
			fmt.Printf("%v\n\n", string(b)) // simplest implementation, doesn't clarify
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
