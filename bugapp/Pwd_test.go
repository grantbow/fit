package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

type tester struct {
	// copied and slightly modified from Bugs.go
	dir string
	bug *bugs.Bug
}

func (t *tester) Setup() {
	config := bugs.Config{}
	gdir, err := ioutil.TempDir("", "issuetest")
	if err == nil {
		os.Chdir(gdir)
		t.dir = gdir
		os.Unsetenv("FIT")
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		//gdir, _ = os.Getwd() // gdir not used later
	} else {
		panic("Failed creating temporary directory")
	}
	// Make sure we get the right directory from the top level
	os.Mkdir("issues", 0755)
	b, err := bugs.New("Test Bug", config)
	if err != nil {
		panic("Unexpected error creating Test Bug")
	}
	t.bug = b
}
func (t *tester) Teardown() {
	os.RemoveAll(t.dir)
}

func TestPwd(t *testing.T) {
	config := bugs.Config{}
	test := tester{} // from Bug_test.go
	test.Setup()
	defer test.Teardown()

	stdout, _ := captureOutput(func() {
		Pwd(config)
	}, t)
	re := regexp.MustCompile("issues")
	matched := re.MatchString(stdout)
	if !matched {
		t.Error("Unexpected output on STDOUT for bugapp/Pwd_test")
		fmt.Printf("Expected to match: %s\nGot: %s\n", "issues", stdout)
	}
}
