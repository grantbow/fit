package fitapp

import (
	bugs "github.com/grantbow/fit/issues"
)

// Milestone is a subcommand to assign a milestone to an issue.
func Milestone(args argumentList, config bugs.Config) {
	fieldHandler("milestone", args, bugs.Issue.SetMilestone, bugs.Issue.Milestone, config)
}
