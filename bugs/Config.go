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
	// aka RootDir
	// storage of location of root dir with .bug.yml
	// if we are reading the config file we already found the root
	// auto-find root dir
	// overridden by PMIT environment variable? ** runtime only
	BugDir string `json:"BugDir"`
	// BugDir+"/.bug.yml" * if present ** runtime only
	BugYml string `json:"BugYml"`
	// Description contents for new issue or empty file (default)
	DefaultDescriptionFile string `json:"DefaultDescriptionFile"`
	// saves raw json files of import (true) or don't save (false default)
	ImportXmlDump bool `json:"ImportXmlDump"`
	// import comments together (true) or separate files (false default)
	ImportCommentsTogether bool `json:"ImportCommentsTogether"`
	// append to the program version ** runtime + append
	ProgramVersion string `json:"ProgramVersion"`
	// file name (Description default)
	DescriptionFileName string `json:"DescriptionFileName"`
	// tag_key_value (true) or tag subdir (false default)
	TagKeyValue bool `json:"TagKeyValue"`
	// tag_Field_value (true) or Field file and contents (false default)
	NewFieldAsTag bool `json:"NewFieldAsTag"`
	// tag_field_value (true) or tag_Field_value (false default)
	NewFieldLowerCase bool `json:"NewFieldLowerCase"`
}

// ErrNoConfig
var ErrNoConfig = errors.New("No .bug.yml provided")

// ConfigRead assigns values to the Config type from .bug.yml.
func ConfigRead(bugYmls string, c *Config, progVersion string) (err error) {
	temp := Config{}
	if fileinfo, err := os.Stat(bugYmls); err == nil && fileinfo.Mode().IsRegular() {
		dat, _ := ioutil.ReadFile(bugYmls)
		// didn't work well
		////wd, _ := os.Getwd()
		////c.BugYmls = wd + "/" + bugYmls
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
		return nil
	}
	return ErrNoConfig
}
