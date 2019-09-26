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
	// overridden by FIT/PMIT environment variable? ** runtime only
	BugDir string `json:"BugDir"`
	// BugDir+"/.bug.yml" * if present ** runtime only
	BugYml string `json:"BugYml"`
	// Description contents for new issue or empty file (default)
	DefaultDescriptionFile string `json:"DefaultDescriptionFile"`
	// saves raw json files of import (true) or don't save (false, default)
	ImportXmlDump bool `json:"ImportXmlDump"`
	// import comments together (true) or separate files (false, default)
	ImportCommentsTogether bool `json:"ImportCommentsTogether"`
	// append to the program version ** runtime + append
	ProgramVersion string `json:"ProgramVersion"`
	// file name (Description is the default) set in main.go
	DescriptionFileName string `json:"DescriptionFileName"`
	// tag_key_value (true) or tag subdir (false, default)
	TagKeyValue bool `json:"TagKeyValue"`
	// tag_Field_value (true) or Field file and contents (false, default)
	NewFieldAsTag bool `json:"NewFieldAsTag"`
	// tag_field_value (true) or tag_Field_value (false, default)
	NewFieldLowerCase bool `json:"NewFieldLowerCase"`
	// github.com/settings/tokens
	GithubPersonalAccessToken string `json:"GithubPersonalAccessToken"`
	//* twilio.com/console "Dashboard" has the "account sid" public acct identifier
	TwilioAccountSid string `json:"TwilioAccountSid"`
	//* twilio "Auth Token" is the "Rest API Key" is for access
	TwilioAuthToken string `json:"TwilioAuthToken"`
	//* your twilio number
	TwilioPhoneNumberFrom string `json:"TwilioPhoneNumberFrom"`
	//* your issues site url for notifications
	IssuesSite string `json:"IssuesSite"`
}

/*
create list of places for help:
    * bugapp/Help.go   // case lines
                       // alias line at bottom of each long description
                       // help output at bottom of the file
    * main.go          // case lines
    * README.md        // includes help output generated from bottom of Help.go file, collected by running program
    * FIT.md
    * FAQ.md

create list of places for config:
    * bugs/Config.go   // bugs.Config struct
                       // ConfigRead for reading values of config file
    * README.md        // includes config descriptions
                       // like comments from bugs.Config struct file

synchronize

Put the list in a good place for when new configs are added

notes:

		//* NewFieldAsTag: true or false,
		//      Default Field file

		//* NewFieldLowerCase: true or false,
		//      Default Field as given
*/

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
		//* github.com/settings/tokens for import of projects and private repos
		if temp.GithubPersonalAccessToken != "" {
			c.GithubPersonalAccessToken = temp.GithubPersonalAccessToken
		} else {
			c.GithubPersonalAccessToken = ""
		}
		//* twilio.com/console "Dashboard" has the "account sid" public acct identifier
		if temp.TwilioAccountSid != "" {
			c.TwilioAccountSid = temp.TwilioAccountSid
		} else {
			c.TwilioAccountSid = ""
		}
		//* twilio "Auth Token" is the "Rest API Key" is for access
		if temp.TwilioAuthToken != "" {
			c.TwilioAuthToken = temp.TwilioAuthToken
		} else {
			c.TwilioAuthToken = ""
		}
		//* your twilio number
		if temp.TwilioPhoneNumberFrom != "" {
			c.TwilioPhoneNumberFrom = temp.TwilioPhoneNumberFrom
		} else {
			c.TwilioPhoneNumberFrom = ""
		}
		//* your issues site url for notifications
		if temp.IssuesSite != "" {
			c.IssuesSite = temp.IssuesSite
		} else {
			c.IssuesSite = ""
		}
		return nil
	}
	return ErrNoConfig
}
