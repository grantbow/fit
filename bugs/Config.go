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
	BugDir                 string `json:"BugDir"`
	// new bug Description text template
	DefaultDescriptionFile string `json:"DefaultDescriptionFile"`
	// saves raw json files of import
	ImportXmlDump          bool   `json:"ImportXmlDump"`
}

var NoConfigError = errors.New("No .bug.yml provided")

func check(e error) () {
	if e != nil {
	//	fmt.Fprintf(os.Stderr, "err: %s", err.Error())
	//	return NoConfigError
		panic (e)
	}
}

func (c *Config) Read() (err error) {
	dat, err := ioutil.ReadFile(string(c.BugDir)+"/.bug.yml")
	check(err)
	err = yaml.Unmarshal(dat, &c)
	check(err)
    //env := os.Getenv("PMIT")
    //if env != "" {
    //    c.BugDir = env
    //}

	return nil
}

