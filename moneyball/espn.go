package main

/**
Copyright (c) 2013 The go-github AUTHORS. All rights reserved.
Copyright (c) 2020 DXC Technology - Dan Hushon. All rights reserved

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc., DXC Technology nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// scaffolding: https://medium.com/@marcus.olsson/writing-a-go-client-for-your-restful-api-c193a2f4998c
// which led to https://github.com/google/go-github which is really quite like our fetch/version process

// ESPN - NBA Scores from secret API
//
// NBA
//
//Scores: https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard
//News: http://site.api.espn.com/apis/site/v2/sports/basketball/nba/news
//All Teams: http://site.api.espn.com/apis/site/v2/sports/basketball/nba/teams
//Specific Team: http://site.api.espn.com/apis/site/v2/sports/basketball/nba/teams/:team
//http://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard/:eventId
//
//WNBA
//
//Scores: http://site.api.espn.com/apis/site/v2/sports/basketball/wnba/scoreboard
//News: http://site.api.espn.com/apis/site/v2/sports/basketball/wnba/news
//All Teams: http://site.api.espn.com/apis/site/v2/sports/basketball/wnba/teams
//Specific Team: http://site.api.espn.com/apis/site/v2/sports/basketball/wnba/teams/:team
//
//Women's College Basketball
//Scores: http://site.api.espn.com/apis/site/v2/sports/basketball/womens-college-basketball/scoreboard
//News: http://site.api.espn.com/apis/site/v2/sports/basketball/womens-college-basketball/news
//All Teams: http://site.api.espn.com/apis/site/v2/sports/basketball/womens-college-basketball/teams
//Specific Team: http://site.api.espn.com/apis/site/v2/sports/basketball/womens-college-basketball/teams/:team
//
//Men's College Basketball
//Scores: http://site.api.espn.com/apis/site/v2/sports/basketball/mens-college-basketball/scoreboard
//News: http://site.api.espn.com/apis/site/v2/sports/basketball/mens-college-basketball/news
//All Teams: http://site.api.espn.com/apis/site/v2/sports/basketball/mens-college-basketball/teams
//Specific Team: http://site.api.espn.com/apis/site/v2/sports/basketball/mens-college-basketball/teams/:team
//
import (
	"bytes"
	"context"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	urlPrefix = "apis/site/v2/sports/basketball/"
)

// StatsService handles communication with the Statistics related methods of the ESPN API.
type StatsService service

//ScoreBoard ...
type ScoreBoard struct {
	Leagues []League    `json:"leagues"`
	Season  SeasonShort `json:"season"`
	Day     Date        `json:"day"`
	Events  []Event     `json:"events"`
}

//SeasonType ... definition of a Season type.. regular, ...
type SeasonType struct {
	ID           string `json:"id"`
	Type         int    `json:"type"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

// SeasonDef ... definition of Season from ESPN.com
type SeasonDef struct {
	Year      int        `json:"year"`
	StartDate espnTime   `json:"startDate"`
	EndDate   espnTime   `json:"endDate"`
	Type      SeasonType `json:"type"`
}

// League ... definition of league JSON from ESPN.com
type League struct {
	ID                  string     `json:"id" binding:"required"`
	UID                 string     `json:"uid"`
	Name                string     `json:"name"`
	Abbreviation        string     `json:"abbreviation"`
	Slug                string     `json:"slug,omitempty"`
	Season              SeasonDef  `json:"season"`
	CalendarType        string     `json:"calendarType"`
	CalendarIsWhiteList bool       `json:"calendarIsWhitelist"`
	CalendarStartDate   espnTime   `json:"calendarStartDate"`
	CalendarEndDate     espnTime   `json:"calendarEndDate"`
	Calendar            []espnTime `json:"calendar"`
}

//SeasonShort ...
type SeasonShort struct {
	Year int `json:"year"`
	Type int `json:"type"`
}

//Date ...
type Date struct {
	Date string `json:"date"`
}

//Event ...
type Event struct {
	Extracted    time.Time     `json:"extract_time,omitempty"`
	ExtractedSrc string        `json:"extract_src,omitempty"`
	ID           string        `json:"id" binding:"required"`
	UID          string        `json:"uid" binding:"required"`
	Date         espnTime      `json:"date"`
	Name         string        `json:"name"`
	ShortName    string        `json:"shortName"`
	Season       SeasonShort   `json:"season"`
	Competitions []Competition `json:"competitions"`
	Links        []Link        `json:"links"`
	Status       GameStatus    `json:"status"`
}

//GameStatus ...
// "status":{"clock":0.0,"displayClock":"0.0","period":0,"type":{"id":"1","name":"STATUS_SCHEDULED","state":"pre","completed":false,"description":"Scheduled","detail":"Thu, December 26th at 7:30 PM EST","shortDetail":"12/26 - 7:30 PM EST"}}}
type GameStatus struct {
	Clock        float32        `json:"clock"`
	DisplayClock string         `json:"displayClock"`
	Period       int            `json:"period"`
	StatusType   GameStatusType `json:"type"`
}

//GameStatusType ...
type GameStatusType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Completed   bool   `json:"completed"`
	Description string `json:"description,omitempty"`
	Detail      string `json:"detail,omitempty"`
	ShortDetail string `json:"shortDetail,omitempty"`
}

//Address ...
type Address struct {
	City  string `json:"city"`
	State string `json:"state"`
}

//Venue ...
type Venue struct {
	ID           string    `json:"id" binding:"required"`
	FullName     string    `json:"fullName,omitempty"`
	Address      Address   `json:"address,omitempty"`
	Capacity     int       `json:"capacity"`
	IsIndoor     bool      `json:"indoor"`
}

//Competition ...
type Competition struct {
	ID                    string          `json:"id" binding:"required"`
	UID                   string          `json:"uid" binding:"required"`
	Date                  espnTime        `json:"date"`
	Attendance            int             `json:"Addendance"`
	Type                  CompetitionType `json:"type"`
	TimeValid             bool            `json:"timeValid"`
	NeutralSite           bool            `json:"neutralSite"`
	ConferenceCompetition bool            `json:"conferenceCompeition"`
	Recent                bool            `json:"recent"`
	Venue                 Venue           `json:"venue"`
	Competitors           []Competitor    `json:"competitors"`
	Notes                 []string        `json:"notes"`
	GameStatus            GameStatus      `json:"status"`
	Broadcasts            []Broadcast     `json:"broadcasts"`
	//Tickets
	StartDate     espnTime       `json:"startDate"`
	GeoBroadcasts []GeoBroadcast `json:"geoBroadcasts"`
	Odds          []Odd          `json:"odds"`
}

//GeoBroadcast ...
//"geoBroadcasts":[
//	{	 "type":{"id":"1","shortName":"TV"},
//	"market":{"id":"2","type":"Home"},
//	"media":{"shortName":"FSDT"},
//	"lang":"en","region":"us"}],
type GeoBroadcast struct {
	Type     GBType   `json:"type"`
	Market   GBMarket `json:"market"`
	Media    GBMedia  `json:"media"`
	Language string   `json:"lang,omitempty"`   // "en"
	Region   string   `json:"region,omitempty"` // "us"
}

//GBType ...
type GBType struct {
	ID        string `json:"id"`        // "1"
	ShortName string `json:"shortName,omitempty"` // "TV"
}

//GBMarket ...
type GBMarket struct {
	ID   string `json:"id"`   // "2"
	Type string `json:"type,omitempty"` // "Home"
}

//GBMedia ...
type GBMedia struct {
	ShortName string `json:"shortName,omitempty"` // "FSDT"
}

//Broadcast ...
type Broadcast struct {
	Market string   `json:"market"`
	Names  []string `json:"names,omitempty"`
}

//Odd ...
type Odd struct {
	Provider  OddProvider `json:"provider"`
	Details   string      `json:"details,omitempty"`
	OverUnder float32     `json:"overUnder"`
}

//OddProvider ...
type OddProvider struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name,omitempty"`
	Priority int    `json:"priority"`
}

//CompetitionType ...
type CompetitionType struct {
	ID string `json:"id" binding:"required"`
}

//Competitor ...
type Competitor struct {
	ID         string       `json:"id" binding:"required"`
	UID        string       `json:"uid" binding:"required"`
	Type       string       `json:"type"` // examples "team"
	Order      int          `json:"order"`
	HomeAway   string       `json:"homeAway"` // example "home"
	Winner     bool         `json:"winner"`   //?
	Team       Team         `json:"team"`
	Score      string       `json:"score"`
	Linescores []Linescore  `json:"linescores"`
	Statistics []Statistic  `json:"statistics"`
	Records    []Record     `json:"records"`
	Leaders    []StatLeader `json:"leaders"`
}

//Linescore ...
type Linescore struct {
	Value float32 `json:"value"`
}

//Statistic ...
type Statistic struct {
	Name             string `json:"name"`
	Abbreviation     string `json:"abbreviation"`
	DisplayValue     string `json:"displayValue"`
	RankDisplayValue string `json:"rankDisplayValue,omitempty"`
}

//Record ...
type Record struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation,omitempty"`
	Type         string `json:"type"`
	Summary      string `json:"summary"`
}

//Team ...
type Team struct {
	ID               string      `json:"id" binding:"required"`
	UID              string      `json:"uid,omitempty"`
	Location         string      `json:"location,omitempty"`         //"Toronto",
	Name             string      `json:"name,omitempty"`             // "Raptors"
	Abbreviation     string      `json:"abbreviation,omitempty"`     // "TOR"
	DisplayName      string      `json:"displayName,omitempty"`      // "Toronto Raptors"
	ShortDisplayName string      `json:"shortDisplayName,omitempty"` // Raptors
	Color            string      `json:"color,omitempty"`            //"CEOF41"
	AlternateColor   string      `json:"alternateColor,omitempty"`   //"061922"
	IsActive         bool        `json:"isActive"`
	Venue            Venue       `json:"venue"`
	Links            []Link      `json:"links"`
	Logo             string      `json:"logo,omitempty"`
	Score            string      `json:"score,omitempty"`
	Linescores       []Linescore `json:"linescores"`
}

//Link ...
type Link struct {
	Language   string   `json:"language,omitonempty"`
	Rel        []string `json:"rel"`            // ["clubhouse","desktop","team"],
	HRef       string   `json:"href"`           //"http://www.espn.com/nba/team/_/name/tor/toronto-raptors",
	Text       string   `json:"text,omitempty"` // "Clubhouse"
	Logo       string   `json:"logo,omitempty"` //"https://a.espncdn.com/i/teamlogos/nba/500/scoreboard/tor.png"
	IsExternal bool     `json:"isExternal"`
	IsPremium  bool     `json:"isPremium"`
}

//StatLeader ..
type StatLeader struct {
	Name             string      `json:"name"`             // e.g. "pointsPerGame"
	DisplayName      string      `json:"displayName,omitonempty"`      // e.g. "Points Per Game"
	ShortDisplayName string      `json:"shortDisplayName,omitonempty"` // e.g. "PPG"
	Abbreviation     string      `json:"abbreviation,omitonempty"`     // e.g. "PPG"
	Leaders          []AthLeader `json:"leaders"`
}

//AthLeader ...
type AthLeader struct {
	DisplayValue string  `json:"displayValue"` // "32"
	Value        float32 `json:"value"`        // 32
	Athlete      Athlete `json:"athlete"`
	Team         Team    `json:"team" binding:"required"`
}

//Position ...
type Position struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

//Athlete ...]
type Athlete struct {
	ID          string   `json:"id" binding:"required"` // e.g. "id":"3012",
	FullName    string   `json:"fullName,omitonempty"`              // e.g. "fullName":"Kyle Lowry",
	DisplayName string   `json:"displayName,omitonempty"`           // e.g. "displayName":"Kyle Lowry",
	ShortName   string   `json:"shortName,omitonempty"`             // e.g."K. Lowry",
	Links       []Link   `json:"links"`
	Jersey      string   `json:"jersey,omitonempty"`   // e.g. "jersey":"7",
	Headshot    string   `json:"headshot"` // e.g. "headshot":"https://a.espncdn.com/i/headshots/nba/players/full/3012.png",
	Position    Position `json:"position"`
	Team        Team     `json:"team" binding:"required"`
	Active      bool     `json:"active"`
}

func getRequest(url string) (*http.Request, error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return req, err
	}
	//req.Header.Add("Accept-Encoding", "gzip")

	// get useragent
	agent, exists := os.LookupEnv("ESPN_USERAGENT")
	if !exists {
		fmt.Println("ESPN_USERAGENT not found, should include your website or email credential")
		agent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 Safari/605.1.15"
	}
	// per instructions set user-information to prevent robot blocking
	req.Header.Set("User-Agent", agent)

	//req.Header.Set("Host", "domain.tld")
	return req, err
}

//espnTime is a custom Time parser
type espnTime time.Time

// UnmarshalJSON ... Custom unxmarshall side effect of time.Time not parsing RFC3339
//
func (espnt *espnTime) UnmarshalJSON(bs []byte) error {
	var s string

	if err := json.Unmarshal(bs, &s); err != nil {
		return err
	}

	//TODO: reset string to be a consistent RFC3339 component
	// shift "2019-09-28T07:00Z" to "2019-09-28T00:00:00Z07:00"
	sa := strings.Split(s, "Z")
	s = sa[0] + ":00Z"

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*espnt = espnTime(t)
	return nil
}

//ScoreBoardService will, for a http client, return a ScoreBoard JSON object
//
func (s *StatsService) ScoreBoardService(ctx context.Context) (*ScoreBoard, *Response, error) {

	req, err := s.client.NewRequest("GET", urlPrefix+"/nba/scoreboard", nil)
	
	//to support gzip encoding uncomment... should probably default to true
	//req.Header.Add("Accept-Encoding", "gzip")

	// get useragent from OS Environment Variables -> often needed to prevent robot blocking or API access with lower DoS thresholds
	agent, exists:= os.LookupEnv("ESPN_USERAGENT")
	if (exists) { 
		req.Header.Set("User-Agent", agent)
	}

	sb := &ScoreBoard{}
	resp, err := s.client.Do(ctx, req, sb)
	if err != nil {
		fmt.Printf("Error on new request: %s\n", err)
		return nil, resp, err
	}
	return sb, resp, err
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

func doGet(baseurl string, query string) (*json.Decoder, error) {
	req, err := getRequest(baseurl + query)
	if err != nil {
		fmt.Printf("The HTTP request header building failed with error %s\n", err)
		return nil, err
	}

	// Send req using http Client
	client := &http.Client{}
	fmt.Println("doing HTTP GET")
	resp, err := client.Do(req)

	// received an error on the HTTP request
	if err != nil {
		fmt.Printf("The HTTP request header building failed with error %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("The HTTP request header building failed with error %s : %s\n", resp.Status, data)
		return nil, errors.New(string(data))
	} else if resp.Header.Get("Content-Encoding") == "gzip" {
		fmt.Println("parsing HTTP GZIP-response")
		//resp.Header.Del("Content-Length")
		gz, err := gzip.NewReader(resp.Body)
		defer gz.Close()
		decoder := json.NewDecoder(gz)
		if err != nil {
			fmt.Printf("Error in gzip response decoding %s\n", err)
			return nil, err
		}
		//resp.Body = gzreadCloser{zr, resp.Body}
		return decoder, nil
	}
	// for initial debugging of data structures uncomment
	//data, _ := ioutil.ReadAll(resp.Body)
	//var sb ScoreBoard
	//err = json.Unmarshal(data, &sb)
	//if err != nil {
	//	fmt.Printf("Error in json error unmarshalling %s\n\n", err)
	//}
	//fmt.Printf("The HTTP request header coded %s : %s\n\n", resp.Status, data)
	//fmt.Printf(fmt.Sprintf("Scoreboard: %#v\n", sb))
	decoder := json.NewDecoder(resp.Body)
	return decoder, nil
}

func decodeScoreboard(decoder *json.Decoder) (*ScoreBoard, error) {

	var sb ScoreBoard

	// Decode the response into our Events struct
	err := decoder.Decode(&sb)
	if err != nil {
		fmt.Printf("error caught: %s", err)
		return nil, err
	}
	//setup provenance
	//ev.Extracted = extractTime
	//ev.ExtractedSrc = extractSrc
	//TODO: dig deep into nested structure to set time/src
	return &sb, nil
}

type gzreadCloser struct {
	*gzip.Reader
	io.Closer
}

func (gz gzreadCloser) Close() error {
	return gz.Closer.Close()
}

/*func main() {
	fmt.Println("Loading environment configuration constants")

	//construct the event header
	baseurl, exists := os.LookupEnv("ESPN_STATS_URL")
	if !exists {
		fmt.Println("ESPN_STATS_URL not found, should include your website or email credential")
		baseurl = "https://site.api.espn.com/apis/site/v2/sports/"
	}
	// test fetching BoxScores
	decoder, err := doGet(baseurl, "/basketball/nba/scoreboard")
	if (err) != nil {
		fmt.Printf("error caught: %s", err)
		return
	}
	client := NewClient(nil)
	

	scoreboard, _ := decodeScoreboard(decoder)
	fmt.Printf(fmt.Sprintf("Scoreboard: %#v\n", scoreboard))

	fmt.Println("Terminating the application normally...")
}*/
