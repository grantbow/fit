package issues

import (
	_ "fmt" // for debugging
	"io/ioutil"
	"os"
)

var dops = Directory(os.PathSeparator)
var sops = string(os.PathSeparator)

func check(e error) {
	if e != nil {
		//	fmt.Fprintf(os.Stderr, "err: %s", err.Error())
		//	return NoConfigError
		panic(e)
	}
}

// also in fitapp/utils.go
func removeFi(slice []os.FileInfo, i int) []os.FileInfo {
	//flag.Parse()
	//Debug("debug len " + string(len(slice)) + " i " + string(i) + " slice[0].Name() " + string(slice[0].Name()) + "\n")
	//fmt.Printf("%+v\n", flag.Args()) // didn't seem to help, needs more work to make it active
	//
	//fmt.Printf("debug ok 01 \n")
	//fmt.Printf("debug len " + string(len(slice)) + " i " + string(i) + " slice[0].Name() " + string(slice[0].Name()) + "\n")
	//fmt.Printf("debug removeFi args len " + string(len(slice)) + " i " + string(i) + "\n")
	if (len(slice) == 1) && (i == 0) {
		return []os.FileInfo{}
	} else if i < len(slice)-2 {
		copy(slice[i:], slice[i+1:])
	}
	return slice[:len(slice)-1]
}

// also in fitapp/utils.go
func readIssues(dirname string) []os.FileInfo {
	//var issueList []os.FileInfo
	fis, _ := ioutil.ReadDir(string(dirname))
	issueList := fis
	for idx, fi := range issueList {
		//Debug("debug fi " + string(fi.Name()) + "idx " + string(idx) + "\n")
		//fmt.Printf("debug readIssues loop fi " + string(fi.Name()) + "idx " + string(idx) + "\n")
		if fi.IsDir() != true {
			//fmt.Printf("debug before removeFi name " + fi.Name() + " idx " + string(idx) + "\n")
			issueList = removeFi(issueList, idx)
		}
	}
	return issueList
}

//byDir allows sort.Sort(byDir(issues))
// type and three functions are needed - also see Bug.go for type byBug
// rather than a custom Len function for os.FileInfo, Len is calculated in Less
type byDir []os.FileInfo

func (t byDir) Len() int {
	return len(t) // time.Format(time.UnixNano(t.modtime).UnixNano())
}
func (t byDir) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t byDir) Less(i, j int) bool {
	return (t[i]).ModTime().Unix() < (t[j]).ModTime().Unix()
}
