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
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const ()

var timeN = time.Now()
var team1 = Team{EntityID{"WAS", &timeN, "sampleData"}, "1610612745", "", "WAS", "Washington Wizards", "Washintgon", nil, nil, nil, nil}
var team2 = Team{EntityID{"DET", &timeN, "sampleData"}, "", "26", "DET", "Detroit Pistons", "Detroit", nil, nil, nil, nil}
var team3 = Team{EntityID{"TOR", &timeN, "sampleData"}, "", "25", "TOR", "Toronto Raptors", "Toronto", nil, nil, nil, nil}
var team4 = Team{EntityID{"BOS", &timeN, "sampleData"}, "", "17", "BOS", "Boston Celtics", "Boston", nil, nil, nil, nil}

//type Team struct {
//Records []*TeamSeasonRecords `json:"records"`
//Rosters []*TeamSeasonRoster  `json:"rosters"` // roster is copied to Competitor for a given game}
var bss = ScoreBoard{
	[]Event{
		Event{EntityID{"2019-12-28.WSH.DET", nil, ""}, "2019-12-28.WSH.DET", "NBA", Season{2019, 1},
			&Competitor{EntityID{"DET-NBA-2019", nil, ""}, "Detroit Pistons", "DET", &team2, Record{0, 1, []Item{}}, 0, &[]Score{}, "Detroit", "0x0000", "0xffff", true, false, nil},
			&Competitor{EntityID{"WAS-NBA-2019", nil, ""}, "Washington Wizards", "WAS", &team1, Record{1, 0, []Item{}}, 0, &[]Score{}, "Washington", "0E3764", "e31837", true, false, nil},
			&Venue{EntityID{}, "", "Little Caesars Arena", &Address{}, 10000, true},
			&GameStatus{0.0, 0, "Final", "Thu, December 28th at 7:00 PM EST"},
			&[]Link{
				Link{"http://www.espn.com/nba/team/roster/_/name/det/detroit-pistons",
					[]string{"roster"}, "roster", nil, false},
			},
			&GameDetail{},
		},
		Event{EntityID{"2017-02-03.TOR.BOS", nil, ""}, "2017-02-03.TOR.BOS", "NBA", Season{2017, 1},
			&Competitor{EntityID{"TOR-NBA-2017", nil, ""}, "Toronto Raptors", "TOR", &team3, Record{1, 0, []Item{}}, 109, &[]Score{}, "Toronto", "0x0000", "0xffff", true, false, nil},
			&Competitor{EntityID{"BOS-NBA-2017", nil, ""}, "Boston Celtics", "BOS", &team4, Record{0, 1, []Item{}}, 104, &[]Score{}, "Boston", "0x0000", "0xffff", true, false, nil},
			&Venue{EntityID{}, "", "TD Garden", &Address{}, 10000, true},
			&GameStatus{0.0, 0, "Final", "Thu, February 3rd at 7:00 PM EST"},
			&[]Link{},
			&GameDetail{},
		},
	},
}

var ven = Venue{EntityID{}, "", "Little Caesars Arena", &Address{Street: "2645 Woodward Avenue", City: "Detroit", State: "MI", Country: "US", GeoLoc: ""}, 20332, true}

//Testing routine for entity mastering - need to figure out how to cary
//google credentials else it will always fail.
func TestMasterEntity(t *testing.T) {
	const GeoCode = "86JR8WRV+CV" // real GeoCode for sample venue
	s, err := ven.MasterEntity()
	assert.Nil(t, err, "Venue.MasterEntity returned error", err)
	assert.EqualValues(t, s, GeoCode, "Venue.MasterEntity returned wrong geocode")
}

//NDJSONService ... test the NDJSON Encoding of a struct
func TestMarshalNBJSON(t *testing.T) {
	//marshall to ndjson so that we can push to/towards bigquery
	var b bytes.Buffer // for testing/development we can use bytes.Buffer as writer
	//marshall with json newline
	bss.marshalNBJSON(&b)
	//log.Printf("ndjson Test result %s\n",b.String())
	assert.NotZero(t, b.Len, "ndjson should not be length 0 %s", b)
	//we could try and re-read the JSON as another assert.. but too hard? as must have an Unmarshaller that supports
	//NBJSON vs. JSON... below throws error
	//var e = []Event{}
	//err := json.Unmarshal(b.Bytes(),&e)
	//assert.Nil(t, err, "Unmarshalling of NBJSON failed", err)
	log.Println("exiting")
}

func TestEntityIDExtract(t *testing.T) {

}
