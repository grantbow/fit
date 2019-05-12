// bug manages plain text issue files.
package main

import (
	"fmt"
	"github.com/driusan/bug/bugapp"
	"github.com/driusan/bug/bugs"
	"github.com/driusan/bug/scm"
	"os"
	"strings"
)

func main() {

	config := bugs.Config{}
	config.ProgramVersion = bugapp.ProgramVersion()
	config.DescriptionFileName = "Description"

	rootPresent := false
	bugYml := ".bug.yml"
	if bugsgetrootdir := bugs.GetRootDir(config); bugsgetrootdir != "" {
		rootPresent = true
		config.BugDir = string(bugsgetrootdir)
		// now try to read config
		cerr := bugs.ConfigRead(bugYml, &config, bugapp.ProgramVersion())
		if cerr == nil {
			config.BugYml = config.BugDir + "/" + bugYml
		}
	}

	if bugapp.SkipRootCheck(&os.Args) && !rootPresent {
		//bugapp.PrintVersion()
		fmt.Printf("Error: Could not find `issues` directory. You probably want to create one.\n")
		fmt.Printf("    Make sure the current directory or a parent directory has an issues folder\n")
		fmt.Printf("    or set the PMIT environment variable.\n")
		fmt.Printf("Aborting.\n")
		os.Exit(2)
	}

	bugs.GetIssuesDir(config) // from bugs/Directory.go, uses config.BugDir

	scmoptions := make(map[string]bool)
	handler, _, herr := scm.DetectSCM(scmoptions, config)
	if herr != nil {
		if _, uerr := handler.GetSCMIssuesUpdates(); uerr != nil {
			if _, cerr := handler.GetSCMIssuesCached(); cerr != nil {
				fmt.Printf("Warn: %s\n", cerr)
			} else {
				fmt.Printf("Warn: %s\n", uerr)
			}
		} else {
			fmt.Printf("Warn: %s\n", herr)
		}
	}

	// flags package is nice but this seemed more direct at the time
	// because of subcommands and
	// arguments that are space separated names
	osArgs := os.Args // TODO: use an env var and assign to osArgs to setup for testing
	//fmt.Printf("A %s %#v\n", "osArgs: ", osArgs)
	if len(osArgs) <= 1 {
		if rootPresent {
			bugapp.List([]string{}, config)
		} else {
			fmt.Printf("Usage: " + os.Args[0] + " <command> [options]\n")
			fmt.Printf("\nUse \"bug help\" or \"bug help <command>\" for details.\n")
		}
	} else if len(osArgs) == 2 &&
		(osArgs[1] == "--help" ||
			osArgs[1] == "help" ||
			osArgs[1] == "-h") {
		bugapp.Help("help") // just bug help
	} else if len(osArgs) >= 3 &&
		(osArgs[2] == "--help" ||
			osArgs[2] == "help" ||
			osArgs[2] == "-h") { // 'bug cmd --help' just like 'bug help cmd'
		//fmt.Printf("B %s %#v\n", "osArgs: ", osArgs)
		bugapp.Help(osArgs[1])
	} else {
		switch osArgs[1] {
		// subcommands without osArgs first
		case "tagslist", "tagsassigned", "tags":
			bugapp.TagsAssigned(config)
		case "notags":
			bugapp.TagsNone(config)
		case "idslist", "idsassigned", "ids", "identifiers":
			bugapp.IdsAssigned(config)
		case "noids", "noidentifiers":
			bugapp.IdsNone(config)
		case "env":
			bugapp.Env(config)
		case "pwd", "dir", "cwd":
			bugapp.Pwd(config)
		case "version", "about", "--version", "-v":
			bugapp.PrintVersion()
		case "staging", "staged", "cached", "cache", "index":
			if b, err := handler.GetSCMIssuesUpdates(); err != nil {
				fmt.Printf("Files in issues/ need committing, see $ git status --porcelain -u -- :/issues\nand for files already in index see $ git diff --name-status --cached HEAD -- :/issues\n")
				if _, errc := handler.GetSCMIssuesCached(); errc != nil {
					for _, bline := range strings.Split(string(b), "\n") {
						//if bline in c {
						//} else {
						fmt.Printf("%v\n", string(bline))
						//}
					}
				} else {
					fmt.Printf("%v\n", string(b))
				}
			} else {
				fmt.Printf("No files in issues/ need committing, see $ git status --porcelain -u issues \":top\"\n")
			}
		case "purge":
			bugapp.Purge(config)
		// subcommands with osArgs next
		case "list", "view", "show", "display", "ls":
			// bug list with no parameters shouldn't autopage,
			// bug list with bugs to view should. So the original
			// stdout is passed as a parameter.
			bugapp.List(osArgs[2:], config)
		case "find":
			bugapp.Find(osArgs[2:], config)
		case "create", "add", "new":
			//fmt.Printf("%s %#v\n", "osArgs: ", len(osArgs))
			bugapp.Create(osArgs[2:], config)
			//fmt.Printf("%s %#v\n", "osArgs: ", len(osArgs))
		case "edit":
			bugapp.Edit(osArgs[2:], config)
		case "retitle", "mv", "rename", "relabel":
			bugapp.Relabel(osArgs[2:], config)
		case "close", "rm":
			bugapp.Close(osArgs[2:], config)
		case "tag":
			bugapp.Tag(osArgs[2:], config)
		case "id", "identifier":
			bugapp.Identifier(osArgs[2:], config)
		case "status":
			bugapp.Status(osArgs[2:], config)
		case "priority":
			bugapp.Priority(osArgs[2:], config)
		case "milestone":
			bugapp.Milestone(osArgs[2:], config)
		case "import":
			bugapp.Import(osArgs[2:], config)
		case "commit", "save":
			bugapp.Commit(osArgs[2:], config)
		case "roadmap":
			bugapp.Roadmap(osArgs[2:], config)
		case "help", "--help", "-h":
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
