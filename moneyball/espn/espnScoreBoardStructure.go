package espn

/**
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
	"encoding/json"
	"fmt"
	"moneyball/go-moneyball/moneyball/ms"
	"strings"
	"time"
)

const (
	//EspnBaseURL is the URL basis for calls to ESPN API's
	EspnBaseURL = "https://site.api.espn.com/"
	//EspnURLPrefix is the URL filepath prefix for calls to v2 of ESPN API's
	EspnURLPrefix = "apis/site/v2/sports/basketball/"
)

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
	UID                 string     `json:"uid" binding:"required"`
	Name                string     `json:"name,omitempty"`
	Abbreviation        string     `json:"abbreviation,omitempty"`
	Slug                string     `json:"slug,omitempty"`
	Season              SeasonDef  `json:"season,omitempty"`
	CalendarType        string     `json:"calendarType,omitempty"`
	CalendarIsWhiteList bool       `json:"calendarIsWhitelist,omitempty"`
	CalendarStartDate   espnTime   `json:"calendarStartDate,omitempty"`
	CalendarEndDate     espnTime   `json:"calendarEndDate,omitempty"`
	Calendar            []espnTime `json:"calendar,omitempty"`
	Teams               []Team     `json:"teams,omitempty"`
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
	Extracted    *time.Time    `json:"extract_time,omitempty"`
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
	ID       string  `json:"id" binding:"required"`
	FullName string  `json:"fullName,omitempty"`
	Address  Address `json:"address,omitempty"`
	Capacity int     `json:"capacity"`
	IsIndoor bool    `json:"indoor"`
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
	ID        string `json:"id"`                  // "1"
	ShortName string `json:"shortName,omitempty"` // "TV"
}

//GBMarket ...
type GBMarket struct {
	ID   string `json:"id"`             // "2"
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
	ID               string        `json:"id" binding:"required"`
	UID              string        `json:"uid" binding:"required"`
	Slug             string        `json:"slug,omitempty"`
	Location         string        `json:"location,omitempty"`         //"Toronto",
	Name             string        `json:"name,omitempty"`             // "Raptors"
	Abbreviation     string        `json:"abbreviation,omitempty"`     // "TOR"
	DisplayName      string        `json:"displayName,omitempty"`      // "Toronto Raptors"
	ShortDisplayName string        `json:"shortDisplayName,omitempty"` // Raptors
	Color            string        `json:"color,omitempty"`            //"CEOF41"
	AlternateColor   string        `json:"alternateColor,omitempty"`   //"061922"
	IsActive         bool          `json:"isActive"`
	IsAllStar        bool          `json:"isAllStar"`
	Venue            Venue         `json:"venue"`
	Links            []Link        `json:"links,omitempty"`
	Logos            []Link        `json:"logos,omitempty"`
	Logo             string        `json:"logo,omitempty"`
	Score            string        `json:"score,omitempty"`
	Linescores       []Linescore   `json:"linescores,omitempty"`
	Record           []RecordItems `json:"record,omitempty"`
}

// RecordItems ...
type RecordItems struct {
	Items []Item `json:"items,omitempty"`
}

// TeamStatistic ...
type TeamStatistic struct {
	Name  string  `json:"name"`  //"name":"playoffSeed",
	Value float32 `json:"value"` //"value":15.0},
}

// Item ...
type Item struct {
	Summary string          `json:"summary"` //"summary":"7-27",
	Stats   []TeamStatistic `json:"stats"`   //"stats":[
}

//Link ...
type Link struct {
	Language   string   `json:"language,omitempty"`
	Rel        []string `json:"rel"`            // ["clubhouse","desktop","team"],
	HRef       string   `json:"href"`           //"http://www.espn.com/nba/team/_/name/tor/toronto-raptors",
	Text       string   `json:"text,omitempty"` // "Clubhouse"
	Logo       string   `json:"logo,omitempty"` //"https://a.espncdn.com/i/teamlogos/nba/500/scoreboard/tor.png"
	IsExternal bool     `json:"isExternal"`
	IsPremium  bool     `json:"isPremium"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
}

//StatLeader ..
type StatLeader struct {
	Name             string      `json:"name"`                       // e.g. "pointsPerGame"
	DisplayName      string      `json:"displayName,omitempty"`      // e.g. "Points Per Game"
	ShortDisplayName string      `json:"shortDisplayName,omitempty"` // e.g. "PPG"
	Abbreviation     string      `json:"abbreviation,omitempty"`     // e.g. "PPG"
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
	FullName    string   `json:"fullName,omitempty"`    // e.g. "fullName":"Kyle Lowry",
	DisplayName string   `json:"displayName,omitempty"` // e.g. "displayName":"Kyle Lowry",
	ShortName   string   `json:"shortName,omitempty"`   // e.g."K. Lowry",
	Links       []Link   `json:"links"`
	Jersey      string   `json:"jersey,omitempty"` // e.g. "jersey":"7",
	Headshot    string   `json:"headshot"`         // e.g. "headshot":"https://a.espncdn.com/i/headshots/nba/players/full/3012.png",
	Position    Position `json:"position"`
	Team        Team     `json:"team" binding:"required"`
	Active      bool     `json:"active"`
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

//Sport ...
type Sport struct {
	ID      string   `json:"id" binding:"required"`  //"id":"40",
	UID     string   `json:"uid" binding:"required"` //	"uid":"s:40",
	Name    string   `json:"name,omitempty"`         //	"name":"Basketball",
	Slug    string   `json:"slug,omitempty"`         //	"slug":"basketball",
	Leagues []League `json:"leagues,omitempty"`      // "leagues":[
}

//TeamSport ... array of teams for json depacking
type TeamSport struct {
	Sport []Sport `json:"sports"`
}

//MarshalMS marshalls espn.Scoreboard structures to ms.Scoreboard structures, can return partial results
//in the case of one event causing an error deep in the array
func (s *ScoreBoard) MarshalMS() (*ms.ScoreBoard, error) {
	sb := ms.ScoreBoard{}
	bs := []ms.Event{}
	for _, event := range s.Events {
		evented, err := (&event).MarshalMSEvent(s.Leagues[0])
		bs = append(bs, *evented)
		if err != nil {
			sb.Events = bs
			return &sb, err
		}
	}
	sb.Events = bs
	return &sb, nil
}

//MarshalMSEvent marshals espn.Event to ms.Event
func (e *Event) MarshalMSEvent(l League) (*ms.Event, error) {
	bs := ms.Event{}
	/* Event
	EntityID
	GameID     GameID      `json:"gameId"`
	League     League      `json:"league"`
	Season     Season      `json:"season"`
	HomeTeam   *Competitor `json:"homeTeam"`
	VisitTeam  *Competitor `json:"visitTeam"`
	Venue      *Venue      `json:"location,omitempty"`
	Status     *GameStatus  `json:"status,omitempty"`
	Links      *[]Link     `json:"link,omitempty"`
	GameDetail *GameDetail `json:"gameDetail,omitempty"`
	*/
	eID := ms.EntityID{}
	eID.Extracted = e.Extracted
	eID.ExtractedSrc = e.ExtractedSrc
	bs.EntityID = eID
	bs.GameID = ms.GameID(e.ID)
	bs.League = ms.League(l.Abbreviation)
	bs.Season = ms.Season{e.Season.Year, e.Season.Type}

	if len(e.Competitions) > 1 {
		// set error
		fmt.Printf("error: compeitions should be 1 %d", len(e.Competitions))
	}

	for _, ref := range e.Competitions[0].Competitors {
		switch ref.HomeAway {
		case "home":
			bs.HomeTeam, _ = ref.marshalMSCompetitor()
		case "away":
			bs.VisitTeam, _ = ref.marshalMSCompetitor()
		default:
			//throw error...
			fmt.Printf("error: compeition should be home or away... found %s", ref.HomeAway)
		}
	}

	venue := e.Competitions[0].Venue
	bs.Venue = &ms.Venue{ms.EntityID{"", nil, ""}, venue.ID, venue.FullName, marshalMSAddress(venue.Address), venue.Capacity, venue.IsIndoor}
	bs.Status, _ = marshalMSGameStatus(e.Status)

	links := []ms.Link{}
	for _, link := range e.Links {
		l, _ := marshalMSLink(link)
		links = append(links, *l)
	}
	bs.Links = &links

	//TODO: GameDetail
	/*
		//GameDetail .. extra detail about the game including things like startTime...
		type GameDetail struct {
			StartTime           *time.Time  `json:"startTimeUTC,omitempty"`     //"startTimeUTC":"2019-10-01T00:00:00.000Z",
			StartDateEastern    string      `json:"startDateEastern,omitempty"` //"startDateEastern":"20190930",
			StartTimeEastern    string      `json:"startTimeEastern,omitempty"`
			Period              *GamePeriod `json:"period,omitempty,omitempty"` // "period": {}
			Attendance          string      `json:"attendance,omitempty"`       //"attendance":"18624",
			GameDurationMinutes int         `json:"gameDuration,omitempty"`
		} */
	gd := ms.GameDetail{}
	refTime := time.Time(e.Date)
	fmt.Printf("timeRef %s\n", refTime)
	gd.StartTime = &refTime
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		// set error
		fmt.Printf("error: timezone conversion %#v\n", err)
	}
	est := refTime.In(location)
	fmt.Printf("time: UTC %s, EST %s\n", refTime, est)

	gd.StartDateEastern = est.Format("2006-01-02")
	gd.StartTimeEastern = est.Format("15:04:05")
	/*
		type GamePeriod struct {
		Current       int  `json:"current"`       //"current":4,
		Type          int  `json:"type"`          //"type":0,
		MaxRegular    int  `json:"maxRegular"`    //"maxRegular":4
		IsHalftime    bool `json:"isHalftime"`    //`"isHalftime":false,
		IsEndOfPeriod bool `json:"isEndOfPeriod"` //"isEndOfPeriod":false
		}*/
	//gd.Period = &ms.GamePeriod{0,0,4,false,false}
	gd.Attendance = e.Competitions[0].Attendance
	//gd.GameDurationMinutes =
	bs.GameDetail = &gd

	/* Event
	Extracted    *time.Time     `json:"extract_time,omitempty"`
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
	*/
	ms.MasterIdentity(&bs)
	return &bs, nil
}

