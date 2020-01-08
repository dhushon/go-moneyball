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

//navigation link with php source: http://nbasense.com/nba-api/Stats/Data/
//http://data.nba.net/json/bios/player_201935.json

//http://data.nba.net/json/cms/noseason/game/{gameDate}/{gameId}/boxscore.json
//http://data.nba.net/json/cms/noseason/game/20170201/0021600732/boxscore.json where game/`json:"date"`/`json:"ID"`/boxscore.json`
//REF: http://nbasense.com/nba-api/Data/Cms/Game/Boxscore

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	dataNBABaseURLv2       = "https://data.nba.net/"
	dataNBAProdURLPrefixv2 = "prod/v2/"
)

//CMSProdv2Schedule ... based upon this structure http://data.nba.net/prod/v2/2019/schedule.json
type CMSProdv2Schedule struct {
	InternalStuff  InternalProdv2   `json:"_internal"` //"_internal":{}
	LeagueSchedule LeagueSchedulev2 `json:"league"`    //"league":{}
}

//CMSProdv1BoxScore ...
type CMSProdv1BoxScore struct {
	InternalStuff InternalProdv2  `json:"_internal"`     //"_internal":{}
	Game          ScheduledGamev2 `json:"basicGameData"` //"basicGameData":{}
	PrevMatchup   GamePointerv2   `json:"previoudMatchup"`
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
	Standard []ScheduledGamev2 `json:"standard"` //"standard":{}
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

func nbaPathModifier(orig string, modifier map[string]string) (string, error) {
	for param := range modifier { //go thru find/replace on map
		//note that source string de-mark is '{' '}' eted.
		search := "{" + (strings.ToLower(param)) + "}"
		orig = strings.Replace(orig, search, modifier[param], 1)
	}
	if strings.Contains(orig, "{") {
		return orig, fmt.Errorf("new Path: %s continues to include variables not satisfied", orig)
	}
	return orig, nil
}

//NBABoxScoreServicev2 will, for a http client, provide a ScheduledGame ( note that this is not yet normalized to structures)
//		boxscorev1 http://data.nba.net/prod/v1/{gameDate}/{gameId}_boxscore.json e.g. http://data.nba.net/prod/v1/20170201/0021600732_boxscore.json
func (s *ScoreService) NBABoxScoreServicev2(ctx context.Context, modifier map[string]string) (*ScheduledGamev2, *Response, error) {

	s.client.BaseURL, _ = url.Parse(dataNBABaseURLv2)
	path := "prod/v1/{gamedate}/{gameid}_boxscore.json"
	suffix, err := nbaPathModifier(path, modifier)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", suffix, nil)
	if err != nil {
		return nil, nil, err
	}

	//to support gzip encoding uncomment... should probably default to true
	//req.Header.Add("Accept-Encoding", "gzip")

	// get useragent from OS Environment Variables -> often needed to prevent robot blocking or API access with lower DoS thresholds
	//agent, exists := os.LookupEnv("NBA_USERAGENT")
	//if exists {
	//	req.Header.Set("User-Agent", agent)
	//}
	event := &CMSProdv1BoxScore{}
	resp, err := s.client.Do(ctx, req, event, true)
	if err != nil {
		fmt.Printf("Error on new request: %s\n", err)
		return nil, resp, err
	}
	return &event.Game, resp, err
}

//NBAScheduleServicev2 is an updated nba feed for NBA Schefule information
//http://data.nba.net/prod/v2/{year}/schedule.json e.g. http://data.nba.net/prod/v2/2019/schedule.json
func (s *ScheduleService) NBAScheduleServicev2(ctx context.Context, modifier map[string]string) (*[]ScheduledGamev2, *Response, error) {

	s.client.BaseURL, _ = url.Parse(dataNBABaseURLv2)
	path := "prod/v2/{year}/schedule.json"
	suffix, err := nbaPathModifier(path, modifier)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", suffix, nil)
	if err != nil {
		return nil, nil, err
	}

	event := &CMSProdv2Schedule{}
	resp, err := s.client.Do(ctx, req, event, true)
	if err != nil {
		fmt.Printf("Error caught: %s\n", err)
	}
	return &event.LeagueSchedule.Standard, resp, err
}