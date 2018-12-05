package main

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	args := ArgumentList(os.Args)
	config := bugs.Config{}
	temp := bugs.Config{}
	bugs.GetIssuesDir(config)
	bug_yml := ".bug.yml"
	//fmt.Fprint(os.Stdout, bug_yml)
	if fileinfo, err := os.Stat(bug_yml); err == nil && fileinfo.Mode().IsRegular() {
		dat, _ := ioutil.ReadFile(bug_yml)
		//fmt.Fprint(os.Stdout, dat)
		err := yaml.Unmarshal(dat, &temp); if err == nil {
			//fmt.Fprint(os.Stdout, config)
			//fmt.Fprint(os.Stdout, temp)
			//os.Exit(0)
			//config = temp // overwrites
			if temp.ImportXmlDump {
				config.ImportXmlDump = true
			}
			if temp.DefaultDescriptionFile != "" {
				config.DefaultDescriptionFile = temp.DefaultDescriptionFile
			}
		}
	}
	if githubRepo := args.GetArgument("--github", ""); githubRepo != "" {
		if strings.Count(githubRepo, "/") != 1 {
			fmt.Fprintf(os.Stderr, "Invalid GitHub repo: %s\n", githubRepo)
			os.Exit(2)
		}
		pieces := strings.Split(githubRepo, "/")
		githubImport(pieces[0], pieces[1], config)

	} else if args.GetArgument("--be", "") != "" {
		beImport(config)
	} else {
		if strings.Count(githubRepo, "/") != 1 {
			fmt.Fprintf(os.Stderr, "Usage: %s --github user/repo\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "       %s --be\n", os.Args[0])
			fmt.Fprintf(os.Stderr, `
Use this tool to import an external bug database into the local
issues/ directory.

Either "--github user/repo" is required to import GitHub issues,
from GitHub, or "--be" is required to import a local BugsEverywhere
database.
`)
			os.Exit(2)
		}
	}
}
