package bugapp

import (
	"fmt"
	"os"
	"runtime"
)

func ProgramVersion() string {
	return "0.4"
}

// Version is a subcommand to output the command name and golang runtime.Version().
func PrintVersion() {
	fmt.Printf("%s version %s built using %s\n", os.Args[0], ProgramVersion(), runtime.Version())
}
