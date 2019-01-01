package bugapp

import (
	"github.com/driusan/bug/bugs"
)

func Priority(args ArgumentList, config bugs.Config) {
	fieldHandler("priority", args, bugs.Bug.SetPriority, bugs.Bug.Priority, config)
}
