package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
)

func Pwd(config bugs.Config) {
	fmt.Printf("%s", bugs.GetIssuesDir(config))
}
