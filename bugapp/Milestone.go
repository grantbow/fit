package bugapp

import (
	"github.com/driusan/bug/bugs"
)

func Milestone(args ArgumentList, config bugs.Config) {
	fieldHandler("milestone", args, bugs.Bug.SetMilestone, bugs.Bug.Milestone, config)
}
