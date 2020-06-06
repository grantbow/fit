# fit
manages plain text issues with git or mercurial

[![GoDoc](https://godoc.org/github.com/grantbow/fit?status.svg)](https://godoc.org/github.com/grantbow/fit) [![Build Status](https://travis-ci.org/grantbow/fit.svg?branch=master)](https://travis-ci.org/grantbow/fit) [![Test Coverage](https://codecov.io/gh/grantbow/fit/branch/master/graphs/badge.svg)](https://codecov.io/gh/grantbow/fit) [![GoReportCard](https://goreportcard.com/badge/github.com/grantbow/fit)](https://goreportcard.com/report/github.com/grantbow/fit) [![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2820/badge)](https://bestpractices.coreinfrastructure.org/projects/2820) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/fit-issue/community)

**TOC:**

<!-- toc -->

- [Goal](#goal)
- [History](#history)
- [Background](#history)
- [Getting Started](#getting-started)
  * [Layout](#layout)
  * [Installation](#installation)
  * [Configuration](#configuration)
  * [Hooks](#hooks)
  * [Example Use](#example-use)
- [Next Steps](#next-steps)
  * [Feedback](#feedback)

<!-- tocstop -->

## Goal

Standard coding practices improve project outcomes. Using as few tools and
systems as possible increases productivity. Using fit minimizes context
switching between coding and issue tracking systems.

The fit implementation is (almost) the simplest issue system that can still
work. Human readable plain text files are intuitively understood, viewable
and editable.

## History

fit is the golang program first developed as "bug" by Dave MacFarlane (driusan).
Filesystem Issue Tracker ([FIT.md](FIT.md)) is the new name for the Poor Man's
Issue Tracker (PMIT) storage system also first developed by driusan. For his
demo from 2016, see [driusan's
talk](https://www.youtube.com/watch?v=ysgMlGHtDMo) at the first
GolangMontreal.org conference, GoMTL-01. The program and storage system have
evolved while trying to remain backward compatible. See the [FAQ.md](FAQ.md)
for more information.

## Background

A limited but sufficient number of conventions with just enough organization
can quickly capture issues using human readable issue directories and files.
fit can be the primary system if no other system is provided or supplement
other issue/bug systems to quickly capture issues and their context as close
to the issue as possible.

Using fit helps implementers streamline working with
[issues](https://en.wikipedia.org/wiki/Issue_tracking_system) and [version
control](https://en.wikipedia.org/wiki/Version_control). fit works with
both git and mercurial distributed version control though the git features are
more well exercised.

fit is designed to adapt to your processes using issue key/value pair metadata.

The fit too manages issues using conventions/format of
Filesystem Issue Tracker (see [FIT.md](FIT.md)). A `fit/` or `issues/`
directory holds one descriptively titled directory per issue. Each directory 
holds a Description (name is configurable) text file and anything else needed.

At first people may naturally try to keep track of issues in a single text
file and/or spreadsheet but these can fail to meet project needs.
(see [FAQ.md](FAQ.md))

Issue context is valuable to coders and may be difficult for others to
understand. fit also supports multiple `fit/` directories in the
repository's tree for stronger coordination of coding and issue tracking.

An alternative in IT projects is all too common: implementers are not given the
tools needed to record code issues because issue systems take resources to
setup and maintain. These separate isssue systems are often focused on user
facing issues so valuable implementation details are often lost. IT hopes
problems will not attract attention. Valuable project budget, time, scope,
quality or other resources are focused on new features. Beyond managing
contentious flat files or spreadsheets there is FIT.

Important issues can be captured, surfaced and addressed, whether they are
actual problems, questions, possible features or ideas by those most familiar
with the project. Less code savvy project collaborators are not distracted by
implementation details, code reviews or operational facing features.

[Software Development Life Cycles](https://en.wikipedia.org/wiki/Software_development_process) (SDLCs) involve more than just the source code.
Over time needs may change from hacking/coding, just getting something working,
to implementing more disciplined software engineering best practices. Code can
start small and grow gradually as users, use cases and developers are added.
The FIT issue system can adapt to each stage.

While one issue set used for one git repository may be enough the use of
recursive fit directories are supported. As complexity increases adding
multiple `fit/` directories in different parts of your git repo may help
project coders keep focused.

There are some choices for how to handle closed issues. As the number of
issues grows closed can simply be deleted or an archive can hold the
inactive issues. While deleting issues helps keep things uncluttered issues
still have value over time and may be difficult to find using only version
control history.

fit can be aliased as a git subcommand "git fit ..."

Security concerns are handled using standard git repository practices.

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

## Getting Started

### Layout

Filesystem Issue Tracker ([FIT.md](FIT.md)) conventions/format are a set of
suggestions for storing issues, one directory/folder per issue with plain text
file details.

A `fit/` directory holds one (descriptively titled) directory per issue.
The "Description" file is the only text needed providing the details. Optional
tag_key_value files assign meta data.

fit maintains the nearest `fit/` directory to your current working directory
or it's parent directories. There can be more than one. fit can commit (or
remove) issues from versioning or this can be done manually without fit. Unlike
many other issue systems, fit issues naturally branch and merge along with the
rest of your versioned files. Using branches or related repos is optional.

Some support is available to import and/or reference other issue trackers.

### Installation

After you have [go installed](https://golang.org/doc/install), install the
latest version of fit with:

`GO111MODULE=on go install github.com/grantbow/fit`

If that does not work in one command then:
    export GO111MODULE=on
    cd $GOPATH/src/github.com/grantbow/fit/cmd/fit
    go install

This will create the binary $GOPATH/src/github.com/grantbow/fit/cmd/fit/fit(.exe)
and move it to $GOPATH/bin/fit

Make sure `$GOPATH/bin` or `$GOBIN` are in your path or you can copy
the "fit" binary somewhere that is in your path.

The environment variable GO111MODULE changes how your golang works by enabling
golang 1.11+ module support required for this version of fit. The default in
golang 1.12, 1.13 and 1.14 are still "auto".

Working with fit and git via the command line can be simplified. You can run
fit as a git subcommand like `git fit` or `git issue`. This [chapter about git
aliases](https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases) describes how
to set them up very well. It is part of the Pro Git book available for free
online. You can quickly add the alias to your .gitconfig:

`git config --global alias.issue !/path/to/fit`

This will add to your $HOME/.gitconfig or you can edit it manually:

`[alias]
    issue = !/path/to/fit`

### Configuration

Settings can be read from an optional config file .fit.yml placed next to the
fit directory. Current options include:

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
    * MultipleFitDirs: true or false
          Default is false.
          always recursive when possible
    * CloseStatusTag: true or false
          Default is false.
          fit close will not add a tag (false) or
          add tag_status_closed (true)
    * CloseMove: true or false
          Default is false.
          fit close will use ClosePreventDelete (false) or
          fit close will move issues to CloseDir (true) and
              sets ClosePreventDelete=true
    * FitClosedDirName: string
          Default is "closed"
          this is the name used inside the fit directory
    * ClosePreventDelete: true or false
          Default is false.
          fit close will delete (false) or
          fit close will not delete (true)
              implies CloseStatusTag and/or CloseMove
    * IdAbbreviate: true or false
          Default is false.
          Use Identifier.
    * IdAutomatic: true or false
          Default is false.
          Use Identifier.

Other issue systems may use databases, hidden directories or hidden branches.
While these may be useful techniques in certain circumstances this seems to
unnecessarily obfuscate access.

### Hooks

Event based automation can be added through git or mercurial. We created a
hooks directory and look forward to seeing what code teams use and contribute.
Work to help adapt hooks to both git and hg would be appreciated.

### Example Use

To get started in the top of an existing git repo simply `mkdir fit` then
add and commit to the repo fit/_issue name_/Description files. Create and
maintain issues by editing the files inside fit/_issue names_/.

If an environment variable named FIT is set that value will be used as a
directory name used to find the 'fit' or 'issues' directory. All fit commands
will use FIT no matter your present working directory.

If an 'fit' directory/folder is not found fit will walk up the tree toward
the root until it finds a "fit" subdirectory similar to how git looks for
.git or hg looks for .hg. A warning is provided if no directory is found.

```
$ mkdir foo && cd foo
$ git init
$ mkdir fit
$ fit help
cmd: help
Usage: help <command>

Use "fit help <command>" or "fit <command> help" for
more information about any command below.
fit version 0.6 built using go1.12.6 GOOS android
executable: -rwx------ 8749028 Sat Jun 22 10:19:54 PDT 2019 /data/data/com.termux/files/home/go/bin/fit

Status/reading commands:
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

Editing commands:
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

Version control commands:
    commit     Commit any new, changed or deleted issues
    purge      Remove all issues not tracked

Processing commands:
    roadmap    Print list of open issues sorted by milestone

aliases for help: --help -h

$ fit create Need better help
(<your editor> Description)
(save and quit)
Created issue: Need better help

$ fit list

===== list /...
Issue 1: Need better help

$ fit list 1

===== list /...
Title: Need better help
Description:
<the entered description>

$ fit create -n Need better formating for README
(no editor, default to empty Description)
Created issue: Need better formatting for README

$ fit list
Issue 1: Need better help
Issue 2: Need better formating for README
```

## Next Steps

* [FIT.md](FIT.md)
* [FAQ.md](FAQ.md)
* [CONTRIBUTING.md](CONTRUBUTING.md)
* [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
* [SUPPORT.md](SUPPORT.md)
* [wiki](https://github.com/grantbow/fit/wiki)
* [gitter](https://gitter.im/fit-issue/community)

Your system is the beginning, not the end. Much has been written about how to
use and setup systems to track or manage issues, software bugs, trouble
tickets, support tickets, incident tickets or requests. See the FAQ.md

### Feedback

We would like to hear about how you use this system.

I would like to work with others and would appreciate feedback at
grantbow+fit@gmail.com.

Since the original project is not very active I have gone ahead and continuted
development on my fork. I encourage discussion. Submitting can be done with a
pull request, to our upstream project or using [git
remotes](https://stackoverflow.com/questions/36628859/git-how-to-merge-a-pull-request-into-a-fork).

Anyone thinking of [CONTRIBUTING.md](CONTRIBUTING.md) is encouraged to do so.
As this is an issue tracking system a pull request with an issue seems logical
enough and good working practice. Development Guidelines are included.

As mentioned in [SUPPORT.md](SUPPORT.md) questions are encouraged via email, issues or pull
requests for now.

The [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) is the standard recommended by github
offered by contributor-covenant.org.

