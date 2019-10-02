package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os"
	"strings"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// Import is a subcommand to read from a bugsEverywhere.org or github.com systems and create identical issues.
func Import(args argumentList, config bugs.Config) {
	if len(args) < 1 {
		fmt.Printf("Usage: %s import {--github,--be} <repo>\n", os.Args[0])
		//fmt.Printf("Usage: %s import {<github.com/user/repo>,--be}\n", os.Args[0])
		return
	}
	switch args[0] {
	case "--github":
		if githubRepo := args.GetArgument("--github", ""); githubRepo != "" {
			numStrings := strings.Count(githubRepo, sops)
			pieces := strings.Split(githubRepo, sops)
			//fmt.Printf("ns %v\np %v\n", numStrings, pieces)
			if numStrings == 1 {
				githubImportIssues(pieces[0], pieces[1], config)
			} else if numStrings == 2 &&
				pieces[2] == "projects" &&
				config.GithubPersonalAccessToken != "" {
				githubImportProjects(pieces[0], pieces[1], config)
			} else {
				fmt.Fprintf(os.Stderr, "GitHub invalid: %s\n", githubRepo)
				return
			}
		}
	case "--be":
		if len(args) > 1 {
			fmt.Fprintf(os.Stderr, "BugsEverywhere repo ignored: %s\n", args[1:])
		}
		beImport(config)
	default:
		fmt.Fprintf(os.Stderr, "Usage: %s import --github user/repo\n", os.Args[0])
		//fmt.Fprintf(os.Stderr, "Usage: %s import github.com/user/repo\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s import --be\n", os.Args[0])
		fmt.Fprintf(os.Stderr, `
Use this command to import an external bug database into the local
issues/ directory.

Either "--github user/repo" is required to import GitHub issues
or "--be" is required to import a local BugsEverywhere database.
`)
	}
	return
}
