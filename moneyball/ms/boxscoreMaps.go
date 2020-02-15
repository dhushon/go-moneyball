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
	"context"
	"log"
	"os"

	"googlemaps.github.io/maps"
)

func getMapCreds() maps.ClientOption {
	key, exists := os.LookupEnv("GOOGLE_APP_API")
	if !(exists) {
		return nil
	}
	return maps.WithAPIKey(key)
}

func (v *Venue) toString() string {
	//TODO: get venue address
	a := v.Address
	str := v.FullName + "," + a.Street + ", " + a.City + ", " + a.State + ", " + a.Country
	return str
}

// GetGeoCodeAddress ... gets a geolocplaceID/pluscode for an address
func GetGeoCodeAddress(v *Venue) (string, error) {
	c, err := maps.NewClient(getMapCreds())
	if err != nil {
		log.Fatalf("fatal error ensure that Google Mapping API credential were provided: %s", err)
		return "", err
	}
	r := &maps.GeocodingRequest{Address: v.toString()}
	resp, err := c.Geocode(context.Background(), r)

	if len(resp) != 1 {
		log.Printf("Expected length of response is 1, was %+v", len(resp))
		return "", err
	}
	if err != nil {
		log.Printf("r.Get returned non nil error: %v", err)
		return "", err
	}
	log.Printf("geoloc: %#v", resp)
	// resp[0].PlaceID = provides a PlaceID
	// resp[0].PlusCode.GlobalCode = provides a new Global Code
	return resp[0].PlusCode.GlobalCode, err
}
