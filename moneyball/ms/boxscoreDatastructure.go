package ms

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

import (
	"time"
)

//GameID ...
type GameID string

//League ...
type League string

//Season ...
type Season struct {
	SeasonYear  int `json:"seasonYear,omitempty"`
	SeasonStage int `json:"seasonStageId,omitempty"`
}

//Competitor ...
type Competitor struct {
	EntityID
	Name           string   `json:"name,omitempty"`
	Abbreviation   string   `json:"abbreviation"`
	Record         Record   `json:"record,omitempty"`
	LineScore      *[]Score `json:"linescore,omitempty"` //"linescore":[{"score":"30"},{"score":"32"},{"score":"23"},{"score":"19"}]},
	Location       string   `json:"location"`
	Color          string   `json:"color"`
	AlternateColor string   `json:"alternateColor"`
	IsActive       bool     `json:"isActive"`
	IsAllStar      bool     `json:"isAllStar"`
	Logo           *Logo    `json:"logos"`
}

//Logo is an image...
type Logo struct {
	Image
}

//Image is a link that has dimensions
type Image struct {
	Link
	Dimension LinkDimensions `json:"dimensions"`
}

//Score ... used in linescore to show period score for a team/competitor
type Score struct {
	Score int `json:"score,omitempty"`
}

//Record ... win/loss record for team
type Record struct {
	Win   int    `json:"win"`
	Loss  int    `json:"loss"`
	Items []Item `json:"items"`
}

//Item is a stats element that includes a summary plus a name/value pair
type Item struct {
	Summary string `json:"summary"`
	Stats   []struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	} `json:"stats"`
}

//Venue is the data around the sports venu
type Venue struct {
	EntityID
	LocalID  string   `json:"id"`
	FullName string   `json:"fullName,omitempty"`
	Address  *Address `json:"address,omitempty"`
	Capacity int      `json:"capacity"`
	IsIndoor bool     `json:"indoor"`
}

//Address is the street address of the venue
type Address struct {
	Street  string `json:"street,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
	GeoLoc  string `json:"geoloc,omitempty"`
}

//GameStatus ...
type GameStatus string

//GameScore ...
type GameScore struct {
	HomeScore  int `json:"home"`
	VisitScore int `json:"away"`
}

//Link ...
type Link struct {
	HRef string   `json:"href"`          //"http://www.espn.com/nba/team/_/name/tor/toronto-raptors",
	Rel  []string `json:"rel,omitempty"` // ["clubhouse","desktop","team"],
	Alt  string   `json:"alt,omitempty"` // "Clubhouse"
}

//LinkDimensions ...
type LinkDimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// EntityID provides the Monumental Foreign key resolution for key types, like Games, Players, Teams that help to resolve
// across a variety of source API's and data bases
type EntityID struct {
	ID string `json:"id"`
}

//BoxScore ...
type BoxScore struct {
	EntityID
	GameID     GameID      `json:"gameId"`
	League     League      `json:"league"`
	Season     Season      `json:"season"`
	HomeTeam   *Competitor `json:"homeTeam"`
	VisitTeam  *Competitor `json:"visitTeam"`
	Venue      *Venue      `json:"location,omitempty"`
	Status     GameStatus  `json:"status,omitempty"`
	Score      *GameScore  `json:"gamescore,omitempty"`
	Links      *[]Link     `json:"link,omitempty"`
	GameDetail *GameDetail `json:"gameDetail,omitempty"`
}

//GameDetail .. extra detail about the game including things like startTime...
type GameDetail struct {
	StartTime           *time.Time  `json:"startTimeUTC,omitempty"`     //"startTimeUTC":"2019-10-01T00:00:00.000Z",
	StartDateEastern    string      `json:"startDateEastern,omitempty"` //"startDateEastern":"20190930",
	StartTimeEastern    string      `json:"startTimeEastern,omitempty"`
	Period              *GamePeriod `json:"period,omitempty,omitempty"` // "period": {}
	Attendance          string      `json:"attendance,omitempty"`       //"attendance":"18624",
	GameDurationMinutes int         `json:"gameDuration,omitempty"`
}

//GamePeriod provides a structure that holds information about the period/quarter/half... that can be used to show game progession
type GamePeriod struct {
	Current       int  `json:"current"`       //"current":4,
	Type          int  `json:"type"`          //"type":0,
	MaxRegular    int  `json:"maxRegular"`    //"maxRegular":4
	IsHalftime    bool `json:"isHalftime"`    //`"isHalftime":false,
	IsEndOfPeriod bool `json:"isEndOfPeriod"` //"isEndOfPeriod":false
}

//Scoreboard ... holding structure for a set of BoxScores
type Scoreboard struct {
	BoxScores []BoxScore
}

// MasterIdentity will provide a basic "soure->target" mapping of different data sets against a
// set of common table keys... things like events, players, and even locations need to be mastered
func asMasterIdentity(v interface{}) string {
	// test if interface isA EntityID struct
	return ""
}
