package issues

import (
	"errors"
	//"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

// Config type holds .fit.yml configured values.
type Config struct {
	// FitDir aka RootDir or RepoDir
	// storage of location of dir containing:
	//     issues directory
	//     .fit.yml
	//     likely .git
	// overridden by FIT/PMIT environment variable    ** runtime only
	FitDir string `json:"FitDir"`
	// overridden by FIT/PMIT environment variable    ** runtime only
	FitDirName string `json:"FitDirName"`
	// save the detected directory name               ** runtime only
	ScmDir string `json:"ScmDir"`
	// save the detected scm type                     ** runtime only
	ScmType string `json:"ScmType"`
	// FitYmlDir                         * if present ** runtime only
	// Now important because this could be FitDir or ScmDir
	FitYmlDir string `json:"FitYmlDir"`
	// FitYmlDir+"/.fit.yml" or .bug.yml * if present ** runtime only
	FitYml string `json:"FitYml"`
	// Description contents for new issue or empty file (empty default)
	// relative to FitYmlDir, aka FitDir or ScmDir
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
	//* base url for notifications
	FitSite string `json:"FitSite"`
	// fit directories always recursive (true) or need -r cli option (false, default)
	MultipleFitDirs bool `json:"MultipleFitDirs"`
	// close will add tag_status_close (true) or deletes issue (false, default)
	CloseStatusTag bool `json:"CloseStatusTag"`
	// Abbreviate Identifier as Id (true) or use Identifier (false, default)
	IdAbbreviate bool `json:"IdAbbreviate"`
	// Identifier Automatic assignment (true) or not (false, default)
	IdAutomatic bool `json:"IdAutomatic"`
}

/*
list of places for
creating new commands
and where the help output is modeled:
    * bugapp/Help.go   // case lines
                       // alias line at bottom of each long description
                       // help output at bottom of the file
    * cmd/fit/main.go  // case lines
    * README.md        // includes help output generated from bottom of Help.go file, collected by running program
    * docs/FIT.md
    * docs/FAQ.md

list of places for
creating new configs:
    * issues/Config.go   // issues.Config struct
                       // ConfigRead for reading values of config file
    * issues/Config_test.go
    * cmd/fit/main_test.go
    * README.md        // includes config descriptions
                       // like comments from issues.Config struct file

synchronize

Put the list in a good place for when new configs are added

notes:

		//* NewFieldAsTag: true or false,
		//      Default Field file

		//* NewFieldLowerCase: true or false,
		//      Default Field as given
*/

// ErrNoConfig
var ErrNoConfig = errors.New("No .fit.yml provided")

// ConfigRead assigns values to the Config type from .fit.yml.
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
		//      Default false, saves raw xml as a file
		if temp.ImportXmlDump {
			c.ImportXmlDump = true
		} else {
			c.ImportXmlDump = false
		}
		//* ImportCommentsTogether: true or false,
		//      Default false, commments save to one file or many files
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
		//      Default false, use tags subdir
		if temp.TagKeyValue {
			c.TagKeyValue = true
		} else {
			c.TagKeyValue = false
		}
		//* NewFieldAsTag: true or false,
		//      Default false use Field file
		if temp.NewFieldAsTag {
			c.NewFieldAsTag = true
		} else {
			c.NewFieldAsTag = false
		}
		//* NewFieldLowerCase: true or false,
		//      Default false, use Field as given
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
		//* base url for notifications
		if temp.FitSite != "" {
			c.FitSite = temp.FitSite
		} else {
			c.FitSite = ""
		}
		//* MultipleFitDirs: true or false,
		//      Default false, need to use -r cli option
		if temp.MultipleFitDirs {
			c.MultipleFitDirs = true
		} else {
			c.MultipleFitDirs = false
		}
		//* CloseStatusTag: true or false,
		//      Default false, delete
		if temp.CloseStatusTag {
			c.CloseStatusTag = true
		} else {
			c.CloseStatusTag = false
		}
		//* IdAbbreviate: true or false,
		//      Default false, Identifier
		if temp.IdAbbreviate {
			c.IdAbbreviate = true
		} else {
			c.IdAbbreviate = false
		}
		//* IdAutomatic: true or false,
		//      Default false
		if temp.IdAutomatic {
			c.IdAutomatic = true
		} else {
			c.IdAutomatic = false
		}
		return nil // success
	} else {
		return ErrNoConfig
	}
}

func ConfigWrite(bugYmls string) (err error) {

	if fileinfo, err := os.Stat(bugYmls); err != nil && fileinfo.Mode().IsRegular() {
		err = ioutil.WriteFile(bugYmls, []byte(`
DefaultDescriptionFile: fit/DescriptionTemplate.txt
ImportXmlDump: false
ImportCommentsTogether: false
ProgramVersion:
DescriptionFileName: Description
TagKeyValue: false
NewFieldAsTag: false
NewFieldLowerCase: false
GithubPersonalAccessToken:
TwilioAccountSid:
TwilioAuthToken:
TwilioPhoneNumberFrom:
FitSite: https://github.com/<you>/<proj>/tree/master/<proj>/
MultipleFitDirs: false
CloseStatusTag: false
IdAbbreviate: false
IdAutomatic: true
`), 0644)
		// check error
		return nil
	}
	return nil
}
