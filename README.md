# Bug

[![GoDoc](https://godoc.org/github.com/driusan/bug?status.svg)](https://godoc.org/github.com/driusan/bug) [![Build Status](https://travis-ci.org/driusan/bug.svg?branch=master)](https://travis-ci.org/driusan/bug) [![Test Coverage](https://codecov.io/gh/driusan/bug/branch/master/graphs/badge.svg)](https://codecov.io/gh/driusan/bug)

bug writes code problem reports to plain text files.

bug requires Go version 1.9 or greater.

Bug is an implementation of a distributed issue tracker using
git (or mercurial/hg) to manage issues on the filesystem following [poor man's
issue tracker](https://github.com/driusan/PoormanIssueTracker) conventions.

Though issues can be edited and saved to the issues directory of a git or hg
repository bug helps create, read, update and delete them.

Using Bug with git you can run Bug as a git subcommand like `git bug` or `git
issue`. As described in this [chapter about
Aliases](https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases) of the Pro Git
book available online. Simple add it to your .gitconfig manually or:

`git config --global alias.issue '!bug'`

The goal is to use the filesystem in a human readable way, similar to
how an organized person without any bug tracking software might, 
by keeping track of bugs in an `issues/` directory, one (descriptive)
subdirectory per issue. bug provides a tool to maintain the nearest 
`issues/` directory to your current working directory and provides hooks 
to commit (or remove) the issues from source control.

An optional config file next to the issues directory named .bug.yml may
specify options. Current options include:
    * DefaultDescriptionFile: string,
          create bug template file name
    * ImportXmlDump: true or false, 
          during import, save raw xml files
    * ImportCommentsTogether: true or false,
          during import, commments save together as one file
          instead of one comment per file.

Some other distributed bug tracking tools store a database in a hidden
directory. A separation of bug related information and bug fix code seems
artifical, unnecessary and/or even harmful to clearly understanding issues and
resolutions. Using files allows easily viewing, editing and understanding
issues even without access to the bug tool. We hope you find the bug program
streamlines your process.

A beneficial consequence of file level storage is allowing multiple `issues/`
directories at different places in your directory tree to, for instance, keep
separate bug repositories for different submodules or packages contained in a
single git repository.

Because issues are stored as human readable plain text files, they branch
and merge along with the rest of your code, and you can resolve conflicts 
using your standard tools.

For a demo, see my talk at [GoMTL-01](https://www.youtube.com/watch?v=ysgMlGHtDMo)

# Installation

If you have go installed, install the latest released version with:

`go get github.com/driusan/bug`

Make sure `$GOPATH/bin` or `$GOBIN` are in your path (or copy
the "bug" binary somewhere that is.)

Otherwise, you can download a 64-bit release for OS X or Linux on the 
[releases](https://github.com/driusan/bug/releases/) page.

(The latest development version is on the latest v0.x-dev branch)

# Example Usage

If an environment variable named PMIT is set, that directory will be
used to create and maintain issues, otherwise the bug command will
walk up the tree until it finds somewhere with a subdirectory named
"issues".  Examples assume you are already in a directory tracked by
git. To get started simply `mkdir issues`.

Example usage:

```
$ bug help
Usage: bug command [options]

Use "bug help [command]" or "bug [command] help" for
more information about any command below.

Valid Commands

Status/reading commands:
	list       List existing bugs
	find       Search bugs for a tag, status, priority, or milestone
	env        Show settings that bug will use if invoked from this directory
	pwd        Prints the issues directory to stdout (useful subcommand in the shell)
	version    Print the version of this software
	help       Show this screen

Issue editing commands:
	create     File a new bug
	edit       Edit an existing bug
	tag        Tag a bug with a category
	identifier Set a stable identifier for the bug
	relabel    Rename the title of a bug
	close      Delete an existing bug
	status     View or edit a bug's status
	priority   View or edit a bug's priority
	milestone  View or edit a bug's milestone
	import     Create local bugs from a github repository

Source control commands:
	commit     Commit any new, changed or deleted bug to git
	purge      Remove all issues not tracked by git

Other commands:
	roadmap    Print list of open issues sorted by milestone

$ bug create Need better help
(Your editor opens here to enter a description)

$ bug list
Issue 1: Need better help

$ bug list 1
Title: Need better help

Description:
The description that I entered

$ bug purge
Removing issues/Need-better-help

$ bug create -n Need better formating for README

$ bug list
Issue 1: Need better formating for README

$ bug commit
$ git push
```

# Feedback

Currently, there aren't enough users to set up a mailing list, but 
I'd nonetheless appreciate any feedback at driusan+bug@gmail.com. 

You can report any bugs either by email, via GitHub issues, or by sending
a pull request to this repo.
