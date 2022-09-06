# fit
filesystem issue tracker: manages plain text issues with git or hg

**TOC:**

<!-- toc -->

- [Prerequisites](#prerequisites)
- [Installation](#installation)

<!-- tocstop -->

## Prerequisites

git or hg (mercurial)

golang

linux, mac or windows OS.

## Installation

First [install go](https://go.dev/doc/install) completely.

make sure you
have `go/bin` (or equivalent) for golang and `$HOME/go/bin` for installed binaries in your path.

`GOPATH` should be set to something like `$HOME/go`.

Install the latest version of fit with:

`go install github.com/grantbow/fit@latest`

If that does not work in one command then:
```
    $ mkdir -p $GOPATH/src/github.com/grantbow/fit
    $ git clone https://github.com/grantbow/fit $GOPATH/src/github.com/grantbow/fit
    $ cd $GOPATH/src/github.com/grantbow/fit/cmd/fit 
      # that's the fit/cmd/fit dir
    $ go install
```

This will create the binary `$GOPATH/src/github.com/grantbow/fit/cmd/fit/fit(.exe)`
and move it to `$GOPATH/bin/fit(.exe)`

The environment variable set using `export GO111MODULE=on` changed how old
golang versions work by enabling golang 1.11+ module support required by fit.
The defaults in golang 1.13 and 1.14 and 1.15 were still "auto".
The defaults in golang 1.16 and above are "on" so this setup is no longer required.

You can use fit directly or use fit as a git subcommand like `git fit`.
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

This [chapter about git aliases](https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases) describes how
to set them up very well. It is part of the Pro Git book available for free
online. 

Your system is just the beginning, not the end. Much has been written about
how to use and setup systems to track or manage issues, software bugs, trouble
tickets, support tickets, incident tickets or requests. See the docs/FAQ.md

