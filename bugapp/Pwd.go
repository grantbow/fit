package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
)

// Pwd is a subcommand to output the issues directory.
func Pwd(config bugs.Config) {
	fmt.Printf("%s", bugs.GetIssuesDir(config))
}
