package issues

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//var dops = Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

// An Issue
type Issue struct {
	Dir                 Directory
	modtime             int
	descFile            *os.File
	DescriptionFileName string
	TagArray            []TagKeyValue
}

// TagBoolTrue only has a string key.
// Implied values are true/present and false/absent.
type TagBoolTrue string

// TagKeyValue were added after type Tag (renamed to TagBoolTrue)
type TagKeyValue struct {
	key  string
	file string
}

// Four hard coded "fields" were allowed to have the four specific keys
// and values. These files were originally kept in a "tags" subdirectory.
//     []string{"Status", "Priority", "Milestone", "Identifier"}
//     (optionally Id for Identifier)

// Comment is the struct type of a unit of discussion about an issue.
type Comment struct {
	Author string
	Time   time.Time
	Body   string
	Order  int
	Xml    []byte
}

// ErrNoDescription defines a new error.
var ErrNoDescription = errors.New("No description provided")

// ErrNotFound defines a new error.
var ErrNotFound = errors.New("Could not find issue")

// TitleToDirString converts string from title to directory
func TitleToDirString(title string) string {
	// replace non-matching valid characters with _
	// for user entered strings
	re := regexp.MustCompile("[^a-zA-Z0-9_ -]+")
	s := re.ReplaceAllString(title, "_")

	replaceWhitespaceWithUnderscore := func(match string) string {
		return strings.Replace(match, " ", "_", -1)
	}
	//replaceDashWithMore := func(match string) string {
	//	if strings.Count(match, " ") > 0 {
	//		return match
	//	}
	//	return "-" + match
	//}

	// Replace sequences of dashes with 1 more dash,
	// as long as there's no whitespace around them
	//commented because of unexpected behavior
	//re = regexp.MustCompile("([\\s]*)([-]+)([\\s]*)")
	//s = re.ReplaceAllStringFunc(s, replaceDashWithMore)

	// If there are dashes with whitespace around them,
	// replace the whitespace with underscores
	// This is a two step process, because the whitespace
	// can independently be on either side, so it's difficult
	// to do with 1 regex..
	re = regexp.MustCompile("([\\s]+)([-]+)")
	s = re.ReplaceAllStringFunc(s, replaceWhitespaceWithUnderscore)
	re = regexp.MustCompile("([-]+)([\\s]+)")
	s = re.ReplaceAllStringFunc(s, replaceWhitespaceWithUnderscore)

	s = strings.Replace(s, " ", "-", -1)
	return s
}

// TitleToDir returns a Directory from a string argument.
func TitleToDir(title string) Directory {
	return Directory(TitleToDirString(title))
}

// ShortTitleToDir truncates a title to 25 characters.
func ShortTitleToDir(title string) Directory {
	if len(title) > 25 {
		return TitleToDir(title[:25]) // TODO: remove leading or trailing _ or -
    }
    return TitleToDir(title)
}

// Direr returns the directory of an issue.
func (b Issue) Direr() Directory {
	return b.Dir
}

// LoadIssue sets an issue's directory, modtime and DescriptionFileName
// and enforces IdAutomatic.
func (b *Issue) LoadIssue(dir Directory, config Config) {
	b.Dir = dir
	b.modtime = int((dir.ModTime()).Unix())
	b.DescriptionFileName = config.DescriptionFileName
	if config.IdAutomatic {
		if id := b.Identifier(); id != "" {
			//fmt.Printf("debug id : (%v)\n", id)
			i := getIdNext(config)
			if i != -1 {
				b.SetIdentifier(strconv.Itoa(i), config)
			}
		}
	}
}

