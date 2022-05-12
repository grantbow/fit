package fitapp

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// ProgramVersion returns the base version string
func ProgramVersion() string {
	return "0.7"
}

// PrintVersion describes fit and golang.
func PrintVersion() {
	ex, _ := os.Executable()
	fi, _ := os.Stat(ex)
	fmt.Printf("%s version %s built using %s GOOS %v\n",
		os.Args[0],
		ProgramVersion(),
		runtime.Version(),
		runtime.GOOS)

	var loc *time.Location
	var errl error
	if runtime.GOOS == "android" {
		// workaround for https://golang.org/src/time/zoneinfo_android.go line 25 hard coded UTC
		tz := os.Getenv("TZ")
		if os.Getenv("TZ") == "" {
			fmt.Println("WARN: TZ is empty or not set")
		}
		loc, errl = time.LoadLocation(tz)
		if errl != nil {
			panic(errl)
		}
	} else {
		loc = time.Local
	}

	fmt.Printf("executable: %v %v %s %s\n",
		fi.Mode(),
		fi.Size(),
		fi.ModTime().In(loc).Format(time.UnixDate),
		ex)
	//name, offset := time.Now().Local().Zone()
	//fmt.Printf("z: %+v %+v\n", name, offset)
	//fmt.Printf("tz: %+v\n", os.Getenv("TZ"))
	//fmt.Printf("aha: %+v\n", fi.ModTime().In(loc).Format(time.UnixDate))

	//x := fi.ModTime()
	//x.loc = Local
	//fmt.Printf("tz: %+v\n", os.Setenv("TZ", "America/Los_Angeles"))
	//time.LoadLocation(time.Local)
	//time.LoadLocation("America/Los_Angeles")

}
