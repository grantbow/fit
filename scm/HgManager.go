package scm

import (
	"fmt"
	"github.com/driusan/bug/bugs"
	"os/exec"
)

// HgManager is a struct for a mercurial (hg) software configuration manager.
type HgManager struct{}

// Purge would give the hg command to purge files but this is not supported.
func (a HgManager) Purge(dir bugs.Directory) error {
	return UnsupportedType("Purge is not supported under Hg. Sorry!")
}

// Commit gives the hg command to commit files.
func (a HgManager) Commit(dir bugs.Directory, commitMsg string, config bugs.Config) error {
	cmd := exec.Command("hg", "addremove", string(dir))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Could not add issues to be committed: %s?\n", err.Error())
		return err
	}

	cmd = exec.Command("hg", "commit", string(dir), "-m", commitMsg)
	// stdout and stderr not captured in HgManager_test.go runtestCommitDirtyTree()
	if err := cmd.Run(); err != nil {
		//fmt.Printf("post 2 runtestCommitDirtyTree error %v\n", err) // 255 when $?=1 and stdout text "nothing changed" present
		fmt.Printf("No new issues to commit.\n") // assumes this error, same for GitManager.go
		return err
	}
	return nil
}

// GetSCMType returns hg.
func (a HgManager) GetSCMType() string {
	return "hg"
}
