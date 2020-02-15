package espn

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"go-moneyball/moneyball/ms"

	"github.com/stretchr/testify/assert"
)

const tFilename = "../../examples/json/espnTeams123019-0931.json"

func TestMarshalMSTeam(t *testing.T) {
	b, err := ioutil.ReadFile(tFilename)
	if err != nil {
		log.Printf("ERROR thrown: %s\n", err)
	}
	assert.Nil(t, err, fmt.Errorf("couldn't read file: %s", tFilename))
	ts := TeamSport{}
	err = json.Unmarshal(b, &ts)
	fmt.Printf("bytes: %s",b)
	assert.Nil(t, err, fmt.Errorf("error decoding json file: %s, %#v", tFilename, err))
	mTeam := []*ms.Team{}
	for _, team := range ts.Sport[0].Leagues[0].Teams {
		tm, err := marshalMSTeam(&team.Team)
		assert.Nil(t, err, fmt.Errorf("error %#v: translating %#v, to %#v", err, team, tm))
		mTeam = append(mTeam, tm)
	}
	//should have a certain number of teams
	log.Printf("espnTeam->MSTeam worked\n")
}
