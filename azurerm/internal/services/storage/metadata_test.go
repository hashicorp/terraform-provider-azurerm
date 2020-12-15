package storage

import (
	"testing"
)

func TestValidateMetaDataKeys(t *testing.T) {
	testData := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "Hello",
			Expected: false,
		},
		{
			Input:    "hello",
			Expected: true,
		},
		{
			Input:    "hello",
			Expected: true,
		},
		{
			// C# keyword
			Input:    "using",
			Expected: false,
		},
		{
			Input:    "0hello",
			Expected: false,
		},
		{
			Input:    "heLLo",
			Expected: false,
		},
		{
			Input:    "panda_cycle",
			Expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		value := map[string]interface{}{
			v.Input: "hello",
		}
		warnings, errors := ValidateMetaDataKeys(value, "field")
		if len(warnings) != 0 {
			t.Fatalf("Expected no warnings but got %d", len(warnings))
		}

		actual := len(errors) == 0
		if v.Expected != actual {
			t.Fatalf("Expected %t but got %t", v.Expected, actual)
		}
	}
}
