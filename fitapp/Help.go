package fitapp

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
	//fmt.Printf("cmd: %v\n", cmd) // debug
	switch cmd {
	case "create", "add", "new":
		fmt.Printf("usage: " + os.Args[0] + " create [-n] [options] Issue Title\n\n")
		fmt.Printf(
			`This will create an issue with the title Issue Title.  An editor 
will be opened on Description automatically.  If your EDITOR
environment variable is set, it will be used, otherwise
the default editor is vim.

If the first argument to create is "-n", then %s will not open 
any editor and create an empty Description.

Options take a value and set a field on the issue at the same
time as creating it. Valid options are:
    --status     Sets status to the next argument
    --tag        Adds tag on creation
    --priority   Sets the priority to the next argument
    --milestone  Sets the milestone to the next argument
    --identifier Sets the identifier to the next argument
    --generate-id Automatically generate a stable issue identifier

aliases for create: add new
`, os.Args[0])
	case "list", "view", "show", "display", "ls":
		fmt.Printf("usage: " + os.Args[0] + " list \n")
		fmt.Printf("       " + os.Args[0] + " list <IssueID>...\n")
		fmt.Printf("       " + os.Args[0] + " list <-m|--match> <regex>...\n")
		fmt.Printf("       " + os.Args[0] + " list <-t|--tags> <IssueID>...\n")
		fmt.Printf("       " + os.Args[0] + " list <tag>...\n\n")
		fmt.Printf("       " + os.Args[0] + " list <-r|--recursive>...\n")
		fmt.Printf(
			`This will list the issues found in the current environment

With no arguments, issue number and titles will be printed.
Issue numbers can reference this issue on the command line.

The [-m|--match] option tells list you are providing a regular
expression. Matching issues are listed.

All unix shell arguments that contain special characters, which many
regular expressions use, must be escaped. This prevents the automatic
"filename expansion" or "globbing" performed by the shell before
launching the fit command. A good references with details are at:
http://tldp.org/LDP/abs/html/globbingref.html
https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html

If valid IssueIDs are provided, whole issues with Description
will print.  See "fit help ids" for what makes a IssueID.

If 1 or more <tag>s are provided, matching issues are listed.

Note that IssueIDs may change as you create, edit, and close other
issues. Details are provided by "fit help ids."

The [-r|--recursive] option lists matching issues in subdirectories.

aliases for list: view show display ls
`)

	case "edit":
		fmt.Printf("usage: " + os.Args[0] + " edit <IssueID> <Filename>\n\n")
		fmt.Printf(
			`This will launch your standard editor to edit the Description 
of the issue identified by IssueID.  See "fit help ids" for
what makes an IssueID.

If the Filename option is provided, fit will instead launch an editor
to edit that file name within the issue directory. Files that have
special meaning (Status, Milestone, Priority, Identifier) are treated 
in a case insensitive manner, otherwise the filename is passed directly
to your editor.
`)
	case "status":
		fmt.Printf("usage: " + os.Args[0] + " status <IssueID> <NewStatus>\n\n")
		fmt.Printf(
			`This will edit or display the status of the issue identified by IssueID.
See "fit help ids" for what constitutes a IssueID.
            
If NewStatus is provided, it will update the first line of the Status file
for the issue (creating the file as necessary). If not provided, it will 
display the first line of the Status file to STDOUT.

Note that you can edit the status in your standard editor with the
command "%s edit status <IssueID>". If you provide a longer than 1 line
status with "fit edit status", "fit status" will preserve everything
after the first line when editing a status. You can use this to provide
further context on a status (for instance, why that status is setup.)
`, os.Args[0])
	case "priority":
		fmt.Printf("usage: " + os.Args[0] + " priority <IssueID> <New Priority>\n\n")
		fmt.Printf(
			`This will edit or display the priority of IssueID. See "fit help ids"
for what constitutes an IssueID.

By convention, priorities should be an integer number (higher is more 
urgent), but that is not enforced by this command and <New Priority> can
be any free-form text if you prefer.
            
If <New Priority> is provided, it will update the first line of the Priority
file for the issue (creating the file as necessary). If not provided, it 
will display the first line of the Priority file to STDOUT.

Note that you can manually edit the Priority file in the issues/ directory
by running "%s edit priority <IssueID>", to provide further explanation (for 
instance, why that priority is set.) This command will preserve the 
explanation when updating a priority.
`, os.Args[0])
	case "milestone":
		fmt.Printf("usage: " + os.Args[0] + " milestone <IssueID> <New Milestone>\n\n")
		fmt.Printf(
			`This will edit or display the milestone of the identified by IssueID.
See "%s help ids" for what constitutes an IssueID.

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
milestone is set) with the command "fit edit milestone <IssueID>"

This command will preserve the explanation when updating a priority.
`, os.Args[0], os.Args[0])
	case "retitle", "mv", "rename", "relabel":
		fmt.Printf("usage: " + os.Args[0] + " retitle <IssueID> <New Title>\n\n")
		fmt.Printf(
			`This will change the title of IssueID to <New Title>. Use this
to rename an issue.

aliases for retitle: mv rename relabel
`)
	case "rm", "close":
		fmt.Printf("usage: " + os.Args[0] + " close <IssueID>\n\n")
		fmt.Printf(
			`This will delete the issue identifier by IssueID. See
"help ids" for details on what constitutes a IssueID.

Note that closing an issue may cause existing IssueIDs to change if
they do not have a stable id set (see "help ids",
again.)

Also note that this does not remove the issue from git, but only 
from the file system. You'll need to execute "fit commit" to
remove the issue from version control.

alias for close: rm
`)
	case "find":
		fmt.Printf("usage: " + os.Args[0] + " find tag <value1> [value2 ...]\n")
		fmt.Printf("usage: " + os.Args[0] + " find status <value1> [value2 ...]\n")
		fmt.Printf("usage: " + os.Args[0] + " find priority <value1> [value2 ...]\n")
		fmt.Printf("usage: " + os.Args[0] + " find milestone <value1> [value2 ...]\n\n")
		fmt.Printf(
			`This will search all issues for multiple tags, statuses, priorities, or milestone.
The matching issues will be printed.
`)
	case "purge":
		fmt.Printf("usage: " + os.Args[0] + " purge\n\n")
		fmt.Printf(
			`This will delete any issues that are not currently tracked by
git or hg.
`)
	case "twilio":
		fmt.Printf("usage: " + os.Args[0] + " twilio\n\n")
		fmt.Printf(
			`This will send via twilio notifications of modified issues.
`)
	case "commit", "save":
		fmt.Printf("usage: " + os.Args[0] + " commit [--no-autoclose]\n\n")
		fmt.Printf(`This will commit any new, modified, or removed issues to
git or hg.

Your working tree and staging area should be otherwise
unaffected by using this command.

If the --no-autoclose option is passed to commit, fit will
not include a "Closes #x" line for each issue imported from
"fit import --github." Otherwise, the commit message will
include the list of issues that were closed so that GitHub
will autoclose them when the changes are pushed upstream.

alias for commit: save
`)
	case "env":
		fmt.Printf("usage: " + os.Args[0] + " env\n\n")
		fmt.Printf(`This will print the environment variables used by the command to stdout.

Use this command if you want to see settings what directory "fit create" is
using to store issues, or what editor will be invoked for a create/edit.
`)

	case "pwd", "dir", "cwd":
		fmt.Printf("usage: " + os.Args[0] + " pwd\n\n")
		fmt.Printf(
			`This will print the issue directory to stdout, 
so you can use it as a subcommand for arguments to any 
arbitrary shell commands. For example "cd $(fit pwd)"

aliases for pwd: dir cwd
`)
	case "tag":
		fmt.Printf("usage: " + os.Args[0] + " tag [--rm] <IssueID> <tag>...\n\n")
		fmt.Printf(`This will tag the given IssueID with the tags
given as arguments. At least one tag is required.

Tags can be any string which would make a valid file name.

If the --rm option is provided before the IssueID, all tags provided will
be removed instead of added.
`)
	case "roadmap":
		fmt.Printf("usage: " + os.Args[0] + " roadmap [options]\n\n")
		fmt.Printf(
			`This will print a markdown formatted list of all open
issues, grouped by milestone.

Valid options are:
    --simple      Don't show anything other than the title in the output
    --no-status   Don't show the status of an issue
    --no-priority Don't show the priority of an issue
    --no-identifier Don't include the issue identifier of an issue
    --tags        Include the tags attached to an issue in it's output

    --filter tag           Only show issues matching tag
    --filter tag1,tag2,etc Only show issues matching at least one of
                           the supplied tags

`)
	case "id", "identifier":
		fmt.Printf("usage: " + os.Args[0] + " id <IssueID> [--generate-id] <value>\n\n")
		fmt.Printf(
			`This will either set of retrieve the identifier for the issue
currently identified by IssueID.

If value is provided as an argument, the issue identifier will be set
to the value passed in. You should take care to ensure that any
identifier used has at least 1 non-numeric character, to ensure there
are no conflicts with automatically generated issue numbers used for
an issue that has no explicit identifier set.

If the --generate-id option is passed instead of a static value, a
short identifier will be generated derived from the issue's current
title (however, the identifier will remain unchanged if the issue's title
is changed.)

If only a IssueID is provided, the current identifier will be printed.

alias for id: identifier
`)
	case "identifiers", "ids":
		fmt.Printf(
			`Issues can be referenced in 2 ways on the commandline, either by
an index of where the issue directory is located inside the issues
directory, or by an ID. "IssueID" can be either of these,
and %s will try to intelligently guess which your command is
referencing.

By default, no IDs are set for an issue. This means that
the issue number provided in "%s list" is an index into the directory
sorted by filesystem directory modified time.  It is not completely
stable as issues are created, modified, and closed. The numbers are
easy to reference and remember in the short term.

If you have longer lasting issues that need a stable ID,
they can be created by "%s id <IssueID> <New ID>".
You can then use <New ID> to reference the issue.

There are no rules for what constitutes a valid id, but
you should try and ensure that they have at least 1 non-numeric
character so that they don't conflict with directory indexes.

If you just want an id but don't care what it is, you
can use "%s id <IssueID> --generate-id".

If there are no exact matches for the IssueID provided, %s commands will
also try and look up the issue by a substring match on all the valid 
IDs in the system before giving up.
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	case "version", "about", "--version", "-v":
		fmt.Printf("usage: " + os.Args[0] + " version\n\n")
		fmt.Printf(
			`This will print information about the version of %s being
invoked.

aliases for version: about --version -v
`, os.Args[0])
	case "import":
		fmt.Printf("usage: " + os.Args[0] + " import <--github|--be> <repo>\n\n")
		fmt.Printf(
			`This will read from github <user>/<repository> issues 
or a local BugsEverywhere bug database to the issues/ directory.

Either "--github <user>/repo>" is required to import GitHub issues
or  "--github <user>/<repo>/projects/<num>" to import a GitHub project
or "--be <path>" is required to import a local BugsEverywhere database.
GitHub projects require a configured GithubPersonalAccessToken value.
`)
		//ids aliases: idlist idsassigned identifiers
		//noids alias: noidentifiers
		//tagslist aliases: tagsassigned tags
		//notags
	case "help", "--help", "-h":
		fallthrough
	default:
		fmt.Printf("usage: " + os.Args[0] + ` help <command>

fit manages plain text issues with git or hg.
Use "fit help <command>" or "fit <command> help" for
    more information about any command below.
`)
		PrintVersion()
		fmt.Printf(`
Commands for status/reading:
    list       List issues
    find       Search for tag of fields: id, status, priority, or milestone
    tagslist   List assigned tags
    notags     List issues without tags
    ids        List stable identifiers
    noids      List issues without stable identifiers
    env        Show settings used when invoked from this directory
    pwd        Print the issues directory
    help       Show this screen
    version    Print the version of this software

Commands for editing:
    create     Open new issue
    edit       Edit an issue
    retitle    Rename an issue
    close      Delete an issue
    tag        Tag an issue
    id         View or set a stable identifier
    status     View or set status
    priority   View or set priority
    milestone  View or set milestone
    import     Download from github or bugseverywhere repository

Commands for version control:
    commit     Commit any new, changed or deleted issues
    purge      Remove all issues not tracked

Commands for processing:
    roadmap    Print list of open issues sorted by milestone

aliases for help: --help -h

`)
	}
}
