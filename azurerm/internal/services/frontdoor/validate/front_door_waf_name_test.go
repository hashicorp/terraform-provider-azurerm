package validate

import (
	"testing"
)

func TestAccFrontDoorFirewallPolicy_validateName(t *testing.T) {
	cases := []struct {
		Name        string
		Input       string
		ExpectError bool
	}{
		{
			Name:        "Empty String",
			Input:       "",
			ExpectError: true,
		},
		{
			Name:        "Starts with Numeric",
			Input:       "1WellThisIsAllWrong",
			ExpectError: true,
		},
		{
			Name:        "Has Spaces",
			Input:       "What part of no spaces do you not understand",
			ExpectError: true,
		},
		{
			Name:        "Has Hyphens",
			Input:       "What-part-of-no-hyphens-do-you-not-understand",
			ExpectError: true,
		},
		{
			Name:        "Special Characters",
			Input:       "WellArn`tTheseSpecialCharacters?!",
			ExpectError: true,
		},
		{
			Name:        "Mixed Case Alpha and Numeric",
			Input:       "ThisNameIsAPerfectlyFine1",
			ExpectError: false,
		},
		{
			Name:        "Too Long",
			Input:       "OhMyLordyThisNameIsWayToLooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooog",
			ExpectError: true,
		},
		{
			Name:        "Max Length",
			Input:       "NowThisNameIsThePerfectLengthForAFrontdoorFireWallPolicyDontYouThinkAnyLongerWouldBeJustWayToLoooooooooooooooooooongDontYouThink",
			ExpectError: false,
		},
		{
			Name:        "Minimum Length Upper",
			Input:       "A",
			ExpectError: false,
		},
		{
			Name:        "Minimum Length Lower",
			Input:       "a",
			ExpectError: false,
		},
		{
			Name:        "Mixed Case Alpha no Numeric",
			Input:       "LookMomNoNumbers",
			ExpectError: false,
		},
		{
			Name:        "All Upper Alpha with Numeric",
			Input:       "OU812",
			ExpectError: false,
		},
		{
			Name:        "All lower no Numeric",
			Input:       "heythisisalllowercase",
			ExpectError: false,
		},
	}

	for _, tc := range cases {
		_, errors := FrontDoorWAFName(tc.Input, tc.Name)

		hasError := len(errors) > 0

		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the FrontDoor WAF Name to trigger a validation error for '%s'", tc.Name)
		}
	}
}
