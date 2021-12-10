package validate

import (
	"strings"
	"testing"
)

func TestSnapshotName_validation(t *testing.T) {
	str := strings.Repeat("a", 80)
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "ab",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "cosmosDBAccount1",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 0,
		},
		{
			Value:    "hello+world",
			ErrCount: 1,
		},
		{
			Value:    str,
			ErrCount: 0,
		},
		{
			Value:    str + "a",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := SnapshotName(tc.Value, "azurerm_snapshot")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Snapshot Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}
