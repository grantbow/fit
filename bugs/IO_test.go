package bugs

import (
	_ "fmt"
	"io/ioutil"
	"os"
	"testing"
)

//var dops = Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

func TestBugWrite(t *testing.T) {
	config := Config{}
	config.DescriptionFileName = "Description"
	var b *Bug
	if dir, err := ioutil.TempDir("", "BugWrite"); err == nil {
		os.Chdir(dir)
		b = &Bug{Dir: Directory(dir + sops + "issues" + sops + "Test-bug"), DescriptionFileName: config.DescriptionFileName}
		defer os.RemoveAll(dir)
	} else {
		t.Error("Could not get temporary directory to test bug write()")
		return
	}

	_, err := b.Write([]byte("Hello there, Mr. Test"))
	if err != nil {
		t.Errorf("Error writing to bug at %s.", b.Dir)
	}
	b.Close()

	fp, _ := os.Open("issues" + sops + "Test-bug" + sops + "Description")
	desc, err := ioutil.ReadAll(fp)
	fp.Close()

	if err != nil {
		t.Error("Error reading description file.")
		return
	}

	// Cast the values to strings because []byte complains that
	// slices can only be compared to nil.
	if string(desc) != string("Hello there, Mr. Test") {
		t.Error("Incorrect description file after writing to bug")
	}
}

/*
func ExampleBugWriter() {
	config := Config{}
	config.DescriptionFileName = "Description"
	if b, err := New("Bug Title", config); err != nil {
		fmt.Fprintf(b, "This is a bug report.\n")
		fmt.Fprintf(b, "The bug will be created as necessary.\n")
	}
}
*/
