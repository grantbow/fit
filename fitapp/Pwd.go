package fitapp

import (
	"fmt"
	bugs "github.com/driusan/bug/bugs"
)

// Pwd is a subcommand to output the issues directory.
func Pwd(config bugs.Config) {
	fmt.Printf("%s\n", bugs.IssuesDirer(config))
}