func getIdNext(config Config) int {
	// config.FitDir/.fit_idnext_1001
	files, err := filepath.Glob(config.FitDir + sops + ".fit_idnext_*")
	//fmt.Printf("debug Found %v files: %v\n", len(files), strings.Join(files, ", "))
	if err == nil && len(files) > 1 {
		fmt.Printf("Found %v files: %v\n", len(files), strings.Join(files, ", "))
		os.Exit(1) // error
	} else if err == nil && len(files) == 1 {
		parts := strings.Split(files[0], "_")
		val := parts[len(parts)-1]
		i, _ := strconv.Atoi(val)
		//removeIdNext(i, config)
		//writeIdNext(i+1, config)
		renameIdNext(i, i+1, config)
		return i
	} else if err != nil && len(files) == 0 {
		// missing
		i := 1001
		writeIdNext(i+1, config)
		return i
	} else {
		if err != nil {
			fmt.Printf("Found %v files: %v\nError : %v\n", len(files), strings.Join(files, ", "), err.Error())
		} else {
			fmt.Printf("Found %v files: %v\nNo error returned\n", len(files), strings.Join(files, ", "))
		}
		return -1
		//errors.New("file %s%s.fit_idnext_<i>", config.FitDir, sops)
	}
	return -1
}

func writeIdNext(j int, config Config) {
	content := ""
	ioutil.WriteFile(config.FitDir+sops+".fit_idnext_"+strconv.Itoa(j), []byte(content+"\n"), 0644)
	// TODO: scm.add
}

//func removeIdNext(k int, config Config) {
//	os.Remove(config.FitDir + sops + ".fit_idnext_" + strconv.Itoa(k))
//	// TODO: scm.add
//}

func renameIdNext(i int, j int, config Config) {
	os.Rename(config.FitDir+sops+".fit_idnext_"+strconv.Itoa(i),
		config.FitDir+sops+".fit_idnext_"+strconv.Itoa(j))
	// TODO: scm.add
}

// Title returns a string with the name of an issue and
// optionally present Identifier, Status, Priority and tags.
func (b Issue) Title(options string) string {
	// options indicate what should be formatted and returned with the title.
	var hasOption = func(o string) bool {
		return strings.Contains(options, o)
	}

	title := b.Dir.ShortNamer().ToTitle()

	if id := b.Identifier(); hasOption("identifier") && id != "" {
		title = fmt.Sprintf("(%s) %s", id, title)
	}
	if hasOption("tags") {
		tags := b.StringTags()
		if len(tags) > 0 {
			title += fmt.Sprintf(" (%s)", strings.Join(tags, ", "))
		}
	}

	priority := hasOption("priority") && b.Priority() != ""
	status := hasOption("status") && b.Status() != ""
	if options == "" {
		priority = false
		status = false
	}

	if priority && status {
		title += fmt.Sprintf(" (Status: %s; Priority: %s)", b.Status(), b.Priority())
	} else if priority {
		title += fmt.Sprintf(" (Priority: %s)", b.Priority())
	} else if status {
		title += fmt.Sprintf(" (Status: %s)", b.Status())
	}
	return title
}

// Description returns a string of an issue.
func (b Issue) Description() string {
	//does filepath.FromSlash() really work?
	df := string(b.Dir) + sops + b.DescriptionFileName
	value := ""
	if _, staterr := os.Stat(df); staterr == nil {
		v, readerr := ioutil.ReadFile(df)
		//fmt.Printf("debug %v %v \n", b.DescriptionFileName, v)
		if readerr == nil {
			value = string(v)
		} else if perr, ok := staterr.(*os.PathError); ok {
			switch perr.Err.(syscall.Errno) {
			// os.PathError Op, Path, Err
			case syscall.ENOENT:
				return string(value)
			default:
				panic("Unhandled error " + fmt.Sprint(reflect.TypeOf(readerr)) + " " + readerr.Error())
			}
		}
	}
	//if string(value) == "" {
	//	return "(No description provided.)\n"
	//}
	return string(value)
}

