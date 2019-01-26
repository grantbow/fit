// bug manages plain text issue files.
package main

import (
	"fmt"
	"github.com/driusan/bug/bugapp"
	"github.com/driusan/bug/bugs"
	"os"
)

func main() {

	config := bugs.Config{}
	bugs.GetIssuesDir(config) // bugs/Directory.go
	bugYml := ".bug.yml"
	bugs.ConfigRead(bugYml, &config)

	if bugapp.SkipRootCheck(&os.Args) && bugs.GetRootDir(config) == "" {
		fmt.Printf("Could not find issues directory.\n")
		fmt.Printf("Make sure either the PMIT environment variable is set, or a parent directory of your working directory has an issues folder.\n")
		fmt.Println("(If you just started new repo, you probably want to create directory named `issues`).")
		fmt.Printf("Aborting.\n")
		os.Exit(2)
	}

	// flags package is nice but this seems more direct because of
	// subcommands and arguments that are space separated names
	if len(os.Args) <= 1 {
		bugapp.Help()
	} else if len(os.Args) >= 3 && os.Args[2] == "--help" {
		bugapp.Help(os.Args[1])
	} else {
		switch os.Args[1] {
		case "--version", "version":
			bugapp.Version()
		case "dir", "pwd":
			bugapp.Pwd(config)
		case "env":
			bugapp.Env(config)
		case "purge":
			bugapp.Purge(config)
		case "add", "new", "create":
			bugapp.Create(os.Args[2:], config)
		case "commit":
			bugapp.Commit(os.Args[2:], config)
		case "edit":
			bugapp.Edit(os.Args[2:], config)
		case "find":
			bugapp.Find(os.Args[2:], config)
		case "id", "identifier":
			bugapp.Identifier(os.Args[2:], config)
		case "import":
			bugapp.Import(os.Args[2:], config)
		case "milestone":
			bugapp.Milestone(os.Args[2:], config)
		case "mv", "rename", "retitle", "relabel":
			bugapp.Relabel(os.Args[2:], config)
		case "priority":
			bugapp.Priority(os.Args[2:], config)
		case "roadmap":
			bugapp.Roadmap(os.Args[2:], config)
		case "status":
			bugapp.Status(os.Args[2:], config)
		case "rm", "close":
			bugapp.Close(os.Args[2:], config)
		case "tag":
			bugapp.Tag(os.Args[2:], config)
		case "view", "list":
			// bug list with no parameters shouldn't autopage,
			// bug list with bugs to view should. So the original
			// stdout is passed as a parameter.
			bugapp.List(os.Args[2:], config)
		case "help":
			fallthrough
		default:
			bugapp.Help(os.Args[1:]...)
		}
	}
}
