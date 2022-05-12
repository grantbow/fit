// bug manages plain text issue files.
package main

import (
	"fmt"
	fitapp "github.com/grantbow/fit/fitapp"
	issues "github.com/grantbow/fit/issues"
	"github.com/grantbow/fit/scm"
	"os"
	"strings"
)

func main() {

	config := issues.Config{}
	config.ProgramVersion = fitapp.ProgramVersion()
	config.DescriptionFileName = "Description"
	config.FitDirName = "fit"
	rootPresent := false
	skip := fitapp.SkipRootCheck(&os.Args) // too few args or help or env
	// see fitapp/utils.go

	// detect scm first to determine backup location for .fit.yml
	scmoptions := make(map[string]bool)
	handler, detectedDir, ErrH := scm.DetectSCM(scmoptions, config)
	//a, b, c := scm.DetectSCM(scmoptions, config)
	//fmt.Printf("%+v %+v %+v\n", a, b, c)
	if ErrH != nil {
		fmt.Printf("Warn: no .git or .hg directory. Use \"{git|hg} init\".\n")
		//fmt.Printf("Warn: %s\n", ErrH) // No SCM found
		//a, b := handler.SCMIssuesUpdaters()
		//fmt.Printf("%+v %+v\n", a, b)
		if handler != nil {
			if _, ErrU := handler.SCMIssuesUpdaters(config); ErrU != nil {
				if _, ErrCa := handler.SCMIssuesCacher(config); ErrCa != nil {
					fmt.Printf("Warn: %s\n", ErrCa)
				} else {
					fmt.Printf("Warn: %s\n", ErrU)
				}
			}
		}
	} else {
		config.ScmDir = string(detectedDir)
		config.ScmType = handler.SCMTyper()
	}

	if rd := issues.RootDirer(&config); rd != "" {
		// issues/Directory.go func RootDirer sets config.FitDir runs os.Chdir()
		rootPresent = true
		config.FitYmlDir = config.FitDir // default
		fitYmlFileName := ".fit.yml"
		bugYmlFileName := ".bug.yml"
		if ErrC := issues.ConfigRead(fitYmlFileName, &config, fitapp.ProgramVersion()); ErrC == nil {
			// tried to read FitDir fit config, must try both fit and bug
			config.FitYmlDir = config.FitDir
			config.FitYml = config.FitYmlDir + string(os.PathSeparator) + fitYmlFileName
			//var sops = string(os.PathSeparator) not yet available
			//var dops = Directory(os.PathSeparator)
		} else if ErrC := issues.ConfigRead(bugYmlFileName, &config, fitapp.ProgramVersion()); ErrC == nil {
			// tried to read FitDir bug config, must try both fit and bug
			config.FitYmlDir = config.FitDir
			config.FitYml = config.FitDir + string(os.PathSeparator) + bugYmlFileName
		} else if config.ScmType == "git" {
			// TODO: collapse else if ...git... && else if ...hg... with .git(4 char len, stl ScmTypeLen) and .hg(3 char len)
			s := len(config.ScmDir)
			//fmt.Printf("debug s01 %v\n scmdir %v\n dir %v\n", s, config.ScmDir, string(config.ScmDir[:s-4])+fitYmlFileName)
			if ErrC := issues.ConfigRead(string(config.ScmDir[:s-4])+fitYmlFileName, &config, fitapp.ProgramVersion()); ErrC == nil {
				// tried to read .git fit config
				config.FitYmlDir = string(config.ScmDir[:s-5])
				config.FitYml = string(config.ScmDir[:s-4]) + fitYmlFileName
				//fmt.Printf("debug 02\n %v\n", config.FitYml)
			} else if ErrC := issues.ConfigRead(string(config.ScmDir[:s-4])+bugYmlFileName, &config, fitapp.ProgramVersion()); ErrC == nil {
				// tried to read .git bug config
				config.FitYmlDir = string(config.ScmDir[:s-5])
				config.FitYml = string(config.ScmDir[:s-4]) + bugYmlFileName
			}
		} else if config.ScmType == "hg" {
			// try to read .hg fit config
			s := len(config.ScmDir)
			//fmt.Printf("debug s03 %v\n scmdir %v\n dir %v\n", s, config.ScmDir, string(config.ScmDir[:s-3])+fitYmlFileName)
			if ErrC := issues.ConfigRead(string(config.ScmDir[:s-3])+fitYmlFileName, &config, fitapp.ProgramVersion()); ErrC == nil {
				// tried to read .hg fit config
				config.FitYmlDir = string(config.ScmDir[:s-4])
				config.FitYml = string(config.ScmDir[:s-3]) + fitYmlFileName
				//fmt.Printf("debug 04\n %v\n", config.FitYml)
			} else if ErrC := issues.ConfigRead(string(config.ScmDir[:s-3])+bugYmlFileName, &config, fitapp.ProgramVersion()); ErrC == nil {
				// tried to read .hg fit config
				config.FitYmlDir = string(config.ScmDir[:s-4])
				config.FitYml = string(config.ScmDir[:s-3]) + bugYmlFileName
				//fmt.Printf("debug 04\n %v\n", config.FitYml)
			}
		}
		//} else { // no
	}

	if !rootPresent {
		if skip {
			fmt.Printf("Warn: no `" + config.FitDirName + "` directory\n")
		} else { // !skip
			//fitapp.PrintVersion()
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

	issues.FitDirer(config) // from issues/Directory.go, uses config.FitDir from issues/Bug.go

	// flags package is nice but this seemed more direct at the time
	// because of subcommands and
	// arguments that can be space separated names
	// glog requires the use of flag
	osArgs := os.Args // TODO: use an env var and assign to osArgs to setup for testing
	//fmt.Printf("A %s %#v\n", "osArgs: ", osArgs)
	if len(osArgs) <= 1 {
		if rootPresent {
			fitapp.List([]string{}, config, true)
		} else {
			fmt.Printf("Usage: " + os.Args[0] + " <command> [options]\n")
			fmt.Printf("\nUse \"bug help\" or \"bug help <command>\" for details.\n")
		}
	} else if len(osArgs) == 2 &&
		(osArgs[1] == "--help" ||
			osArgs[1] == "help" ||
			osArgs[1] == "-h") {
		fitapp.Help("help") // just bug help
	} else if len(osArgs) >= 3 &&
		(osArgs[2] == "--help" ||
			osArgs[2] == "help" ||
			osArgs[2] == "-h") { // 'bug cmd --help' just like 'bug help cmd'
		//fmt.Printf("B %s %#v\n", "osArgs: ", osArgs)
		fitapp.Help(osArgs[1])
	} else {
		switch osArgs[1] {
		// subcommands without osArgs first
		case "notags", "notag":
			fitapp.TagsNone(config)
		case "idslist", "idsassigned", "ids", "identifiers":
			fitapp.IdsAssigned(config)
		case "noids", "noid", "noidentifiers", "noidentifier":
			fitapp.IdsNone(config)
		case "env", "environment", "config", "settings":
			fitapp.Env(config)
		case "pwd", "dir", "cwd":
			fitapp.Pwd(config)
		case "version", "about", "--version", "-v":
			fitapp.PrintVersion()
		case "purge":
			fitapp.Purge(config)
		case "twilio":
			fitapp.Twilio(config)
		case "staging", "staged", "cached", "cache", "index":
			// TODO: scm/Staged.go
			if b, err := handler.SCMIssuesUpdaters(config); err != nil {
				fmt.Printf("Files in " + config.FitDirName + "/ need committing, see $ git status --porcelain -u -- :/" + config.FitDirName + "\nor if already in the index see     $ git diff --name-status --cached HEAD -- :/" + config.FitDirName + "\n")
				if _, ErrCach := handler.SCMIssuesCacher(config); ErrCach != nil {
					for _, bline := range strings.Split(string(b), "\n") {
						//if c, bline
						//if bline in c {
						//} else {
						fmt.Printf("%v\n", string(bline))
						//}
					}
				} else {
					fmt.Printf("%v\n", string(b))
				}
			} else {
				fmt.Printf("No files in " + config.FitDirName + "/ need committing, see $ git status --porcelain -u :/" + config.FitDirName + " \":top\"\n")
			}
		// subcommands that pass osArgs
		case "tagslist", "taglist", "tagsassigned", "tags":
			fitapp.TagsAssigned(osArgs[2:], config)
		case "list", "view", "show", "display", "ls":
			// bug list with no parameters shouldn't autopage,
			// bug list with issues to view should. So the original
			fitapp.List(osArgs[2:], config, true)
		case "find":
			fitapp.Find(osArgs[2:], config)
		case "create", "add", "new":
			//fmt.Printf("%s %#v\n", "osArgs: ", len(osArgs))
			fitapp.Create(osArgs[2:], config)
			//fmt.Printf("%s %#v\n", "osArgs: ", len(osArgs))
		case "edit":
			fitapp.Edit(osArgs[2:], config)
		case "retitle", "mv", "rename", "relabel":
			fitapp.Relabel(osArgs[2:], config)
		case "close", "rm":
			fitapp.Close(osArgs[2:], config)
		case "tag":
			fitapp.Tag(osArgs[2:], config) // boolean only
		case "id", "identifier":
			fitapp.Identifier(osArgs[2:], config)
		case "status":
			if len(osArgs) == 2 {
				// overview like a git status
				fitapp.Env(config)
			} else {
				// get or set the status of an issue
				fitapp.Status(osArgs[2:], config)
			}
		case "priority":
			fitapp.Priority(osArgs[2:], config)
		case "milestone":
			fitapp.Milestone(osArgs[2:], config)
		case "import":
			fitapp.Import(osArgs[2:], config)
		case "commit", "save":
			fitapp.Commit(osArgs[2:], config)
		case "roadmap":
			fitapp.Roadmap(osArgs[2:], config)
		case "help", "--help", "-h":
			//fmt.Printf("C %s %#v\n", "osArgs: ", osArgs)
			fitapp.Help(osArgs[2])
		default:
			//if
			if len(osArgs) == 2 {
				buglist, _ := issues.LoadIssueByHeuristic(osArgs[1], config)
				//fmt.Printf("%+v\n", buglist)
				if buglist != nil { // || ae, ok := bugerr.(issues.ErrNotFound); ! ok { // bug list when possible, not help
					fitapp.List(osArgs[1:], config, true)
				} else {
					fitapp.Help(osArgs[1:]...)
				}
			}
		}
	}
}