// SetDescription writes the Description file of an issue.
func (b *Issue) SetDescription(val string, config Config) error {
	dir := b.Direr()
	//fmt.Printf("aha %s\n", config.DescriptionFileName)
	b.DescriptionFileName = config.DescriptionFileName

	//return ioutil.WriteFile(filepath.FromSlash(string(dir)+"/"+b.DescriptionFileName), []byte(val+"\n"), 0644)
	return ioutil.WriteFile(string(dir)+sops+b.DescriptionFileName, []byte(val+"\n"), 0644)
}

// RemoveTag deletes a tag file of an issue.
func (b *Issue) RemoveTag(tag TagBoolTrue, config Config) {
	if dir := b.Direr(); dir != "" {
		os.Remove(string(dir) + sops + "tags" + sops + string(tag))
		files, err := filepath.Glob(string(dir) + sops + "tag_" + string(tag) + "*")
		if err == nil {
			for _, x := range files {
				os.Remove(x)
			}
		}
	} else {
		// no b.Dir - should not happen any more
		// still good to check just in case
		fmt.Printf("Error removing tag: %s", tag)
	}
}

// TagIssue writes an empty *boolean* tag file: key, no value
func (b *Issue) TagIssue(tag TagBoolTrue, config Config) {
	var key string
	if dir := b.Direr(); dir != "" {
		if config.NewFieldLowerCase {
			key = strings.ToLower(string(tag))
		} else {
			key = string(tag)
		}
		if config.TagKeyValue == true {
			ioutil.WriteFile(string(dir)+sops+"tag_"+key, []byte(""), 0644)
		} else {
			os.Mkdir(string(dir)+sops+"tags"+sops, 0755)
			ioutil.WriteFile(string(dir)+sops+"tags"+sops+key, []byte(""), 0644)
		}
	} else {
		fmt.Printf("Error tagging issue: %s", key)
	}
}

// RemoveComment deletes a comment file of an issue.
func (b *Issue) RemoveComment(comment Comment) {
	if dir := b.Direr(); dir != "" {
		//os.Remove(filepath.FromSlash(string(dir) + "/comment-" + string(ShortTitleToDir(string(comment.Body)))))
		os.Remove(string(dir) + sops + "comment-" + string(ShortTitleToDir(string(comment.Body))))
	} else {
		fmt.Printf("Error removing comment: %s", comment.Body)
	}
}

// CommentIssue writes a text file for an issue.
func (b *Issue) CommentIssue(comment Comment, config Config) {
	if dir := b.Direr(); dir != "" {
		//os.Mkdir(filepath.FromSlash(string(dir)+"/"), 0755)
		commenttext := []byte(comment.Body + "\n")
		if config.ImportCommentsTogether { // not efficient but ok for now
			data, err := ioutil.ReadFile(string(dir) + sops + "comments")
			check(err)
			commentappend := []byte(fmt.Sprintf("%s%s%s", data, "\n", commenttext))
			werr := ioutil.WriteFile(string(dir)+sops+"comments", commentappend, 0644)
			check(werr)
		} else {
			werr := ioutil.WriteFile(string(dir)+sops+"comment-"+string(ShortTitleToDir(string(comment.Body))), commenttext, 0644)
			check(werr)
		}
	} else {
		fmt.Printf("Error commenting issue: %s", comment.Body)
	}
}

// ViewIssue outputs an issue.
func (b Issue) ViewIssue() {
	// Fields and tags could be more general if architected differently.
	if identifier := b.Identifier(); identifier != "" {
		fmt.Printf("Identifier: %s\n", identifier)
	}

	fmt.Printf("Title: %s\n", b.Title(""))
	fmt.Printf("Description: %s\n", b.Description())

	if status := b.Status(); status != "" {
		fmt.Printf("Status: %s\n", status)
	}
	if priority := b.Priority(); priority != "" {
		fmt.Printf("Priority: %s\n", priority)
	}
	if milestone := b.Milestone(); milestone != "" {
		fmt.Printf("Milestone: %s\n", milestone)
	}
	if tags := b.StringTags(); len(tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join([]string(tags), ", "))
	}

}

