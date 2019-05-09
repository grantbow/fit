package bugapp

import (
	"fmt"
	"os"
)

// Help is a subcommand to describe the program and it's subcommands.
func Help(args ...string) {
	var cmd string
	if args == nil || len(args) == 0 {
		cmd = "help"
	} else if len(args) == 1 {
		//cmd = "help"
		cmd = args[0]
	} else if len(args) == 2 && (args[1] == "help" || args[1] == "--help") {
		cmd = args[0]
	} else {
		cmd = args[1]
	}
	fmt.Printf("cmd %v\n", cmd)
	switch cmd {
	case "add", "new", "create":
		fmt.Printf("Usage: " + os.Args[0] + " create [-n] [options] Issue Title\n\n")
		fmt.Printf(
			`This will create an issue with the title Issue Title.  An editor 
will be opened on Description automatically.  If your EDITOR
environment variable is set, it will be used, otherwise
the default editor is vim.

If the first argument to create is "-n", then %s will not open 
any editor and create an empty Description.

Options take a value and set a field on the bug at the same
time as creating it. Valid options are:
    --status     Sets the bug status to the next parameter
    --tag        Tags the bug with a tag on creation
    --priority   Sets the priority to the next parameter
    --milestone  Sets the milestone to the next parameter
    --identifier Sets the identifier to the next parameter
    --generate-id Automatically generate a stable bug identifier
`, os.Args[0])
	case "list":
		fmt.Printf("Usage: " + os.Args[0] + " list \n\n")
		fmt.Printf("       " + os.Args[0] + " list <BugID>...\n\n")
		fmt.Printf("       " + os.Args[0] + " list <-m|--match> <regex>...\n\n")
		fmt.Printf("       " + os.Args[0] + " list <-t|--tags> <BugID>...\n\n")
		fmt.Printf("       " + os.Args[0] + " list <tag>...\n\n")
		fmt.Printf(
			`This will list the issues found in the current environment

With no arguments, issue number and titles will be printed.
Issue numbers can reference this issue on the command line.

The [-m|--match] option tells list you are providing  a regular
expression. Matching issues are listed.

All unix shell parameters that contain special characters, which many
regular expressions use, must be escaped. This prevents the automatic
"filename expansion" or "globbing" performed by the shell before
launching the bug command. A good references with details are at:
http://tldp.org/LDP/abs/html/globbingref.html
https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html

If valid BugIDs are provided, whole issues with Description
will print.  See "bug help ids" for what makes a BugID.

If 1 or more <tag]s are provided, matching issues are listed.

Note that BugIDs may change as you create, edit, and close other
issues. Details are provided by "bug help ids."

An alias for "list" is  "view".
`)

	case "edit":
		fmt.Printf("Usage: " + os.Args[0] + " edit <Filename> <BugID>\n\n")
		fmt.Printf(
			`This will launch your standard editor to edit the Description 
of the bug identified by BugID.  See "bug help ids" for
what makes a BugID.

If the Filename option is provided, bug will instead launch an editor
to edit that file name within the bug directory. Files that have
special meaning to bug (Status, Milestone, Priority, Identifier) are
treated in a case insensitive manner, otherwise the filename is passed
directly to your editor.
`)
	case "status":
		fmt.Printf("Usage: " + os.Args[0] + " status <BugID> <NewStatus>\n\n")
		fmt.Printf(
			`This will edit or display the status of the bug identified by BugID.
See "bug help ids" for what constitutes a BugID.
            
If NewStatus is provided, it will update the first line of the Status file
for the issue (creating the file as necessary). If not provided, it will 
display the first line of the Status file to STDOUT.

Note that you can edit the status in your standard editor with the
command "%s edit status BugID". If you provide a longer than 1 line
status with "bug edit status", "bug status" will preserve everything
after the first line when editing a status. You can use this to provide
further context on a status (for instance, why that status is setup.)
`, os.Args[0])
	case "priority":
		fmt.Printf("Usage: " + os.Args[0] + " priority <BugID> <New Priority>\n\n")
		fmt.Printf(
			`This will edit or display the priority of BugID. See "bug help ids"
for what constitutes a BugID.

By convention, priorities should be an integer number (higher is more 
urgent), but that is not enforced by this command and <New Priority> can
be any free-form text if you prefer.
            
If <New Priority> is provided, it will update the first line of the Priority
file for the issue (creating the file as necessary). If not provided, it 
will display the first line of the Priority file to STDOUT.

Note that you can manually edit the Priority file in the issues/ directory
by running "%s edit priority BugID", to provide further explanation (for 
instance, why that priority is set.) This command will preserve the 
explanation when updating a priority.
`, os.Args[0])
	case "milestone":
		fmt.Printf("Usage: " + os.Args[0] + " milestone <BugID> <New Milestone>\n\n")
		fmt.Printf(
			`This will edit or display the milestone of the identified by BugID.
See "%s help ids" for what constitutes a BugID.

There are no restrictions on how milestones must be named, but
semantic versioning is a good convention to adopt. Failing that,
it's a good idea to use milestones that collate properly when
sorted as strings so that they appear properly in "%s roadmap".

If <New Milestone> is provided, it will update the first line of the
Milestone file for the issue (creating the file as necessary). 
If not provided, it will display the first line of the Milestone 
file to STDOUT.

Note that you can manually edit the Milestone file in the issues/
directory to provide further explanation (for instance, why that 
milestone is set) with the command "bug edit milestone BugID"

This command will preserve the explanation when updating a priority.
`, os.Args[0], os.Args[0])
	case "retitle", "mv", "rename", "relabel":
		fmt.Printf("Usage: " + os.Args[0] + " retitle <BugID> <New Title>\n\n")
		fmt.Printf(
			`This will change the title of BugID to <New Title>. Use this
to rename an issue.

"%s mv", "%s relabel", and "%s rename" are all aliases for "%s retitle".
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	case "rm", "close":
		fmt.Printf("Usage: " + os.Args[0] + " close <BugID>\n")
		fmt.Printf("       " + os.Args[0] + " rm <BugID>\n\n")
		fmt.Printf(
			`This will delete the issue identifier by BugID. See
"%s help ids" for details on what constitutes a BugID.

Note that closing a bug may cause existing BugIDs to change if
they do not have a stable id set (see "%s help ids",
again.)

Also note that this does not remove the issue from git, but only 
from the file system. You'll need to execute "bug commit" to
remove the bug from version control.

"%s rm" is an alias for this "%s close"
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	case "find":
		fmt.Printf("Usage: %s find tag <value1> [value2 ...]\n", os.Args[0])
		fmt.Printf("Usage: %s find status <value1> [value2 ...]\n", os.Args[0])
		fmt.Printf("Usage: %s find priority <value1> [value2 ...]\n", os.Args[0])
		fmt.Printf("Usage: %s find milestone <value1> [value2 ...]\n\n", os.Args[0])
		fmt.Printf(
			`This will search all bugs for multiple tags, statuses, priorities, or milestone.
The matching bugs will be printed.
`)
	case "purge":
		fmt.Printf("Usage: " + os.Args[0] + " purge\n\n")
		fmt.Printf(
			`This will delete any bugs that are not currently tracked by
git.
`)
	case "commit":
		fmt.Printf("Usage: " + os.Args[0] + " commit [--no-autoclose]\n\n")
		fmt.Printf(`This will commit any new, modified, or removed issues to
git or hg.

Your working tree and staging area should be otherwise
unaffected by using this command.

If the --no-autoclose option is passed to commit, bug will
not include a "Closes #x" line for each issue imported from
"bug-import --github." Otherwise, the commit message will
include the list of issues that were closed so that GitHub
will autoclose them when the changes are pushed upstream.
`)
	case "env":
		fmt.Printf("Usage: " + os.Args[0] + " env\n\n")
		fmt.Printf(`This will print the environment used by the bug command to stdout.

Use this command if you want to see what directory bug create is
using to store bugs, or what editor will be invoked by bug create/edit.
`)

	case "dir", "pwd":
		fmt.Printf("Usage: " + os.Args[0] + " dir\n\n")
		fmt.Printf(
			`This will print the undecorated bug directory to stdout, 
so you can use it as a subcommand for arguments to any 
arbitrary shell commands. For example "cd $(bug dir)"

"%s dir" is an alias for "%s pwd"
`, os.Args[0], os.Args[0])
	case "tag":
		fmt.Printf("Usage: " + os.Args[0] + " tag [--rm] <BugID> <tag>...\n\n")
		fmt.Printf(`This will tag the given BugID with the tags
given as parameters. At least one tag is required.

Tags can be any string which would make a valid file name.

If the --rm option is provided before the BugID, all tags provided will
be removed instead of added.
`)
	case "roadmap":
		fmt.Printf("Usage: " + os.Args[0] + " roadmap [options]\n\n")
		fmt.Printf(
			`This will print a markdown formatted list of all open
issues, grouped by milestone.

Valid options are:
    --simple      Don't show anything other than the title in the output
    --no-status   Don't show the status of an issue
    --no-priority Don't show the priority of an issue
    --no-identifier Don't include the bug identifier of an issue
    --tags        Include the tags attached to a bug in it's output

    --filter tag           Only show bugs matching tag
    --filter tag1,tag2,etc Only show issues matching at least one of
                           the supplied tags

`)
	case "id", "identifier":
		fmt.Printf("Usage: " + os.Args[0] + " id <BugID> [--generate-id] <value>\n\n")
		fmt.Printf(
			`This will either set of retrieve the identifier for the bug
currently identified by BugID.

If value is provided as an argument, the bug identifier will be set
to the value passed in. You should take care to ensure that any
identifier used has at least 1 non-numeric character, to ensure there
are no conflicts with automatically generated issue numbers used for
a bug that has no explicit identifier set.

If the --generate-id option is passed instead of a static value, a
short identifier will be generated derived from the issue's current
title (however, the identifier will remain unchanged if the bug's title
is changed.)

If only a BugID is provided, the current identifier will be printed.

"%s id" is an alias for "%s identifier"
`, os.Args[0], os.Args[0])
	case "about", "version":
		fmt.Printf("Usage: " + os.Args[0] + " version\n\n")
		fmt.Printf(
			`This will print information about the version of %s being
invoked.

"%s about" is an alias for "version".
`, os.Args[0], os.Args[0])
	case "identifiers", "ids":
		fmt.Printf(
			`Bugs can be referenced in 2 ways on the commandline, either by
an index of where the bug directory is located inside the issues
directory, or by an ID. "BugID" can be either of these,
and %s will try to intelligently guess which your command is
referencing.

By default, no IDs are set for an issue. This means that
the issue number provided in "%s list" is an index into the directory
sorted by filesystem directory modified time.  It is not completely
stable as bugs are created, modified, and closed. The numbers are
easy to reference and remember in the short term.

If you have longer lasting issues that need a stable ID,
they can be created by "%s id <BugID> <New ID>".
You can then use <New ID> to reference the issue.

There are no rules for what constitutes a valid id, but
you should try and ensure that they have at least 1 non-numeric
character so that they don't conflict with directory indexes.

If you just want an id but don't care what it is, you
can use "%s id BugID --generate-id" to generate a new
ID for BugID.

If there are no exact matches for the BugID provided, %s commands will
also try and look up the bug by a substring match on all the valid 
IDs in the system before giving up.
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	case "import":
		fmt.Printf("Usage: %s import <--github|--be> <repo>\n\n", os.Args[0])
		fmt.Printf(
			`This will read from github <user/repository> issues 
or a local BugsEverywhere bug database to the issues/ directory.

Either "--github user/repo" is required to import GitHub issues
or "--be" is required to import a BugsEverywhere database
in the current directory.
`)
	case "help":
		fallthrough
	default:
		fmt.Printf("Usage: " + os.Args[0] + " <command>\n")

		fmt.Printf("\nUse \"bug help <command>\" or \"bug <command> help\" for\n")
		fmt.Printf("more information about any command below.\n")
		PrintVersion()
		fmt.Printf("\n\nStatus/reading commands:\n")
		fmt.Printf("\tlist       List existing bugs\n")
		fmt.Printf("\tfind       Search bugs for a tag, status, priority, or milestone\n")
		fmt.Printf("\tenv        Show settings that bug will use if invoked from this directory\n")
		fmt.Printf("\tpwd        Prints the issues directory to stdout (useful subcommand in the shell)\n")
		fmt.Printf("\tversion    Print the version of this software\n")
		fmt.Printf("\thelp       Show this screen\n")

		fmt.Printf("\nIssue editing commands:\n")
		fmt.Printf("\tcreate     File a new bug\n")
		fmt.Printf("\tedit       Edit an existing bug\n")
		fmt.Printf("\ttag        Tag a bug\n")
		fmt.Printf("\tid         Set a stable identifier for the bug\n")
		fmt.Printf("\tretitle    Rename the title of a bug\n")
		fmt.Printf("\tclose      Delete an existing bug\n")
		fmt.Printf("\tstatus     View or edit a bug's status\n")
		fmt.Printf("\tpriority   View or edit a bug's priority\n")
		fmt.Printf("\tmilestone  View or edit a bug's milestone\n")
		fmt.Printf("\timport     Create local bugs from a github repository\n")

		fmt.Printf("\nVersion control commands:\n")
		fmt.Printf("\tcommit     Commit any new, changed or deleted bug to git\n")
		fmt.Printf("\tpurge      Remove all issues not tracked by git\n")

		fmt.Printf("\nOther commands:\n")
		fmt.Printf("\troadmap    Print list of open issues sorted by milestone\n")
	}
}
