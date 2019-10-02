package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// filecp copies files.
// see https://opensource.com/article/18/6/copying-files-go
func filecp(sourceFile string, destinationFile string) {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}
}

// Create is a subcommand to open a new issue.
func Create(Args argumentList, config bugs.Config) {
	//fmt.Print("a\n")
	if len(Args) < 1 || (len(Args) < 2 && Args[0] == "-n") {
		//fmt.Print("b\n")
		fmt.Fprintf(os.Stderr, "Usage: %s create [-n] <Bug Description>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNo Bug Description provided.\n")
		return
	}
	var noDesc bool = false

	if Args.HasArgument("-n") {
		noDesc = true
		Args = Args[1:]
	}

	Args, argVals := Args.GetAndRemoveArguments([]string{"--tag", "--status", "--priority", "--milestone", "--identifier", "--id"})
	tag := argVals[0]
	status := argVals[1]
	priority := argVals[2]
	milestone := argVals[3]
	identifier := argVals[4] + argVals[5]

	if Args.HasArgument("--generate-id") {
		for i, token := range Args {
			if token == "--generate-id" {
				if i+1 < len(Args) {
					Args = append(Args[:i], Args[i+1:]...)
					break
				} else {
					Args = Args[:i]
					break
				}
			}
		}
		identifier = generateID(strings.Join(Args, " "))
	}

	// It's possible there were arguments provided, but still no title
	// included. Do another check before trying to create the bug.
	if strings.TrimSpace(strings.Join(Args, " ")) == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s create [-n] <Bug Description>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNo Bug Description provided.\n")
		return
	}
	var bgid = bugs.IssuesDirer(config)
	if bgid == "" {
		os.MkdirAll("issues", 0700)
		bgid = bugs.IssuesDirer(config)
	}
	var bug = bugs.Bug{
		Dir:                 bgid + dops + bugs.TitleToDir(strings.Join(Args, " ")),
		DescriptionFileName: config.DescriptionFileName,
	}

	dir := bug.Direr()

	var mode os.FileMode
	mode = 0775
	err := os.Mkdir(string(dir), mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%s error: mkdir\n", os.Args[0])
		log.Fatal(err)
	}
	DescriptionFile := string(dir) + sops + config.DescriptionFileName
	if noDesc {
		txt := []byte("")
		if config.DefaultDescriptionFile != "" {
			filecp(config.DefaultDescriptionFile, DescriptionFile)
		} else {
			//fmt.Printf("here %s\n", config.DescriptionFileName)
			ioutil.WriteFile(DescriptionFile, txt, 0644)
		}
	} else {
		if config.DefaultDescriptionFile != "" {
			filecp(config.DefaultDescriptionFile, DescriptionFile)
		}
		cmd := exec.Command(getEditor(), DescriptionFile)

		//osi := os.Stdin
		//oso := os.Stdout
		//ose := os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	if tag != "" {
		bug.TagBug(bugs.TagBoolTrue(tag), config)
	}
	if status != "" {
		bug.SetStatus(status, config)
	}
	if priority != "" {
		bug.SetPriority(priority, config)
	}
	if milestone != "" {
		bug.SetMilestone(milestone, config)
	}
	if identifier != "" {
		bug.SetIdentifier(identifier, config)
	}
	fmt.Printf("Created issue: %s\n", bug.Title(""))
}
