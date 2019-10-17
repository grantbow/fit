package bugs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"syscall"
	"time"
)

//var dops = Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

type TagKeyValue struct {
	key  string
	file string
}

// Bug is the type of an issue.
// The fields are Dir and descFile.
type Bug struct {
	Dir                 Directory
	modtime             int
	descFile            *os.File
	DescriptionFileName string
	TagArray            []TagKeyValue
}

// Tag is the first type of an issue identifier.
// There is only a string key.
// Values were not supported originally
// so there is an implied true/present false/absent value.
type TagBoolTrue string

// Comment is the struct type of a unit of discussion about an issue.
// The fields are Author, Time, Body, Order and Xml.
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
var ErrNotFound = errors.New("Could not find bug")

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
	} else {
		return TitleToDir(title)
	}
}

// Direr returns the directory of an issue.
func (b Bug) Direr() Directory {
	return b.Dir
}

// LoadBug assigns a directory to an issue.
func (b *Bug) LoadBug(dir Directory, config Config) {
	b.Dir = dir
	b.modtime = int((dir.ModTime()).Unix())
	b.DescriptionFileName = config.DescriptionFileName
}

// Title returns a string with the name of an issue and optionally present Status and Priority.
func (b Bug) Title(options string) string {
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
func (b Bug) Description() string {
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
func (b *Bug) SetDescription(val string, config Config) error {
	dir := b.Direr()
	//fmt.Printf("aha %s\n", config.DescriptionFileName)
	b.DescriptionFileName = config.DescriptionFileName

	//return ioutil.WriteFile(filepath.FromSlash(string(dir)+"/"+b.DescriptionFileName), []byte(val+"\n"), 0644)
	return ioutil.WriteFile(string(dir)+sops+b.DescriptionFileName, []byte(val+"\n"), 0644)
}

// RemoveTag deletes a tag file of an issue.
func (b *Bug) RemoveTag(tag TagBoolTrue, config Config) {
	if dir := b.Direr(); dir != "" {
		//os.Remove(filepath.FromSlash(string(dir) + "/tags/" + string(tag)))
		os.Remove(string(dir) + sops + "tags" + sops + string(tag))
		files, err := filepath.Glob(string(dir) + "tag_" + string(tag) + "")
		if err == nil {
			for _, x := range files {
				os.Remove(x)
			}
		}
	} else {
		fmt.Printf("Error removing tag: %s", tag)
		// no b.Dir?
		// this was added during debugging when this happened sometimes
		// it's still a good idea to check just in case
	}
}

// TagBug writes an empty tag file.
func (b *Bug) TagBug(tag TagBoolTrue, config Config) {
	if dir := b.Direr(); dir != "" {
		if config.TagKeyValue == true {
			//ioutil.WriteFile(filepath.FromSlash(string(dir)+"/tag_"+string(tag)), []byte(""), 0644)
			ioutil.WriteFile(string(dir)+sops+"tag_"+string(tag), []byte(""), 0644)
		} else {
			//os.Mkdir(filepath.FromSlash(string(dir)+"/tags/"), 0755)
			os.Mkdir(string(dir)+sops+"tags"+sops, 0755)
			//ioutil.WriteFile(filepath.FromSlash(string(dir)+"/tags/"+string(tag)), []byte(""), 0644)
			ioutil.WriteFile(string(dir)+sops+"tags"+sops+string(tag), []byte(""), 0644)
		}
	} else {
		fmt.Printf("Error tagging bug: %s", tag)
	}
}

// RemoveComment deletes a comment file of an issue.
func (b *Bug) RemoveComment(comment Comment) {
	if dir := b.Direr(); dir != "" {
		//os.Remove(filepath.FromSlash(string(dir) + "/comment-" + string(ShortTitleToDir(string(comment.Body)))))
		os.Remove(string(dir) + sops + "comment-" + string(ShortTitleToDir(string(comment.Body))))
	} else {
		fmt.Printf("Error removing comment: %s", comment.Body)
	}
}

// CommentBug writes a text file for an issue.
func (b *Bug) CommentBug(comment Comment, config Config) {
	if dir := b.Direr(); dir != "" {
		//os.Mkdir(filepath.FromSlash(string(dir)+"/"), 0755)
		commenttext := []byte(comment.Body + "\n")
		if config.ImportCommentsTogether { // not efficient but ok for now
			//data, err := ioutil.ReadFile(filepath.FromSlash(string(dir) + "/comments"))
			data, err := ioutil.ReadFile(string(dir) + sops + "comments")
			check(err)
			commentappend := []byte(fmt.Sprintf("%s%s%s", data, "\n", commenttext))
			//werr := ioutil.WriteFile(filepath.FromSlash(string(dir)+"/comments"), commentappend, 0644)
			werr := ioutil.WriteFile(string(dir)+sops+"comments", commentappend, 0644)
			check(werr)
		} else {
			//werr := ioutil.WriteFile(filepath.FromSlash(string(dir)+"/comment-")+string(ShortTitleToDir(string(comment.Body))), commenttext, 0644)
			werr := ioutil.WriteFile(string(dir)+sops+"comment-"+string(ShortTitleToDir(string(comment.Body))), commenttext, 0644)
			check(werr)
		}
	} else {
		fmt.Printf("Error commenting bug: %s", comment.Body)
	}
}

// ViewBug outputs an issue.
func (b Bug) ViewBug() {
	// This could be more generalized and shortened if the bug design is changed.
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
func (b Bug) StringTags() []string {
	tags := b.Tags()
	tagout := []string{}
	for _, tag := range tags {
		tagout = append(tagout, string(tag))
	}
	sort.Strings(tagout)
	return tagout
}

// HasTag returns if an issue is assigned a tag.
func (b Bug) HasTag(tag TagBoolTrue) bool {
	allTags := b.Tags()
	for _, bugTag := range allTags {
		if bugTag == tag {
			return true
		}
	}
	return false
}

// created b.tager for similar needs of {bugs/Bug.go:Tags, bugs/Bug.go:setField}, also bugs/Bug.go:liners

// tager takes a file name, returns the key, value,
// bool if value in name,
// bool if value in file contents, error
func (b Bug) tager(abspath string) (string, string, bool, bool, error) {
	dir := b.Direr()
	//hit := withtagsubdirfile.Name() // simple for tags subdir
	//   aka abspath
	key := ""
	value := ""
	tag_name := false
	tag_contents := false
	//var presentLines []string
	segments := strings.Split(abspath, string(os.PathSeparator)) // path separator
	parts := strings.Split(segments[len(segments)-1], "_")
	if len(parts) <= 1 {
		return key, value, tag_name, tag_contents, errors.New("tag has no key or value")
	} else if len(parts) == 2 {
		key = parts[1]
		//field, err := ioutil.ReadFile(filepath.FromSlash(string(dir) + "/tag_"))
		field, err := ioutil.ReadFile(string(dir) + sops + "tag_")
		if err == nil {
			value = ([]string(strings.Split(string(field), "\n")))[0] // tag_Status file contents overrides "Status" file contents
			// assumes value is ok, not false
			tag_contents = true
		}
	} else if len(parts) >= 3 {
		key = parts[1] // tag_Status_ file overrides "Status" file contents
		//presentLines = append(presentLines, strings.Join(parts[2:], "_")) // tag_Status_ file overrides "Status" file contents
		value = strings.Join(parts[2:], "_") // tag_Status_ file overrides "Status" file contents
		// assumes value is ok, not false
		tag_name = true
	}
	if value == "false" {
		return "", "", false, false, errors.New("tag has no key or value")
	} else {
		return strings.ToLower(key), strings.ToLower(value), tag_name, tag_contents, nil // key, value, tag_name, tag_contents, err
	}
}

// Tags returns a bug's array of tags.
func (b Bug) Tags() []TagBoolTrue {
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
	//withtagsubdir, errsubdir := ioutil.ReadDir(filepath.FromSlash(string(dir) + "/tags/")) // returns []os.FileInfo
	withtagsubdir, errsubdir := ioutil.ReadDir(string(dir) + sops + "tags" + sops) // returns []os.FileInfo
	withtagfile, errtagfile := filepath.Glob(string(dir) + sops + "tag_*")         // returns []string
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

//byBug allows sort.Sort(byBug(bugs)). Requirements are
//     func (bug) Len() int,
//     type
//     and the (byBug) {Len, Swap, Less} functions
// see also List.go for type byDir

func (t Bug) Len() int {
	return t.modtime // time.Format(time.UnixNano(t.modtime).UnixNano())
}

type byBug []Bug

func (t byBug) Len() int {
	return len(t) // time.Format(time.UnixNano(t.modtime).UnixNano())
}
func (t byBug) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t byBug) Less(i, j int) bool {
	return (t[i]).Len() < (t[j]).Len()
}

// findArrayString returns a bool if looking is an element of arr.
func findArrayString(arr []string, looking string) bool {
	for loop, _ := range arr {
		if arr[loop] == looking {
			return true
		}
	}
	return false
}

// fielder reads and returns the string value from the file of an issue.
func (b Bug) fielder(fieldName string) string {
	lines := b.liners(fieldName)
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	} else {
		return ""
	}
}

// created b.liners for similar needs of {bugs/Bug.go:fielder, bugs/Bug.go:setField}
// liners does the work for fielder with extra lines.
func (b Bug) liners(fieldName string) []string {
	dirr := b.Direr()
	dir := string(dirr)
	lines := []string{}
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
	// try tag_(K)ey_value
	_, value, _, _, err := b.tager(dir + sops + "tag_" + fieldName)
	if err == nil {
		lines = []string{value}
		return lines
	}
	// try tag_(k)ey_value
	_, value, _, _, err = b.tager(dir + sops + "tag_" + strings.ToLower(fieldName))
	if err == nil {
		lines = []string{value}
		return lines
	}
	// try tag_(K)ey file contents
	field, err = ioutil.ReadFile(dir + sops + "tag_" + fieldName)
	if err == nil {
		lines = strings.Split(string(field), "\n")
		return lines
	}
	// try tag_(k)ey file contents
	field, err = ioutil.ReadFile(dir + sops + "tag_" + strings.ToLower(fieldName))
	if err == nil {
		lines = strings.Split(string(field), "\n")
		return lines
	}
	return lines
}

//key, value, _, _, err := fielder(withtagfile[withtagfilefile], string(dir))

// setField writes the string value to the file of an issue.
func (b Bug) setField(fieldName string, value string, config Config) error { // TODO: complete func for config tag files : paused with tag_name, tag_contents, file_contents
	// using Status for fielName string example in comments
	dir := b.Direr()
	//possible locations
	tag_name := false
	tag_contents := false
	file_contents := false
	// try "Status" file
	presentLines := b.liners(fieldName) // var presentLines []string
	if len(presentLines) > 0 {
		file_contents = true
	}
	// try tag_Status* files
	withtagfile, errtagfile := filepath.Glob(string(dir) + sops + "tag_" + fieldName + "*") // returns []string
	errfind := errtagfile
	// two cases, ie tag_Status_closed or tag_Status contains closed
	if errtagfile == nil {
		for _, withtagfilefile := range withtagfile {
			presentvalue := ""
			_, presentvalue, tag_name, tag_contents, errfind = b.tager(withtagfilefile)
			if errfind == nil {
				presentLines = []string{presentvalue}
			}
			//segments := strings.Split(withtagfile[f], "/") // path separator
			//parts := strings.Split(segments[len(segments)-1], "_")
			//if len(parts) == 2 {
			//	field, errpresent = ioutil.ReadFile(string(dir) + "/tag_" + fieldName)
			//	if errpresent == nil {
			//		presentLines = strings.Split(string(field), "\n") // tag_ file contents overrides "Status" file contents
			//		tag_contents = true
			//	}
			//} else if len(parts) >= 3 {
			//	presentLines = append(presentLines, strings.Join(parts[2:], "_")) // tag_ file overrides "Status" file contents
			//	tag_name = true
			//}
		}
	}
	_ = tag_name
	_ = tag_contents
	_ = file_contents

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
			err = ioutil.WriteFile(string(dir)+sops+"tag_"+fieldName+"_"+strings.ToLower(TitleToDirString(newValue)), []byte(""), 0644)
		} else {
			err = ioutil.WriteFile(string(dir)+sops+"tag_"+fieldName+"_"+TitleToDirString(newValue), []byte(""), 0644)
		}
	} else {
		err = ioutil.WriteFile(string(dir)+sops+fieldName, []byte(newValue), 0644)
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Status returns the string from the Status file of an issue.
func (b Bug) Status() string {
	return b.fielder("Status")
}

// SetStatus writes the Status file to an issue.
func (b Bug) SetStatus(newStatus string, config Config) error {
	return b.setField("Status", newStatus, config)
}

// Priority returns the string from the Priority file of an issue.
func (b Bug) Priority() string {
	return b.fielder("Priority")
}

// SetPriority writes the Priority file to an issue.
func (b Bug) SetPriority(newValue string, config Config) error {
	return b.setField("Priority", newValue, config)
}

// Milestone returns the string from the Milestone file of an issue.
func (b Bug) Milestone() string {
	return b.fielder("Milestone")
}

// SetMilestone writes the Milestone file to an issue.
func (b Bug) SetMilestone(newValue string, config Config) error {
	return b.setField("Milestone", newValue, config)
}

// Identifier returns the string from the Identifier file of an issue.
func (b Bug) Identifier() string {
	i := b.fielder("Identifier")
	if i != "" {
		return i
	} else {
		return b.fielder("Id")
	}
}

// SetIdentifier writes the Identifier file to an issue.
func (b Bug) SetIdentifier(newValue string, config Config) error {
	return b.setField("Identifier", newValue, config)
}

// New assigns and writes an issue.
func New(title string, config Config) (*Bug, error) {
	expectedDir := IssuesDirer(config) + Directory(os.PathSeparator) + TitleToDir(title)
	err := os.Mkdir(string(expectedDir), 0755)
	if err != nil {
		return nil, err
	}
	return &Bug{Dir: expectedDir}, nil
}
