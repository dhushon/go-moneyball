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

// scaffolding: https://medium.com/@marcus.olsson/writing-a-go-client-for-your-restful-api-c193a2f4998c
// which led to https://github.com/google/go-github which is really quite like our fetch/version process

//navigation link with php source: http://nbasense.com/nba-api/Stats/Data/
//http://data.nba.net/json/bios/player_201935.json

//http://data.nba.net/json/cms/noseason/game/{gameDate}/{gameId}/boxscore.json
//http://data.nba.net/json/cms/noseason/game/20170201/0021600732/boxscore.json where game/`json:"date"`/`json:"ID"`/boxscore.json`
//REF: http://nbasense.com/nba-api/Data/Cms/Game/Boxscore

import (
	"context"
	"fmt"
	"log"
	"go-moneyball/moneyball/nba"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func scheduleServiceURLModifier(modifier map[string]string) (string, error) {
	//http://data.nba.net/json/cms/2016/league/nba_games.json
	suffix := ""
	for param := range modifier {
		switch param {
		case "period", "Period", "PERIOD":
			//legal values are years, dates [yyyymmdd], today or week
			p := modifier[param]
			if len(p) > 4 {
				// go ahead and parse extended date
				suffix = "2018/" //TODO parse extended date
				break
			} else if p == "ALL" {
				break
			} else {
				if strings.HasPrefix(p, "20") || strings.HasPrefix(p, "19") {
					// test to see if can convert to int... if yes
					if _, err := strconv.Atoi(p); err == nil {
						suffix = suffix + p + "/"
						break
					}
					//return "", err
					break
				}
				break
			}
		case "team", "Team", "TEAM":
			//see if we have an ID
			//see if we have an abbreviation
			//if all just zero out
			//else ignore or trigger an error?
			break
		default:
			//parameter unknown trigger an error?
		}
	}
	return suffix, nil
}

/*NBAScheduleService ...
  NBA schedule for [year,Today,Week] for Team [Team-UID]
  - this is done upstream for the different league services? ?league = ["NBA","WNBA"] absent defaults to NBA
  ?period = [$yyyy, $yyyymmdd, "Today","Week"] absent defaults to Today
  ?team = [$teamID or $teamAbbr.] absent returns all teams
*/
func (s *ScheduleService) NBAScheduleService(ctx context.Context, modifier map[string]string) (*[]nba.ScheduledGame, *Response, error) {
	s.client.BaseURL, _ = url.Parse(nba.DataNBABaseURL)
	//http://data.nba.net/json/cms/2016/league/nba_games.json
	suffix, err := scheduleServiceURLModifier(modifier)
	//+"league/nba_games.json"
	if !(strings.HasSuffix(suffix, "/")) {
		suffix = suffix + "/"
	}
	req, err := s.client.NewRequest("GET", "json/cms/"+suffix+"league/nba_games.json", nil)
	if err != nil {
		return nil, nil, err
	}
	event := &nba.SportsEvent{}
	resp, err := s.client.Do(ctx, req, event, true)
	if err != nil {
		log.Printf("Error caught: %s\n", err)
	}
	return nil, resp, err
	//return &event.Event.Schedule.Games, resp, err
}

//BoxScoreService will, for a http client, return a StatsTLN JSON object ( note that this is not yet normalized to structures)
func (s *ScoreService) BoxScoreService(ctx context.Context) (*nba.SportsEvent, *Response, error) {

	s.client.BaseURL, _ = url.Parse(nba.DataNBABaseURL)
	req, err := s.client.NewRequest("GET", "json/cms/noseason/game/20170201/0021600732/boxscore.json", nil)
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
	event := &nba.SportsEvent{}
	resp, err := s.client.Do(ctx, req, event, true)
	if err != nil {
		log.Printf("Error on new request: %s\n", err)
		return nil, resp, err
	}
	//TODO... extract meta and bring backscore of game
	return event, resp, err
}

//NBAPlayerMovementStatsService will, for a http client, return a StatsTLN JSON object ( note that this is not yet normalized to structures)
func (s *StatsService) NBAPlayerMovementStatsService(ctx context.Context) (*nba.StatsTLN, *Response, error) {

	s.client.BaseURL, _ = url.Parse(nba.NBAStatsBaseURL)
	req, err := s.client.NewRequest("GET", nba.NBAStatsURLPrefix+nba.PlayerMovementPath, nil)

	//to support gzip encoding uncomment... should probably default to true
	//req.Header.Add("Accept-Encoding", "gzip")

	// get useragent from OS Environment Variables -> often needed to prevent robot blocking or API access with lower DoS thresholds
	agent, exists := os.LookupEnv("NBA_USERAGENT")
	if exists {
		req.Header.Set("User-Agent", agent)
	}
	tln := &nba.StatsTLN{}
	resp, err := s.client.Do(ctx, req, tln, true)
	if err != nil {
		log.Printf("Error on new request: %s\n", err)
		return nil, resp, err
	}
	return tln, resp, err
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
func (s *ScoreService) NBABoxScoreServicev2(ctx context.Context, modifier map[string]string) (*nba.ScheduledGamev2,
	*Response, error) {

	s.client.BaseURL, _ = url.Parse(nba.DataNBABaseURLv2)
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
	event := &nba.CMSProdv1BoxScore{}
	resp, err := s.client.Do(ctx, req, event, true)
	if err != nil {
		log.Printf("Error on new request: %s\n", err)
		return nil, resp, err
	}
	return event.Game, resp, err
}

//NBAScheduleServicev2 is an updated nba feed for NBA Schefule information
//http://data.nba.net/prod/v2/{year}/schedule.json e.g. http://data.nba.net/prod/v2/2019/schedule.json
func (s *ScheduleService) NBAScheduleServicev2(ctx context.Context, modifier map[string]string) (*[]nba.ScheduledGamev2, *Response, error) {

	s.client.BaseURL, _ = url.Parse(nba.DataNBABaseURLv2)
	path := "prod/v2/{year}/schedule.json"
	suffix, err := nbaPathModifier(path, modifier)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", suffix, nil)
	if err != nil {
		return nil, nil, err
	}

	event := &nba.CMSProdv2Schedule{}
	resp, err := s.client.Do(ctx, req, event, true)
	if err != nil {
		log.Printf("Error caught: %s\n", err)
	}
	return &event.LeagueSchedule.Events, resp, err
}
