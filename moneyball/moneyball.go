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

import (
	"context"
	"fmt"
	"moneyball/go-moneyball/moneyball/espn"
	"time"

	"github.com/davecgh/go-spew/spew"
	//"golang.org/x/oauth2"
)

/**
 * ability to test OAUTH2 access for sso
 */
func main() {
	/** uncomment for oauth2
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "... your access token ..."},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := NewClient(tc)
	*/
	// if using oath2 comment rest
	client := NewClient(nil)
	ctx := context.Background()

	schedParams := map[string]string{
		"year": "2019", //2019 season (current)
	}
	schedule, _, err := client.Schedule.NBAScheduleServicev2(ctx, schedParams)
	if err != nil {
		fmt.Printf("ScheduleService: Error %s\n", err)
	}
	fmt.Printf("NBAScheduleService: %d with values %#v retrieved\n", len(*schedule), (*schedule)[0])
	//getTodayGames(schedule)
	todayStart := time.Now()
	tomorrow := todayStart.AddDate(0, 0, 1)
	counter := 0
	// look for games today... get box scores...
	for i, game := range *schedule {
		if game.StartTime.After(todayStart) && game.StartTime.Before(tomorrow) {
			//have a game I care about
			counter = counter + 1
			fmt.Printf("today game id: %s, start: %s %s, url: %s\n", game.GameID, game.StartTimeEastern, game.StartDateEastern, game.GameURLCode)
			// go get details.
			params := map[string]string{
				"gamedate": game.StartDateEastern,
				"gameid":   game.GameID}
			temp, _, err := client.Score.NBABoxScoreServicev2(ctx, params)
			if err != nil {
				fmt.Printf("BoxScoreService: Error %s\n", err)
			} else {
				// replace existing game with the detailed box.
				fmt.Printf("orig_game %s", game.GameURLCode)
				(*schedule)[i] = *temp
				fmt.Printf("new game %#v", (*schedule)[i])
				// could build independent array of games or add detail or...

			}
			fmt.Printf("next\n")
		}
	}
	//todo: getYesterdayBoxes(schedule)

	// test script for old BoxScore Service using old NBA API's - this is kinda a mess
	boxscore, _, err := client.Score.BoxScoreService(ctx)
	if err != nil {
		fmt.Printf("BoxScoreService: Error %s\n", err)
	}
	fmt.Printf("NBABoxScoreService: %s with values %#v retrieved\n", "test", boxscore) //boxscore.Event.Game.GameID, boxscore)

	// tests for PlayerMovement service from nba... this is used to show player roster changes (but seems to be non-authoritative)
	statstln, _, err := client.Stats.NBAPlayerMovementStatsService(ctx)
	if err != nil {
		fmt.Printf("ScoreBoardService: Error: %s\n", err)
	}
	fmt.Printf("NBAPlayerMovementStatsService: %s StatName with values of %#v retrieved\n", statstln.StatGroupName, statstln.StatGroup)

	//get the current scoreboard
	scoreboard, _, err := client.Score.ESPNBoxScoreService(ctx)
	if err != nil {
		fmt.Printf("ScoreBoardService: Error: %s\n", err)
	}
	sb, _ := espn.MarshalMS(scoreboard)
	spew.Printf("espn.ScoreBoard: %v\n\n", scoreboard)
	spew.Printf("ms.ScoreBoard: %v\n\n", sb)

	teams, _, err := client.Stats.ESPNTeamsService(ctx)
	if err != nil {
		fmt.Printf("TeamsService: Error %s\n", err)
	}
	_ = teams
	//fmt.Printf("TeamsService: %d teams for date retrieved\n", len(teams.Sport[0].Leagues[0].Teams))

}
