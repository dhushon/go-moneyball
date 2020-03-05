package nba

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

//navigation link with php source: http://nbasense.com/nba-api/Stats/Data/
//http://data.nba.net/json/bios/player_201935.json

//http://data.nba.net/json/cms/noseason/game/{gameDate}/{gameId}/boxscore.json
//http://data.nba.net/json/cms/noseason/game/20170201/0021600732/boxscore.json where game/`json:"date"`/`json:"ID"`/boxscore.json`
//REF: http://nbasense.com/nba-api/Data/Cms/Game/Boxscore

import (
	"go-moneyball/moneyball/ms"
	"log"
	"strconv"
	"time"
)

const (
	//DataNBABaseURLv2 ...
	DataNBABaseURLv2 = "https://data.nba.net/"
	//DataNBAProdURLPrefixv2 ...
	DataNBAProdURLPrefixv2 = "prod/v2/"
)

//CMSProdv2Schedule ... based upon this structure http://data.nba.net/prod/v2/2019/schedule.json
type CMSProdv2Schedule struct {
	InternalStuff  InternalProdv2   `json:"_internal"` //"_internal":{}
	LeagueSchedule LeagueSchedulev2 `json:"league"`    //"league":{}
}

//CMSProdv1BoxScore ...
type CMSProdv1BoxScore struct {
	InternalStuff *InternalProdv2  `json:"_internal"`     //"_internal":{}
	Game          *ScheduledGamev2 `json:"basicGameData"` //"basicGameData":{}
	PrevMatchup   *GamePointerv2   `json:"previoudMatchup"`
	BoxStats	  *BoxStats  		`json:"stats"`

}

//GamePointerv2 ...
type GamePointerv2 struct {
	GameID           string `json:"gameId"`   //"gameId":"0021600572",
	StartDateEastern string `json:"gameDate"` //"gameDate":"20170110"},

}

//InternalProdv2 ... source provenance information, what we care about is pubDate
type InternalProdv2 struct {
	PubDateTime string `json:"pubDateTime"` //"pubDateTime":"2020-01-06 10:20:39.221 EST",
}

//LeagueSchedulev2 ... we don't understand what happens with non-standard?
type LeagueSchedulev2 struct {
	Events []ScheduledGamev2 `json:"standard"` //"standard":{}
}

//ScheduledGamev2 ... based upon this structure
type ScheduledGamev2 struct {
	GameID           string    `json:"gameId"`                //"gameId":"0011900001",
	SeasonStageID    int       `json:"seasonStageId"`         //"seasonStageId":1,
	SeasonYear       FlexInt   `json:"seasonYear"`            //"seasonYear":"2016",
	Arena            Arena     `jsonm:"arena"`                //"arena": {	"name":"TD Garden","isDomestic":false,"city":"Boston","stateAbbr":"MA","country":""},
	IsGameActivated  bool      `json:"isGameActivated"`       //"isGameActivated":false,
	GameURLCode      string    `json:"gameUrlCode,omitempty"` //"gameUrlCode":"20190930/SDSHOU", //NOTE: this is Date(Eastern)/vTeamhTeam
	StatusNum        int       `json:"statusNum"`             //"statusNum":3,
	ExtStatusNum     int       `json:"extendedStatusNum"`     //"extendedStatusNum":0,
	IsStartTimeTBD   bool      `json:"isStartTimeTBD"`        //"isStartTimeTBD":false,
	StartTime        time.Time `json:"startTimeUTC"`          //"startTimeUTC":"2019-10-01T00:00:00.000Z",
	StartDateEastern string    `json:"startDateEastern"`      //"startDateEastern":"20190930",
	StartTimeEastern string    `json:"startTimeEastern"`      //"startTimeEastern":"8:00 PM ET",
	IsBuzzerBeater   bool      `json:"isBuzzerBeater"`        //"isBuzzerBeater":false,
	//"isPreviewArticleAvail":false,
	//"isRecapArticleAvail":false,
	//"tickets": {"mobileApp":"https://a.data.nba.com/tickets/single/2016/0021600732/APP_TIX","desktopWeb":"https://a.data.nba.com/tickets/single/2016/0021600732/TEAM_SCH","mobileWeb":"https://a.data.nba.com/tickets/single/2016/0021600732/WEB_MWEB"},
	//"hasGameBookPdf":true,
	Period GamePeriodv2 `json:"period"` // "period": {}
	//"nugget": {"text":""},
	Attendance   string        `json:"attendance,omitempty"`   //"attendance":"18624",
	GameDuration *GameDuration `json:"gameDuration,omitempty"` //"gameDuration":{"hours":"2","minutes":"33"},
	HomeTeam     GameTeamv2    `json:"hTeam"`                  //"hTeam":{"teamId":"1610612745","score":"140","win":"1","loss":"0"},
	VisitingTeam GameTeamv2    `json:"vTeam"`                  //"vTeam":{"teamId":"12329","score":"71","win":"0","loss":"1"},
	//Watch        json.RawMessage `json:"watch"` //"watch":{"broadcast":{"video":{"regionalBlackoutCodes":"","isLeaguePass":true,"isNationalBlackout":false,"isTNTOT":false,"canPurchase":false,"isVR":false,"isNextVR":false,"isNBAOnTNTVR":false,"isMagicLeap":false,"isOculusVenues":false,"national":{"broadcasters":[{"shortName":"NBA TV","longName":"NBA TV"}]},"canadian":[{"shortName":"NBAC","longName":"NBA TV Canada"}],"spanish_national":[]}}}},
}

