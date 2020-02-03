package main

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
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestESPNScoreBoardService(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()
	//get the current scoreboard
	scoreboard, _, err := client.Score.ESPNBoxScoreService(ctx)
	assert.Nil(t, err, err)
	//fmt.Printf("ScoreBoardService: %d scores for date %s retrieved\n", len(scoreboard.Events), scoreboard.Day.Date)
	assert.NotZero(t, len(scoreboard.Events), t.Name()+" espnScoreboard positive length response")
	//fmt.Printf("Response: %#v\n", scoreboard)
}

func TestESPNScoreBoardMarshalService(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()
	//get the current scoreboard
	scoreboard, _, err := client.Score.ESPNBoxScoreService(ctx)
	assert.Nil(t, err, err)
	sb, err := scoreboard.MarshalMS()
	spew.Printf("scoreboard: %#v \n ms.scoreboard: %#+v\n", scoreboard, sb)
	fmt.Printf("what'd we get %#v", sb)

}

func TestESPNTeamService(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()
	teams, _, err := client.Schedule.client.Stats.ESPNTeamsService(ctx)
	assert.Nil(t, err, err)
	assert.NotZero(t, len(teams.Sport[0].Leagues[0].Teams) > 0, "espnTeams should be a positive length response")
	//fmt.Printf("TeamsService: %d teams for date retrieved %#v\n", len(teams.Sport[0].Leagues[0].Teams), teams.Sport[0].Leagues[0].Teams)
}
