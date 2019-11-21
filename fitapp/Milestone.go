package fitapp

import (
	bugs "github.com/driusan/bug/bugs"
)

// Milestone is a subcommand to assign a milestone to an issue.
func Milestone(args argumentList, config bugs.Config) {
	fieldHandler("milestone", args, bugs.Bug.SetMilestone, bugs.Bug.Milestone, config)
}
