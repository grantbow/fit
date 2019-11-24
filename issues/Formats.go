package issues

import (
	"encoding/json"
)

// ToJSONString encodes an issue. A string and an error are returned.
func (i Issue) ToJSONString() (string, error) {
	iJSONStruct := struct {
		Identifier  string `json:",omitempty"`
		Title       string
		Description string
		Status      string   `json:",omitempty"`
		Priority    string   `json:",omitempty"`
		Milestone   string   `json:",omitempty"`
		Tags        []string `json:",omitempty"`
	}{
		Identifier:  i.Identifier(),
		Title:       i.Title(""),
		Description: i.Description(),
		Status:      i.Status(),
		Priority:    i.Priority(),
		Milestone:   i.Milestone(),
		Tags:        i.StringTags(),
	}

	iJSON, err := json.Marshal(iJSONStruct)
	if err != nil {
		return "", err
	}
	return string(iJSON), nil
}
