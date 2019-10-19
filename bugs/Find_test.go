package bugs

import (
	"fmt"
	"testing"
)

func TestFindBugsByTag(t *testing.T) {
	//tests: func FindBugsByTag(tags []string, config Config) []Bug
	config := Config{}
	test := tester{} // from Bug_test.go
	test.Setup()
	defer test.Teardown()

	//b := test.bug

	c, err := New("find me", config)
	if err != nil {
		panic("Unexpected error creating bug find me")
	}
	c.TagBug("hit", config)
	c.SetIdentifier("special", config)

	d, errd := New("hiding", config)
	if errd != nil {
		panic("Unexpected error creating bug hiding")
	}
	d.TagBug("miss", config)

	a := FindBugsByTag([]string{"tag"}, config)
	if len(a) != 0 {
		t.Error(fmt.Sprintf("Failed %s: expected %v but got %v\n%#v", "FindBugsByTag 0", 0, len(a), a))
	}
	a = FindBugsByTag([]string{"hit"}, config)
	if len(a) != 1 {
		t.Error(fmt.Sprintf("Failed %s: expected %v but got %v\n%#v", "FindBugsByTag 1", 1, len(a), a))
	}

	//tests: func (b BugNotFoundError) Error() string
	out := BugNotFoundError("foo")
	if out != "foo" {
		t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "Find BugNotFoundError", "foo", out))
	}

	// LOAD
	//tests: func LoadBugByDirectory(dir string, config Config) (*Bug, error)
	var x *Bug
	x, errx := LoadBugByDirectory("find-me", config) // String
	if errx != nil {
		y, booly := errx.(BugNotFoundError)
		if booly {
			t.Error(fmt.Sprintf("Failed %s: got type %v %v", "LoadBugByDirectory", "BugNotFoundError", errx))
		} else {
			t.Error(fmt.Sprintf("Failed %s: got %v", "LoadBugByDirectory", y))
		}
	} else {
		if string(x.Dir) != string(c.Dir) {
			t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "LoadBugByDirectory", c.Dir, x.Dir))
		}
	}
	// LOAD
	//tests: func LoadBugByHeuristic(id string, config Config) (*Bug, error)
	x, errx = LoadBugByHeuristic("2", config) // c is 2 String
	if errx != nil {
		y, booly := errx.(BugNotFoundError)
		if booly {
			t.Error(fmt.Sprintf("Failed %s: got type %v %v", "LoadBugByHeuristic", "BugNotFoundError", errx))
		} else {
			t.Error(fmt.Sprintf("Failed %s: got %v", "LoadBugByHeuristic", y))
		}
	} else {
		if string(x.Dir) != string(c.Dir) {
			t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "LoadBugByHeuristic", c.Dir, x.Dir))
		}
	}
	//// LOAD
	////tests: func LoadBugByStringIndex(i string, config Config) (*Bug, error)
	//x, errx = LoadBugByStringIndex("2", config) // c is 2 String
	//if errx != nil {
	//	y, booly := errx.(BugNotFoundError)
	//	if booly {
	//		t.Error(fmt.Sprintf("Failed %s: got type %v %v", "LoadBugByStringIndex", "BugNotFoundError", errx))
	//	} else {
	//		t.Error(fmt.Sprintf("Failed %s: got %v", "LoadBugByStringIndex", y))
	//	}
	//} else {
	//	if string(x.Dir) != string(c.Dir) {
	//		t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "LoadBugByStringIndex", c.Dir, x.Dir))
	//	}
	//}
	// LOAD
	//tests: func LoadBugByIndex(idx int, config Config) (*Bug, error)
	x, errx = LoadBugByIndex(2, config) // c is 2 Int
	if errx != nil {
		y, booly := errx.(BugNotFoundError)
		if booly {
			t.Error(fmt.Sprintf("Failed %s: got type %v %v", "LoadBugByIndex", "BugNotFoundError", errx))
		} else {
			t.Error(fmt.Sprintf("Failed %s: got %v", "LoadBugByIndex", y))
		}
	} else {
		if string(x.Dir) != string(c.Dir) {
			t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "LoadBugByIndex", c.Dir, x.Dir))
		}
	}
	// LOAD
	// tests: func LoadBugByIdentifier(id string, config Config) (*Bug, error)
	x, errx = LoadBugByIdentifier("special", config) // c is 2
	if errx != nil {
		y, booly := errx.(BugNotFoundError)
		if booly {
			t.Error(fmt.Sprintf("Failed %s: got type %v %v", "LoadBugByIdentifier", "BugNotFoundError", errx))
		} else {
			t.Error(fmt.Sprintf("Failed %s: got %v", "LoadBugByIdentifier", y))
		}
	} else {
		if string(x.Dir) != string(c.Dir) {
			t.Error(fmt.Sprintf("Failed %s: expected %v but got %v", "LoadBugByIdentifier", c.Dir, x.Dir))
		}
	}

	//tests: func GetAllBugs(config Config) []Bug
	bugarray := GetAllBugs(config)
	if len(bugarray) != 3 { // minimal
		t.Error(fmt.Sprintf("Failed %s: got %v", "GetAllBugs", bugarray))
	}
}
