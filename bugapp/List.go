package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
)

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
		b.LoadBug(bugs.Directory(bugs.IssuesDirer(config) + "/" + bugs.Directory(files[idx].Name())))

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
func List(args argumentList, config bugs.Config) {
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

	if matchRegex && (len(args) > 1) {
		for i, length := 0, len(args); i < length; i += 1 {
			// TODO for _, tag := range args { // idx not needed
			if args[i] == "--match" || args[i] == "-m" ||
				args[i] == "--tags" || args[i] == "-t" {
				continue
			}
			fmt.Printf("\n===== matching (?i)%s\n", args[i])
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
	} else if len(args) == 0 || (wantTags && len(args) == 1) { // --regex alone makes no sense
		// No parameters, print a list of all bugs
		//os.Stdout = stdout
		for idx, issue := range issues {
			if issue.IsDir() != true {
				continue
			}
			printIssueByDir(idx, issue, issuesroot, config, wantTags)
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
				args[i] == "--tags" || args[i] == "-t" {
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
}

func printIssueByDir(idx int, issue os.FileInfo, issuesroot bugs.Directory, config bugs.Config, wantTags bool) {
	// TODO: same next eight lines func (idx, issue)
	var dir bugs.Directory = issuesroot + "/" + bugs.Directory(issue.Name())
	b := bugs.Bug{Dir: dir, DescriptionFileName: config.DescriptionFileName} // assumes DescriptionFileName
	name := bugNamer(b, idx)                                                 // Issue idx: b.Title
	if wantTags == false {
		fmt.Printf("%s: %s\n", name, b.Title(""))
	} else {
		fmt.Printf("%s: %s\n", name, b.Title("tags"))
	}
}
