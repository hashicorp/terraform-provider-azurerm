package validate

import (
	"testing"
)

func TestLogAnalyticsClusterName(t *testing.T) {
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
			Input:    "invalid_Clusters_Name",
			Expected: false,
		},
		{
			Name:     "Invalid characters space",
			Input:    "invalid Clusters Name",
			Expected: false,
		},
		{
			Name:     "Invalid name starts with hyphen",
			Input:    "-invalidClustersName",
			Expected: false,
		},
		{
			Name:     "Invalid name ends with hyphen",
			Input:    "invalidClustersName-",
			Expected: false,
		},
		{
			Name:     "Invalid name too long",
			Input:    "thisIsToLoooooooooooooooooooooooooooooooooooooongForAClusterName",
			Expected: false,
		},
		{
			Name:     "Valid name",
			Input:    "validClustersName",
			Expected: true,
		},
		{
			Name:     "Valid name with hyphen",
			Input:    "validClustersName-2",
			Expected: true,
		},
		{
			Name:     "Valid name max length",
			Input:    "thisIsTheLooooooooooooooooooooooooongestValidClusterNameThereIs",
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

		_, errors := LogAnalyticsClusterName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
