package main

import (
	"encoding/json"
	"fmt"
)

type _StatsTLN StatsTLN // preventing recursion

//StatsTLN topLevel fireld decoding
type StatsTLN struct {
	StatGroupName string                  	`json:"statGroupName"` // map of category/group name infered from structure
	StatGroup	  []StatsRow			    `json:"statGroup"` // map of string/value generices to hold json
	TLN       	  map[string]interface{} 	`json:"-"` // initial map to hold the improperly formated stats
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

func main() {
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
}
