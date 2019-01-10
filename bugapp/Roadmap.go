package bugapp

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/driusan/bug/bugs"
	"sort"
	"strconv"
	"strings"
)

// BugListByMilestone is a pkg global to hold a list of issues.
type BugListByMilestone [](bugs.Bug)

// Len, Swap and Less sort issues.
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

// Roadmap is a subcommand to output issues by milestone.
func Roadmap(args ArgumentList, config bugs.Config) {
	var bgs []bugs.Bug

	if args.HasArgument("--filter") {
		tags := strings.Split(args.GetArgument("--filter", ""), ",")
		fmt.Printf("%s", tags)
		bgs = bugs.FindBugsByTag(tags, config)
	} else {
		bgs = bugs.GetAllBugs(config)
	}
	sort.Sort(BugListByMilestone(bgs))

	fmt.Printf("# Roadmap for %s\n", bugs.GetRootDir(config).GetShortName().ToTitle())
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
