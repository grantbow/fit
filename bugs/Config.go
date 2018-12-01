package bugs

import (
	"errors"
	//"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

type Config struct {
	Dir                    Directory // location of root dir and .bug.yml
	PMIT                   string    // overrides auto-find root dir
	                                 // overridden by PMIT environment variable
	DefaultDescriptionFile string    // new bug Description text template
	ImportXmlDump          bool      // saves raw json files of import
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
	dat, err := ioutil.ReadFile(string(c.Dir)+"/.bug.yml")
	check(err)
	err = yaml.Unmarshal(dat, &c)
	check(err)
    env := os.Getenv("PMIT")
    if env != "" {
        c.PMIT = env
    }

	return nil
}

