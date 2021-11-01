# fit
filesystem issue tracker: manages plain text issues with git or mercurial

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
  * [Hooks](#hooks)
  * [Example Script](#example-script)
- [History](#history)
- [Background](#background)
- [Next Steps](#next-steps)
  * [Feedback](#feedback)

<!-- tocstop -->

## Prerequisites

git or hg (mercurial)

golang 1.13 or higher

linux, mac or windows OS.

## Goal

Standard coding practices improve outcomes. Using fit minimizes switching
between coding and issue tracking systems which increases productivity,
maintains code context and stores issue histories.

The fit implementation is (almost) the simplest issue system that can still
work. See ([Background](#background)) for an explanation.

## Getting Started

### Layout

Filesystem Issue Tracker conventions/format ([Filesystem_Issues.md](docs/Filesystem_Issues.md)) are a set of
suggestions for storing issues, one directory/folder per issue with plain text
file details.

A `fit/` directory holds one (descriptively titled) directory per issue.
The "Description" file is the only fiel needed in each directory. Optional
tag_key_value files assign meta data. A minimal issue looks like:

    fit/name_of_issue/Description

One fit directory may be used for a whole git/mercurial repository or
subfolders may contain their own fit directory.

fit maintains the nearest `fit/` directory to your current working directory or
it's parent directories. fit can commit (or remove) issues from versioning or
this can be done manually without fit. Unlike many other issue systems, fit
issues naturally branch and merge along with the rest of your versioned files.

Some support is available to import and/or reference other issue trackers.
Usage reports via email or gitter are encouraged.

### Example Use

To get started in the top of an existing git repo simply `mkdir fit`.
Then `mkdir fit/<issue_name>` and edit `fit/issue_name/Description`.

`fit list` shows your directory or repo's issues.

Add and commit the Description file like any other file in your repository.

If an environment variable named FIT is set that value will be used as a
directory name used to find the 'fit' or 'issues' directory instead of your
present working directory. All fit commands use the FIT environment variable
if present.

If a 'fit' directory/folder is not found fit will walk up your filesystem tree
until it finds a "fit" subdirectory similar to how git looks for
.git or hg looks for .hg. A warning is provided if no directory is found.

fit uses subcommands like git. For a list of commands use `fit help`

### Installation

After you have [go installed](https://golang.org/doc/install) make sure you
have `/usr/local/go/bin` (or equivalent) and `$HOME/go/bin` in your path and
`GOPATH` is set to something like `$HOME/go`. Install the
latest version of fit with:

`GO111MODULE=on go install github.com/grantbow/fit`

If that does not work in one command then:
```
    $ export GO111MODULE=on
    $ mkdir -p $GOPATH/src/github.com/grantbow/fit
    $ git clone https://github.com/grantbow/fit $GOPATH/src/github.com/grantbow/fit
    $ cd $GOPATH/src/github.com/grantbow/fit/cmd/fit
    $ go install
```

This will create the binary `$GOPATH/src/github.com/grantbow/fit/cmd/fit/fit(.exe)`
and move it to `$GOPATH/bin/fit(.exe)`

Make sure `$GOPATH/bin` or `$GOBIN` are in your path or you can copy
the "fit" binary somewhere that is in your path.

The environment variable set using `export GO111MODULE=on` changes how old
golang versuibs work by enabling golang 1.11+ module support required by fit.
The defaults in golang 1.13 and 1.14 and 1.15 were still "auto".
The defaults in golang 1.16 and 1.17 are "on" so this setup is no longer
required any more.

Working with fit and git via the command line can be simplified. You can run
fit as it's own command or as a git subcommand like `git fit`.
You can quickly add the alias to your .gitconfig:

```
 git config --global alias.fit \!/home/<user>/go/bin/fit`  
 git config --global alias.issue \!/home/<user>/go/bin/fit`  
 git config --global alias.bug \!/home/<user>/go/bin/fit`  
```

Note: cygwin users use !/cygdrive/c/Users/\<user\>/go/bin/fit.exe

This will add to your $HOME/.gitconfig or you can edit it manually:

```
[alias]  
    fit = !/home/<user>/go/bin/fit
    issue = !/home/<user>/go/bin/fit
    bug = !/home/<user>/go/bin/fit
```

This[chapter about git aliases](https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases) describes how
to set them up very well. It is part of the Pro Git book available for free
online. 

### Configuration

The environment variable EDITOR is used to execute your preferred editor
when needed.

An important choice is what to do with closed issues. They can be deleted
(the historical default), moved to a subdirectory "closed" or
add a tag\_status\_closed.

Settings can be read from .fit.yml next to the fit directory.
This is an optional config file. Defaults are backwards compatible with
the original bug program so far. Current options include:

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
While these may be useful techniques in certain circumstances this seems to
unnecessarily obfuscate access.

Tags have been significantly enhanced since they were originally implemented.
Using the above options the default behavior of boolean present/not present
file names in a "tags" subdirectory can instead be filenames like Status that
contain the values or even simple tag\_key\_value filenames with empty contents
or comment contents. The last option enables great flexibility. A few keys are
hard coded in the program with special features: Identifier, Priority, Status,
Milestone and Tag. Newer tag\_key\_value filenames are recommended.

As every bug system operates within the context of a number of people that use
the system many efforts to support as many choices of system use that
are reasonably possible. Comments and suggestions are welcomed. Pull requests
are even better but are not required to participate in this project.

### Hooks

Event based automation can be added through git or mercurial. We created a
hooks directory and look forward to seeing what teams use and contribute.
Work to help adapt hooks to both git and hg are appreciated.

### Example

```
$ mkdir foo && cd foo
$ git init
$ mkdir fit
$ fit help
cmd: help
Usage: help <command>

Use "fit help <command>" or "fit <command> help" for
more information about any command below.
fit version 0.7 built using go1.17.2 GOOS linux
executable: -rwxrwxr-x 7931914 Sun Oct 31 23:43:08 PDT 2021 /home/grantbow/go/bin/fit

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

$ fit create Need better help
(<your editor> Description)
(save and quit)
Created issue: Need better help

$ fit list

===== list /.../foo/fit
Issue 1: Need better help

$ fit list 1

===== list /.../foo/fit
Title: Need better help
Description:
<the entered description>

$ fit create -n Need better formatting for README
(no editor launched, defaults to empty Description file)
Created issue: Need better formatting for README

$ fit list

===== list /.../foo/fit
Issue 1: Need better help
Issue 2: Need better formatting for README
```

## History

fit is the golang program first developed as "bug" by Dave MacFarlane (driusan).
Filesystem Issue Tracker ([Filesystem_Issues.md](docs/Filesystem_Issues.md)) is the new name for the Poor Man's
Issue Tracker (PMIT) storage system also first developed by driusan. See the
2016 demo video of [driusan's
talk](https://www.youtube.com/watch?v=ysgMlGHtDMo) at the first
GolangMontreal.org conference, GoMTL-01. The program and storage system have incrementally
evolved while trying to remain backward compatible. See the docs/[FAQ.md](docs/FAQ.md)
for even more information.

## Background

A limited but sufficient number of conventions with just enough organization
can quickly capture issues using human readable issue directories and files.
fit can be the primary system if no other system is provided or supplement
other issue/bug systems to quickly capture issues and their context as close
to the code as possible.

Using fit helps implementers streamline working with
[issues](https://en.wikipedia.org/wiki/Issue_tracking_system) and [version
control](https://en.wikipedia.org/wiki/Version_control). fit works with
both git and mercurial distributed version control though the git features are
more well exercised.

fit is designed to adapt to your processes using issue key/value pair metadata.

The fit tool manages issues using conventions/format of
Filesystem Issue Tracker (see [Filesystem_Issues.md](docs/Filesystem_Issues.md)). A `fit/` or `issues/`
directory holds one descriptively titled directory per issue. Each directory 
holds a file Description (name is configurable) which is a text file.
Issue directories hold anything else needed about the issue.

Issue systems typically evolve from the most simple systems that work to slightly
more complex systems that work better when working with others.

At first people may naturally try to keep track of issues in a single text
file and/or spreadsheet but these can fail to meet project needs.
(see docs/[FAQ.md](docs/FAQ.md))

Issue context is valuable to coders and may be difficult for others to
understand, especially without the context of the code they describe.
fit can support multiple `fit/` directories in a
repository's tree for stronger coordination of coding and issue tracking.

Projects in IT environments face all too common circumstances: implementers
may not be given the tools needed (or given bad tools) to record code issues.
Other issue systems typically take some resources to setup and maintain which
can be difficult to justify until long after a system is desperately needed
by implementers. Separate issue systems are often focused on higher volume
user facing streams of problem reports. These systems may or may not be meet
project needs to capture valuable implementation details so valuable details
are often poorly documented or completely lost. Some IT groups hope that
problems will not attract attention and see obfuscation as a way to reduce
perhaps already oversized workloads. It can be reasonable to focus valuable
project budget, time, scope, quality or other resources on new features but
code level or low volume streams of issues can be ignored. Beyond managing
contentious flat files or spreadsheets there is FIT.

Important issues can be captured, surfaced and addressed, whether they are
actual problems, questions, possible features or ideas by those most familiar
with the project. It is hoped that all code savvy project collaborators can
capture implementation details of varying importance quickly and easily using
fit compared to using larger, possibly distracting systems best designed for
other uses like code reviews or operational facing streams of issues. Regardless
of other available issue systems that might be available, using a fit system
might advantageously complement project workflows.

[Software Development Life Cycles](https://en.wikipedia.org/wiki/Software_development_process) (SDLCs) involve more than just the source code.
Over time needs of a project may change from hacking/coding, just getting something working,
to implementing more disciplined software engineering best practices. Code can
start small and grow as users, TODO comments, use cases and developers are added.
The FIT issue system was designed to adapt to each stage of needs.

While one issue set used for one git repository may be enough the use of
recursive fit directories are now supported. As complexity increases adding
multiple `fit/` directories in different parts of your git repo may help
project coders keep focused.

There are some choices each project can make for how to handle closed
issues. As the number of issues grows closed issues can simply be deleted or
an archive can hold the inactive issues. While deleting issues helps keep things
uncluttered issues still have value over time and may be difficult to find
using only version control history.

fit can be aliased as a git subcommand "git fit ..." It is intended that similar
subcommands perform similarly expected functions.

fit software is written using [golang](https://golang.org) to make things easy,
simple and reliable. Why go? [This video](https://vimeo.com/69237265) from a
2013 Ruby conference by Andrew Gerrand of Google seems a good explanation. It
was linked from the golang.org home page.

Engineers know that there is more to code than the source itself. For some rare
individuals the code is enough context. For most people new to a code base or
distracted by other concerns any recorded context can be extremely helpful.
Notes about refactoring history, code reviews or feature ideas can be
important to grok a code base more quickly. This context may originate from
researching a user reported problem or may arise any time while coding.

## Next Steps

* docs/[Filesystem_Issues.md](docs/Filesystem_Issues.md)
* docs/[FAQ.md](docs/FAQ.md)
* [CONTRIBUTING.md](CONTRUBUTING.md)
* [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
* [SUPPORT.md](SUPPORT.md)
* [wiki](https://github.com/grantbow/fit/wiki)
* [gitter](https://gitter.im/fit-issue/community)
* [SECURITY.md](SECURITY.md)

Your system is just the beginning, not the end. Much has been written about
how to use and setup systems to track or manage issues, software bugs, trouble
tickets, support tickets, incident tickets or requests. See the docs/FAQ.md

### Feedback

We would very much like to hear about how you use this system.

I would like to work with others and would appreciate feedback at
grantbow+fit@gmail.com.

Since the original bug project is not very active I have gone ahead and continuted
development. I encourage discussion. Submissions can be done with a
pull request or using [git remotes](https://stackoverflow.com/questions/36628859/git-how-to-merge-a-pull-request-into-a-fork).

Anyone thinking of [CONTRIBUTING.md](CONTRIBUTING.md) is encouraged to do so
and development guidelines are included in that file. As this is an issue
tracking system a pull request with an issue seems logical enough and a good
practice to exercise our own tool.

As mentioned in [SUPPORT.md](SUPPORT.md) questions are encouraged via email, issues or pull
requests.

The [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) is the standard recommended by github
offered by contributor-covenant.org.

As mentioned in [SECURITY.md](SECURITY.md) vulnerabilities are encouraged via email.
Security concerns are generally handled using standard git repository practices.
