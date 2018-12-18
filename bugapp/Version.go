package bugapp

import (
	"fmt"
	"os"
	"runtime"
)

func Version() {
	fmt.Printf("%s version 0.4 built using %s\n", os.Args[0], runtime.Version())

}
