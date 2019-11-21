package issues

import (
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

// also in bugapp/utils.go
func removeFi(slice []os.FileInfo, i int) []os.FileInfo {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

// also in bugapp/utils.go
func readIssues(dirname string) []os.FileInfo {
	//var issueList []os.FileInfo
	fis, _ := ioutil.ReadDir(string(dirname))
	issueList := fis
	for idx, fi := range issueList {
		if fi.IsDir() != true {
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
