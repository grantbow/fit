package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// issueNamer takes a bug and an int index then outputs a string.
func issueNamer(b bugs.Issue, idx int) string {
	id := b.Identifier()
	//fmt.Printf("debug b.Dir %v id %v\n", b.Dir, id)
	if id != "" {
		return fmt.Sprintf("Issue %s", id)
	}
	return fmt.Sprintf("Issue %d", idx+1)
}

// listTags takes an array of os.FileInfo directories and prints bugs.
func listTags(files []os.FileInfo, args argumentList, config bugs.Config) {
	b := bugs.Issue{}
	for idx := range files {
		b.LoadIssue(bugs.Directory(bugs.FitDirer(config)+dops+bugs.Directory(files[idx].Name())), config)

		for _, tag := range args {
			if b.HasTag(bugs.TagBoolTrue(tag)) {
				fmt.Printf("%s: %s\n", issueNamer(b, idx), b.Title("tags"))
			}
		}
	}
}

// List is a subcommand to print lists and individual issues.
func List(args argumentList, config bugs.Config, topRecurse bool) {
	fitdir := bugs.FitDirer(config)
	issues := readIssues(string(fitdir))
	sort.Sort(byDir(issues))

	var wantTags bool = false
	if args.HasArgument("--tags") || args.HasArgument("-t") {
		wantTags = true
	}
	var matchRegex bool = false
	if args.HasArgument("--match") || args.HasArgument("-m") {
		matchRegex = true
	}
	var wantRecursive bool = false
	if args.HasArgument("--recursive") || args.HasArgument("-r") {
		wantRecursive = true
	}

	//fmt.Printf("debug topRecurse %v wantRecursive %v config.MultipleDirs %v \n", topRecurse, wantRecursive, config.MultipleDirs == true)
	if topRecurse == true && (wantRecursive == true || config.MultipleDirs == true) {
		// print warning if below the .git directory
		//fmt.Printf("debug config.FitDir %v len %v config.ScmDir %v len %v\n",
		//	config.FitDir, len(config.FitDir),
		//	config.ScmDir, len(config.ScmDir))
		if config.ScmType == "git" &&
			len(config.ScmDir) != len(config.FitDir)+5 {
			fmt.Printf("\n===== WARNING, path from .git to %s: %s\n", config.FitDirName, config.FitDir[len(config.ScmDir)-5:])
		} else if config.ScmType == "hg" &&
			len(config.ScmDir) != len(config.FitDir)+4 {
			fmt.Printf("\n===== WARNING, path from .hg to %s: %s\n", config.FitDirName, config.FitDir[len(config.ScmDir)-4:])
		}
	}

	fmt.Printf("\n===== list %s\n", config.FitDir+sops+config.FitDirName)
	if matchRegex && (len(args) > 1) {
		for i, length := 0, len(args); i < length; i += 1 {
			// TODO for _, tag := range args { // idx not needed
			if args[i] == "--match" || args[i] == "-m" ||
				args[i] == "--tags" || args[i] == "-t" ||
				args[i] == "--recursive" || args[i] == "-r" {
				continue
			}
			fmt.Printf("===== matching (?i)%s\n", args[i])
			for idx, issue := range issues {
				if issue.IsDir() != true {
					continue
				}
				re, err := regexp.Compile("(?i)" + args[i])
				if err == nil {
					s := re.Find([]byte(issue.Name()))
					if s != nil {
						printIssueByDir(idx, issue, fitdir, config, wantTags)
					} // else { continue }
				} // else { continue }
			}
		}
	} else if len(args) == 0 || (wantTags && len(args) == 1) || // --regex alone makes no sense
		(wantRecursive && len(args) == 1) ||
		(wantTags && wantRecursive && len(args) == 2) ||
		config.MultipleDirs == true {
		// No parameters, print a list of all bugs
		//os.Stdout = stdout
		foundsome := 0
		for idx, issue := range issues {
			if issue.IsDir() != true {
				continue
			}
			printIssueByDir(idx, issue, fitdir, config, wantTags)
			foundsome += 1
		}
		if foundsome == 0 {
			fmt.Printf("   << found no issues >>\n")
		}
		if topRecurse == true && (wantRecursive || config.MultipleDirs == true) {
			fi, _ := os.Stat(config.FitDir)
			//if fierr != nil {
			//	panic(fierr)
			//}
			checkDirTreeDown(args, config, fi, false)
		}
		return
	} else {
		// Get a list of tags, so that when we encounter
		// an error we can check if it's because the user
		// provided a tagname instead of a IssueID. If they
		// did, then list bugs matching that tag instead
		// of full descriptions
		//tags := getAllTags(config) // defined in Tag.go
		tags := uniqueTagList(config)
		// There were parameters, so show the full description of each
		// of those issues
		for i, length := 0, len(args); i < length; i += 1 {
			// TODO for _, tag := range args { // idx not needed
			if args[i] == "--match" || args[i] == "-m" ||
				args[i] == "--tags" || args[i] == "-t" ||
				args[i] == "--recursive" || args[i] == "-r" {
				continue
			}
			b, err := bugs.LoadIssueByHeuristic(args[i], config)
			if err != nil {
				for _, tagname := range tags {
					if tagname == args[i] && matchRegex == false {
						listTags(issues, args, config)
						return
					}
				}
				fmt.Printf("%s\n", err.Error())
				continue
			}

			// err == nil so issue loaded
			b.ViewIssue()
			if i < length-1 {
				fmt.Printf("\n--\n\n")
			}
		}
	}
	//fmt.Printf("\n")

	//if wantRecursive && topRecurse == true {
	// the first readdir is special so as to not find the same/regular fit dir
	//	fi, _ := os.Stat(config.FitDir)
	//if fierr != nil {
	//	panic(fierr)
	//}
	//	checkDirTreeDown(args, config, fi, false)
	//fileinfos, _ := ioutil.ReadDir(config.FitDir) // IssueRootDir
	//for _, node := range fileinfos {
	//	// search for a non-"issues" dir containing an "issues" subdir
	//	if node.Name() != "issues" &&
	//		node.IsDir() == true {
	//		os.Chdir(node.Name())
	//		checkDirTreeDown(args, config, node, true)
	//		os.Chdir(wd) // go back
	//	}
	//}
	//}
}

