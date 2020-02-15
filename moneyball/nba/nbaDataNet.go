package nba

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

//navigation link with php source: http://nbasense.com/nba-api/Stats/Data/
//http://data.nba.net/json/bios/player_201935.json

//http://data.nba.net/json/cms/noseason/game/{gameDate}/{gameId}/boxscore.json
//http://data.nba.net/json/cms/noseason/game/20170201/0021600732/boxscore.json where game/`json:"date"`/`json:"ID"`/boxscore.json`
//REF: http://nbasense.com/nba-api/Data/Cms/Game/Boxscore

import (
	"encoding/json"
	"log"
	"strconv"
)

const (
	//DataNBABaseURL ...
	DataNBABaseURL = "https://data.nba.net/"
	//DataNBAURLPathPrefix ...
	DataNBAURLPathPrefix = "json/cms/"
	//PlayerBioPath ...
	PlayerBioPath = "bios/player_201935.json"
	//BoxScorePath ...
	BoxScorePath = "noseason/game/"
)

//SeasonMeta ... //TODO.. nba is inconsistent in structure of SeasonMeta... sometimes strings, someimtes ints...
type SeasonMeta struct {
	CalendarDate        FlexInt `json:"calendar_date"`         //`"calendar_date":20160908, 20190522
	SeasonYear          FlexInt `json:"season_year"`           //"season_year":2016, 2018
	StatsSeasonYear     FlexInt `json:"stats_season_year"`     //"stats_season_year":2015,2018
	StatsSeasonID       FlexInt `json:"stats_season_id"`       //"stats_season_id":42015,42018
	StatsSeasonStage    FlexInt `json:"stats_season_stage"`    //"stats_season_stage":4, 4
	RosterSeasonYear    FlexInt `json:"roster_season_year"`    //"roster_season_year":2016,
	ScheduleSeasonYear  FlexInt `json:"schedule_season_year"`  //"schedule_season_year":2016,
	StandingsSeasonYear FlexInt `json:"standings_season_year"` //"standings_season_year":2016,
	SeasonID            FlexInt `json:"season_id"`             //"season_id":22016,
	DisplayYear         string  `json:"display_year"`          //"display_year":"2016-17",
	DisplaySeason       string  `json:"display_season"`        //"display_season":"Regular Season", "Post Season"
	SeasonStage         FlexInt `json:"season_stage"`          //"season_stage":2},
	LeagueID            FlexInt `json:"league_id,omitempty"`   //"league_id":"00"
}

//SportsSchedule ...
type SportsSchedule struct {
	Games []ScheduledGame `json:"game"`
}

//ScheduledGame ...
type ScheduledGame struct {
	HomeAbbreviation    string  `json:"h_abrv"` //"h_abrv":"TOR",
	VisitorAbbreviation string  `json:"v_abrv"` //"v_abrv":"GSW",
	GameID              FlexInt `json:"id"`     //"id":"0011600001", //inconsistent string vs. id
	DateTime            string  `json:"dt"`     //"dt":"2016-10-01 19:30:00.0", //TODO normalize to datetime
	RReg                string  `json:"r_reg"`  //"r_reg":"", have no idea what this is
	IsLP                bool    `json:"is_lp"`  //"is_lp":true,
	SG                  bool    `json:"sg"`     //"sg":false
}

//SportsEvent ...
type SportsEvent struct {
	Event SportsContent `json:"sports_content"`
}

//SportsContent ...
type SportsContent struct {
	Meta     SportsMeta     `json:"sports_meta"`
	Game     SportsGame     `json:"game"`
	Schedule SportsSchedule `json:"schedule"`
}

// SportsMeta ...
type SportsMeta struct {
	DateTime   string     `json:"date_time"`   //TODO fix time date for strong temporal alignment "20170510 1438"
	SeasonMeta SeasonMeta `json:"season_meta"` //TODO rebuild parser as SeasonMeta is inconsistent from NBA
	Next       CMSTarget  `json:"next"`
}

//CMSTarget ...
type CMSTarget struct {
	URL string `json:"url"`
}

