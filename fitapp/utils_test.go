package fitapp

import (
	"testing"
)

func TestGetArgument(t *testing.T) {
	args := argumentList{"foo", "--github", "bar"}
	if "bar" != args.GetArgument("--github", "") {
		t.Error("GetArgument should return the string value")
	}
	args = argumentList{"foo", "other", "bar"}
	if "" != args.GetArgument("--github", "") {
		t.Error("GetArgument should return the default value")
	}
}

func TestGetAndRemoveArguments(t *testing.T) {
	Args := argumentList{"--tag", "bar"}
	argout, argVals := Args.GetAndRemoveArguments([]string{"--tag", "--status", "--priority", "--milestone", "--identifier"})

	if len(argout) != 0 {
		t.Error("--tag not removed from Args")
	}
	if argVals[0] != "bar" {
		t.Errorf("--tag value not set, expected %s got %s", Args[0], argVals[0])
	}
}

func TestSkipRootCheck(t *testing.T) {
	list := []string{}
	if !SkipRootCheck(&list) {
		t.Error("SkipRootCheck 0 should return true")
	}
	list = []string{"bin"}
	if !SkipRootCheck(&list) {
		t.Error("SkipRootCheck 1 should return true")
	}
	list = []string{"bin", "help"}
	if !SkipRootCheck(&list) {
		t.Error("SkipRootCheck 2 help should return true")
	}
	list = []string{"bin", "nohelp"}
	if SkipRootCheck(&list) {
		t.Error("SkipRootCheck 2 nohelp should return false")
	}
	list = []string{"bin", "ver", "--help"}
	if !SkipRootCheck(&list) {
		t.Error("SkipRootCheck 3 --help should return true")
	}
}
