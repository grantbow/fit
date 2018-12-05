package main

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/google/go-github/github" // handles json
	"io/ioutil"
	"encoding/json"
	"os"
	"context"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func githubImport(user, repo string, config bugs.Config) {
	client := github.NewClient(nil)
	i := 0
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// https://api.github.com/repos/<user>/<repo>/issues
	issues, resp, err := client.Issues.ListByRepo(context.Background(), user, repo, opt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for lastPage := false; lastPage != true; {
		i = 0
		for _, issue := range issues {
			i += 1
			// issues includes pull requests, so skip each pull request
			if issue.PullRequestLinks == nil {
				b := bugs.Bug{Dir: bugs.Directory(config.BugDir+"issues/" + string(bugs.TitleToDir(*issue.Title)))}
				if dir := b.GetDirectory(); dir != "" {
					os.Mkdir(string(dir), 0755)
				}
				if issue.Body != nil {
					b.SetDescription(*issue.Body)
				}
				if issue.Milestone != nil {
					b.SetMilestone(*issue.Milestone.Title)
				}
				if config.ImportXmlDump == true {
					// b.SetXml()
					xml, _ := json.MarshalIndent(issue, "", "    ")
					err = ioutil.WriteFile(string(b.GetDirectory())+"/issue.xml",append(xml,'\n'),0644)
					check(err)
				}
				// Don't set a bug identifier, but put an empty line and
				// then a GitHub identifier, so that bug commit can include
				// "Closes ..." in the commit message.
				b.SetIdentifier(fmt.Sprintf("\n\nGitHub:%s/%s%s%d\n", user, repo, "#", *issue.Number))
				for _, l := range issue.Labels {
					b.TagBug(bugs.Tag(*l.Name))
				}
				fmt.Printf("Importing %s\n", *issue.Title)
			}
		}
		if resp.NextPage == 0 {
			lastPage = true
		} else {
			opt.ListOptions.Page = resp.NextPage
			issues, resp, err = client.Issues.ListByRepo(context.Background(), user, repo, opt)
		}
	}
}
