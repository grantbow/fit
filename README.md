# fit
filesystem issue tracker: manages plain text issues with git or hg 

[![GoDoc](https://godoc.org/github.com/grantbow/fit?status.svg)](https://godoc.org/github.com/grantbow/fit) [![Build Status](https://travis-ci.com/grantbow/fit.svg?branch=master)](https://app.travis-ci.com/github/grantbow/fit) [![Test Coverage](https://codecov.io/gh/grantbow/fit/branch/master/graphs/badge.svg)](https://codecov.io/gh/grantbow/fit) [![GoReportCard](https://goreportcard.com/badge/github.com/grantbow/fit)](https://goreportcard.com/report/github.com/grantbow/fit) [![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2820/badge)](https://bestpractices.coreinfrastructure.org/projects/2820) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/fit-issue/community)

**TOC:**

<!-- toc -->

- [Prerequisites](#prerequisites)
- [Goal](#goal)
- [Getting Started](#getting-started)
  * [Layout](#layout)
  * [Example Use](#example-use)
  * [Installation](#installation)
  * [Configuration](#configuration)
  * [Example](#example)
- [Hooks](#hooks)
- [Governance](#governance)
- [Feedback](#feedback)
- [Next Steps](#next-steps)

<!-- tocstop -->
<!-- older topics
# fit init, fit open, fit close, fit archive, fit list
# what's an issue
# fit codecov, 3 x scan issues,
# fit summary
# characteristics
-->

## Prerequisites

The only prerequisite is a filesystem. Issues can easily be created with new
directories and a text editor.

To list and manage the issues with the fit tool you need:
- git or hg (mercurial)
- go1.19
- any operating system [supported by go](https://go.dev/doc/install/source) and
  [supported by
  git](https://git.wiki.kernel.org/index.php/Interfaces,_frontends,_and_tools)
  including Windows, MacOS, Linux, \*BSD, Android, etc..

## Goal

Capture issues fast.

fit can be a project's primary issue system or a secondary system for code related issues.

Standard coding practices improve project outcomes. Using fit minimizes
switching between coding and issue tracking systems which increases
productivity and maintains code context.

The fit implementation is (almost) the simplest issue system that can still
work. The intent is to make fit conventions natural for users of
[git](https://en.wikipedia.org/wiki/Git) and
[golang](https://en.wikipedia.org/wiki/Go_(programming_language)). Contention
inherent in simpler issue systems is minimized and maintenance inherent in more
complex issue systems is avoided. See ([Background](#background)) for an
explanation.

## Getting Started

### Layout

Filesystem Issue Tracker conventions/format
([Filesystem_Issues.md](docs/Filesystem_Issues.md)) store issues, one
directory/folder per issue, inside a "fit" directory. That document also
contains information about what a good issue looks like. Inside each issue
directory is a plain text "Description" file. 

Optional tag\_<key>\_<value> files assign meta data. A minimal issue looks like:

    fit/better_docs/Description
    fit/better_docs/tag_id_1

The fit directory is located at the top of a git/mercurial repository and
subfolders may contain their own fit directory.

The config file is named `.fit.yml`,

The EDITOR environment variable is the editor used by default.

fit maintains the nearest `fit/` directory to your current working directory or
it's parent directories.

Setting an environment variable FIT overrides the default search for a location.

Unlike other issue systems, fit
issues naturally branch and merge along with the rest of your versioned files.

Some support is available to import and/or reference other issue trackers.
Usage reports via email or gitter are encouraged.

### Example Use

To get started in the top of an existing git repo simply
`mkdir -p fit/<issue_name>` and edit `fit/<issue_name>/Description` with
your editor set in the EDITOR environment variable.

`fit list` shows your issues.

Add and commit the Description file like any other file in your repository.

`git add fit/<issue_name>/Description && git commit -m "first issue"`

If an environment variable named FIT is set that value will be used as a
directory name used to find the 'fit' or 'issues' directory instead of your
present working directory. All fit commands use the FIT environment variable
if present.

If a 'fit' directory/folder is not found and you enable recursion
fit will walk up your filesystem tree
until it finds a "fit" subdirectory similar to how git looks for
.git or hg looks for .hg. A warning is provided if no directory is found.

fit uses subcommands like git. For a list of commands use `fit help`

### Installation

Briefly, you need git and [go installed](https://golang.org/doc/install).

`go install github.com/grantbow/fit@latest`

You can run fit as it's own command or as a git subcommand.

For details see ([INSTALL.md](INSTALL.md)). The software was developed with go1.13 but only the current and one previous version are supported before being declared [end of life](https://go.dev/doc/devel/release#policy).

### Configuration

The environment variable EDITOR is used to execute your preferred editor
when needed.

An important choice is what to do with closed issues. They can be deleted
(the historical default), moved to a subdirectory "closed" or
add a tag\_status\_closed.

Settings can be read from .fit.yml next to the fit directory.  Defaults are
backwards compatible with the original bug program so far. Current options
include:

    * DescriptionFileName: string
          Default is "Description".
          The name is one of the few imposed limitations.
          This configuration allows overriding the name used in your system.
    * DefaultDescriptionFile: string,
          Default is ""
          when doing fit {add|new|create}
          first copy this file name to Description
          recommended: fit/DescriptionDefault.txt
    * ImportXmlDump: true or false, 
          Default is false.
          during import, save raw xml files
    * ImportCommentsTogether: true or false,
          Default is false.
          during import, commments save together as one file
          instead of one comment per file.
    * ProgramVersion: string
          Default is "". Set to 0.6 when run. This
          string is appended to identify local customizations.
    * TagKeyValue: true or false
          Default is false.
          writes tags in tag_key_value format (true)
          rather than .../tag/key without values.
          tags in both forms are read automatically.
    * NewFieldAsTag: true or false
          Default is false.
          writes fields in tag_key_value format (true)
          rather than .../Key file with value in the text.
          fields in both forms are read automatically.
    * NewFieldLowerCase: true or false
          Default is false.
          writes fields as tag_key_value format (true)
          rather than .../tag_Key_value.
          fields in both forms are read automatically.
    * GithubPersonalAccessToken: string
          Default is empty.
          Set one at github.com/settings/tokens
          for import of private projects that need authentication
    * TwilioAccountSid: string
          Default is empty.
          Needed for twilio use.
    * TwilioAuthToken: string
          Default is empty.
          Needed for twilio use.
    * TwilioPhoneNumberFrom: string
          Default is empty.
          Needed for twilio use.
    * FitSite: string
          Default is empty.
          base url used in notifications
    * MultipleDirs: true or false
          Default is false.
          Set to always look recursive.
    * CloseStatusTag: true or false
          Default is false.
          fit close will not add a tag (false) or
          add tag_status_closed (true)
    * CloseMove: true or false
          Default is false.
          fit close will use ClosePreventDelete (false) or
          fit close will move issues to CloseDir (true) and
              sets ClosePreventDelete=true
    * ClosedDirName: string
          Default is "closed"
          this is the name used inside the fit directory
    * ClosePreventDelete: true or false
          Default is false.
          fit close will delete (false) or
          fit close will not delete (true)
              implies CloseStatusTag and/or CloseMove
    * IdAbbreviate: true or false
          Default is false.
          Use Identifier. True uses Id.
    * IdAutomatic: true or false
          Default is false.
          Generates an Identifier.

Other issue systems may use databases, hidden directories or hidden branches.
While these may be useful techniques in certain circumstances they obfuscate
access to the data.

Tags have been significantly enhanced since they were originally implemented.
Using the above options the default behavior of boolean present/not present
file names in a "tags" subdirectory can instead be filenames like Status that
contain the values or even simple tag\_key\_value filenames with empty contents
or comment contents. The last option enables great flexibility. A few keys are
hard coded in the program with special features: Identifier, Priority, Status,
Milestone and Tag. Newer tag\_key\_value filenames are recommended.

Comments and suggestions are welcomed. Pull requests are even better but are
not required to participate in this project.

### Example

```
$ mkdir foo && cd foo
$ git init
$ mkdir fit
$ fit help
usage: fit help <command>

fit manages plain text issues with git or hg.
Use "fit help <command>" or "fit <command> help" for
    more information about any command below.
fit version 0.7 built using go1.19 GOOS windows
executable: -rw-rw-rw- 8411136 Tue Sep  6 04:44:31 PDT 2022 fit.exe

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

$ fit create better docs
(<your editor> Description)
(save and quit)
Created issue: better docs

$ fit list

===== list /.../foo/fit
Issue 1: better docs

$ fit list 1

===== list /.../foo/fit
Title: better docs
Description:
<the entered description>

$ fit create -n better README formatting
(no editor launched, defaults to empty Description file)
Created issue: better README formatting

$ fit list

===== list /.../foo/fit
Issue 1: better help
Issue 2: better README formatting
```

## Hooks

Event based automation can be added through git or mercurial hooks. Our hooks
directory contains some examples. We look forward to seeing what teams use and
contribute. Adapting hooks to both git and hg would be appreciated.

## Governance

There are two key roles in the fit project: project owner & lead contributors.
Collaboration is encouraged. See docs/FAQ.md for details.

## Feedback

We would very much like to hear about how you use this system.

I would like to work with others and would appreciate feedback at
grantbow+fit@gmail.com.

Since the original bug project is not very active I have gone ahead and continuted
development. I encourage discussion. Submissions can be done with a
pull request or using [git remotes](https://stackoverflow.com/questions/36628859/git-how-to-merge-a-pull-request-into-a-fork).

## Next Steps

* docs/[Filesystem_Issues.md](docs/Filesystem_Issues.md)
* docs/[FAQ.md](docs/FAQ.md)
* docs/[Background.md](Background.md)
* [CONTRIBUTING.md](CONTRIBUTING.md)
* [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) is the standard recommended by github offered by contributor-covenant.org.
* [SUPPORT.md](SUPPORT.md)
* [wiki](https://github.com/grantbow/fit/wiki)
* [gitter](https://gitter.im/fit-issue/community)
* [SECURITY.md](SECURITY.md)

Your system is just the beginning, not the end. Much has been written about
how to use and setup systems to track or manage issues, software bugs, trouble
tickets, support tickets, incident tickets or requests. See the docs/FAQ.md