//SportsGame ...
type SportsGame struct {
	GameID            string           `json:"id"`                 //"id":"0021600732",
	GameURL           string           `json:"game_url"`           //"game_url":"20170201\/TORBOS",
	SeasonID          string           `json:"season_id"`          //"season_id":"22016",
	Date              string           `json:"date"`               //"date":"20170201",
	Time              string           `json:"time"`               //"time":"1930",
	Arena             string           `json:"arena"`              //"arena":"TD Garden",
	City              string           `json:"city"`               //"city":"Boston",
	State             string           `json:"state"`              //"state":"MA",
	Country           string           `json:"country"`            //"country":"",
	HomeStartDate     string           `json:"home_start_date"`    //:"20170201",
	HomeStartTime     string           `json:"home_start_time"`    //	"home_start_time":"1930",
	VisitorStartDate  string           `json:"visitor_start_date"` //"visitor_start_date":"20170201",
	VisitorStartTime  string           `json:"visitor_start_time"` //"visitor_start_time":"1930",
	PreviewAvailable  string           `json:"previewAvailabile"`  //"previewAvailable":"0", //TODO: convert 0 to false (boolean)
	RecapAvailable    string           `json:"recapAvailable"`     //"recapAvailable":"0", //TODO: convert 0 to false (boolean)
	NotebookAvailable string           `json:"notebookAvailable"`  //"notebookAvailable":"0",//TODO: convert 0 to false (boolean)
	TNTOT             string           `json:"tnt_ot"`             //	"tnt_ot":"0", //TODO: convert 0 to false (boolean)
	Attendance        string           `json:"attendance"`         //	"attendance":"18624", //TODO: convert string to int
	Officials         []OfficialPerson `json:"officials"`          //"officials":[{"person_id":"1146","first_name":"Tony","last_name":"Brothers","jersey_number":"25"},
	Ticket            json.RawMessage  `json:"ticket,omitempty"`
	Broadcast         json.RawMessage  `json:"broadcasters,omitempty"`
	PeriodTime        TimePeriod       `json:"period_time,omitempty"`
	Visitor           WorkingTeam      `json:"visitor"`
	Home              WorkingTeam      `json:"home"`
}

//Bio is the tag for a player bio
type Bio struct {
	Player Player `json:"Bio"`
}

//OfficialPerson ... is a person that is a game official
type OfficialPerson Person

//Player is some basic information on player
type Player struct {
	PlayerID    string `json:"id"`                    //"id":"201935",
	Type        string `json:"type"`                  //"type":"player",
	DisplayName string `json:"display_name"`          //"display_name":"Harden, James",
	Abstract    string `json:"professional"`          //"professional":"html bio here"
	College     string `json:"college,omitempty"`     //"college":"",
	Highschool  string `json:"highschool,omitempty"`  //"highschool":"",
	Twitter     string `json:"twitter,omitempty"`     //"twitter":"",
	OtherLabel  string `json:"other_label,omitempty"` //"other_label":"",
	OtherText   string `json:"other_text,omitempty"`  //"other_text":""}
}

