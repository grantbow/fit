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
	config.IssuesDirName = "issues"

	rootPresent := false
	bugYmlFileName := ".bug.yml"
	fitYmlFileName := ".fit.yml"
	skip := bugapp.SkipRootCheck(&os.Args) // too few args or help or env
	if rd := bugs.RootDirer(&config); rd != "" {
		// bugs/Directory.go func RootDirer sets config.BugDir does os.Chdir()
		rootPresent = true
		// now try to read config
		if ErrC := bugs.ConfigRead(fitYmlFileName, &config, bugapp.ProgramVersion()); ErrC == nil {
			config.BugYml = config.BugDir + string(os.PathSeparator) + fitYmlFileName
			//var sops = string(os.PathSeparator) not yet available
			//var dops = Directory(os.PathSeparator)
		} else if ErrC := bugs.ConfigRead(bugYmlFileName, &config, bugapp.ProgramVersion()); ErrC == nil {
			config.BugYml = config.BugDir + string(os.PathSeparator) + bugYmlFileName
		}
	}

	if !rootPresent {
		if skip {
			fmt.Printf("Warn: no `issues` directory\n")
		} else { // !skip
			//bugapp.PrintVersion()
			fmt.Printf("bug manages plain text issues with git or hg.\n")
			fmt.Printf("Error: Could not find `fit` or `issues` directory.\n")
			fmt.Printf("    Check that the current or a parent directory has a fit directory\n")
			fmt.Printf("    or set the FIT environment variable.\n")
			//fmt.Printf("Each issues directory contains issues.\n")
			//fmt.Printf("Each issue directory contains a text file named Description.\n")
			fmt.Printf("Use \"bug help\" or \"bug help help\" for details.\n")
			fmt.Printf("Aborting.\n")
			os.Exit(2)
		}
	}

	bugs.IssuesDirer(config) // from bugs/Directory.go, uses config.BugDir from bugs/Bug.go

	scmoptions := make(map[string]bool)
	handler, _, ErrH := scm.DetectSCM(scmoptions, config)
	//a, b, c := scm.DetectSCM(scmoptions, config)
	//fmt.Printf("%+v %+v %+v\n", a, b, c)
	if ErrH != nil {
		fmt.Printf("Warn: no .git or .hg directory. Use \"{git|hg} init\".\n")
		//fmt.Printf("Warn: %s\n", ErrH) // No SCM found
		//a, b := handler.SCMIssuesUpdaters()
		//fmt.Printf("%+v %+v\n", a, b)
		if handler != nil {
			if _, ErrU := handler.SCMIssuesUpdaters(); ErrU != nil {
				if _, ErrCa := handler.SCMIssuesCacher(); ErrCa != nil {
					fmt.Printf("Warn: %s\n", ErrCa)
				} else {
					fmt.Printf("Warn: %s\n", ErrU)
				}
			}
		}
	}

	// flags package is nice but this seemed more direct at the time
	// because of subcommands and
	// arguments that can be space separated names
	osArgs := os.Args // TODO: use an env var and assign to osArgs to setup for testing
	//fmt.Printf("A %s %#v\n", "osArgs: ", osArgs)
	if len(osArgs) <= 1 {
		if rootPresent {
			bugapp.List([]string{}, config, true)
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
		case "notags", "notag":
			bugapp.TagsNone(config)
		case "idslist", "idsassigned", "ids", "identifiers":
			bugapp.IdsAssigned(config)
		case "noids", "noid", "noidentifiers", "noidentifier":
			bugapp.IdsNone(config)
		case "env":
			bugapp.Env(config)
		case "pwd", "dir", "cwd":
			bugapp.Pwd(config)
		case "version", "about", "--version", "-v":
			bugapp.PrintVersion()
		case "purge":
			bugapp.Purge(config)
		case "twilio":
			bugapp.Twilio(config)
		case "staging", "staged", "cached", "cache", "index":
			if b, err := handler.SCMIssuesUpdaters(); err != nil {
				fmt.Printf("Files in issues/ need committing, see $ git status --porcelain -u -- :/issues\nand for files already in index see $ git diff --name-status --cached HEAD -- :/issues\n")
				if _, ErrCach := handler.SCMIssuesCacher(); ErrCach != nil {
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
		// subcommands that pass osArgs
		case "tagslist", "taglist", "tagsassigned", "tags":
			bugapp.TagsAssigned(osArgs[2:], config)
		case "list", "view", "show", "display", "ls":
			// bug list with no parameters shouldn't autopage,
			// bug list with bugs to view should. So the original
			bugapp.List(osArgs[2:], config, true)
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
			bugapp.Tag(osArgs[2:], config) // boolean only
		case "id", "identifier":
			bugapp.Identifier(osArgs[2:], config)
		case "status":
			if len(osArgs) == 2 {
				// overview like a git status
				bugapp.Env(config)
			} else {
				// get or set the status of an issue
				bugapp.Status(osArgs[2:], config)
			}
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
					bugapp.List(osArgs[1:], config, true)
				} else {
					bugapp.Help(osArgs[1:]...)
				}
			}
		}
	}
}
