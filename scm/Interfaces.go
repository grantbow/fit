package scm

import bugs "github.com/grantbow/fit/issues"

// SCMHandler interface defines how to call Commit, Purge and SCMTyper.
type SCMHandler interface {
	Commit(dir bugs.Directory, commitMsg string, config bugs.Config) error
	Purge(bugs.Directory) error
	SCMTyper() string
	SCMIssuesUpdaters(config bugs.Config) ([]byte, error)
	SCMIssuesCacher(config bugs.Config) ([]byte, error)
}

// FileStatus type holds information about a file.
type FileStatus struct {
	Filename      string
	IndexStatus   string
	WorkingStatus string
}
