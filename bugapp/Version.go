package bugapp

import (
	"fmt"
	"os"
	"runtime"
)

// Version is a subcommand to output the command name and golang runtime.Version().
func Version() {
	fmt.Printf("%s version 0.4 built using %s\n", os.Args[0], runtime.Version())

}
