package scm

import "github.com/driusan/bug/bugs"

// SCMHandler interface defines how to call Commit, Purge and GetSCMType.
type SCMHandler interface {
	Commit(dir bugs.Directory, commitMsg string) error
	Purge(bugs.Directory) error
	GetSCMType() string
}

// FileStatus type holds information about a file.
type FileStatus struct {
	Filename      string
	IndexStatus   string
	WorkingStatus string
}
