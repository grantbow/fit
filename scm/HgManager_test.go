	"flag"
	"log"
		//c.Fail()
		log.Fatal(err)
var hg bool

func init() {
	flag.BoolVar(&hg, "hg", true, "Mercurial (hg) presence")
	flag.Parse()
	_, err := runCmd("hg")
	if err != nil {
		hg = false
	}
}

	if hg == false {
		t.Skip("hg executable not found")
	}