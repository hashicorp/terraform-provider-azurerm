package validate

import "testing"

func TestApplicationSubdomain(t *testing.T) {
	testData := []struct {
		Value string
		Error bool
	}{
		{
			Value: "a1",
			Error: false,
		},
		{
			Value: "11",
			Error: false,
		},
		{
			Value: "1a",
			Error: false,
		},
		{
			Value: "aa",
			Error: false,
		},
		{
			Value: "1-1",
			Error: false,
		},
		{
			Value: "a-a",
			Error: false,
		},
		{
			Value: "a1-1",
			Error: false,
		},
		{
			Value: "a1-a",
			Error: false,
		},
		{
			Value: "1a-1",
			Error: false,
		},
		{
			Value: "1a-a",
			Error: false,
		},
		{
			Value: "a1-11",
			Error: false,
		},
		{
			Value: "aa-11",
			Error: false,
		},
		{
			Value: "11-1a",
			Error: false,
		},
		{
			Value: "11-a1",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde123",
			Error: false,
		},
		{
			Value: "a",
			Error: true,
		},
		{
			Value: "1",
			Error: true,
		},
		{
			Value: "1-",
			Error: true,
		},
		{
			Value: "a-",
			Error: true,
		},
		{
			Value: "a1-",
			Error: true,
		},
		{
			Value: "1a-",
			Error: true,
		},
		{
			Value: "aa-",
			Error: true,
		},
		{
			Value: "-",
			Error: true,
		},
		{
			Value: "-1",
			Error: true,
		},
		{
			Value: "-a",
			Error: true,
		},
		{
			Value: "AA",
			Error: true,
		},
		{
			Value: "AA-1",
			Error: true,
		},
		{
			Value: "AA-a",
			Error: true,
		},
		{
			Value: "A1-",
			Error: true,
		},
		{
			Value: "AA-A",
			Error: true,
		},
		{
			Value: "AA-aA",
			Error: true,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Value)

		_, err := ApplicationSubdomain(v.Value, "unit test")
		if err != nil && !v.Error {
			t.Fatalf("Expected pass but got an error: %s", err)
		}
	}
}
