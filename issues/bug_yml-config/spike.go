package main

import (
	"fmt"
	"github.com/ghodss/yaml"
)

type Person struct {
	// location of root dir and .bug.yml
	BugDir string `json:"BugDir"`
	// overrides auto-find root dir
	// overridden by PMIT environment variable
	PMIT string `json:"PMIT"`
	// new bug Description text template
	DefaultDescriptionFile string `json:"DefaultDescriptionFile"`
	// saves raw json files of import
	ImportXmlDump bool `json:"ImportXmlDump"`
}

func main() {
	// Marshal a Person struct to YAML.
	p := Person{"goodbugdir", "pain", "path", true}
	y, err := yaml.Marshal(p)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(y))
	/* Output:
	age: 30
	name: John
	*/

	// Unmarshal the YAML back into a Person struct.
	var p2 Person
	err = yaml.Unmarshal(y, &p2)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(p2)
	/* Output:
	{John 30}
	*/
}
