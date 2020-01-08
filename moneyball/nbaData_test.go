package main

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, testName string, a interface{}, b interface{}) {
	t.Logf("testing: %s", testName)
	fmt.Printf("testing: %s =>", testName)
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func AssertTrue(t *testing.T, testName string, a bool) {
	t.Logf("testing: %s", testName)
	fmt.Printf("testing: %s =>", testName)
	if !a {
		t.Errorf("Received %v (type %v), expected true", a, reflect.TypeOf(a))
	}
}

func TestPathModifier(t *testing.T) {
	// test both correct and error modifier
	testString := []string{"http://google.com/{param1}/{param2}", "http://google.com/{param1}/{param3}"}
	modifier := map[string]string{
		"param1": "12435",
		"param2": "67890",
	}
	// test positive case - both found
	_, err := nbaPathModifier(testString[0], modifier)
	AssertTrue(t, t.Name()+" positive effect", (err == nil))
	_, err = nbaPathModifier(testString[1], modifier)
	AssertTrue(t, t.Name()+" error correct effect", (err != nil))
}

func TestNBAScheduleServicev2(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	schedParams := map[string]string{
		"year": "2019", //2019 season (current)
	}
	schedule, _, err := client.Schedule.NBAScheduleServicev2(ctx, schedParams)
	AssertTrue(t, t.Name()+" no error returned", (err == nil))
	AssertTrue(t, t.Name()+" returns schedule", len(*schedule) > 0)
	fmt.Printf("NBAScheduleService: %d with values %#v retrieved\n", len(*schedule), (*schedule)[0])
}

func TestPlayerMovementStatsService(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()
	// tests for PlayerMovement service from nba... this is used to show player roster changes (but seems to be non-authoritative)
	statstln, _, err := client.Stats.PlayerMovementStatsService(ctx)
	AssertTrue(t, t.Name()+" no error returned", (err == nil))
	AssertTrue(t, t.Name()+" returns schedule", len(statstln.StatGroup) > 0)
	fmt.Printf("NBAPlayerMovementStatsService: %s StatName with values of %#v retrieved\n", statstln.StatGroupName, statstln.StatGroup)
}