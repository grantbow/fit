package bugs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

// ErrNoDescription defines a new error.
var ErrNoDescription = errors.New("No description provided")

// ErrNotFound defines a new error.
var ErrNotFound = errors.New("Could not find bug")

// Bug is the type of an issue.
// The fields are Dir and descFile.
type Bug struct {
	Dir      Directory
	descFile *os.File
}

// Tag is the type of an issue identifier.
// There is only a string key.
// Values are not supported so there is an implied true/present false/absent value.
type Tag string

// Comment is the struct type of a unit of discussion about an issue.
// The fields are Author, Time, Body, Order and Xml.
type Comment struct {
	Author string
	Time   time.Time
	Body   string
	Order  int
	Xml    []byte
}

// TitleToDir returns a Directory from a string argument.
func TitleToDir(title string) Directory {
	// replace non-matching valid characters with _
	// for user entered strings
	re := regexp.MustCompile("[^a-zA-Z0-9_ -]+")
	s := re.ReplaceAllString(title, "_")

	replaceWhitespaceWithUnderscore := func(match string) string {
		return strings.Replace(match, " ", "_", -1)
	}
	replaceDashWithMore := func(match string) string {
		if strings.Count(match, " ") > 0 {
			return match
		}
		return "-" + match
	}

	// Replace sequences of dashes with 1 more dash,
	// as long as there's no whitespace around them
	re = regexp.MustCompile("([\\s]*)([-]+)([\\s]*)")
	s = re.ReplaceAllStringFunc(s, replaceDashWithMore)
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
	return Directory(s)
}

// ShortTitleToDir truncates a title to 25 characters.
func ShortTitleToDir(title string) Directory {
	return TitleToDir(title[:25]) // TODO: remove leading or trailing _ or -
}

// GetDirectory returns the directory of an issue.
func (b Bug) GetDirectory() Directory {
	return b.Dir
}

// LoadBug assigns a directory to an issue.
func (b *Bug) LoadBug(dir Directory) {
	b.Dir = dir
}

// Title returns a string with the name of an issue and optionally present Status and Priority.
func (b Bug) Title(options string) string {
	var hasOption = func(o string) bool {
		return strings.Contains(options, o)
	}

	title := b.Dir.GetShortName().ToTitle()

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
	value, err := ioutil.ReadAll(&b)

	if err != nil {
		if err == ErrNoDescription {
			return "No description provided.\n"
		}
		panic("Unhandled error" + err.Error())
	}

	if string(value) == "" {
		return "No description provided.\n"
	}
	return string(value)
}

// SetDescription writes the Description file of an issue.
func (b Bug) SetDescription(val string) error {
	dir := b.GetDirectory()

	return ioutil.WriteFile(string(dir)+"/Description", []byte(val+"\n"), 0644)
}

// RemoveTag deletes a tag file in the tags subdirectory of an issue.
func (b *Bug) RemoveTag(tag Tag) {
	if dir := b.GetDirectory(); dir != "" {
		os.Remove(string(dir) + "/tags/" + string(tag))
	} else {
		fmt.Printf("Error removing tag: %s", tag)
	}
}

// TagBug writes an empty tag file in the tags subdirectory of an issue.
func (b *Bug) TagBug(tag Tag) {
	if dir := b.GetDirectory(); dir != "" {
		os.Mkdir(string(dir)+"/tags/", 0755)
		ioutil.WriteFile(string(dir)+"/tags/"+string(tag), []byte(""), 0644)
	} else {
		fmt.Printf("Error tagging bug: %s", tag)
	}
}

// RemoveComment deletes a comment file of an issue.
func (b *Bug) RemoveComment(comment Comment) {
	if dir := b.GetDirectory(); dir != "" {
		os.Remove(string(dir) + "/comment-" + string(ShortTitleToDir(string(comment.Body))))
	} else {
		fmt.Printf("Error removing comment: %s", comment.Body)
	}
}

