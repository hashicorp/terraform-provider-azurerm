package validate

import (
	"testing"
)

func TestLogAnalyticsClustersName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "inv",
			Expected: false,
		},
		{
			Input:    "invalid_Cluster_Name",
			Expected: false,
		},
		{
			Input:    "invalid Cluster Name",
			Expected: false,
		},
		{
			Input:    "-invalidClusterName",
			Expected: false,
		},
		{
			Input:    "invalidClusterName-",
			Expected: false,
		},
		{
			Input:    "validClusterName",
			Expected: true,
		},
		{
			Input:    "validClusterName-2",
			Expected: true,
		},
		{
			Input:    "thisIsTheLooooooooooooooooooooooooongestValidClusterNameThereIs",
			Expected: true,
		},
		{
			Input:    "vali",
			Expected: true,
		},
	}
	for _, v := range testCases {
		_, errors := LogAnalyticsClustersName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
