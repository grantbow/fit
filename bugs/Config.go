package bugs

import (
	"errors"
	"fmt"
    "github.com/tkanos/gonfig"
	"os"
)

type Config struct{
	Dir                    Directory
	DefaultDescriptionFile string
	PMIT                   string
}

var NoConfigError = errors.New("No config provided")

/*
func (c *Config) GetDirectory() Directory {
	c.Dir := Directory.GetRootDir()
	return c.Dir
}
*/

func (c *Config) Read() (err error) {
	err = gonfig.GetConf(string(c.Dir)+"/.bug.yml" , &c); if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err.Error())
		return NoConfigError
	}
	return nil
}

