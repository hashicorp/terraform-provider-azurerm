package validate

import (
	"testing"
)

func TestLogAnalyticsDataExportName(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected bool
	}{
		{
			Name:     "Too short",
			Input:    "inv",
			Expected: false,
		},
		{
			Name:     "Invalid characters underscores",
			Input:    "invalid_Exports_Name",
			Expected: false,
		},
		{
			Name:     "Invalid characters space",
			Input:    "invalid Exports Name",
			Expected: false,
		},
		{
			Name:     "Invalid name starts with hyphen",
			Input:    "-invalidExportsName",
			Expected: false,
		},
		{
			Name:     "Invalid name ends with hyphen",
			Input:    "invalidExportsName-",
			Expected: false,
		},
		{
			Name:     "Invalid name too long",
			Input:    "thisIsToLooooooooooooooooooooooooooooooooooooooongForAExportName",
			Expected: false,
		},
		{
			Name:     "Valid name",
			Input:    "validExportsName",
			Expected: true,
		},
		{
			Name:     "Valid name with hyphen",
			Input:    "validExportsName-2",
			Expected: true,
		},
		{
			Name:     "Valid name max length",
			Input:    "thisIsTheLoooooooooooooooooooooooooongestValidExportNameThereIs",
			Expected: true,
		},
		{
			Name:     "Valid name min length",
			Input:    "vali",
			Expected: true,
		},
	}
	for _, v := range testCases {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		_, errors := LogAnalyticsDataExportName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %v but got %v (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
