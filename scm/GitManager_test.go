	commit string
	log    string
func (c GitCommit) CommitID() string {
	return c.commit
}
func (c GitCommit) LogMsg() string {
	return c.log
}
	return runCmd("git", "show", "--pretty=format:%b", c.CommitID())
}

func (c GitCommit) CommitMessage() (string, error) {
	return runCmd("git", "show", "--pretty=format:%B", "--quiet", c.CommitID())
}

type GitTester struct {
	handler SCMHandler
	workdir string

func (t GitTester) GetLogs() ([]Commit, error) {
	logs, err := runCmd("git", "log", "--oneline", "--reverse", "-z")
		wd, _ := os.Getwd()
		fmt.Fprintf(os.Stderr, "Error retrieving git logs: %s in directory %s\n", logs, wd)
	logMsgs := strings.Split(logs, "\000")
	commits := make([]Commit, len(logMsgs)-1)
func (g GitTester) AssertStagingIndex(t *testing.T, f []FileStatus) {
	for _, file := range f {
		out, err := runCmd("git", "status", "--porcelain", file.Filename)
		if err != nil {
			t.Error("Could not run git status")
		}
		expected := file.IndexStatus + file.WorkingStatus + " " + file.Filename + "\n"
		if out != expected {
			t.Error("Incorrect file status")
			t.Error("Got" + out + " not " + expected)
		}
	}
func (g GitTester) StageFile(file string) error {
	_, err := runCmd("git", "add", file)
	return err
}
func (t *GitTester) Setup() error {
		t.workdir = dir
		os.Chdir(t.workdir)
	out, err := runCmd("git", "init")
		fmt.Fprintf(os.Stderr, "Error initializing git: %s", out)
func (t GitTester) TearDown() {
	os.RemoveAll(t.workdir)
}
func (t GitTester) GetWorkDir() string {
	return t.workdir
func (m GitTester) AssertCleanTree(t *testing.T) {
		t.Error("Unexpected Output from git status (expected nothing):\n" + out)
func (m GitTester) GetManager() SCMHandler {
	return m.handler
}

	gm := GitTester{}
	gm.handler = GitManager{}
	expectedDiffs := []string{
		`
diff --git a/issues/Test-bug/Description b/issues/Test-bug/Description
new file mode 100644
index 0000000..e69de29
`, `
diff --git a/issues/Test-bug/Description b/issues/Renamed-bug/Description
similarity index 100%
rename from issues/Test-bug/Description
rename to issues/Renamed-bug/Description
`}

	runtestRenameCommitsHelper(&gm, t, expectedDiffs)
}
func TestGitFilesOutsideOfBugNotCommited(t *testing.T) {
	gm := GitTester{}
	gm.handler = GitManager{}
	runtestCommitDirtyTree(&gm, t)
}
func TestGitManagerGetType(t *testing.T) {
	manager := GitManager{}
	if getType := manager.GetSCMType(); getType != "git" {
		t.Error("Incorrect SCM Type for GitManager. Got " + getType)
}
func TestGitManagerPurge(t *testing.T) {
	gm := GitTester{}
	gm.handler = GitManager{}
	runtestPurgeFiles(&gm, t)
}

func TestGitManagerAutoclosingGitHub(t *testing.T) {
	// This test is specific to gitmanager, since GitHub
	// only supports git..
	tester := GitTester{}
	tester.handler = GitManager{Autoclose: true}
	err := tester.Setup()
		panic("Something went wrong trying to initialize git:" + err.Error())
	defer tester.TearDown()
	m := tester.GetManager()
	if m == nil {
		t.Error("Could not get manager")
		return
	}
	os.Mkdir("issues", 0755)
	runCmd("bug", "create", "-n", "Test", "bug")
	runCmd("bug", "create", "-n", "Test", "Another", "bug")
	if err = ioutil.WriteFile("issues/Test-bug/Identifier", []byte("\n\nGitHub:#TestBug"), 0644); err != nil {
		t.Error("Could not write Identifier file")
		return
	}
	if err = ioutil.WriteFile("issues/Test-Another-bug/Identifier", []byte("\n\nGITHuB:  #Whitespace   "), 0644); err != nil {
		t.Error("Could not write Identifier file")
		return
	// Commit the file, so that we can close it..
	m.Commit(bugs.Directory(tester.GetWorkDir()+"/issues"), "Adding commit")
	// Delete the bug
	os.RemoveAll(tester.GetWorkDir() + "/issues/Test-bug")
	os.RemoveAll(tester.GetWorkDir() + "/issues/Test-Another-bug")
	m.Commit(bugs.Directory(tester.GetWorkDir()+"/issues"), "Removal commit")

	commits, err := tester.GetLogs()
	if len(commits) != 2 || err != nil {
		t.Error("Error getting git logs while attempting to test GitHub autoclosing")
		return
	if msg, err := commits[1].(GitCommit).CommitMessage(); err != nil {
		t.Error("Error getting git logs while attempting to test GitHub autoclosing")
	} else {
		closing := func(issue string) bool {
			return strings.Contains(msg, "Closes #"+issue) ||
				strings.Contains(msg, ", closes #"+issue)
		}
		if !closing("Whitespace") || !closing("TestBug") {
			fmt.Printf("%s\n", msg)
			t.Error("GitManager did not autoclose Github issues")
		}