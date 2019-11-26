package scm

import (
	bugs "github.com/grantbow/fit/issues"
	"os"
)

var dops = bugs.Directory(os.PathSeparator)
var sops = string(os.PathSeparator)
