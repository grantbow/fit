package fitapp

import (
	bugs "github.com/grantbow/fit/issues"
)

// Priority is a subcommand to assign a priority to an issue.
func Priority(args argumentList, config bugs.Config) {
	fieldHandler("priority", args, bugs.Issue.SetPriority, bugs.Issue.Priority, config)
}
