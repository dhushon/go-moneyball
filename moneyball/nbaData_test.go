package main

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
	statstln, _, err := client.Stats.PlayerMovementStatsService(ctx)
	assert.Nil(t, err, err)
	assert.NotZero(t, len(statstln.StatGroup) > 0, "StatGroup should not be nil")
	//fmt.Printf("NBAPlayerMovementStatsService: %s StatName with values of %#v retrieved\n", statstln.StatGroupName, statstln.StatGroup)
}