// StringTags gets all Tags and returns []string.
func (b Issue) StringTags() []string {
	tags := b.Tags()
	tagout := []string{}
	for _, tag := range tags {
		tagout = append(tagout, string(tag))
	}
	sort.Strings(tagout)
	return tagout
}

// HasTag returns if an issue is assigned a tag.
func (b Issue) HasTag(tag TagBoolTrue) bool {
	allTags := b.Tags()
	for _, issueTag := range allTags {
		if issueTag == tag {
			return true
		}
	}
	return false
}

// created b.tager for similar needs of {issues/Issue.go:Tags, issues/Issue.go:SetField}, also issues/Issue.go:liners

// tager takes a file name, returns the key, value,
// bool if value is located in the name,
// bool if value is located in file contents, error
func (b Issue) tager(abspath string) (string, string, bool, bool, error) {
	dir := b.Direr()
	//hit := withtagsubdirfile.Name() // simple for tags subdir
	//   aka abspath
	key := ""
	value := ""
	tagName := false
	tagContents := false
	//var presentLines []string
	segments := strings.Split(abspath, string(os.PathSeparator)) // path separator
	// no glob - won't find tag_key_value
	parts := strings.Split(segments[len(segments)-1], "_")
	if len(parts) <= 1 {
		return key, value, tagName, tagContents, errors.New("tag has no key or value")
	} else if len(parts) == 2 {
		key = parts[1]
		field, err := ioutil.ReadFile(string(dir) + sops + "tag_")
		if err == nil {
			value = ([]string(strings.Split(string(field), "\n")))[0] // tag_Status file contents overrides "Status" file contents
			// assumes value is ok, not false
			tagContents = true
		}
	} else if len(parts) >= 3 {
		key = parts[1] // tag_Status_ file overrides "Status" file contents
		//presentLines = append(presentLines, strings.Join(parts[2:], "_")) // tag_Status_ file overrides "Status" file contents
		value = strings.Join(parts[2:], "_") // tag_Status_ file overrides "Status" file contents
		// assumes value is ok, not false
		tagName = true
	}
	if value == "false" {
		return "", "", false, false, errors.New("tag has no key or value")
    }
    return strings.ToLower(key), strings.ToLower(value), tagName, tagContents, nil // key, value, tagName, tagContents, err
}

// Tags returns an issue's array of tags.
func (b Issue) Tags() []TagBoolTrue {
	dir := b.Direr()
	tags := []string{}
	// fields
	for _, k := range []string{"Status", "Priority", "Milestone", "Identifier"} {
		if v := b.fielder(k); v != "" {
			keyvalue := strings.ToLower(k) + ":" + strings.ToLower(v)
			if !findArrayString(tags, keyvalue) {
				//fmt.Printf("keyvalue 1 %v\n", keyvalue)
				tags = append(tags, keyvalue)
			}
		}
	}
	// look in the <issue>/tags subdir
	withtagsubdir, errsubdir := ioutil.ReadDir(string(dir) + sops + "tags" + sops) // returns []os.FileInfo
	// look in the <issue> dir for tag_<key> and tag_<key>_<value>
	withtagfile, errtagfile := filepath.Glob(string(dir) + sops + "tag_*") // returns []string
	if len(tags) == 0 && errsubdir != nil && errtagfile != nil {
		return nil
	}

	//fmt.Printf("withtagsubdir %v\n", withtagsubdir)
	for _, withtagsubdirfile := range withtagsubdir {
		name := strings.ToLower(withtagsubdirfile.Name())
		if !findArrayString(tags, name) {
			//fmt.Printf("keyvalue 2 %v\n", name)
			tags = append(tags, name)
		}
	}
	//fmt.Printf("withtagfile %v\n", withtagfile)
	for _, withtagfilefile := range withtagfile {
		k, v, _, _, err := b.tager(withtagfilefile)
		key := strings.ToLower(k)
		value := strings.ToLower(v)
		//fmt.Printf("key %v value %v e %v\n", key, value, err)
		if err == nil && value != "false" {
			if value != "" && !findArrayString(tags, key+":"+value) {
				// only add unique TagBoolTrue
				//fmt.Printf("keyvalue 3 %v\n", key+":"+value)
				tags = append(tags, key+":"+value)
			} else if value == "" && !findArrayString(tags, key) {
				//fmt.Printf("keyvalue 3.5 %v\n", key+":"+value)
				//fmt.Printf("keyvalue 4 %v\n", key)
				tags = append(tags, key)
			}
		}
	}
	// sort
	sort.Strings(tags)

	tagtags := []TagBoolTrue{}
	for _, x := range tags {
		tagtags = append(tagtags, TagBoolTrue(x))
	}
	return tagtags
}

