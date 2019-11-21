package fitapp

import bugs "github.com/driusan/bug/bugs"

// Status is a subcommand to assign a status to an issue.
func Status(args argumentList, config bugs.Config) {
	fieldHandler("status", args, bugs.Bug.SetStatus, bugs.Bug.Status, config)
}
