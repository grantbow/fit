package issues

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

//var dops = Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// IssueNotFoundError defines a new error.
type IssueNotFoundError string

// Error returns a string of the error.
func (b IssueNotFoundError) Error() string {
	return string(b)
}

//var dops = Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// FindIssuesByTag returns an array of tagged issues.
func FindIssuesByTag(tags []string, config Config) []Issue {
	root := RootDirer(&config)
	issuesroot := FitDirer(config)
	issues := readIssues(string(issuesroot))
	sort.Sort(byDir(issues))

	var bugs []Issue
	for _, issue := range issues { // idx not needed
		if issue.IsDir() == true {
			bug := Issue{}
			bug.LoadIssue(root+dops+Directory(config.FitDirName)+dops+Directory(issue.Name()), config)
			for _, tag := range tags {
				if bug.HasTag(TagBoolTrue(tag)) {
					bugs = append(bugs, bug)
					break
				}
			}
		}
	}
	//fmt.Printf("c %+v\n", bugs)
	sort.Sort(byIssue(bugs)) // allowed by byIssue defined in Issue.go
	//fmt.Printf("d %+v\n", bugs)
	return bugs
}

// LoadIssueByDirectory returns an issue from the directory name.
func LoadIssueByDirectory(dir string, config Config) (*Issue, error) {
	root := RootDirer(&config)
	_, err := ioutil.ReadDir(string(root) + sops + config.FitDirName + sops + dir)
	if err != nil {
		return nil, IssueNotFoundError("Not found " + dir)
	}
	bug := Issue{}
	bug.LoadIssue(FitDirer(config)+dops+Directory(dir), config)
	return &bug, nil
}

// LoadIssueByHeuristic returns an issue.
func LoadIssueByHeuristic(id string, config Config) (*Issue, error) {
	root := RootDirer(&config)
	issuesroot := FitDirer(config)
	issues := readIssues(string(issuesroot))
	sort.Sort(byDir(issues))
	//fmt.Printf("debug Found %v files\n", len(issues)) //, strings.Join(issues, ", "))
	//for _, issue := range issues {                    // idx not needed
	//	fmt.Printf("debug %v\n", issue.Name())
	//}

	if idx, err := strconv.Atoi(id); err == nil { // && idx > 0 && idx <= len(issues) {
		// check for an assigned id before index
		if bugptr, errId := LoadIssueByIdentifier(id, config); errId == nil {
			return bugptr, nil
		}
        // now check for index within range or pass through error
        return LoadIssueByIndex(idx, config)
	}
	// just a string, not an string of integers

	var candidate *Issue
	for _, issue := range issues { // idx not needed
		if issue.IsDir() == true {
			bug := Issue{}
			bug.LoadIssue(root+dops+Directory(config.FitDirName)+dops+Directory(issue.Name()), config)
			if bugid := bug.Identifier(); bugid == id {
				return &bug, nil
			} else if strings.Index(bugid, id) >= 0 {
				candidate = &bug
			}

		}
	}
	if candidate != nil {
		return candidate, nil
	}
	return nil, IssueNotFoundError("Not found " + id)
}

//// LoadIssueByStringIndex returns an issue from a string index.
//func LoadIssueByStringIndex(i string, config Config) (*Issue, error) {
//	root := RootDirer(&config)
//	issues, _ := ioutil.ReadDir(string(root) + sops + config.FitDirName)
//
//	idx, err := strconv.Atoi(i)
//	if err != nil {
//		return nil, IssueNotFoundError("Index not a number")
//	}
//	if idx < 1 || idx > len(issues) {
//		return nil, IssueNotFoundError("Invalid Index")
//	}
//
//	b := Issue{}
//	directoryString := fmt.Sprintf("%s%s%s%s%s", root, sops, config.FitDirName, sops, issues[idx-1].Name())
//	b.LoadIssue(Directory(directoryString), config)
//	return &b, nil
//}

// LoadIssueByIdentifier returns an issue from a string Identifier
func LoadIssueByIdentifier(id string, config Config) (*Issue, error) {
	root := RootDirer(&config)
	issuesroot := FitDirer(config)
	issues := readIssues(string(issuesroot))
	sort.Sort(byDir(issues))

	for _, issue := range issues { // idx not needed
		if issue.IsDir() == true {
			bug := Issue{}
			bug.LoadIssue(root+dops+Directory(config.FitDirName)+dops+Directory(issue.Name()), config)
			if bug.Identifier() == id {
				return &bug, nil
			}
		}
	}
	return nil, IssueNotFoundError("No issue named " + id)
}

// LoadIssueByIndex returns an issue from an int index.
func LoadIssueByIndex(idx int, config Config) (*Issue, error) {
	root := RootDirer(&config)
	issuesroot := FitDirer(config)
	issues := readIssues(string(issuesroot))
	sort.Sort(byDir(issues))
	if idx < 1 || idx > len(issues) {
		return nil, IssueNotFoundError("Invalid issue index")
	}

	b := Issue{}
	directoryString := fmt.Sprintf("%s%s%s%s%s", root, sops, config.FitDirName, sops, issues[idx-1].Name())
	// TODO: fix, can return files that are non-issues located in the issues directory
	b.LoadIssue(Directory(directoryString), config)
	return &b, nil
}

// GetAllIssues returns an array of all issues.
func GetAllIssues(config Config) []Issue {
	root := RootDirer(&config)
	issuesroot := FitDirer(config)
	issues := readIssues(string(issuesroot))
	sort.Sort(byDir(issues))
	//fmt.Printf("%+v\n", issues)

	var bugs []Issue
	for _, issue := range issues { // idx not needed
		if issue.IsDir() == true {
			bug := Issue{}
			bug.LoadIssue(root+dops+Directory(config.FitDirName)+dops+Directory(issue.Name()), config)
			bugs = append(bugs, bug)
		}
	}

	//fmt.Printf("a %+v\n", bugs)
	sort.Sort(byIssue(bugs)) // allowed by byIssue defined in Issue.go
	//fmt.Printf("b %+v\n", bugs)
	return bugs
}
