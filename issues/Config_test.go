package issues

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func doconfigteststring(t *testing.T, rootDir string, bugymlfile string, config *Config, configstr *string, expected string) {
	// write
	err := ioutil.WriteFile(string(rootDir)+sops+".bug.yml", []byte(bugymlfile), 0644)
	if err != nil {
		t.Error(err)
	}
	// read and stored correctly
	err = ConfigRead(".bug.yml", config, "testversion")
	if err != nil {
		t.Error(err)
	}
	if *configstr != expected {
		t.Errorf("%s expected: %v\nGot: %v\n", bugymlfile, expected, *configstr)
	}
}

func doconfigtestbool(t *testing.T, rootDir string, bugymlfile string, config *Config, configbool *bool, expected bool) {
	// write
	err := ioutil.WriteFile(string(rootDir)+sops+".bug.yml", []byte(bugymlfile), 0644)
	if err != nil {
		t.Error(err)
	}
	// read and stored correctly
	err = ConfigRead(".bug.yml", config, "testversion")
	if err != nil {
		t.Error(err)
	}
	if *configbool != expected {
		t.Errorf("%s expected: %v\nGot: %v\n", bugymlfile, expected, *configbool)
	}
}

func TestConfigRead(t *testing.T) {
	//tests: func ConfigRead(bug_yml string, c *Config) (err error) {
	config := Config{} //// clears
	test := tester{}   // from Bug_test.go
	test.Setup()
	defer test.Teardown()
	rootDir := RootDirer(&config)

	doconfigteststring(t, string(rootDir),
		"DefaultDescriptionFile: issues.bug-template.txt\n",
		&config,
		&config.DefaultDescriptionFile,
		"issues.bug-template.txt")
	config = Config{} //// clears
	doconfigtestbool(t, string(rootDir),
		"ImportXmlDump: true\n",
		&config,
		&config.ImportXmlDump,
		true)
	config = Config{} //// clears
	doconfigtestbool(t, string(rootDir),
		"ImportCommentsTogether: false\n",
		&config,
		&config.ImportCommentsTogether,
		false)
	config = Config{} //// clears
	progVersion := "testversion"
	doconfigteststring(t, string(rootDir),
		"ProgramVersion: special\n",
		&config,
		&config.ProgramVersion,
		fmt.Sprintf("%sspecial", progVersion))
	config = Config{} //// clears
	doconfigteststring(t, string(rootDir),
		"DescriptionFileName: description.txt\n",
		&config,
		&config.DescriptionFileName,
		"description.txt")
	config = Config{} //// clears
	doconfigtestbool(t, string(rootDir),
		"\n",
		&config,
		&config.TagKeyValue,
		false)
	config = Config{} //// clears
	doconfigtestbool(t, string(rootDir),
		"TagKeyValue: true\n",
		&config,
		&config.TagKeyValue,
		true)
	config = Config{} //// clears
	doconfigtestbool(t, string(rootDir),
		"\n",
		&config,
		&config.CloseStatusTag,
		false)
	config = Config{} //// clears
	doconfigtestbool(t, string(rootDir),
		"CloseStatusTag: true\n",
		&config,
		&config.CloseStatusTag,
		true)
	config = Config{} //// clears
	doconfigtestbool(t, string(rootDir),
		"\n",
		&config,
		&config.IdAbbreviate,
		false)
	config = Config{} //// clears
	doconfigtestbool(t, string(rootDir),
		"IdAbbreviate: true\n",
		&config,
		&config.IdAbbreviate,
		true)
	config = Config{} //// clears
	doconfigtestbool(t, string(rootDir),
		"\n",
		&config,
		&config.IdAutomatic,
		false)
	config = Config{} //// clears
	doconfigtestbool(t, string(rootDir),
		"IdAutomatic: true\n",
		&config,
		&config.IdAutomatic,
		true)
}
