package bugs

import "os"

var dops = Directory(os.PathSeparator)
var sops = string(os.PathSeparator)

func check(e error) {
	if e != nil {
		//	fmt.Fprintf(os.Stderr, "err: %s", err.Error())
		//	return NoConfigError
		panic(e)
	}
}