func marshalMSAddress(a Address) *ms.Address {
	addr := ms.Address{}
	addr.City = a.City
	addr.State = a.State
	return nil
}

func (comp *Competitor) marshalMSCompetitor() (*ms.Competitor, error) {
	c := ms.Competitor{}
	t := (*comp).Team
	c.Name = t.Name
	c.Abbreviation = t.Abbreviation
	//c.Record = t.
	linescores := []ms.Score{}
	for _, lsc := range t.Linescores {
		linescores = append(linescores, ms.Score{lsc.Value})
	}
	c.LineScore = &linescores
	c.Location = (*comp).Team.Location
	c.Color = (*comp).Team.Color
	c.AlternateColor = (*comp).Team.AlternateColor
	c.IsActive = (*comp).Team.IsActive
	c.IsAllStar = (*comp).Team.IsAllStar
	links := []ms.Link{}
	for _, link := range (*comp).Team.Links {
		l, _ := marshalMSLink(link)
		links = append(links, *l)
	}
	c.Links = &links

	/*
		type Team struct {
		ID               string        `json:"id" binding:"required"`
		UID              string        `json:"uid" binding:"required"`
		Slug             string        `json:"slug,omitempty"`
		Location         string        `json:"location,omitempty"`         //"Toronto",
		Name             string        `json:"name,omitempty"`             // "Raptors"
		Abbreviation     string        `json:"abbreviation,omitempty"`     // "TOR"
		DisplayName      string        `json:"displayName,omitempty"`      // "Toronto Raptors"
		ShortDisplayName string        `json:"shortDisplayName,omitempty"` // Raptors
		Color            string        `json:"color,omitempty"`            //"CEOF41"
		AlternateColor   string        `json:"alternateColor,omitempty"`   //"061922"
		IsActive         bool          `json:"isActive"`
		IsAllStar        bool          `json:"isAllStar"`
		Venue            Venue         `json:"venue"`
		Links            []Link        `json:"links,omitempty"`
		Logos            []Link        `json:"logos,omitempty"`
		Logo             string        `json:"logo,omitempty"`
		Score            string        `json:"score,omitempty"`
		Linescores       []Linescore   `json:"linescores,omitempty"`
		Record           []RecordItems `json:"record,omitempty"`}*/

	/*
		type Competitor struct {
		EntityID
		Name           string   `json:"name,omitempty"`
		Abbreviation   string   `json:"abbreviation"`
		Record         Record   `json:"record,omitempty"`
		Score		   int 		`json:"score"`
		LineScore      *[]Score `json:"linescore,omitempty"` //"linescore":[{"score":"30"},{"score":"32"},{"score":"23"},{"score":"19"}]},
		Location       string   `json:"location"`
		Color          string   `json:"color"`
		AlternateColor string   `json:"alternateColor"`
		IsActive       bool     `json:"isActive"`
		IsAllStar      bool     `json:"isAllStar"`
		Link           *Link    `json:"logos"`}*/
	return &c, nil
}

func marshalMSGameStatus(gs GameStatus) (*ms.GameStatus, error) {
	return &ms.GameStatus{gs.Clock, gs.Period, gs.StatusType.State, gs.StatusType.Detail}, nil
}

func marshalMSLink(l Link) (*ms.Link, error) {
	link := ms.Link{}
	if l.Logo != "" {
		//we have a logo reference... so we need dimensions
		link.HRef = l.Logo
		link.IsLogo = true
	} else {
		link.HRef = l.HRef
	}
	link.Rel = l.Rel
	link.Alt = l.Text
	link.Dimension = &ms.LinkDimensions{l.Width, l.Height}
	return &link, nil
}