//Person ...
type Person struct {
	FirstName string  `json:"first_name"`
	FName     string  `json:"FirstName,omitempty"` //used in game->team->leaders
	LastName  string  `json:"last_name"`
	LName     string  `json:"LastName"`                //used in game->team->leaders
	JerseyStr FlexInt `json:"jersey_number,omitempty"` //TODO unmarshall to int
	//PlayerCode 			   string `json:"player_code,omitempty"` //"PlayerCode":"kyle_lowry" //used in game->team->leaders TODO redefinition of player_code?
	PersonID               string  `json:"person_id"` //TODO unmarshal to ID[int]
	PersonID2              string  `json:"PersonID"`  // TODO reconcile PersonID2 and PersonID
	PositionShort          string  `json:"position_short,omitempty"`
	PositionFull           string  `json:"position_full,omitempty"`
	Minutes                FlexInt `json:"minutes,omitempty"`                  //TODO unmarshal to int and add to event-stats?
	Seconds                FlexInt `json:"seconds,omitempty"`                  //TODO unmarshal to int and add to event-stats?
	Points                 FlexInt `json:"points,omitempty"`                   //TODO unmarshal to int and add to event-stats?
	FieldGoalsMade         FlexInt `json:"field_goals_made,omitempty"`         //TODO unmarshal to int and add to event-stats?
	FieldGoalsAttempted    FlexInt `json:"field_goals_attempted,omitempty"`    //TODO unmarshal to int and add to event-stats?
	PlayerCode             FlexInt `json:"player_code,omitempty"`              //TODO unmarshal to int and add to event-stats?
	FreeThrowsMade         FlexInt `json:"free_throws_made,omitempty"`         //TODO unmarshal to int and add to event-stats?
	FreeThrowsAttempted    FlexInt `json:"free_throws_attempted,omitempty"`    //TODO unmarshal to int and add to event-stats?
	ThreePointersMade      FlexInt `json:"three_pointers_made,omitempty"`      //TODO unmarshal to int and add to event-stats?
	ThreePointersAttempted FlexInt `json:"three_pointers_attempted,omitempty"` //TODO unmarshal to int and add to event-stats?
	ReboundsOffensive      FlexInt `json:"rebounds_offsensive,omitempty"`      //TODO unmarshal to int and add to event-stats?
	ReboundsDefensive      FlexInt `json:"rebounds defensive,omitempty"`       //TODO unmarshal to int and add to event-stats?
	Assists                FlexInt `json:"assists,omitempty"`                  //TODO unmarshal to int and add to event-stats?
	Fouls                  FlexInt `json:"fouls,omitempty"`                    //TODO unmarshal to int and add to event-stats?
	Steals                 FlexInt `json:"steals,omitempty"`                   //TODO unmarshal to int and add to event-stats?
	Turnovers              FlexInt `json:"turnovers,omitempty"`                //TODO unmarshal to int and add to event-stats?
	TeamTurnovers          FlexInt `json:"team_turnovers,omitempty"`           //TODO unmarshal to int and add to event-stats?
	Blocks                 FlexInt `json:"blocks,omitempty"`                   //TODO unmarshal to int and add to event-stats?
	PlusMinus              FlexInt `json:"plus_minus,omitempty"`               //TODO unmarshal to int and add to event-stats?
	OnCourt                FlexInt `json:"on_court,omitempty"`                 //TODO unmarshal to int and add to event-stats?
	StartingPosition       string  `json:"starting_position,omitempty"`        //TODO  add to event-stats?
}

//TimePeriod ... part of game status
type TimePeriod struct {
	PeriodValue  FlexInt `json:"period_value,omitempty"`  // "period_value":"4", //TODO: convert string to int
	PeriodStatus string  `json:"period_status,omitempty"` //"period_status":"Final",
	GameStatus   FlexInt `json:"game_status,omitempty"`   //"game_status":"3",
	GameClock    string  `json:"game_clock,omitempty"`    //"game_clock":"", //TODO: determine wether this is rolling -> live feed and status?
	TotalPeriods FlexInt `json:"total_periods,omitempty"` //"total_periods":"4", //TODO string to int
	PeriodName   string  `json:"period_name,omitempty"`   //"period_name":"Qtr"},
}

//Linescore ...
type Linescore struct {
	Period []PeriodicScore `json:"period"` //{	"period_value":"4","period_name":"Q4","score":"19"}]},
}

//PeriodicScore ...
type PeriodicScore struct {
	PeriodValue string  `json:"period_value"`
	PeriodName  string  `json:"period_name"`
	Score       FlexInt `json:"score"` //TODO string to int
}

//TeamStatLeader ... gamestats TODO: need to put this in namespace, and realign statistics
type TeamStatLeader struct {
	Points   TeamStatistic `json:"Points,omitempty"`
	Assists  TeamStatistic `json:"Assists,omitempty"`
	Rebounds TeamStatistic `json:"Rebounds,omitempty"`
}

//TeamStatistic ... TODO: need to put this in namespace, and realign statistics
type TeamStatistic struct {
	PlayerCount FlexInt  `json:"PlayerCount"` // "PlayerCount":"1", //TODO string to int (derived from sizeof array?)
	StatValue   string   `json:"StatValue"`   // "StatValue":"32", // TODO string to int/float?
	Players     []Person `json:"leader"`      //	"leader":[{"PersonID":"200768","PlayerCode":"kyle_lowry","FirstName":"Kyle","LastName":"Lowry"}],
}

