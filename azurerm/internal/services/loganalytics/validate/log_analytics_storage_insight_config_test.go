package validate

import (
	"testing"
)

func TestLogAnalyticsStorageInsightConfigWorkspaceName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "inv",
			Expected: false,
		},
		{
			Input:    "invalid_Exports_Name",
			Expected: false,
		},
		{
			Input:    "invalid Storage Insight Config Name Name",
			Expected: false,
		},
		{
			Input:    "-invalidStorageInsightConfigName",
			Expected: false,
		},
		{
			Input:    "invalidStorageInsightConfigName-",
			Expected: false,
		},
		{
			Input:    "thisIsToLoooooooooooooooooooooongestForAStorageInsightConfigName",
			Expected: true,
		},
		{
			Input:    "validStorageInsightConfigName",
			Expected: true,
		},
		{
			Input:    "validStorageInsightConfigName-2",
			Expected: true,
		},
		{
			Input:    "thisIsTheLoooooooooooongestValidStorageInsightConfigNameThereIs",
			Expected: true,
		},
		{
			Input:    "vali",
			Expected: true,
		},
	}
	for _, v := range testCases {
		_, errors := LogAnalyticsStorageInsightConfigWorkspaceName(v.Input, "workspace_name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %q but got %q (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
