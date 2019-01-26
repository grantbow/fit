# Bug

[![GoDoc](https://godoc.org/github.com/grantbow/bug?status.svg)](https://godoc.org/github.com/grantbow/bug) [![Build Status](https://travis-ci.org/grantbow/bug.svg?branch=master)](https://travis-ci.org/grantbow/bug) [![Test Coverage](https://codecov.io/gh/grantbow/bug/branch/master/graphs/badge.svg)](https://codecov.io/gh/grantbow/bug) [![GoReportCard](https://goreportcard.com/badge/github.com/grantbow/bug)](https://goreportcard.com/report/github.com/grantbow/bug)

bug manages plain text issue files.
bug works with git and mercurial distributed versioning.
See below for how bug can be aliased as a git subcommand `bug` or `issue`.

The goal is to use a filesystem in a human readable way similar to how (see
FAQ.md) an organized person would keep track of issues beyond text files or
spreadsheets. This program streamlines working with issues and version control.

An `issues/` directory holds one (descriptive) subdirectory per issue. bug
maintains the nearest `issues/` directory to your current working directory
(there can be more than one) and provides hooks to commit (or remove) issues
from versioning. Issues naturally branch and merge along with the rest of your
versioned files.

bug was started by Dave MacFarlane (driusan) and was extended from his
implementation of his "poor man's issue tracker" conventions. For his demo, see
[his talk](https://www.youtube.com/watch?v=ysgMlGHtDMo) at the first
GolangMontreal.org conference, GoMTL-01.

This fork of bug has a renamed path so it requires version 1.11 to properly
redirect building and testing the sub-modules using go.mod files. This is
currently the best way to handle module paths in golang.

# Feedback
I would like to work with others and would appreciate feedback at
grantbow+bug@gmail.com.

Since the original project is not very active I have gone ahead and published
code through this public fork. Instead of submitting a pull request to our
upstream getting code from others can be done using [git
remotes](https://stackoverflow.com/questions/36628859/git-how-to-merge-a-pull-request-into-a-fork).

You can report any bugs either by email or by sending a pull request.

# Configuration
An optional config file next to the closest issues directory named .bug.yml
may specify options. Current options include:
    * DefaultDescriptionFile: string,
          when doing bug {add|new|create}
          first copy this file name to Description
    * ImportXmlDump: true or false, 
          during import, save raw xml files
    * ImportCommentsTogether: true or false,
          during import, commments save together as one file
          instead of one comment per file.

bug is the (almost) simplest system that can still work. It differs from other
distributed, versioned, filesystem issue tracking tools in several ways.  Human
readable plain text files are still easily viewed, edited and understood
without access to the bug tool unlike a database or hidden directory based
system.  Standard tools are used and further minimize context switching between
systems.  bug also supports multiple `issues/` directories throughout the
directory tree.

For a demo, see my talk at [GoMTL-01](https://www.youtube.com/watch?v=ysgMlGHtDMo)

# Installation
If you have [go installed](https://golang.org/doc/install), install the latest version with:

`GO111MODULE=on go get github.com/grantbow/bug`

The environment variable enables golang 1.11 module support.

Make sure `$GOPATH/bin` or `$GOBIN` are in your path (or copy
the "bug" binary somewhere that is.)

Using Bug with git you can run Bug as a git subcommand like `git bug` or `git
issue`. As described in this [chapter about
Aliases](https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases) of the Pro Git
book available online. Simple add it to your .gitconfig manually or:

`git config --global alias.issue !/path/to/bug`

This adds to your .gitconfig:

`[alias]
    issue = !/path/to/bug`

# Sample Usage
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
