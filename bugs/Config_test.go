package bugs

import (
	"io/ioutil"
	"testing"
)

func doconfigteststring(t *testing.T, rootDir string, bugymlfile string, config *Config, configstr *string, expected string) {
	// write
	err := ioutil.WriteFile(string(rootDir)+"/.bug.yml", []byte(bugymlfile), 0644)
	if err != nil {
		t.Error(err)
	}
	// read and stored correctly
	err = ConfigRead(".bug.yml", config, "testversion")
	if err != nil {
		t.Error(err)
	}
	if *configstr != expected {
		t.Errorf("DefaultDescriptionFile expected: %v\nGot: %v\n", expected, config.DefaultDescriptionFile)
	}
}

func doconfigtestbool(t *testing.T, rootDir string, bugymlfile string, config *Config, configbool *bool, expected bool) {
	// write
	err := ioutil.WriteFile(string(rootDir)+"/.bug.yml", []byte(bugymlfile), 0644)
	if err != nil {
		t.Error(err)
	}
	// read and stored correctly
	err = ConfigRead(".bug.yml", config, "testversion")
	if err != nil {
		t.Error(err)
	}
	if *configbool != expected {
		t.Errorf("DefaultDescriptionFile expected: %v\nGot: %v\n", expected, config.DefaultDescriptionFile)
	}
}

func TestConfigRead(t *testing.T) {
	//tests: func ConfigRead(bug_yml string, c *Config) (err error) {
	config := Config{}
	test := tester{} // from Bug_test.go
	test.Setup()
	defer test.Teardown()
	rootDir := GetRootDir(config)

	doconfigteststring(t, string(rootDir),
		"DefaultDescriptionFile: issues/bug-template.txt\n",
		&config,
		&config.DefaultDescriptionFile,
		"issues/bug-template.txt")
	config = Config{}
	doconfigtestbool(t, string(rootDir),
		"ImportXmlDump: true\n",
		&config,
		&config.ImportXmlDump,
		true)
	config = Config{}
	doconfigtestbool(t, string(rootDir),
		"ImportCommentsTogether: false\n",
		&config,
		&config.ImportCommentsTogether,
		false)
}
