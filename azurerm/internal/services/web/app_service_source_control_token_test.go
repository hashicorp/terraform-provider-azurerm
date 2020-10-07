package web

import "testing"

func TestValidateAppServiceSourceControlToken(t *testing.T) {
	testCases := []struct {
		Input string
		Valid bool
	}{
		{Input: "", Valid: false},
		{Input: "BitBucket", Valid: true},
		{Input: "Dropbox", Valid: true},
		{Input: "GitLab", Valid: false},
		{Input: "GitHub", Valid: true},
		{Input: "OneDrive", Valid: true},
	}

	for _, v := range testCases {
		t.Logf("[DEBUG] Testing %q..", v.Input)
		warns, err := ValidateAppServiceSourceControlTokenName()(v.Input, "id")
		if len(warns) > 0 {
			t.Fatalf("Got warnings when they should be errors")
		}
		isValid := len(err) == 0
		if v.Valid != isValid {
			t.Fatalf("Expected %t but got %t", v.Valid, isValid)
		}
	}
}
