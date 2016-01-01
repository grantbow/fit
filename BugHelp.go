package main

import (
	"fmt"
	"os"
)

func (a BugApplication) Help(args ...string) {
	var cmd string
	if args == nil {
		cmd = "help"

	}
	if len(args) <= 1 {
		cmd = "help"
	} else {
		cmd = args[1]
	}
	switch cmd {
	case "create":
		fmt.Printf("Usage: " + os.Args[0] + " create [-n] Issue Title\n\n")
		fmt.Printf("This will create an issue with the title Issue Title\n\n")
		fmt.Printf("An editor will be opened automatically for you to enter\n")
		fmt.Printf("a more detailed description.\n\n")
		fmt.Printf("If your EDITOR environment variable is set, it will be\n")
		fmt.Printf("used, otherwise the default is vim.\n")
		fmt.Printf("If the first argument to create is \"-n\", then " + os.Args[0] + " will not open any editor and create an empty Description\n\n")
	case "list":
		fmt.Printf("Usage: " + os.Args[0] + " list [issue numbers]\n")
		fmt.Printf("       " + os.Args[0] + " list [tags]\n\n")
		fmt.Printf(`This will list the issues found in the current environment

With no arguments, titles will be printed to the screen along with the issue
number that can be used to reference this issue on the command line.

If 1 or more issue numbers are provided, the whole issue including description
will be printed to stdout.

If, instead of issue numbers, you provide list with 1 or more tags, it will
print any issues which have that tag (in short form)
`)

	case "edit":
		fmt.Printf("Usage: " + os.Args[0] + " edit IssueNumber\n\n")
		fmt.Printf(
			`This will launch your standard editor to edit the description of the bug numbered 
IssueNumber, where IssueNumber is a reference to same index provided with a
"bug list" command.
`)
	case "status":
		fmt.Printf("Usage: " + os.Args[0] + " status IssueNumber [NewStatus]\n\n")
		fmt.Printf(
			`This will edit or display the status of the bug numbered IssueNumber.
            
If NewStatus is provided, it will update the first line of the Status file
for the issue (creating the file as necessary). If not provided, it will 
display the first line of the Status file to STDOUT.

Note that you can manually edit the Status file in the issues/ directory
to provide further explanation (for instance, why that status is set.)
This command will preserve the explanation when updating a status.
`)
	case "priority":
		fmt.Printf("Usage: " + os.Args[0] + " priority IssueNumber [NewPriority]\n\n")
		fmt.Printf(
			`This will edit or display the priority of the bug numbered IssueNumber.
By convention, priorities should be an integer number (higher is more 
urgent), but that is not enforced by this command and NewPriority can
be any free-form text if you prefer.
            
If NewPriority is provided, it will update the first line of the Priority
file for the issue (creating the file as necessary). If not provided, it 
will display the first line of the Priority file to STDOUT.

Note that you can manually edit the Priority file in the issues/ directory
to provide further explanation (for instance, why that priority is set.)
This command will preserve the explanation when updating a priority.
`)
	case "rm":
		fallthrough
	case "close":
		fmt.Printf("Usage: " + os.Args[0] + " close IssueNumber\n\n")
		fmt.Printf(`This will delete the issue numbered IssueNumber. IssueNumbers
correspond to the number in the "bug list" command.

Note that closing a bug will cause all existing bugs to be be renumbered and
IssueNumbers are not intended to be stable.

Also note that this does not remove the issue from git, but only from the file
system. If you want to remove an issue that is tracked by git, you'll have to
manually "git rm -r" the directory from the directory that's printed to the
screen if you execute "bug dir".`)
		fmt.Printf("\n\n\"%s rm\" is an alias for this \"%s close\"\n", os.Args[0], os.Args[0])
	case "purge":
		fmt.Printf("Usage: " + os.Args[0] + " purge\n\n")
		fmt.Printf(`This will delete any bugs that are not currently tracked by
git.

It is an alias for "git clean -fd $(bug dir)"
`)
	case "commit":
		fmt.Printf("Usage: " + os.Args[0] + " commit\n\n")
		fmt.Printf(`This will commit any new, modified, or removed issues to
git.

Your working tree and staging area should be otherwise unaffected by using
this command.
`)
	case "env":
		fmt.Printf("Usage: " + os.Args[0] + " env\n\n")
		fmt.Printf(`This will print the environment used by the bug command to stdout.

Use this command if you want to see what directory bug create is
using to store bugs, or what editor will be invoked by bug create/edit.
`)

	case "pwd":
		fallthrough
	case "dir":
		fmt.Printf("Usage: " + os.Args[0] + " dir\n\n")
		fmt.Printf(`This will bug directory to stdout, so you can use it as a subcommand
for arguments to any arbitrary shell commands. For example "cd $(bug dir)" or 
"git rm -r $(bug dir)/Issue-Title"
`)
		fmt.Printf("\n\n\"%s pwd\" is an alias for this \"%s dir\"\n", os.Args[0], os.Args[0])
	case "tag":
		fmt.Printf("Usage: " + os.Args[0] + " tag IssueNumber [tags]\n\n")
		fmt.Printf(`This will tag the given IssueNumber with the tags
given as parameters. At least one tag is required.

Tags can be any string which would make a valid file name.
`)
	case "help":
		fallthrough
	default:
		fmt.Printf("Usage: " + os.Args[0] + " command [options]\n\n")
		fmt.Printf("Use \"bug help [command]\" for more information about any command below\n\n")
		fmt.Printf("Valid commands\n")
		fmt.Printf("\nIssue editing commands:\n")
		fmt.Printf("\tcreate\t File a new bug\n")
		fmt.Printf("\tlist\t List existing bugs\n")
		fmt.Printf("\tedit\t Edit an existing bug\n")
		fmt.Printf("\ttag\t Tag a bug with a category\n")
		fmt.Printf("\tclose\t Delete an existing bug\n")
		fmt.Printf("\trm\t Alias of close\n")
		fmt.Printf("\tstatus\t View or edit a bug's status\n")
		fmt.Printf("\tpriority View or edit a bug's priority\n")

		fmt.Printf("\nSource control commands:\n")
		fmt.Printf("\tcommit\t Commit any new, changed or deleted bug to git\n")
		fmt.Printf("\tpurge\t Remove all issues not tracked by git\n")

		fmt.Printf("\nOther commands:\n")
		fmt.Printf("\tenv\t Show settings that bug will use if invoked from this directory\n")
		fmt.Printf("\tdir\t Prints the issues directory to stdout (useful subcommand in the shell)\n")
		fmt.Printf("\tpwd\t Alias of dir\n")
		fmt.Printf("\thelp\t Show this screen\n")
	}
}