package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
)

// Pwd is a subcommand to output the issues directory.
func Pwd(config bugs.Config) {
	fmt.Printf("%s\n", bugs.FitDirer(config))
}