//GameDuration ...
type GameDuration struct {
	Hours   FlexInt `json:"hours"`
	Minutes FlexInt `json:"minutes"`
}

//GamePeriodv2 is inforamtion about where we are in the game? or where the game was at completion (unk.)
type GamePeriodv2 struct {
	Current       int  `json:"current"`       //"current":4,
	Type          int  `json:"type"`          //"type":0,
	MaxRegular    int  `json:"maxRegular"`    //"maxRegular":4
	IsHalftime    bool `json:"isHalftime"`    //`"isHalftime":false,
	IsEndOfPeriod bool `json:"isEndOfPeriod"` //"isEndOfPeriod":false
}

//GameTeamv2 is som subset of team information associated with the schedule (including boxscore?)
type GameTeamv2 struct {
	TeamID     string    `json:"teamId"`               //"teamId":"1610612745"
	TriCode    string    `json:"triCode,omitempty"`    //"triCode":"TOR",
	Score      FlexInt   `json:"score"`                // "score":"71"
	Win        FlexInt   `json:"win"`                  // "win":"1"
	Loss       FlexInt   `json:"loss"`                 // "loss":"1"
	SeriesWin  FlexInt   `json:"seriesWin,omitempty"`  //"seriesWin":"2",
	SeriesLoss FlexInt   `json:"seriesLoss,omitempty"` //"seriesLoss":"0",
	Linescore  []Scorev2 `json:"linescore"`            //"linescore":[{"score":"30"},{"score":"32"},{"score":"23"},{"score":"19"}]},
}

//Arena - location where game is played TODO: timezone?
type Arena struct {
	Name       string `json:"name"`       //"name":"TD Garden",
	IsDomestic bool   `json:"isDomestic"` //"isDomestic":false,
	City       string `json:"city"`       //"city":"Boston",
	State      string `json:"stateAbbr"`  //"stateAbbr":"MA",
	Country    string `json:"country"`    //"country":""
}

//Scorev2 - used for linescores
type Scorev2 struct {
	Score FlexInt `json:"score"`
}

//MarshalMS marshalls espn.Scoreboard structures to ms.Scoreboard structures, can return partial results
//in the case of one event causing an error deep in the array
func MarshalMS(s *LeagueSchedulev2) (*ms.ScoreBoard, error) {
	sb := ms.ScoreBoard{}
	bs := []ms.Event{}
	for _, event := range s.Events {
		evented, err := (&event).MarshalMSEvent()
		bs = append(bs, *evented)
		if err != nil {
			sb.Events = bs
			return &sb, err
		}
	}
	sb.Events = bs
	return &sb, nil
}

