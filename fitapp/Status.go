package fitapp

import bugs "github.com/grantbow/fit/issues"

// Status is a subcommand to assign a status to an issue.
func Status(args argumentList, config bugs.Config) {
	fieldHandler("status", args, bugs.Issue.SetStatus, bugs.Issue.Status, config)
}