//byIssue allows sort.Sort(byIssue(issues)). Requirements are
//     func (issue) Len() int,
//     type
//     and the (byIssue) {Len, Swap, Less} functions
// see also List.go for type byDir

// Len allows sort.Sort(byIssue(issues))
func (b Issue) Len() int {
	return b.modtime // time.Format(time.UnixNano(t.modtime).UnixNano())
}

type byIssue []Issue

func (b byIssue) Len() int {
	return len(b) // time.Format(time.UnixNano(t.modtime).UnixNano())
}
func (b byIssue) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b byIssue) Less(i, j int) bool {
	return (b[i]).Len() < (b[j]).Len()
}

// findArrayString returns a bool if looking is an element of arr.
func findArrayString(arr []string, looking string) bool {
	for loop := range arr {
		if arr[loop] == looking {
			return true
		}
	}
	return false
}

// fielder reads and returns the string value from the file of an issue.
func (b Issue) fielder(fieldName string) string {
	lines := b.liners(fieldName)
	//fmt.Printf("debug fielder lines %v\n", lines)
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
    return ""
}

// created b.liners for similar needs of {issues/Issue.go:fielder, issues/Issue.go:SetField}
// liners does the work for fielder with tags and extra lines.
func (b Issue) liners(fieldName string) []string {
	dirr := b.Direr()
	dir := string(dirr)
	lines := []string{}
	withtagfile, errtagfile := filepath.Glob(dir + sops + "tag_*") // returns []string
	if errtagfile == nil {
		for _, withtagfilefile := range withtagfile {
			//fmt.Printf("debug liners %v\n", withtagfilefile)
			// tager finds values
			k, v, _, _, err := b.tager(withtagfilefile)
			//fmt.Printf("debug liners k %v v %v \n", k, v)
			// try tag_(K,k)ey_value or tag_(K,k)ey with contents value
			if err == nil &&
				(k == fieldName || k == strings.ToLower(fieldName)) {
				lines = []string{v}
				return lines
			}
		}
	}
	// try (F)ieldName
	field, err := ioutil.ReadFile(dir + sops + fieldName)
	if err == nil {
		lines = strings.Split(string(field), "\n")
		return lines
	}
	// try lower (f)ieldname
	field, err = ioutil.ReadFile(dir + sops + strings.ToLower(fieldName))
	if err == nil {
		lines = strings.Split(string(field), "\n")
		return lines
	}
	return lines
}

