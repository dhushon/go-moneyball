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

// an Event has competitors (home & away),
// a competitor is a team at an event
// a team has a roster of players
// a competitor has a roster of players for that event that play for the team

import (
	"time"
)

//Team ...
type Team struct {
	EntityID
	TeamIDNBA    string `json:"teamIdNBA,omitempty"` //"teamId":"1610612745"
	TeamIDESPN   string `json:"teamIdESPN,omitempty"`
	Abbreviation string `json:"abbreviation"`
	Name         string `json:"name"`
	Location     string  `json:"teamLocation"` // e.g. "Atlanta" Hawks
	Logos        []*Link `json:"logos,omitempty"`
	Links 		 []*Link `json:"links,omitempty"`
	//TODO: how to treat historic record?
	Records []*TeamSeasonRecords `json:"records,omitempty"`
	Rosters []*TeamSeasonRoster  `json:"rosters,omitempty"` // roster is copied to Competitor for a given game
}

//TeamSeasonRecords ...
type TeamSeasonRecords struct {
	Season  *Season `json:"season,omitempty"`
	Summary string  `json:"summary,omitempty"`
	Stats   []*Stat `json:"teamStat,omitempty"`
}

//GamePlayerStats ... may be used in boxScores to do game stats associated with Player
type GamePlayerStats struct {
	EntityID `json:"gamePlayerID"`
	Stats    []*Stat `json:"gamePlayerStat"`
}

//Stat .. a well known stat both short/long verions if exists
type Stat struct {
	Key     string      `json:"key"`
	LongKey string      `json:"longKey"`
	Value   interface{} `json:"value"`
}

//TeamSeasonRoster ...
type TeamSeasonRoster struct {
	Season Season    `json:"season"`
	Roster []*Player `json:"roster"`
}

//Player ...
type Player struct {
	EntityID    EntityID
	IDESPN      string             `json:"idESPN,omitempty"` // e.g. "id":"3012",
	IDNBA       string             `json:"idNBA,omitempty"`
	FullName    string             `json:"fullName,omitempty"`    // e.g. "fullName":"Kyle Lowry",
	DisplayName string             `json:"displayName,omitempty"` // e.g. "displayName":"Kyle Lowry",
	ShortName   string             `json:"shortName,omitempty"`   // e.g."K. Lowry",
	Links       []Link             `json:"links"`
	Jersey      string             `json:"jersey,omitempty"` // e.g. "jersey":"7",
	Headshot    *Link              `json:"headshot"`         // e.g. "headshot":"https://a.espncdn.com/i/headshots/nba/players/full/3012.png",
	Position    *Position          `json:"position,omitempty"`
	Team        *Team              `json:"team" binding:"required"`
	Active      bool               `json:"active"`
	Career      *PlayerTeamsCareer `json:"career,omitempty"`
}

//PlayerAssignment is a record in the history of  a player inclusive of volunteer,
//highchool, college, national, olympic, international, club...
type PlayerAssignment struct {
	Type         string    `json:"leagueType"`
	DateStart    time.Time `json:"start"`
	DateEnd      time.Time `json:"end"`
	Team         *EntityID `json:"teamID"`
	Position     Position  `json:"position"`
	Injured      bool      `json:"wasInjured"`
	Simultaneous bool      `json:"simultaneoud"` // there are times in which multiple teams simultaneously in events like olympics or international demo
}

//PlayerTeamsCareer containst the longitudinal record of a player, it cannot be assumed that the list will be ordered
// so it is highly recommended that the json is sorted on start or end dates to d a start->end->start motion
type PlayerTeamsCareer struct {
	Program []*PlayerAssignment `json:"program,omitempty"`
}

//Position ... somthing like C[enter], P[oint]G[uard]...
type Position struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}