func checkDirTreeDown(args argumentList, config bugs.Config, node os.FileInfo, allowHits bool) {
	// check pwd
	topwd, _ := os.Getwd()
	//fmt.Printf("/////  debug checkDirTreeDown fit dir %s\n", topwd)
	if dirinfo, err := os.Stat(topwd + sops + config.FitDir); err == nil && dirinfo.IsDir() && allowHits {
		//fmt.Printf("\n/////  topwd FitDir %s\n", topwd+sops+config.FitDir)
		newConfig := config
		newConfig.FitDir = string(topwd) // IssueRootDir
		List(args, newConfig, false)     // process
	}

	// recursively check subdirs
	fileinfos, _ := ioutil.ReadDir(topwd)
	//fmt.Printf("/////  debug checkDirTreeDown subdir %s\n", topwd) // +sops+node.Name())
	for _, nodee := range fileinfos {
		if nodee.Name() != config.FitDir &&
			nodee.IsDir() == true {
			wd, _ := os.Getwd() // save for go back
			// candidate
			os.Chdir(nodee.Name())                      // go down
			checkDirTreeDown(args, config, nodee, true) // calls itself
			os.Chdir(wd)                                // go back
		}
	}
}

func printIssueByDir(idx int, issue os.FileInfo, fitdir bugs.Directory, config bugs.Config, wantTags bool) {
	// TODO: same next eight lines func (idx, issue)
	var dir bugs.Directory = fitdir + dops + bugs.Directory(issue.Name())
	b := bugs.Issue{Dir: dir, DescriptionFileName: config.DescriptionFileName} // usually Description
	name := issueNamer(b, idx)                                                 // Issue idx: b.Title
	if wantTags == false {
		fmt.Printf("%s: %s\n", name, b.Title(""))
	} else {
		fmt.Printf("%s: %s\n", name, b.Title("tags"))
	}
}
