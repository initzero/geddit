package geddit 

import (
	"log"
	"encoding/json"
	"strings"
	"math"
)

// nullable string -- to handle json returning null. 
// if encoding/json expects a string, the api returning a null will cause problems 
// (same for bool/int)
type NullString string

func (n *NullString) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(n))
}

// nullable bool
type NullBool bool

func (n *NullBool) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*bool)(n))
}

// nullable float x_x
type NullFloat float64

func (n *NullFloat) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*float64)(n))
}

type EData struct {
	Domain 				NullString 		`json:"domain"`
	BannedBy 			NullString		`json:"banned_by"`
	//MediaEmbed 			[]string	`json:"media_embed"`
	Subreddit 			NullString		`json:"subreddit"`
	SelftextHTML 		NullString		`json:"selftext_html"`
	Selftext 			NullString		`json:"selftext"`
	//Likes 				float64		`json:"likes"`
	LinkFlairText 		NullString		`json:"link_flair_text"`
	Id 					NullString		`json:"id"`
	//Clicked 			bool			`json:"clicked"`
	Title 				NullString		`json:"title"`
	NumComments 		NullFloat		`json:"num_comments"`
	Score 				NullFloat		`json:"score"`
	ApprovedBy 			NullString		`json:"approved_by"`
	//Over18 				bool		`json:"over_18"`
	//Hidden 				bool		`json:"hidden"`
	Thumbnail 			NullString		`json:"thumbnail"`
	SubredditID			NullString		`json:"subreddit_id"`
	//Edited 				bool		`json:"edited"`
	LinkFlairCSSClass 	NullString		`json:"link_flair_css_class"`
	AuthorFlairCSSClass NullString		`json:"author_flair_css_class"`
	Downs 				NullFloat		`json:"downs"`
	//Saved 				bool		`json:"saved"`
	//IsSelf 				bool		`json:"is_self"`
	Permalink 			NullString		`json:"permalink"`
	Name 				NullString		`json:"name"`
	Created 			NullFloat		`json:"created"`
	Url 				NullString		`json:"url"`
	AuthorFlairText 	NullString		`json:"author_flair_text"`
	Author 				NullString		`json:"author"`
	CreatedUTC 			NullFloat		`json:"created_utc"`
	//Media 				NullString		`json:"media"`
	NumReports 			NullFloat		`json:"num_reports"`
	Ups 				NullFloat		`json:"ups"`
}

type Entry struct {
	Kind 				NullString		`json:"kind"`	
	Data 				EData			`json:"data"`
}

type TData struct {
	Modhash 			NullString		`json:"modhash"`
	Children 			[]Entry			`json:"children"`
	After 				NullString		`json:"after"`
	Before 				NullString 		`json:"before"`
}

type Top struct {
	Kind 				NullString		`json:"kind"`
	Data 				TData			`json:"data"`
	
}

func (t Top) String() string {
	var ret []string
	for _,child := range t.Data.Children {
		ret = append(ret, string(child.Data.Title))
	}
	return strings.Join(ret, "\n")
}

// return an array of strings that can be easily iterated through 
// and sent to irc.SendRaw
func (t Top) List() []string {
	var ret []string
	for _, child := range t.Data.Children {
		ret = append(ret, string(child.Data.Title) + " " + string(child.Data.Url))
	}
	return ret
}

// format Top representation as IRC string(s)
// TODO: make the color codes configurable
func (t Top) ToIRCStrings() []string {
	var ret []string
	colorCodeA := "\x034,1"		// red on black
	colorCodeB := "\x0312,1"	// blue on black
	colorCodeC := "\x031,0"		// black on white
	tmpColorCode := colorCodeA
	var i float64 = 0.0

	ret = append(ret, "\x033d-(^_^)z \x037 Loading reddit feed...")
	for _, child := range t.Data.Children {
		if math.Mod(i, 2) != 0 {
			tmpColorCode = colorCodeA
		} else {
			tmpColorCode = colorCodeB
		}
		i += 1
		// build return array
		ret = append(ret, 
			tmpColorCode + string(child.Data.Title) + "\x031,0 >> " + colorCodeC + string(child.Data.Url))
	}
	return ret
}

func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
