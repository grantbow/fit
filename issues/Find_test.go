package issues

import (
	"fmt"
	"testing"
)

func TestFindIssuesByTag(t *testing.T) {
	//tests: func FindIssuesByTag(tags []string, config Config) []Issue
	config := Config{}
	test := tester{} // from Issue_test.go
	test.Setup()
	defer test.Teardown()

	//b := test.bug

	c, err := New("find me", config)
	if err != nil {
		panic("Unexpected error creating issue find me")
	}
	c.TagIssue("hit", config)
	c.SetIdentifier("special", config)

	d, errd := New("hiding", config)
	if errd != nil {
		panic("Unexpected error creating issue hiding")
	}
	d.TagIssue("miss", config)

	a := FindIssuesByTag([]string{"tag"}, config)
	if len(a) != 0 {
		t.Error(fmt.Sprintf("Failed %s: expected %v but got %v\n%#v", "FindIssuesByTag 0", 0, len(a), a))
	}
	a = FindIssuesByTag([]string{"hit"}, config)
	if len(a) != 1 {
		t.Error(fmt.Sprintf("Failed %s: expected %v but got %v\n%#v", "FindIssuesByTag 1", 1, len(a), a))
	}

	//tests: func (b IssueNotFoundError) Error() string
	out := IssueNotFoundError("foo")
	if out != "foo" {
		t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "Find IssueNotFoundError", "foo", out))
	}

	// LOAD
	//tests: func LoadIssueByDirectory(dir string, config Config) (*Issue, error)
	var x *Issue
	x, errx := LoadIssueByDirectory("find-me", config) // String
	if errx != nil {
		y, booly := errx.(IssueNotFoundError)
		if booly {
			t.Error(fmt.Sprintf("Failed %s: got type %v %v", "LoadIssueByDirectory", "IssueNotFoundError", errx))
		} else {
			t.Error(fmt.Sprintf("Failed %s: got %v", "LoadIssueByDirectory", y))
		}
	} else {
		if string(x.Dir) != string(c.Dir) {
			t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "LoadIssueByDirectory", c.Dir, x.Dir))
		}
	}
	// LOAD
	//tests: func LoadIssueByHeuristic(id string, config Config) (*Issue, error)
	x, errx = LoadIssueByHeuristic("2", config) // c is 2 String
	if errx != nil {
		y, booly := errx.(IssueNotFoundError)
		if booly {
			t.Error(fmt.Sprintf("Failed %s: got type %v %v", "LoadIssueByHeuristic", "IssueNotFoundError", errx))
		} else {
			t.Error(fmt.Sprintf("Failed %s: got %v", "LoadIssueByHeuristic", y))
		}
	} else {
		if string(x.Dir) != string(c.Dir) {
			t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "LoadIssueByHeuristic", c.Dir, x.Dir))
		}
	}
	//// LOAD
	////tests: func LoadIssueByStringIndex(i string, config Config) (*Issue, error)
	//x, errx = LoadIssueByStringIndex("2", config) // c is 2 String
	//if errx != nil {
	//	y, booly := errx.(IssueNotFoundError)
	//	if booly {
	//		t.Error(fmt.Sprintf("Failed %s: got type %v %v", "LoadIssueByStringIndex", "IssueNotFoundError", errx))
	//	} else {
	//		t.Error(fmt.Sprintf("Failed %s: got %v", "LoadIssueByStringIndex", y))
	//	}
	//} else {
	//	if string(x.Dir) != string(c.Dir) {
	//		t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "LoadIssueByStringIndex", c.Dir, x.Dir))
	//	}
	//}
	// LOAD
	//tests: func LoadIssueByIndex(idx int, config Config) (*Issue, error)
	x, errx = LoadIssueByIndex(2, config) // c is 2 Int
	if errx != nil {
		y, booly := errx.(IssueNotFoundError)
		if booly {
			t.Error(fmt.Sprintf("Failed %s: got type %v %v", "LoadIssueByIndex", "IssueNotFoundError", errx))
		} else {
			t.Error(fmt.Sprintf("Failed %s: got %v", "LoadIssueByIndex", y))
		}
	} else {
		if string(x.Dir) != string(c.Dir) {
			t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "LoadIssueByIndex", c.Dir, x.Dir))
		}
	}
	// LOAD
	// tests: func LoadIssueByIdentifier(id string, config Config) (*Issue, error)
	x, errx = LoadIssueByIdentifier("special", config) // c is 2
	if errx != nil {
		y, booly := errx.(IssueNotFoundError)
		if booly {
			t.Error(fmt.Sprintf("Failed %s: got type %v %v", "LoadIssueByIdentifier", "IssueNotFoundError", errx))
		} else {
			t.Error(fmt.Sprintf("Failed %s: got %v", "LoadIssueByIdentifier", y))
		}
	} else {
		if string(x.Dir) != string(c.Dir) {
			t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "LoadIssueByIdentifier", c.Dir, x.Dir))
		}
	}

	//tests: func GetAllIssues(config Config) []Issue
	bugarray := GetAllIssues(config)
	if len(bugarray) != 3 { // minimal
		t.Error(fmt.Sprintf("Failed %s: got %v", "GetAllIssues", bugarray))
	}
}