//TeamStats ...
type TeamStats struct {
	Points                  FlexInt `json:"points,omitempty"`                    //TODO unmarshal to int and add to event-stats?
	FieldGoalsMade          FlexInt `json:"field_goals_made,omitempty"`          //TODO unmarshal to int and add to event-stats?
	FieldGoalsAttempted     FlexInt `json:"field_goals_attempted,omitempty"`     //TODO unmarshal to int and add to event-stats?
	FieldGoalsPercentage    string  `json:"field_goals_percentage,omitempty"`    //TODO unmarshal to float32
	FreeThrowsMade          FlexInt `json:"free_throws_made,omitempty"`          //TODO unmarshal to int and add to event-stats?
	FreeThrowsAttempted     FlexInt `json:"free_throws_attempted,omitempty"`     //TODO unmarshal to int and add to event-stats?
	FreeThrowsPercentage    string  `json:"free_throws_percentage,omitempty"`    //TODO unmarshal to float32
	ThreePointersMade       FlexInt `json:"three_pointers_made,omitempty"`       //TODO unmarshal to int and add to event-stats?
	ThreePointersAttempted  FlexInt `json:"three_pointers_attempted,omitempty"`  //TODO unmarshal to int and add to event-stats?
	ThreePointersPercentage string  `json:"three_pointers_percentage,omitempty"` //TODO unmarshal to float32
	ReboundsOffensive       FlexInt `json:"rebounds_offsensive,omitempty"`       //TODO unmarshal to int and add to event-stats?
	ReboundsDefensive       FlexInt `json:"rebounds defensive,omitempty"`        //TODO unmarshal to int and add to event-stats?
	TeamRebounds            FlexInt `json:"team_rebounds,omitempty"`             // "team_rebounds":"15", //TODO unmarshal to int
	Assists                 FlexInt `json:"assists,omitempty"`                   //TODO unmarshal to int and add to event-stats?
	Fouls                   FlexInt `json:"fouls,omitempty"`                     //TODO unmarshal to int and add to event-stats?
	TeamFouls               FlexInt `json:"team_fouls"`                          //"team_fouls":"10", //TODO unmarshal to int
	TechnicalFouls          FlexInt `json:"technical_fouls"`                     //"technical_fouls":"1", //TODO unmarshal to int
	Steals                  FlexInt `json:"steals,omitempty"`                    //TODO unmarshal to int and add to event-stats?
	Turnovers               FlexInt `json:"turnovers,omitempty"`                 //TODO unmarshal to int and add to event-stats?
	TeamTurnovers           FlexInt `json:"team_turnovers,omitempty"`            //TODO unmarshal to int and add to event-stats?
	Blocks                  FlexInt `json:"blocks,omitempty"`                    //TODO unmarshal to int and add to event-stats?
	ShortTimeoutRemaining   FlexInt `json:"short_timeout_remaining"`             //"short_timeout_remaining":"0", //TODO unmarshal to int
	FullTimeoutRemaining    FlexInt `json:"full_timeout_remaining"`              //"full_timeout_remaining":"0"}, //TODO unmarshal to int
}

//WorkingTeam ... visitor/home team info with status...
type WorkingTeam struct {
	TeamID       string         `json:"id"`                     //"id":"1610612761", // convert string to int?
	TeamKey      string         `json:"team_key"`               //"team_key":"TOR",
	City         string         `json:"city"`                   //"city":"Toronto",
	Abbreviation string         `json:"abbreviation,omitempty"` //"abbreviation":"TOR",
	Nickname     string         `json:"nickname,omitempty"`     //"nickname":"Raptors",
	URLName      string         `json:"url_name"`               //"url_name":"raptors",
	TeamCode     string         `json:"team_code"`              //"team_code":"raptors",
	Score        FlexInt        `json:"score"`                  //"score":"104", //TODO string to int
	Linescores   Linescore      `json:"linescores"`
	Leaders      TeamStatLeader `json:"Leaders"`
	TeamStats    TeamStats      `json:"stats"`
	Players      PlayerArray    `json:"players"` //TODO... properly parse this... don't need holding structure
}

//PlayerArray ... TODO... don't need this extra structure associated with team
type PlayerArray struct {
	Player []Person `json:"player"`
}

//FlexInt ... int unmarshalled fro JSON field that is passed as a string, or
// inconsistently (string or int)
type FlexInt int

// UnmarshalJSON implements the json.Unmarshaler interface, which
// allows us to ingest values of any json type as an int and run our custom conversion
func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	if b[0] != '"' { // already an int
		return json.Unmarshal(b, (*int)(fi))
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*fi = FlexInt(-1)
	} else {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("Error strconv.Atoi FlexInt convert, %s", s)
			return err
		}
		*fi = FlexInt(i)
	}
	return nil
}
