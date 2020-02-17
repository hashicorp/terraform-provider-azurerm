package parse

import (
	"testing"
)

func TestParseDataBoxJobID(t *testing.T) {
	testData := []struct {
		input    string
		expected *DataBoxJobID
	}{
		{
			input:    "",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/providers",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/providers/Microsoft.DataBox",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/providers/Microsoft.DataBox/jobs",
			expected: nil,
		},
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/providers/Microsoft.DataBox/jobs/job1",
			expected: &DataBoxJobID{
				Name:          "job1",
				ResourceGroup: "hello",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)
		actual, err := ParseDataBoxJobID(v.input)

		// if we get something there shouldn't be an error
		if v.expected != nil && err == nil {
			continue
		}

		// if nothing's expected we should get an error
		if v.expected == nil && err != nil {
			continue
		}

		if v.expected == nil && actual == nil {
			continue
		}

		if v.expected == nil && actual != nil {
			t.Fatalf("Expected nothing but got %+v", actual)
		}
		if v.expected != nil && actual == nil {
			t.Fatalf("Expected %+v but got nil", actual)
		}

		if v.expected.ResourceGroup != actual.ResourceGroup {
			t.Fatalf("Expected ResourceGroup to be %q but got %q", v.expected.ResourceGroup, actual.ResourceGroup)
		}
		if v.expected.Name != actual.Name {
			t.Fatalf("Expected Name to be %q but got %q", v.expected.Name, actual.Name)
		}
	}
}
