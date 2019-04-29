# bug

[![GoDoc](https://godoc.org/github.com/grantbow/bug?status.svg)](https://godoc.org/github.com/grantbow/bug) [![Build Status](https://travis-ci.org/grantbow/bug.svg?branch=master)](https://travis-ci.org/grantbow/bug) [![Test Coverage](https://codecov.io/gh/grantbow/bug/branch/master/graphs/badge.svg)](https://codecov.io/gh/grantbow/bug) [![GoReportCard](https://goreportcard.com/badge/github.com/grantbow/bug)](https://goreportcard.com/report/github.com/grantbow/bug) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/fit-issue/community)

bug helps manage plain text issue directories working with git or mercurial.

** TOC:**
<!-- toc -->
- [Goal](#goal)
- [Getting Started](#getting-started)
  * [Layout](#layout)
  * [Installation](#installation)
  * [Configuration](#configuration)
  * [Hooks](#hooks)
  * [Sample Usage](#sample-usage)
- [History](#history)
- [Next Steps](#next-steps)
  * [Feedback](#feedback)

<!-- tocstop -->

## Goal

The goal is to use human readable issue directories and files similar to how an
organized person (see FAQ.md) would keep track of issues beyond text files or
spreadsheets. This program streamlines working with
[issues](https://en.wikipedia.org/wiki/Issue_tracking_system) and [version
control](https://en.wikipedia.org/wiki/Version_control)_

bug saves issues using Filesystem Issue Tracker (see FIT.md)
conventions/format. bug works well with both git and mercurial distributed
version control. The git features are now more advanced.

bug can be aliased as a git subcommand such as `bug` or `issue`.

## Getting Started

### Layout

Filesystem Issue Tracker (FIT.md) conventions/format are a set of suggestions for
storing issues, one directory/folder per issue with plain text file details.

An `issues/` directory holds one (descriptively titled) directory per issue.
The "Description" file is the only text needed providing the details. Optional
tag_key_value files assign meta data.

bug maintains the nearest `issues/` directory to your current working directory
or it's parent directories. There can be more than one. bug can commit (or
remove) issues from versioning. Unlike many other issue systems, bug issues
naturally branch and merge along with the rest of your versioned files.

### Installation

After you have [go installed](https://golang.org/doc/install), install the
latest version of bug with:

`GO111MODULE=on go install github.com/grantbow/bug`

Make sure `$GOPATH/bin` or `$GOBIN` are in your path (or copy
the "bug" binary somewhere that is.)

The environment variable GO111MODULE changes how your golang works by enabling
golang 1.11 module support required for this version of bug. The default in
golang 1.12 is still "auto" but golang 1.13 is expected to default to "on".

Working with bug and git via the command line can be simplified. You can run
bug as a git subcommand like `git bug` or `git issue`. This [chapter about git
aliases](https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases) describes how
to set them up. It is part of the Pro Git book available for free online.
Simply add the alias to your .gitconfig or edit your .gitconfig.

`git config --global alias.issue !/path/to/bug`

This adds to your $HOME/.gitconfig:

`[alias]
    issue = !/path/to/bug`

### Configuration

Settings can be read from an optional config file .bug.yml placed next to the
issues directory. Current options include:

    * DefaultDescriptionFile: string,
          Default is ""
          when doing bug {add|new|create}
          first copy this file name to Description
          recommended: issues/Default
    * ImportXmlDump: true or false, 
          Default is false.
          during import, save raw xml files
    * ImportCommentsTogether: true or false,
          Default is false.
          during import, commments save together as one file
          instead of one comment per file.
    * ProgramVersion: string
          Default is ""
          String appended to the version of the program to
          identify customization.
    * DescriptionFileName: string
          Default is "Description".
          The name is one of the few imposed limitations.
          This configuration allows overriding the name used in your system.
	* TagKeyValue: true or false
          Default is false.
          writes tags in tag_key_value format (true)
          rather than .../tag/key without values.
          tags in both forms are read automatically.
          
The bug implementation of FIT (FIT.md) is (almost) the simplest issue system that can
still work. It differs from other distributed, versioned, filesystem issue
tracking tools in several ways. Human readable plain text files are still
easily viewed, edited and understood. Standard tools are used and further
minimize context switching between systems. bug also supports multiple
`issues/` directories throughout the directory tree.

Other issue systems use databases, hidden directories or hidden branches. While
these may be useful techniques in certain circumstances they do not seem
necessary and obfuscate how to access the valuable data.

### Hooks

Event based automation can be added through git or mercurial. We created a
hooks directory and look forward to seeing what code teams use and contribute.
Adapting hooks for both git and hg would be appreciated.

### Sample Usage

If an environment variable named PMIT is set, that directory will be used to
create and maintain issues as an 'issues' directory, otherwise the bug command
will walk up the tree until it finds an "issues" subdirectory. Examples assume
you are already in a directory tracked by git. To get started simply `mkdir
issues`.

Example usage:

```
$ mkdir issues
$ bug help
Usage: bug <command> [options]

Use "bug help <command>" or "bug <command> help" for
more information about any command below.
bug version 0.4 built using go1.12.2

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
	tag        Tag a bug
	identifier Set a stable identifier for the bug
	retitle    Rename the title of a bug
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
(<your editor> Description)
(save)
Created issue: Need better help

$ bug list
Issue 1: Need better help

$ bug list 1
Title: Need better help

Description:
<the entered description>

$ bug create -n Need better formating for README
(no editor, default to empty Description)
Created issue: Need better formatting for README

$ bug list
Issue 1: Need better formating for README
Issue 2: Need better help
```

## History

bug is the program written in Go developed by Dave MacFarlane (driusan).
Filesystem Issue Tracker (FIT.md) is the new name for the Poor Man's Issue Tracker
(PMIT) storage system also developed by driusan. For his demo from 2016, see
[driusan's talk](https://www.youtube.com/watch?v=ysgMlGHtDMo) at the first
GolangMontreal.org conference, GoMTL-01.

## Next Steps

* [FIT.md](FIT.md)
* [FAQ.md](FAQ.md)
* [CONTRIBUTING.md](CONTRUBUTING.md)
* [CODE_OF_CONDUCT.md](CODE_OF_CONDICT.md)
* [SUPPORT.md](SUPPORT.md)
* https://github.com/grantbow/bug/wiki

Your system is the beginning, not the end. Much has been written about how to
use and setup systems to track or manage issues, software bugs, trouble
tickets, support tickets, incident ticket or requests.

How to Report Bugs Effectively by Simon Tatham
    https://www.chiark.greenend.org.uk/~sgtatham/bugs.html

While I hesitate to include this URL and while it uses very colloquial
terminology the content could help many people.
    http://www.catb.org/esr/faqs/smart-questions.html

### Feedback

We would like to hear about how you use this system.

I would like to work with others and would appreciate feedback at
grantbow+bug@gmail.com.

Since the original project is not very active I have gone ahead and continuted
some development on my fork. I encourage discussion. Submitting can be done
with a pull request to me, to our upstream project or using [git
remotes](https://stackoverflow.com/questions/36628859/git-how-to-merge-a-pull-request-into-a-fork).

Anyone thinking of CONTRIBUTING.md please do so. As this is an issue tracking
system a pull request with an issue seems logical enough and good working
practice.

As mentioned in SUPPORT.md questions are encouraged via email, issues or pull
requests for now.

The CODE_OF_CONDUCT.md is the standard offered by github and looks great.

