package main 

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os"
	"time"
	"strings"
	"log"
)

import (
	geddit "github.com/initzero/geddit"
	irc "github.com/thoj/Go-IRC-Client-Library"
)

func slowSendIRC(s string, icon irc.IRCConnection) {
	second := time.Second
	//	index := len(s)
	time.Sleep(second)
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s server:port channel bot-username", os.Args[0])
		os.Exit(1)
	}
	
	// http client
	client := &http.Client{}

	// setup IRC
	icon := irc.IRC(os.Args[3], "hahawtflol")
	err := icon.Connect(os.Args[1])
	geddit.CheckError(err)
	icon.AddCallback("001", func(e *irc.IRCEvent) { icon.Join("#" + os.Args[2]) })
	icon.AddCallback("JOIN", func(e *irc.IRCEvent) { 
		icon.SendRaw("PRIVMSG #" + os.Args[2] + " :SCV ready s(^_^)-b")
		//for _, str := range jRep.List() {		
			//icon.SendRaw("PRIVMSG #geddit :" + str)
		//	time.Sleep(time.Second*1) 
		//}	
	})
	icon.AddCallback("PRIVMSG", func(e *irc.IRCEvent) {
		if strings.HasPrefix(e.Message, "@reddit") {
			req, err := http.NewRequest("GET", "http://reddit.com/r/" + e.Message[8:] + ".json", nil)
			req.Header.Set("User-Agent", "wtf_is_up ~ playing with Go-lang")
	
			// make request	
			resp, err := client.Do(req)	
			geddit.CheckError(err)
			if resp.StatusCode != http.StatusOK {
				log.Println("bad HTTP response: " + resp.Status)
				return
			}
	
			// read response
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			//fmt.Printf("%s\n", string(body))
			geddit.CheckError(err)
	
			if strings.HasPrefix(string(body), "<") {
				log.Println("returned xml; bad")
				icon.SendRaw("PRIVMSG " + e.Arguments[0] + " :ERROR // bad subreddit // q-(v_v)z")
				return
			}
			// parse JSON
			var jRep geddit.Top
			err = json.Unmarshal(body, &jRep)
			geddit.CheckError(err)

			icon.SendRaw("PRIVMSG " + e.Arguments[0] + " :d-(^_^)z check your PMs for /r/" + e.Message[8:])
			for _, str := range jRep.List() {
				icon.SendRaw("PRIVMSG " + e.Nick + " :" + str)
				time.Sleep(time.Second*1)
			}		
		}
		
	})
	icon.Loop()
}
