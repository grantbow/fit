// bug manages plain text issue files.
package main

import (
	"fmt"
	"github.com/driusan/bug/bugapp"
	"github.com/driusan/bug/bugs"
	"github.com/driusan/bug/scm"
	"os"
)

func main() {

	config := bugs.Config{}
	config.ProgramVersion = bugapp.ProgramVersion()
	config.DescriptionFileName = "Description"
	bugs.GetIssuesDir(config) // bugs/Directory.go
	bugYml := ".bug.yml"
	bugs.ConfigRead(bugYml, &config, bugapp.ProgramVersion())

	scmoptions := make(map[string]bool)
	handler, _, herr := scm.DetectSCM(scmoptions, config)
	if herr != nil {
		fmt.Printf("Warn: %s\n", herr.Error())
	} else if _, err := handler.GetSCMIssuesUpdates(); err != nil {
		fmt.Printf("Warn: %s\n", err)
	}

	if bugapp.SkipRootCheck(&os.Args) && bugs.GetRootDir(config) == "" {
		//bugapp.PrintVersion()
		fmt.Printf("Error: Could not find `issues` directory. You probably want to create one.\n")
		fmt.Printf("    Make sure the current directory or a parent directory has an issues folder\n")
		fmt.Printf("    or set the PMIT environment variable.\n")
		fmt.Printf("Aborting.\n")
		os.Exit(2)
	}

	// flags package is nice but this seemed more direct at the time
	// because of subcommands and
	// arguments that are space separated names
	osArgs := os.Args // TODO: use an env var and assign to osArgs to setup for testing
	//fmt.Printf("A %s %#v\n", "osArgs: ", osArgs)
	if len(osArgs) <= 1 {
		fmt.Printf("Usage: " + os.Args[0] + " <command> [options]\n")
		fmt.Printf("\nUse \"bug help\" or \"bug help <command>\" for details.\n")
	} else if len(osArgs) >= 3 && (osArgs[2] == "--help" || osArgs[2] == "help") { // bug cmd --help just like bug help cmd
		//fmt.Printf("B %s %#v\n", "osArgs: ", osArgs)
		bugapp.Help(osArgs[1])
	} else {
		switch osArgs[1] {
		case "--version", "version", "-v": // subcommands without osArgs
			bugapp.PrintVersion()
		case "dir", "pwd":
			bugapp.Pwd(config)
		case "env":
			bugapp.Env(config)
		case "purge":
			bugapp.Purge(config)
		case "tagsassigned":
			bugapp.TagsAssigned(config)
		case "tagsnone":
			bugapp.TagsNone(config)
		case "staging":
			if b, err := handler.GetSCMIssuesUpdates(); err != nil {
				fmt.Printf("Files in issues/ need committing, see $ git status --porcelain -u issues \":top\"\n%v\n", string(b))
			} else {
				fmt.Printf("No files in issues/ need committing, see $ git status --porcelain -u issues \":top\"\n")
			}
		case "add", "new", "create": // subcommands with    osArgs
			//fmt.Printf("%s %#v\n", "osArgs: ", len(osArgs))
			bugapp.Create(osArgs[2:], config)
			//fmt.Printf("%s %#v\n", "osArgs: ", len(osArgs))
		case "commit", "save":
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
		case "view", "list", "show", "display":
			// bug list with no parameters shouldn't autopage,
			// bug list with bugs to view should. So the original
			// stdout is passed as a parameter.
			bugapp.List(osArgs[2:], config)
		case "help", "--help":
			//fmt.Printf("C %s %#v\n", "osArgs: ", osArgs)
			bugapp.Help(osArgs[2])
		default:
			//if
			if len(osArgs) == 2 {
				buglist, _ := bugs.LoadBugByHeuristic(osArgs[1], config)
				//fmt.Printf("%+v\n", buglist)
				if buglist != nil { // || ae, ok := bugerr.(bugs.ErrNotFound); ! ok { // bug list when possible, not help
					bugapp.List(osArgs[1:], config)
				} else {
					bugapp.Help(osArgs[1:]...)
				}
			}
		}
	}
}
