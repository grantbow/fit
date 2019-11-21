package scm

import (
	bugs "github.com/driusan/bug/bugs"
	"os"
)

var dops = bugs.Directory(os.PathSeparator)
var sops = string(os.PathSeparator)
