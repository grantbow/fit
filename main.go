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

	// flags package is nice but this seemed more direct at the time
	// because of subcommands and
	// arguments that are space separated names
	osArgs := os.Args // TODO: use an env var and assign to osArgs to setup for testing
	if len(osArgs) <= 1 {
		bugapp.Help()
	} else if len(osArgs) >= 3 && osArgs[2] == "--help" { // bug cmd --help just like bug help cmd
		bugapp.Help(osArgs[1])
	} else {
		switch osArgs[1] {
		case "--version", "version": // subcommands without osArgs
			bugapp.Version()
		case "dir", "pwd":
			bugapp.Pwd(config)
		case "env":
			bugapp.Env(config)
		case "purge":
			bugapp.Purge(config)
		case "add", "new", "create": // subcommands with    osArgs
			bugapp.Create(osArgs[2:], config)
		case "commit":
			bugapp.Commit(osArgs[2:], config)
		case "edit":
			bugapp.Edit(osArgs[2:], config)
		case "find":
			bugapp.Find(osArgs[2:], config)
		case "id", "identifier":
			bugapp.Identifier(osArgs[2:], config)
		case "import":
			bugapp.Import(osArgs[2:], config)
		case "milestone":
			bugapp.Milestone(osArgs[2:], config)
		case "mv", "rename", "retitle", "relabel":
			bugapp.Relabel(osArgs[2:], config)
		case "priority":
			bugapp.Priority(osArgs[2:], config)
		case "roadmap":
			bugapp.Roadmap(osArgs[2:], config)
		case "status":
			bugapp.Status(osArgs[2:], config)
		case "rm", "close":
			bugapp.Close(osArgs[2:], config)
		case "tag":
			bugapp.Tag(osArgs[2:], config)
		case "view", "list":
			// bug list with no parameters shouldn't autopage,
			// bug list with bugs to view should. So the original
			// stdout is passed as a parameter.
			bugapp.List(osArgs[2:], config)
		case "help":
			fallthrough
		default:
			bugapp.Help(osArgs[1:]...)
		}
	}
}
