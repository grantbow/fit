package main

import (
	"context"
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/google/go-github/github"
	"os"
)

func FetchIssues(owner string, repo string, opt *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error) {
	client := github.NewClient(nil)
	issues, response, err := client.Issues.ListByRepo(context.Background(), owner, repo, opt)
	return issues, response, err
}

func FetchIssueComments(owner string, repo string, comment int, opt *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error) {
	client := github.NewClient(nil)
	comments, response, err := client.Issues.ListComments(context.Background(), owner, repo, comment, opt)
	return comments, response, err
}

func githubImport(user, repo string) {
	issueDir := bugs.GetIssuesDir()
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, resp, err := FetchIssues(user, repo, opt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for lastPage := false; lastPage != true; {
		for _, issue := range issues {
			if issue.PullRequestLinks == nil {
				b := bugs.Bug{Dir: issueDir + bugs.TitleToDir(*issue.Title)}
				if dir := b.GetDirectory(); dir != "" {
					os.Mkdir(string(dir), 0755)
				}
				if issue.Body != nil {
					b.SetDescription(*issue.Body)
				}
				if issue.Milestone != nil {
					b.SetMilestone(*issue.Milestone.Title)
				}
				// Don't set a bug identifier, but put an empty line and
				// then a GitHub identifier, so that bug commit can include
				// "Closes ..." in the commit message.
				b.SetIdentifier(fmt.Sprintf("\n\nGitHub:%s/%s%s%d\n", user, repo, "#", *issue.Number))
				for _, l := range issue.Labels {
					b.TagBug(bugs.Tag(*l.Name))
				}
				if *issue.Comments > 0 {
					comments, _, err := FetchIssueComments(user, repo, *issue.Number, nil)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return
					}
					for _, l := range comments {
						b.CommentBug(bugs.Comment(*l.Body))
					}
				}
				fmt.Printf("Importing %s\n", *issue.Title)
			}
		}
		if resp.NextPage == 0 {
			lastPage = true
		} else {
			opt.ListOptions.Page = resp.NextPage
			issues, resp, err = FetchIssues(user, repo, opt)
		}
	}
}