// CommentBug writes a text file for an issue.
func (b *Bug) CommentBug(comment Comment, config Config) {
	if dir := b.GetDirectory(); dir != "" {
		//os.Mkdir(string(dir)+"/", 0755)
		commenttext := []byte(comment.Body + "\n")
		if config.ImportCommentsTogether { // not efficient but ok for now
			data, err := ioutil.ReadFile(string(dir) + "/comments")
			check(err)
			commentappend := []byte(fmt.Sprintf("%s%s%s", data, "\n", commenttext))
			werr := ioutil.WriteFile(string(dir)+"/comments", commentappend, 0644)
			check(werr)
		} else {
			werr := ioutil.WriteFile(string(dir)+"/comment-"+string(ShortTitleToDir(string(comment.Body))), commenttext, 0644)
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

	fmt.Printf("Title: %s\n\n", b.Title(""))
	fmt.Printf("Description:\n%s", b.Description())

	if status := b.Status(); status != "" {
		fmt.Printf("\nStatus: %s", status)
	}
	if priority := b.Priority(); priority != "" {
		fmt.Printf("\nPriority: %s", priority)
	}
	if milestone := b.Milestone(); milestone != "" {
		fmt.Printf("\nMilestone: %s", milestone)
	}
	if tags := b.StringTags(); tags != nil {
		fmt.Printf("\nTags: %s", strings.Join([]string(tags), ", "))
	}

}

// StringTags outputs tags of an issue.
func (b Bug) StringTags() []string {
	dir := b.GetDirectory()
	dir += "/tags/"
	issues, err := ioutil.ReadDir(string(dir))
	if err != nil {
		return nil
	}

	tags := make([]string, 0, len(issues))
	for _, issue := range issues {
		tags = append(tags, issue.Name())
	}
	return tags
}

// HasTag returns if an issue is assigned a tag.
func (b Bug) HasTag(tag Tag) bool {
	allTags := b.Tags()
	for _, bugTag := range allTags {
		if bugTag == tag {
			return true
		}
	}
	return false
}

// Tags returns an array of assigned tags.
func (b Bug) Tags() []Tag {
	dir := b.GetDirectory()
	dir += "/tags/"
	issues, err := ioutil.ReadDir(string(dir))
	if err != nil {
		return nil
	}

	tags := make([]Tag, 0, len(issues))
	for _, issue := range issues {
		tags = append(tags, Tag(issue.Name()))
	}
	return tags

}

// getField reads and returns the string value from the file of an issue.
func (b Bug) getField(fieldName string) string {
	dir := b.GetDirectory()
	field, err := ioutil.ReadFile(string(dir) + "/" + fieldName)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(field), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

// setField writes the string value to the file of an issue.
func (b Bug) setField(fieldName, value string) error {
	dir := b.GetDirectory()
	oldValue, err := ioutil.ReadFile(string(dir) + "/" + fieldName)
	var oldLines []string
	if err == nil {
		oldLines = strings.Split(string(oldValue), "\n")
	}

	newValue := ""
	if len(oldLines) >= 1 {
		// If there were 0 or 1 old lines, overwrite them
		oldLines[0] = value
		newValue = strings.Join(oldLines, "\n")
	} else {
		newValue = value
	}

	err = ioutil.WriteFile(string(dir)+"/"+fieldName, []byte(newValue), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Status returns the string from the Status file of an issue.
func (b Bug) Status() string {
	return b.getField("Status")
}

// SetStatus writes the Status file to an issue.
func (b Bug) SetStatus(newStatus string) error {
	return b.setField("Status", newStatus)
}

// Priority returns the string from the Priority file of an issue.
func (b Bug) Priority() string {
	return b.getField("Priority")
}

// SetPriority writes the Priority file to an issue.
func (b Bug) SetPriority(newValue string) error {
	return b.setField("Priority", newValue)
}

// Milestone returns the string from the Milestone file of an issue.
func (b Bug) Milestone() string {
	return b.getField("Milestone")
}

// SetMilestone writes the Milestone file to an issue.
func (b Bug) SetMilestone(newValue string) error {
	return b.setField("Milestone", newValue)
}

// Identifier returns the string from the Identifier file of an issue.
func (b Bug) Identifier() string {
	return b.getField("Identifier")
}

// SetIdentifier writes the Identifier file to an issue.
func (b Bug) SetIdentifier(newValue string) error {
	return b.setField("Identifier", newValue)
}

// New assigns and writes an issue.
func New(title string, config Config) (*Bug, error) {
	expectedDir := GetIssuesDir(config) + TitleToDir(title)
	err := os.Mkdir(string(expectedDir), 0755)
	if err != nil {
		return nil, err
	}
	return &Bug{Dir: expectedDir}, nil
}
