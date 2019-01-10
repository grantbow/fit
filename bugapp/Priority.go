package bugapp

import (
	"github.com/driusan/bug/bugs"
)

// Priority is a subcommand to assign a priority to an issue.
func Priority(args ArgumentList, config bugs.Config) {
	fieldHandler("priority", args, bugs.Bug.SetPriority, bugs.Bug.Priority, config)
}
