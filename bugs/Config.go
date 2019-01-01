package bugs

import (
	"errors"
	//"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

type Config struct {
	// storage of location of root dir with .bug.yml
	// if we are reading the config file we already found the root
	// auto-find root dir or
	// overridden by PMIT environment variable
	BugDir string `json:"BugDir"`
	// new bug Description text template
	DefaultDescriptionFile string `json:"DefaultDescriptionFile"`
	// saves raw json files of import
	ImportXmlDump bool `json:"ImportXmlDump"`
	// import comments together or separate files
	ImportCommentsTogether bool `json:"ImportCommentsTogether"`
}

var NoConfigError = errors.New("No .bug.yml provided")

func ConfigRead(bug_yml string, c *Config) (err error) {
	temp := Config{}
	if fileinfo, err := os.Stat(bug_yml); err == nil && fileinfo.Mode().IsRegular() {
		dat, _ := ioutil.ReadFile(bug_yml)
		//fmt.Fprint(os.Stdout, dat)
		err := yaml.Unmarshal(dat, &temp)
		check(err)
		//* DefaultDescriptionFile: string,
		//      create bug template file name
		if temp.DefaultDescriptionFile != "" {
			c.DefaultDescriptionFile = temp.DefaultDescriptionFile
		}
		//* ImportXmlDump: true or false,
		//      saves raw xml as a file
		if temp.ImportXmlDump {
			c.ImportXmlDump = true
		}
		//* ImportCommentsTogether: true or false,
		//      commments save to one file or many files
		if temp.ImportCommentsTogether {
			c.ImportCommentsTogether = true
		}
	}
	return nil
}