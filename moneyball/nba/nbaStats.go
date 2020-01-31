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
// PlayerMovement: https://stats.nba.com/js/data/playermovement/NBA_Player_Movement.json

import (
	"encoding/json"
)

const (
	//NBAStatsBaseURL ...
	NBAStatsBaseURL = "https://stats.nba.com/"
	//NBAStatsURLPrefix ...
	NBAStatsURLPrefix = "js/data/"
	//PlayerMovementPath ...
	PlayerMovementPath = "playermovement/NBA_Player_Movement.json"
)

/*{
"NBA_Player_Movement": {
  "rows": [
	{
	  "Transaction_Type": "Signing",
	  "TRANSACTION_DATE": "2019-12-27T00:00:00",
	  "TRANSACTION_DESCRIPTION": "Houston Rockets signed guard Chris Clemons to a Rest-of-Season Contract.",
	  "TEAM_ID": 1610612745.0,
	  "PLAYER_ID": 1629598.0,
	  "Additional_Sort": 0.0,
	  "GroupSort": "Signing 1025079"
	},
*/

// seems stats are all database driven... which means they are pretty much coming out in table structure
// so we should probably do something with a dictionary here to lookup key, value[type], and tranlated / lookup strategy
// Declared an empty interface of type Array

type _StatsTLN StatsTLN // preventing recursion

//StatsTLN topLevel fireld decoding
type StatsTLN struct {
	StatGroupName string                 `json:"statGroupName"` // map of category/group name infered from structure
	StatGroup     []StatsRow             `json:"statGroup"`     // map of string/value generices to hold json
	TLN           map[string]interface{} `json:"-"`             // initial map to hold the improperly formated stats
}

//StatsRow is a custom Map parser (pulling generic row table structure from JSON)
type StatsRow map[string]interface{}

// UnmarshalJSON -- custom json Unmarshal
func (tln *StatsTLN) UnmarshalJSON(bs []byte) (err error) {
	t := _StatsTLN{}

	// try to parse... but unlikely
	if err = json.Unmarshal(bs, &t); err == nil {
		// make sure we initiate
		*tln = StatsTLN(t)
	}

	//build a map to support navigation
	mp := make(map[string]interface{})

	//unmarshal into the map[string]interface{} generic
	if err = json.Unmarshal(bs, &mp); err == nil {
		for sgn := range mp {
			tln.StatGroupName = sgn
			// now navigate to "rows"
			rows := mp[sgn].(map[string]interface{})
			for r := range rows { // could just search for "rows" in the map[string]
				ary := rows[r].([]interface{})
				// need to cast the []interface to []map[string]interface{}
				// start by building the holding variable
				tln.StatGroup = make([]StatsRow, len(ary))
				// copy/assign the exisitng maps to the array (each slide must be treated independently)
				for i := range ary {
					tln.StatGroup[i] = StatsRow(ary[i].(map[string]interface{}))
				}
			}
		}
	}
	return err
}

//PlayerMovement ...
type PlayerMovement struct {
}

/*func main() {
	playerMovement := `
			{ "NBA_Player_Movement":
				{ "rows": [
					  {	"Transaction_Type": "Signing",
						"TRANSACTION_DATE": "2019-12-27T00:00:00",
						"TRANSACTION_DESCRIPTION": "Houston Rockets signed guard Chris Clemons to a Rest-of-Season Contract.",
						"TEAM_ID": 1610612745.0,
						"PLAYER_ID": 1629598.0,
						"Additional_Sort": 0.0,
						"GroupSort": "Signing 1025079"},
					{	"Transaction_Type": "Signing",
						"TRANSACTION_DATE": "2019-12-26T00:00:00",
						"TRANSACTION_DESCRIPTION": "Washington Wizards signed forward Johnathan Williams to a Rest-of-Season Contract.",
						"TEAM_ID": 1610612764.0,
						"PLAYER_ID": 1629140.0,
						"Additional_Sort": 0.0,
						"GroupSort": "Signing 1025040"}]
				}
			}`

	//var results map[string]interface{}
	tln := StatsTLN{}
	// try and detect the Top Level Node from NBAStats
	if err := json.Unmarshal([]byte(playerMovement), &tln); err != nil {
		panic(err)
	}
	fmt.Printf("TLN: %#v \n", tln)
}*/
