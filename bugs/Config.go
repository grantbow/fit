package bugs

import (
	"errors"
	//"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

// Config type holds .bug.yml configured values.
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
	// append to the program version
	ProgramVersion string `json:"ProgramVersion"`
	// override the default file name
	DescriptionFileName string `json:"DescriptionFileName"`
	// tag_key_value or tag subdir (default)
	TagKeyValue bool `json:"TagKeyValue"`
	// tag_Field_value or Field file and contents (default)
	NewFieldAsTag bool `json:"NewFieldAsTag"`
	// tag_field_value or tag_Field_value (default)
	NewFieldLowerCase bool `json:"NewFieldLowerCase"`
}

// ErrNoConfig is a new error.
var ErrNoConfig = errors.New("No .bug.yml provided")

// ConfigRead assigns values to the Config type from .bug.yml.
func ConfigRead(bugYml string, c *Config, progVersion string) (err error) {
	temp := Config{}
	if fileinfo, err := os.Stat(bugYml); err == nil && fileinfo.Mode().IsRegular() {
		dat, _ := ioutil.ReadFile(bugYml)
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
		} else {
			c.ImportXmlDump = false
		}
		//* ImportCommentsTogether: true or false,
		//      commments save to one file or many files
		if temp.ImportCommentsTogether {
			c.ImportCommentsTogether = true
		} else {
			c.ImportCommentsTogether = false
		}
		//* ProgramVersion
		c.ProgramVersion = progVersion
		if temp.ProgramVersion != "" {
			c.ProgramVersion = c.ProgramVersion + temp.ProgramVersion
		}
		//* DescriptionFileName
		if temp.DescriptionFileName != "" {
			c.DescriptionFileName = temp.DescriptionFileName
		} else {
			c.DescriptionFileName = "Description"
		}
		//* TagKeyValue: true or false,
		//      Default tags subdir
		if temp.TagKeyValue {
			c.TagKeyValue = true
		} else {
			c.TagKeyValue = false
		}
		//* NewFieldAsTag: true or false,
		//      Default Field file
		if temp.NewFieldAsTag {
			c.NewFieldAsTag = true
		} else {
			c.NewFieldAsTag = false
		}
		//* NewFieldLowerCase: true or false,
		//      Default Field as given
		if temp.NewFieldLowerCase {
			c.NewFieldLowerCase = true
		} else {
			c.NewFieldLowerCase = false
		}
	}
	return nil
}
