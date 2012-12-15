package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

import (
	geddit "github.com/initzero/geddit"
	irc "github.com/thoj/Go-IRC-Client-Library"
)

// send strings to IRC
func sendIRC(s []string, i *irc.IRCConnection, e *irc.IRCEvent, ready chan bool) {
	i.SendRaw("PRIVMSG " + e.Arguments[0] + " :\x031,9d-(^_^)z \x039,1 check your PMs for /r/" + e.Message[8:])
	for _, str := range s {
		i.SendRaw("PRIVMSG " + e.Nick + " :" + str)
		time.Sleep(time.Second * 1)
	}
	ready <- true
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s server:port channel bot-username", os.Args[0])
		os.Exit(1)
	}

	ready := make(chan bool, 1)
	ready <- true
	// http client
	client := &http.Client{}

	// setup IRC
	icon := irc.IRC(os.Args[3], "hahawtflol")
	err := icon.Connect(os.Args[1])
	geddit.CheckError(err)

	icon.AddCallback("001", func(e *irc.IRCEvent) { icon.Join("#" + os.Args[2]) })
	//icon.AddCallback("JOIN", func(e *irc.IRCEvent) { 
	//icon.SendRaw("PRIVMSG #" + os.Args[2] + " :SCV ready s(^_^)-b")
	//})

	icon.AddCallback("PRIVMSG", func(e *irc.IRCEvent) {
		if strings.HasPrefix(e.Message, "@reddit") {
			// check for message longer than '@reddit sub'
			if len(e.Message) < 8 {
				return
			}
			req, err := http.NewRequest("GET", "http://reddit.com/r/"+e.Message[8:]+".json", nil)
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
			geddit.CheckError(err)

			// hackish way of dealing with XML returned from reddit redirect caused 
			// by subreddit not existing	
			if strings.HasPrefix(string(body), "<") {
				log.Println("returned xml; bad")
				icon.SendRaw("PRIVMSG " + e.Arguments[0] + " :\x031,4q-(v_v)z \x034,1 bad subreddit")
				return
			}
			// parse JSON
			var jRep geddit.Top
			err = json.Unmarshal(body, &jRep)
			geddit.CheckError(err)
			select {
			case <-ready:
				go sendIRC(jRep.ToIRCStrings(), icon, e, ready)
			default:
				icon.SendRaw("PRIVMSG " + e.Arguments[0] + " :\x031,4q-(v_v)z \x034,1 wait ur turn (flood control, l0l) ")
			}
		}
		if strings.Contains(e.Message, "ur mom") {
			icon.SendRaw("PRIVMSG " + e.Arguments[0] + " :\x031,4http://en.wikipedia.org/wiki/List_of_burn_centers_in_the_United_States")
		}
	})
	icon.Loop()
}
