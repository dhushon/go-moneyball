package infoservice

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

// ESPN - NBA Scores from secret API
//
// NBA
//
//Scores: https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard
//News: http://site.api.espn.com/apis/site/v2/sports/basketball/nba/news
//All Teams: http://site.api.espn.com/apis/site/v2/sports/basketball/nba/teams
//Specific Team: http://site.api.espn.com/apis/site/v2/sports/basketball/nba/teams/:team
//http://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard/:eventId
//
//WNBA
//
//Scores: http://site.api.espn.com/apis/site/v2/sports/basketball/wnba/scoreboard
//News: http://site.api.espn.com/apis/site/v2/sports/basketball/wnba/news
//All Teams: http://site.api.espn.com/apis/site/v2/sports/basketball/wnba/teams
//Specific Team: http://site.api.espn.com/apis/site/v2/sports/basketball/wnba/teams/:team
//
//Women's College Basketball
//Scores: http://site.api.espn.com/apis/site/v2/sports/basketball/womens-college-basketball/scoreboard
//News: http://site.api.espn.com/apis/site/v2/sports/basketball/womens-college-basketball/news
//All Teams: http://site.api.espn.com/apis/site/v2/sports/basketball/womens-college-basketball/teams
//Specific Team: http://site.api.espn.com/apis/site/v2/sports/basketball/womens-college-basketball/teams/:team
//
//Men's College Basketball
//Scores: http://site.api.espn.com/apis/site/v2/sports/basketball/mens-college-basketball/scoreboard
//News: http://site.api.espn.com/apis/site/v2/sports/basketball/mens-college-basketball/news
//All Teams: http://site.api.espn.com/apis/site/v2/sports/basketball/mens-college-basketball/teams
//Specific Team: http://site.api.espn.com/apis/site/v2/sports/basketball/mens-college-basketball/teams/:team
//
//ScoreBoardService will, for a http client, return a ScoreBoard JSON object
//

import (
	"context"
	"fmt"
	"moneyball/go-moneyball/moneyball/espn"
	"net/url"
	"os"
)

//BoxScoreService provides a fetcher for ESPN's scoreboard API that will pull the latest scoreboard (todays games & results)
func (s *ScoreService) ESPNBoxScoreService(ctx context.Context) (*(espn.ScoreBoard), *Response, error) {

	s.client.BaseURL, _ = url.Parse(espn.EspnBaseURL)
	req, err := s.client.NewRequest("GET", espn.EspnURLPrefix+"nba/scoreboard", nil)

	//to support gzip encoding uncomment... should probably default to true
	//req.Header.Add("Accept-Encoding", "gzip")

	// get useragent from OS Environment Variables -> often needed to prevent robot blocking or API access with lower DoS thresholds
	agent, exists := os.LookupEnv("ESPN_USERAGENT")
	if exists {
		req.Header.Set("User-Agent", agent)
	}

	sb := &espn.ScoreBoard{}
	resp, err := s.client.Do(ctx, req, sb, false)
	if err != nil {
		fmt.Printf("Error on new request: %s\n", err)
		return nil, resp, err
	}
	return sb, resp, err
}

//TeamsService will, for a http client, return a ScoreBoard JSON object
//
func (s *StatsService) ESPNTeamsService(ctx context.Context) (*espn.TeamSport, *Response, error) {

	s.client.BaseURL, _ = url.Parse(espn.EspnBaseURL)
	req, err := s.client.NewRequest("GET", espn.EspnURLPrefix+"nba/teams", nil)

	//to support gzip encoding uncomment... should probably default to true
	//req.Header.Add("Accept-Encoding", "gzip")

	// get useragent from OS Environment Variables -> often needed to prevent robot blocking or API access with lower DoS thresholds
	agent, exists := os.LookupEnv("ESPN_USERAGENT")
	if exists {
		req.Header.Set("User-Agent", agent)
	}
	teams := &espn.TeamSport{}
	resp, err := s.client.Do(ctx, req, teams, false)
	if err != nil {
		fmt.Printf("Error on new request: %s\n", err)
		return nil, resp, err
	}
	return teams, resp, err
}
