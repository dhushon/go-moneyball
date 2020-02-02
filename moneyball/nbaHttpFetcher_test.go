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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathModifier(t *testing.T) {
	// test both correct and error modifier
	testString := []string{"http://google.com/{param1}/{param2}", "http://google.com/{param1}/{param3}"}
	modifier := map[string]string{
		"param1": "12435",
		"param2": "67890",
	}
	// test positive case - both found
	_, err := nbaPathModifier(testString[0], modifier)
	assert.Nil(t, err, err)
	_, err = nbaPathModifier(testString[1], modifier)
	assert.NotNil(t, err, err)
}

func TestNBAScheduleServicev2(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	schedParams := map[string]string{
		"year": "2019", //2019 season (current)
	}
	schedule, _, err := client.Schedule.NBAScheduleServicev2(ctx, schedParams)
	assert.Nil(t, err, err)
	assert.NotZero(t, len(*schedule) > 0)
	//fmt.Printf("NBAScheduleService: %d with values %#v retrieved\n", len(*schedule), (*schedule)[0])
}

func TestPlayerMovementStatsService(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()
	// tests for PlayerMovement service from nba... this is used to show player roster changes (but seems to be non-authoritative)
	statstln, _, err := client.Stats.NBAPlayerMovementStatsService(ctx)
	assert.Nil(t, err, err)
	assert.NotZero(t, len(statstln.StatGroup) > 0, "StatGroup should not be nil")
	//fmt.Printf("NBAPlayerMovementStatsService: %s StatName with values of %#v retrieved\n", statstln.StatGroupName, statstln.StatGroup)
}