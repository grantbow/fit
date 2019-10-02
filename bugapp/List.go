package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// bugNamer takes a bug and an int index then outputs a string.
func bugNamer(b bugs.Bug, idx int) string {
	if id := b.Identifier(); id != "" {
		return fmt.Sprintf("Issue %s", id)
	}
	return fmt.Sprintf("Issue %d", idx+1)
}

// listTags takes an array of os.FileInfo directories and prints bugs.
func listTags(files []os.FileInfo, args argumentList, config bugs.Config) {
	b := bugs.Bug{}
	for idx := range files {
		b.LoadBug(bugs.Directory(bugs.IssuesDirer(config) + dops + bugs.Directory(files[idx].Name())))

		for _, tag := range args {
			if b.HasTag(bugs.TagBoolTrue(tag)) {
				fmt.Printf("%s: %s\n", bugNamer(b, idx), b.Title("tags"))
			}
		}
	}
}

//byDir allows sort.Sort(byDir(issues))
// type and three functions are needed - also see Bug.go for type byBug
// rather than a custom Len function for os.FileInfo, Len is calculated in Less
type byDir []os.FileInfo

func (t byDir) Len() int {
	return len(t) // time.Format(time.UnixNano(t.modtime).UnixNano())
}
func (t byDir) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t byDir) Less(i, j int) bool {
	return (t[i]).ModTime().Unix() < (t[j]).ModTime().Unix()
}

// List is a subcommand to print issues.
func List(args argumentList, config bugs.Config, topRecurse bool) {
	issuesroot := bugs.IssuesDirer(config)
	issues, _ := ioutil.ReadDir(string(issuesroot))
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

	fmt.Printf("\n===== list %s\n", config.BugDir+sops+"issues")
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
						printIssueByDir(idx, issue, issuesroot, config, wantTags)
					} // else { continue }
				} // else { continue }

			}
		}
	} else if len(args) == 0 || (wantTags && len(args) == 1) || // --regex alone makes no sense
		(wantRecursive && len(args) == 1) ||
		(wantTags && wantRecursive && len(args) == 2) ||
		config.MultipleIssuesDirs == true {
		// No parameters, print a list of all bugs
		//os.Stdout = stdout
		for idx, issue := range issues {
			if issue.IsDir() != true {
				continue
			}
			printIssueByDir(idx, issue, issuesroot, config, wantTags)
		}
		if topRecurse == true && (wantRecursive || config.MultipleIssuesDirs == true) {
			fi, _ := os.Stat(config.BugDir)
			//if fierr != nil {
			//	panic(fierr)
			//}
			checkDirTree(args, config, fi, false)
		}
		return
	} else {
		// getAllTags() is defined in Tag.go
		// Get a list of tags, so that when we encounter
		// an error we can check if it's because the user
		// provided a tagname instead of a BugID. If they
		// did, then list bugs matching that tag instead
		// of full descriptions
		tags := getAllTags(config)
		// There were parameters, so show the full description of each
		// of those issues
		for i, length := 0, len(args); i < length; i += 1 {
			// TODO for _, tag := range args { // idx not needed
			if args[i] == "--match" || args[i] == "-m" ||
				args[i] == "--tags" || args[i] == "-t" ||
				args[i] == "--recursive" || args[i] == "-r" {
				continue
			}
			b, err := bugs.LoadBugByHeuristic(args[i], config)
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

			// err == nil so bug loaded
			b.ViewBug()
			if i < length-1 {
				fmt.Printf("\n--\n\n")
			}
		}
	}
	//fmt.Printf("\n")

	//if wantRecursive && topRecurse == true {
	// the first readdir is special so as to not find the same/regular issues dir
	//	fi, _ := os.Stat(config.BugDir)
	//if fierr != nil {
	//	panic(fierr)
	//}
	//	checkDirTree(args, config, fi, false)
	//fileinfos, _ := ioutil.ReadDir(config.BugDir) // BugRootDir
	//for _, node := range fileinfos {
	//	// search for a non-"issues" dir containing an "issues" subdir
	//	if node.Name() != "issues" &&
	//		node.IsDir() == true {
	//		os.Chdir(node.Name())
	//		checkDirTree(args, config, node, true)
	//		os.Chdir(wd) // go back
	//	}
	//}
	//}
}

func checkDirTree(args argumentList, config bugs.Config, node os.FileInfo, allowHits bool) {
	// check pwd
	topwd, _ := os.Getwd()
	//fmt.Printf("/////  debug checkDirTree issues dir %s\n", topwd)
	if dirinfo, err := os.Stat(topwd + sops + "issues"); err == nil && dirinfo.IsDir() && allowHits {
		//fmt.Printf("\n/////  issues in dir %s\n", topwd+sops+"issues")
		newConfig := config
		newConfig.BugDir = string(topwd) // BugRootDir
		List(args, newConfig, false)     // process
	}

	// recursively check subdirs
	fileinfos, _ := ioutil.ReadDir(topwd)
	//fmt.Printf("/////  debug checkDirTree subdir %s\n", topwd) // +sops+node.Name())
	for _, nodee := range fileinfos {
		if nodee.Name() != "issues" &&
			nodee.IsDir() == true {
			wd, _ := os.Getwd() // save for go back
			// candidate
			os.Chdir(nodee.Name()) // go down
			checkDirTree(args, config, nodee, true)
			os.Chdir(wd) // go back
		}
	}
}

func printIssueByDir(idx int, issue os.FileInfo, issuesroot bugs.Directory, config bugs.Config, wantTags bool) {
	// TODO: same next eight lines func (idx, issue)
	var dir bugs.Directory = issuesroot + dops + bugs.Directory(issue.Name())
	b := bugs.Bug{Dir: dir, DescriptionFileName: config.DescriptionFileName} // usually Description
	name := bugNamer(b, idx)                                                 // Issue idx: b.Title
	if wantTags == false {
		fmt.Printf("%s: %s\n", name, b.Title(""))
	} else {
		fmt.Printf("%s: %s\n", name, b.Title("tags"))
	}
}
