package bugapp

import (
	"encoding/json"
	"fmt"
	"github.com/driusan/bug/bugs"
	"github.com/driusan/bug/scm"
	"net/http"
	"net/url"
	"strings"
)

var TwilioUrlHttp = "https://api.twilio.com/2010-04-01/Accounts/"
var TwilioUrlMessages = "/Messages.json"

/*
// getAllTags returns all the tags
func getAllTags(config bugs.Config) []string {
	bugs := bugs.GetAllBugs(config)
	//fmt.Printf("%+v\n", bugs)
	tagMap := make(map[string]int, 0)
	// Put all the tags in a map, then iterate over
	// the tags so that only unique tags are included
	for _, bug := range bugs {
		for _, tag := range bug.Tags() {
			tagMap[strings.ToLower(string(tag))] += 1
		}
	}
	var tags []string
	for k := range tagMap {
		tags = append(tags, k)
	}
	sort.Strings(tags)
	return tags
}
*/

// TwilioSend is a subcommand to send to the assigned tags.
func Twilio(config bugs.Config) {
	var hasPart = func(target string, part string) bool {
		return strings.Contains(target, part)
	}

	buglist := bugs.GetAllBugs(config)
	//tagMap := make(map[string]int, 0)
	if len(buglist) > 0 {
		// from buglist and
		// from scm get added, changed and removed
		//fmt.Printf("debug 1\n")
		scmoptions := map[string]bool{}
		handler, _, ErrH := scm.DetectSCM(scmoptions, config)
		if ErrH == nil {
			//fmt.Printf("debug 1.5\n")
			if b, err := handler.SCMIssuesUpdaters(); err != nil {
				//fmt.Printf("debug 2\n")
				if _, ErrCach := handler.SCMIssuesCacher(); ErrCach != nil {
					//fmt.Printf("debug 3\n")
					updatedissues := map[string]bool{}
					twiliorecipients := map[string]string{}
					fmt.Printf("updated issues with twilio:\n")
					for _, bline := range strings.Split(string(b), "\n") {
						if len(bline) > 0 {
							i := strings.Split(string(bline), "/")
							if len(i) > 2 {
								updatedissues[i[1]] = true
							}
						}
					}
					// updated issues exist
					//for key, _ := range updatedissues {
					//	fmt.Printf("bug dirname: %v\n", key)
					//}
					//fmt.Printf("allbugs: %q\n", buglist)
					bug := bugs.Bug{}

					// twilio tag must exist
					for key, _ := range updatedissues {
						//fmt.Printf("twilio bug dirname: %v\n", key)
						//bug := buglist[string(bugs.IssuesDirer(config))+"/"+key]
						expectedbugdir := string(bugs.IssuesDirer(config)) + "/" + key
						bug.LoadBug(bugs.Directory(expectedbugdir))
						tags := bug.Tags()
						//fmt.Printf("tags: %v\n", tags)
						for _, k := range tags {
							//fmt.Printf("k: %v\n", k)
							if hasPart(string(k), "twilio") {
								a := strings.Split(string(k), ":")
								fmt.Printf("twilio bug dirname: %v tag %v : %v\n", key, a[0], a[1])
								recip := a[1]
								//	if strings.ToLower(string(tag)) == "twilio" {
								if _, ok := twiliorecipients[recip]; ok {
									twiliorecipients[recip] = twiliorecipients[recip] + ", " + key
								} else {
									// new recipient
									twiliorecipients[recip] = "site " + config.IssuesSite + "\nupdated " + key
								}
								//fmt.Printf("twilio %v tags %v\n", key, )
								//fmt.Printf("twilio %v tags %v\n", key, bug.Tags())
								//		// get value for sending
							}
						}
					}
					fmt.Printf("result: %v\n", twiliorecipients)
					for msg := range twiliorecipients {
						TwilioDoSend(config, msg, twiliorecipients[msg])
					}
				} else {
					fmt.Printf("No updated and staged issues.\n")
				}
			} else {
				fmt.Printf("No updated or staged issues.\n")
			}
		}
	} else {
		fmt.Print("<no twilio tags yet>\n")
		return
	}
}

/*
	issuesroot := bugs.IssuesDirer(config)
	issues, _ := ioutil.ReadDir(string(issuesroot)) // TODO: should be a method elsewhere
	sort.Sort(byDir(issues))
	var wantTags bool = false

	allbugs := bugs.GetAllBugs(config)
	tagMap := make(map[string]int, 0)
	for _, bug := range allbugs {
		if len(bug.Tags()) == 0 {
			title := bug.Dir.ShortNamer()
			tagMap[string(title)] += 1
		}
	}
*/

// data comes from tags
//var PNTo = "+14153752752" // Phone Number To
//var BodyStr = "more new issues"

func TwilioDoSend(config bugs.Config, PNTo string, BodyStr string) {
	urlStr := TwilioUrlHttp + config.TwilioAccountSid + TwilioUrlMessages
	fmt.Println(urlStr)
	msgData := url.Values{}
	msgData.Set("To", PNTo)
	msgData.Set("From", config.TwilioPhoneNumberFrom)
	msgData.Set("Body", BodyStr)
	msgDataReader := *strings.NewReader(msgData.Encode())
	//
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(config.TwilioAccountSid, config.TwilioAuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil { // less than perfect, == nil perfect
			fmt.Println(data["sid"])
		} else {
			fmt.Printf("%v\n", resp) // print success
		}
	} else {
		fmt.Printf("%v\n", resp) // print failure
	}
}