// SetField writes the string value to the file of an issue.
// NewFieldAsTag and NewFieldLowerCase are respected
func (b Issue) SetField(fieldName string, value string, config Config) error { // TODO: complete func for config tag files : paused with tagName, tagContents, fileContents
	// using Status for fielName string example in comments
	dir := b.Direr()
	//possible locations
	tagName := false
	tagContents := false
	fileContents := false
	// try "Status" file
	presentLines := b.liners(fieldName) // var presentLines []string
	if len(presentLines) > 0 {
		fileContents = true
	}
	// try tag_Status* files
	withtagfile, errtagfile := filepath.Glob(string(dir) + sops + "tag_" + fieldName + "*") // returns []string
	errfind := errtagfile
	// two cases, ie tag_Status_closed or tag_Status contains closed
	if errtagfile == nil {
		for _, withtagfilefile := range withtagfile {
			presentvalue := ""
			_, presentvalue, tagName, tagContents, errfind = b.tager(withtagfilefile)
			if errfind == nil {
				presentLines = []string{presentvalue}
			}
			//segments := strings.Split(withtagfile[f], "/") // path separator
			//parts := strings.Split(segments[len(segments)-1], "_")
			//if len(parts) == 2 {
			//	field, errpresent = ioutil.ReadFile(string(dir) + "/tag_" + fieldName)
			//	if errpresent == nil {
			//		presentLines = strings.Split(string(field), "\n") // tag_ file contents overrides "Status" file contents
			//		tagContents = true
			//	}
			//} else if len(parts) >= 3 {
			//	presentLines = append(presentLines, strings.Join(parts[2:], "_")) // tag_ file overrides "Status" file contents
			//	tagName = true
			//}
		}
	}
	_ = tagName
	_ = tagContents
	_ = fileContents

	newValue := ""
	if len(presentLines) >= 1 {
		// If there were 0 or 1 present lines, overwrite 1 and maintain the rest
		presentLines[0] = value
		newValue = strings.Join(presentLines, "\n")
	} else {
		newValue = value
	}

	var err error
	if config.NewFieldAsTag == true {
		if config.NewFieldLowerCase == true {
			err = ioutil.WriteFile(string(dir)+sops+"tag_"+strings.ToLower(fieldName)+"_"+strings.ToLower(TitleToDirString(newValue)), []byte(""), 0644)
		} else {
			err = ioutil.WriteFile(string(dir)+sops+"tag_"+fieldName+"_"+TitleToDirString(newValue), []byte(""), 0644)
		}
	} else {
		err = ioutil.WriteFile(string(dir)+sops+fieldName, []byte(newValue), 0644)
	}
	if err != nil {
		return err
	}
    return nil
}

// Status returns the string from the Status file of an issue.
func (b Issue) Status() string {
	return b.fielder("Status")
}

// SetStatus writes the Status file to an issue.
func (b Issue) SetStatus(newStatus string, config Config) error {
	return b.SetField("Status", newStatus, config)
}

// Priority returns the string from the Priority file of an issue.
func (b Issue) Priority() string {
	return b.fielder("Priority")
}

// SetPriority writes the Priority file to an issue.
func (b Issue) SetPriority(newValue string, config Config) error {
	return b.SetField("Priority", newValue, config)
}

// Milestone returns the string from the Milestone file of an issue.
func (b Issue) Milestone() string {
	return b.fielder("Milestone")
}

// SetMilestone writes the Milestone file to an issue.
func (b Issue) SetMilestone(newValue string, config Config) error {
	return b.SetField("Milestone", newValue, config)
}

// Identifier returns the string from the Identifier of an issue.
func (b Issue) Identifier() string {
	// try to read both
	i := b.fielder("Id")
	//fmt.Printf("debug b.fielder %v\n", i)
	if i != "" {
		return i
	} else if i = b.fielder("Identifier"); i != "" {
		return i
	}
	return ""
}

// SetIdentifier writes the Identifier file to an issue.
func (b Issue) SetIdentifier(newValue string, config Config) error {
	if config.IdAbbreviate {
		return b.SetField("Id", newValue, config)
	}
    return b.SetField("Identifier", newValue, config)
}

// New prepares an issue directory.
func New(title string, config Config) (*Issue, error) {
	expectedDir := FitDirer(config) + Directory(os.PathSeparator) + TitleToDir(title)
	err := os.Mkdir(string(expectedDir), 0755)
	if err != nil {
		return nil, err
	}
	return &Issue{Dir: expectedDir}, nil
}
