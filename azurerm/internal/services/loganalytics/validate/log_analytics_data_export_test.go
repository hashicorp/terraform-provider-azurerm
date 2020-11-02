package validate

import (
	"testing"
)

func TestLogAnalyticsDataExportName(t *testing.T) {
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
			Input:    "invalid Exports Name",
			Expected: false,
		},
		{
			Input:    "-invalidExportsName",
			Expected: false,
		},
		{
			Input:    "invalidExportsName-",
			Expected: false,
		},
		{
			Input:    "validExportsName",
			Expected: true,
		},
		{
			Input:    "validExportsName-2",
			Expected: true,
		},
		{
			Input:    "thisIsTheLooooooooooooooooooooooooongestValidExportsNameThereIs",
			Expected: true,
		},
		{
			Input:    "vali",
			Expected: true,
		},
	}
	for _, v := range testCases {
		_, errors := LogAnalyticsDataExportName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
