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

// Twilio is a subcommand to send to changed isssues with tag_twilio_4155551212
func Twilio(config bugs.Config) {
	var hasPart = func(target string, part string) bool {
		return strings.Contains(target, part)
	}

	buglist := bugs.GetAllBugs(config)
	//fmt.Printf("getallbugs: %q\n", buglist)
	if len(buglist) > 0 {
		// from buglist and
		// from scm get added, changed and removed
		//fmt.Printf("debug 1\n")
		scmoptions := map[string]bool{}
		handler, _, ErrH := scm.DetectSCM(scmoptions, config)
		if ErrH == nil {
			// scm exists
			if b, err := handler.SCMIssuesUpdaters(); err != nil {
				// uncommitted files including staged AND working directory
				//fmt.Printf("debug 2\n")
				if _, ErrCach := handler.SCMIssuesCacher(); ErrCach != nil {
					// uncommitted files staged only NOT working directory
					//fmt.Printf("debug 3\n")
					updatedissues := map[string]bool{}      // issues staged no duplicates
					twiliorecipients := map[string]string{} // one or more per issue staged
					//fmt.Printf("updated issues:\n")
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
					bug := bugs.Bug{}

					// build message for each recipient from updated issues and twilio tags
					for key, _ := range updatedissues {
						//fmt.Printf("twilio bug dirname: %v\n", key)
						expectedbugdir := string(bugs.IssuesDirer(config)) + "/" + key
						bug.LoadBug(bugs.Directory(expectedbugdir))
						tags := bug.Tags()
						//fmt.Printf("debug %v tags %v\n", key, tags)
						for _, k := range tags {
							//fmt.Printf("k: %v\n", k)
							if hasPart(string(k), "twilio") { // local function returns bool
								a := strings.Split(string(k), ":") // : separated from bug.Tags
								recip := a[1]
								//fmt.Printf("twilio bug dirname: %v tag %v : %v\n", key, a[0], recip)
								//	if strings.ToLower(string(tag)) == "twilio" {
								if _, ok := twiliorecipients[recip]; ok {
									// recipient exists, append
									twiliorecipients[recip] = twiliorecipients[recip] + ", " + key
								} else {
									// new recipient
									twiliorecipients[recip] = "site " + config.IssuesSite + "\nupdated " + key
								}
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

func TwilioDoSend(config bugs.Config, PNTo string, BodyStr string) {
	urlStr := TwilioUrlHttp + config.TwilioAccountSid + TwilioUrlMessages
	fmt.Println(urlStr)
	msgData := url.Values{}
	msgData.Set("To", PNTo) // Phone Number To
	msgData.Set("From", config.TwilioPhoneNumberFrom)
	msgData.Set("Body", BodyStr) // text message body
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
