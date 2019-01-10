package bugapp

import (
	"github.com/driusan/bug/bugs"
)

// Milestone is a subcommand to assign a milestone to an issue.
func Milestone(args ArgumentList, config bugs.Config) {
	fieldHandler("milestone", args, bugs.Bug.SetMilestone, bugs.Bug.Milestone, config)
}