//MarshalMSEvent marshals nba.Event to ms.BoxScore
func (e *ScheduledGamev2) MarshalMSEvent() (*ms.Event, error) {
	bs := ms.Event{}
	/* ms.Event
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
	//eID.Extracted = e.Extracted
	//eID.ExtractedSrc = e.ExtractedSrc
	bs.EntityID = eID
	bs.GameID = ms.GameID(e.GameID)
	bs.League = ms.League("NBA")
	bs.Season = ms.Season{SeasonYear: int(((*e).SeasonYear)), SeasonStage: (*e).SeasonStageID}
	bs.HomeTeam, _ = (*e).HomeTeam.marshalMSCompetitor()
	bs.VisitTeam, _ = (*e).VisitingTeam.marshalMSCompetitor()
	bs.Venue, _ = (*e).Arena.marshalMSVenue()
	bs.GameDetail = e.marshalMSGameDetail()

	ms.MasterIdentity(&bs)
	return &bs, nil
}

func (e *ScheduledGamev2) marshalMSGameDetail() *ms.GameDetail {
	/*
		type GameDetail struct {
		StartTime           *time.Time  `json:"startTimeUTC,omitempty"`     //"startTimeUTC":"2019-10-01T00:00:00.000Z",
		StartDateEastern    string      `json:"startDateEastern,omitempty"` //"startDateEastern":"20190930",
		StartTimeEastern    string      `json:"startTimeEastern,omitempty"`
		Period              *GamePeriod `json:"period,omitempty"`     // "period": {}
		Attendance          int         `json:"attendance,omitempty"` //"attendance":"18624",
		GameDurationMinutes int         `json:"gameDuration,omitempty"`
		}*/
	gd := ms.GameDetail{}
	refTime := time.Time((*e).StartTime)
	//fmt.Printf("timeRef %s\n", refTime)
	gd.StartTime = &refTime
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		// set error
		log.Printf("error: timezone conversion %#v\n", err)
	}
	est := refTime.In(location)
	//fmt.Printf("time: UTC %s, EST %s\n", refTime, est)

	gd.StartDateEastern = est.Format("2006-01-02")
	gd.StartTimeEastern = est.Format("15:04:05")
	//gd.StartDateEastern = (*e).StartDateEastern  // ignore initial - wrong formatting
	//gd.StartTimeEastern = (*e).StartTimeEastern  // ignore initial - wrong formatting
	//TODO: Period
	if ((*e).Attendance == "") {
		gd.Attendance = 0 
	} else {
		gd.Attendance, err = strconv.Atoi((*e).Attendance)
		if (err != nil ) {
			log.Printf("Error strconv.Atoi e.Attendance, %s set to zero", (*e).Attendance)
		}
		gd.Attendance = 0
	}	
	gd.GameDurationMinutes = ((int((*e).GameDuration.Hours) * 60) + int((*e).GameDuration.Minutes))
	return &gd
}

func (t *GameTeamv2) marshalMSCompetitor() (*ms.Competitor, error) {
	c := ms.Competitor{}
	c.ID = t.TeamID
	c.Abbreviation = t.TriCode
	//c.Record = t.
	linescores := []ms.Score{}
	for _, lsc := range t.Linescore {
		linescores = append(linescores, ms.Score{Score: float32(lsc.Score)})
	}
	c.LineScore = &linescores
	c.Score = int(t.Score)
	return &c, nil
}

func (a *Arena) marshalMSVenue() (*ms.Venue, error) {
	v := ms.Venue{}
	v.FullName = a.Name
	v.Address = &ms.Address{Street: "", City: a.City, State: a.State, Country: a.Country}
	_, err := ms.GetGeoCodeAddress(&v)
	//TODO: Setup EntityID..
	return &v, err
}

//marshalMSGamePlayerStat ... long name, but need to map NBA BoxScore Player stats to MS.GamePlayerStats
func (ps *PlayerStats) marshalMSGamePlayerStat(gameID string, teamID string) (*ms.GamePlayersStats, error) {
	
	return nil, nil
}
