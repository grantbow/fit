package bugapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/google/go-github/github" // handles json
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
)

func fetchIssues(owner string, repo string, opt *github.IssueListByRepoOptions, client *github.Client) ([]*github.Issue, *github.Response, error) {
	issues, response, err := client.Issues.ListByRepo(context.Background(), owner, repo, opt)
	return issues, response, err
}

func fetchIssueComments(owner string, repo string, comment int, opt *github.IssueListCommentsOptions, client *github.Client) ([]*github.IssueComment, *github.Response, error) {
	comments, response, err := client.Issues.ListComments(context.Background(), owner, repo, comment, opt)
	return comments, response, err
}

// githubImportIssues downloads issues and comments from a github repository
func githubImportIssues(user, repo string, config bugs.Config) {
	i := 0
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	client := github.NewClient(nil)
	// https://api.github.com/repos/<user>/<repo>/issues
	issues, resp, err := fetchIssues(user, repo, opt, client)
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
				ititle := string(bugs.TitleToDir(fmt.Sprintf("%s%s%v", *issue.Title, "-", *issue.Number)))
				fmt.Printf("Importing issue %s\n", ititle)
				// add issue.Number to title
				//b := bugs.Bug{Dir: bugs.Directory(config.BugDir+"issues/" + string(bugs.TitleToDir(*issue.Title)))}
				b := bugs.Bug{Dir: bugs.Directory(config.BugDir + "/issues/" + ititle)}
				if dir := b.Direr(); dir != "" {
					os.Mkdir(string(dir), 0755)
				}
				if issue.Body != nil {
					b.SetDescription(*issue.Body, config)
				} else {
					b.SetDescription("", config)
				}
				if issue.Milestone != nil {
					b.SetMilestone(*issue.Milestone.Title, config)
				}
				if config.ImportXmlDump == true {
					// b.SetXml()
					xml, _ := json.MarshalIndent(issue, "", "    ")
					err = ioutil.WriteFile(string(b.Direr())+"/issue.xml", append(xml, '\n'), 0644)
					check(err)
				}
				// Don't set a bug identifier, but put an empty line and
				// then a GitHub identifier, so that bug commit can include
				// "Closes ..." in the commit message.
				b.SetIdentifier(fmt.Sprintf("GitHub:%s/%s%s%d", user, repo, "#", *issue.Number), config)
				for _, lab := range issue.Labels {
					b.TagBug(bugs.TagBoolTrue(*lab.Name), config)
				}
				j := 1
				if *issue.Comments > 0 {
					comments, _, err := fetchIssueComments(user, repo, *issue.Number, nil, client)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return
					}
					for _, co := range comments {
						xml, err := json.MarshalIndent(co, "", "    ")
						check(err)
						x := bugs.Comment{
							Author: *co.User.Login,
							Time:   *co.CreatedAt,
							Body:   *co.Body,
							Order:  j,
							Xml:    xml}
						b.CommentBug(x, config)
						if config.ImportXmlDump == true {
							// b.SetXml()
							comname := "comment-" + string(bugs.ShortTitleToDir(string(*co.Body))) + "-" + fmt.Sprintf("%v", j)
							err = ioutil.WriteFile(string(b.Direr())+"/"+comname+".xml", append(xml, '\n'), 0644)
							check(err)
						}
						j += 1
					}
				}
			}
		}
		if resp.NextPage == 0 {
			lastPage = true
		} else {
			opt.ListOptions.Page = resp.NextPage
			issues, resp, err = fetchIssues(user, repo, opt, client)
			check(err)
		}
	}
}

func fetchProjectCards(columnid int64, opt *github.ProjectCardListOptions, client *github.Client) ([]*github.ProjectCard, *github.Response, error) {
	projectcards, response, err := client.Projects.ListProjectCards(context.Background(), columnid, opt)
	return projectcards, response, err
}

// uninteresting, no new info
//func fetchProjectColumn(columnid int64, opt *github.ListOptions, client *github.Client) ([]*github.ProjectColumn, *github.Response, error) {
//projectcolumn, response, err := client.Projects.GetProjectColumn(context.Background(), projectid)
//return projectcolumn, response, err
//}

func fetchProjectColumns(projectid int64, opt *github.ListOptions, client *github.Client) ([]*github.ProjectColumn, *github.Response, error) {
	projectcolumns, response, err := client.Projects.ListProjectColumns(context.Background(), projectid, opt)
	return projectcolumns, response, err
}

func fetchProjects(owner string, repo string, opt *github.ProjectListOptions, client *github.Client) ([]*github.Project, *github.Response, error) {
	projects, response, err := client.Repositories.ListProjects(context.Background(), owner, repo, opt)
	return projects, response, err
}

func githubImportProjects(user, repo string, config bugs.Config) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GithubPersonalAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc) // oauthClient
	i := 0
	opt := &github.ProjectListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// https://api.github.com/repos/<user>/<repo>/projects
	projects, resp, err := fetchProjects(user, repo, opt, client)
	//fmt.Printf("%v\n", projects)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for lastPage := false; lastPage != true; {
		i = 0
		for _, project := range projects {
			i += 1
			projname := "proj-" + string(bugs.TitleToDir(fmt.Sprintf("%v%s%s", *project.Number, "-", *project.Name)))
			fmt.Printf("Importing %s\n", projname)
			b := bugs.Bug{Dir: bugs.Directory(config.BugDir + "/issues/" + projname)}
			if dir := b.Direr(); dir != "" {
				os.Mkdir(string(dir), 0755)
			}
			if project.Body != nil {
				b.SetDescription(*project.Body, config)
			} else {
				b.SetDescription("", config)
			}
			b.SetIdentifier(fmt.Sprintf("GitHub:%s/%s%s%v", user, repo, "/projects/", *project.Number), config)

			j := 1
			projectcolumns, _, err := fetchProjectColumns(int64(*project.ID), nil, client)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			for _, pc := range projectcolumns {
				xmlbytes := &bytes.Buffer{}
				xml, err := json.MarshalIndent(pc, "", "    ")
				check(err)
				for i := 0; i < len(xml); i++ {
					xmlbytes.WriteByte(xml[i])
				}

				projectcards, _, err := fetchProjectCards(int64(*pc.ID), nil, client)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}
				for _, card := range projectcards {
					cardxml, err := json.MarshalIndent(card, "", "    ")
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return
					}
					xmlbytes.WriteByte("\n"[0])
					for i := 0; i < len(cardxml); i++ {
						xmlbytes.WriteByte(cardxml[i])
					}
				}
				xmlbytes.WriteByte("\n"[0])

				if config.ImportXmlDump == true {
					colname := "col-" + fmt.Sprintf("%v", j) + "-" + string(bugs.ShortTitleToDir(string(*pc.Name)))
					fmt.Printf("\nImporting %v\n", colname)
					err = ioutil.WriteFile(string(b.Direr())+"/"+colname+".xml", xmlbytes.Bytes(), 0644)
					check(err)
				}
				j += 1
			}
			if config.ImportXmlDump == true {
				// b.SetXml()
				xml, _ := json.MarshalIndent(*project, "", "    ")
				err = ioutil.WriteFile(string(b.Direr())+"/project.xml", append(xml, '\n'), 0644)
				check(err)
			}

		}
		if resp.NextPage == 0 {
			lastPage = true
		} else {
			opt.ListOptions.Page = resp.NextPage
			projects, resp, err = fetchProjects(user, repo, opt, client)
			check(err)
		}
	}
}
