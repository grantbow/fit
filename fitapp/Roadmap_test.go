package fitapp

import (
	"fmt"
	bugs "github.com/grantbow/fit/issues"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
	"time"
)

//var dops = bugs.Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

func TestRoadmapLess(t *testing.T) {
	// func (a BugListByMilestone) Less(i, j int) bool {
	config := bugs.Config{}
	config.FitDirName = "fit"
	dir, err := ioutil.TempDir("", "roadmaptest")
	if err != nil {
		t.Error("Could not create temporary dir for test")
		return
	}
	pwd, _ := os.Getwd()
	os.Chdir(dir)
	os.MkdirAll(config.FitDirName, 0700)
	defer os.RemoveAll(dir)
	// On MacOS, /tmp is a symlink, which causes GetDirectory() to return
	// a different path than expected in these tests, so make the issues
	// directory explicit with an environment variable
	err = os.Setenv("FIT", dir)
	if err != nil {
		t.Error("Could not set environment variable: " + err.Error())
		return
	}
	// bug 1
	expected := "Created issue: Test1bug\n"
	stdout, stderr := captureOutput(func() {
		Create(argumentList{"-n", "Test1bug"}, config)
	}, t)
	if stderr != "" || stdout != expected {
		t.Error("Unexpected output/err create 1")
		fmt.Printf("Expected stdout: %s\nGot: %s\n", expected, stdout)
		fmt.Printf("Expected stderr: %s\nGot: %s\n", "", stderr)
	}
	// bug 2
	expected = "Created issue: Test2bug\n"
	stdout, stderr = captureOutput(func() {
		Create(argumentList{"-n", "Test2bug"}, config)
	}, t)
	if stderr != "" || stdout != expected {
		t.Error("Unexpected output/error create 2")
		fmt.Printf("Expected stdout: %s\nGot: %s\n", expected, stdout)
		fmt.Printf("Expected stderr: %s\nGot: %s\n", "", stderr)
	}
	// milestones
	runmiles(argumentList{"1", "v1.0"}, "", t)
	time.Sleep(3 * time.Second)
	runmiles(argumentList{"2", "v2.0"}, "", t)
	// roadmap
	expected = `# Roadmap for .*

## v2.0:
- Test2bug

## v1.0:
- Test1bug
`
	re := regexp.MustCompile(expected)
	//d, _   := os.Open(dir)
	//dd, _  := d.Readdir(0)
	//dd, _ := ioutil.ReadDir(dir+sops+config.FitDirName)
	//fmt.Printf("dirdump %s\n", dirDump(fmt.Sprintf("%s%s%s",dir, sops,config.FitDirName)))
	stdout, stderr = captureOutput(func() {
		Roadmap(argumentList{}, bugs.Config{})
	}, t)
	matched := re.MatchString(stdout)
	if stderr != "" || !matched {
		t.Error("Unexpected out/error of roadmap")
        if !matched {
            fmt.Printf("Expected stdout:\n%s\nGot:\n%s\n", expected, stdout)
        }
        if stderr != "" {
            fmt.Printf("Expected stderr:\n%s\nGot:\n%s\n", "", stderr)
        }
	}
	os.Chdir(pwd)
}

//func TestRoadmap(t *testing.T) {
// Roadmap(argumentList{}, bugs.Config{})
//}

/*
type BugListByMilestone [](bugs.Bug)
func (a BugListByMilestone) Len() int      { return len(a) }
func (a BugListByMilestone) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BugListByMilestone) Less(i, j int) bool {
	iMS := a[i].Milestone()
	jMS := a[j].Milestone()
	// If there's a "v" at the start, strip it out
	// before doing any comparisons of semantic
	// versions
	if len(iMS) > 1 && iMS[0] == "v"[0] {
		iMS = iMS[1:]
	}
	if len(jMS) > 1 && jMS[0] == "v"[0] {
		jMS = jMS[1:]
	}
	// First try semantic versioning comparison
	iVer, iVerErr := semver.Make(iMS)
	jVer, jVerErr := semver.Make(jMS)
	if iVerErr == nil && jVerErr == nil {
		return iVer.LT(jVer)
	}

	// Next try floating point comparison as an
	// approximation of real number comparison..
	iFloat, iVerErr := strconv.ParseFloat(iMS, 32)
	jFloat, jVerErr := strconv.ParseFloat(jMS, 32)
	if iVerErr == nil && jVerErr == nil {
		return iFloat < jFloat
	}

	// Finally, just use a normal string collation
	return iMS < jMS
}

func Roadmap(args argumentList, config bugs.Config) {
	var bgs []bugs.Bug

	if args.HasArgument("--filter") {
		tags := strings.Split(args.GetArgument("--filter", ""), ",")
		fmt.Printf("%s", tags)
		bgs = bugs.FindBugsByTag(tags, config)
	} else {
		bgs = bugs.GetAllBugs(config)
	}
	sort.Sort(BugListByMilestone(bgs))

	fmt.Printf("# Roadmap for %s\n", bugs.RootDirer(&config).ShortNamer().ToTitle())
	milestone := ""
	for i := len(bgs) - 1; i >= 0; i -= 1 {
		b := bgs[i]
		newMilestone := b.Milestone()
		if milestone != newMilestone {
			if newMilestone == "" {
				fmt.Printf("\n## No milestone set:\n")
			} else {
				fmt.Printf("\n## %s:\n", newMilestone)
			}
		}
		if args.HasArgument("--simple") {
			fmt.Printf("- %s\n", b.Title(""))
		} else {
			options := ""
			if !args.HasArgument("--no-status") {
				options += "status"
			}
			if !args.HasArgument("--no-priority") {
				options += " priority"
			}
			if !args.HasArgument("--no-identifier") {
				options += " identifier"
			}

			if args.HasArgument("--tags") {
				options += "tags"
			}
			fmt.Printf("- %s\n", b.Title(options))
		}
		milestone = newMilestone

	}
}
*/
