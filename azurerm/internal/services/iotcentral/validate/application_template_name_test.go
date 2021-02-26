package validate

import (
	"testing"
)

func TestApplicationTemplateName(t *testing.T) {
	testData := []struct {
		Value string
		Error bool
	}{
		{
			Value: "a",
			Error: false,
		},
		{
			Value: "A",
			Error: false,
		},
		{
			Value: "1",
			Error: false,
		},
		{
			Value: "1-",
			Error: false,
		},
		{
			Value: "a-",
			Error: false,
		},
		{
			Value: "A-",
			Error: false,
		},
		{
			Value: "-",
			Error: false,
		},
		{
			Value: "-1",
			Error: false,
		},
		{
			Value: "_-a",
			Error: false,
		},
		{
			Value: "#$%$#!",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
			Error: false,
		},
		{
			Value: "",
			Error: true,
		},
		{
			Value: "abcdeabcdeabcdeabcde@$#%abcdeabcdeadeabcdeabcdeabcdeabcde-1a",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Value)

		_, err := ApplicationTemplateName(v.Value, "unit test")
		if err != nil && !v.Error {
			t.Fatalf("Expected pass but got an error: %s", err)
		}
	}
}
