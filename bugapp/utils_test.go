package bugapp

import (
	"testing"
)

func TestGetArgument(t *testing.T) {
	Args := ArgumentList{"--tag", "bar"}
	argout, argVals := Args.GetAndRemoveArguments([]string{"--tag", "--status", "--priority", "--milestone", "--identifier"})

	if len(argout) != 0 {
		t.Error("--tag not removed from Args")
	}
	if argVals[0] != "bar" {
		t.Errorf("--tag value not set, expected %s got %s", Args[0], argVals[0])
	}
}
