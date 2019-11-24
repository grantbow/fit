package issues

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

type tester struct {
	dir   string
	issue *Issue
	pwd   string
}

func (t *tester) Setup() {
	config := Config{}
	config.DescriptionFileName = "Description"
	config.FitDirName = "issues"
	gdir, err := ioutil.TempDir("", "issuetest")
	pwd, _ := os.Getwd()
	t.pwd = pwd
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
	os.Mkdir(config.FitDirName, 0755)
	b, err := New("Test Issue", config)
	if err != nil {
		panic("Unexpected error creating Test Issue")
	}
	t.issue = b
}
func (t *tester) Teardown() {
	os.Chdir(t.pwd)
	os.RemoveAll(t.dir)
}

func TestTitleToDirectory(t *testing.T) {
	var assertDirectory = func(title, directory string) {
		titleStr := TitleToDir(title)
		dirStr := Directory(directory).ShortNamer()

		if titleStr != dirStr {
			t.Error(fmt.Sprintf("Failed on %s: got %s but expected %s\n", title, titleStr, dirStr))
		}
	}

	assertDirectory("Test", "Test")
	assertDirectory("Test Space", "Test-Space")
	//assertDirectory("Test-Dash", "Test--Dash")
	//assertDirectory("Test--TripleDash", "Test--TripleDash")
	assertDirectory("Test --WithSpace", "Test_--WithSpace")
	assertDirectory("Test - What", "Test_-_What")
	assertDirectory("Test : What", "Test-_-What")
	assertDirectory("Test ? What", "Test-_-What")
	assertDirectory("Test / What", "Test-_-What")
	assertDirectory("Test . What", "Test-_-What")
}

func TestShortTitleToDir(t *testing.T) {
	var assertDirectory = func(title, directory string) {
		titleStr := ShortTitleToDir(title)
		dirStr := Directory(directory).ShortNamer()

		if titleStr != dirStr {
			t.Error(fmt.Sprintf("Failed on %s: got %s but expected %s\n", title, titleStr, dirStr))
		}
	}
	assertDirectory("123456789012345678901234567", "1234567890123456789012345")

}

func TestNewIssue(t *testing.T) {
	var gdir string
	config := Config{}
	config.DescriptionFileName = "Description"
	config.FitDirName = "fit"
	gdir, err := ioutil.TempDir("", "newissue")
	pwd, _ := os.Getwd()
	if err == nil {
		os.Chdir(gdir)
		// Hack to get around the fact that /tmp is a symlink on
		// OS X, and it causes the directory checks to fail..
		gdir, _ = os.Getwd()
		defer os.RemoveAll(gdir)
	} else {
		t.Error("Failed creating temporary directory for detect")
		return
	}
	os.Mkdir(config.FitDirName, 0755)
	b, err := New("I am a test", config)
	if err != nil || b == nil {
		t.Error("Unexpected error when creating New issue" + err.Error())
	}
	if b.Dir != FitDirer(config)+Directory(os.PathSeparator)+TitleToDir("I am a test") {
		t.Error("Unexpected directory when creating New issue")
	}
	os.Chdir(pwd)
}

func TestSetDescription(t *testing.T) {
	config := Config{}
	config.DescriptionFileName = "Description"
	test := tester{}
	test.Setup()
	defer test.Teardown()

	b := test.issue

	b.SetDescription("Hello, I am an issue.", config)
	val, err := ioutil.ReadFile(string(b.Direr()) + sops + config.DescriptionFileName)
	if err != nil {
		t.Error("Could not read Description file")
	}

	if string(val) != "Hello, I am an issue.\n" {
		t.Error("Unexpected description after SetDescription")
	}
}

func TestTitle(t *testing.T) {
	config := Config{}
	config.DescriptionFileName = "Description"
	test := tester{}
	test.Setup()
	defer test.Teardown()

	b := test.issue

	expected := "Test Issue"
	val := b.Title("")
	if string(val) != expected {
		t.Error(fmt.Sprintf("Failed on %s: got %s but expected %s\n", "TestTitle", val, expected))
	}
}

//type Comment struct {
//	Author string
//	Time   time.Time
//	Body   string
//	Order  int
//	Xml    []byte
//}

func TestCommentStatusPriorityMilestone(t *testing.T) {
	config := Config{}
	config.DescriptionFileName = "Description"
	test := tester{}
	test.Setup()
	defer test.Teardown()

	b := test.issue

	expected := "Test Issue Comment"
	//b.CommentIssue(Comment("Author", time.Now(), expected, 0, []byte("")), config)
	b.CommentIssue(Comment{Author: "Author", Time: time.Now(), Body: expected, Order: 0, Xml: []byte("")}, config)
	b.RemoveComment(Comment{Author: "Author", Time: time.Now(), Body: expected, Order: 0, Xml: []byte("")})
	_ = b.SetStatus("do", config)
	_ = b.Status()
	_ = b.SetPriority("urgent", config)
	_ = b.Priority()
	_ = b.SetMilestone("release", config)
	_ = b.Milestone()
}

func TestDescription(t *testing.T) {
	config := Config{}
	config.DescriptionFileName = "Description"
	test := tester{}
	test.Setup()
	defer test.Teardown()

	b := test.issue
	b.DescriptionFileName = config.DescriptionFileName

	desc := "I am yet another issue.\nWith Two Lines."
	b.SetDescription(desc, config)

	if b.Description() != desc+"\n" {
		title := "TestDescription"
		t.Error(fmt.Sprintf("Failed on %s:\ngot:\n%s\nbut expected:\n%s\n", title, b.Description(), desc))
	}
}
