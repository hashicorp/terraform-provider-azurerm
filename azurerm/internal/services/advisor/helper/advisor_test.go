package helper

import "testing"

func TestConvertToAdvisorSuppresionTTL(t *testing.T) {
	testCases := []struct {
		input       int
		expect string
	}{
		{-1, ""},
		{1800, "0.0:30:0"},
		{90000, "1.1:0:0"},
		{604800, "7.0:0:0"},
		{1, "0.0:0:1"},
		{60, "0.0:1:0"},
		{3600, "0.1:0:0"},
		{259200000, "3000.0:0:0"},
	}

	for _, test := range testCases {
		es := ConvertToAdvisorSuppresionTTL(test.input)

		if es!=test.expect {
			t.Fatalf("Expected %q convert to %q", test.input,test.expect)
		}
	}
}

func TestParseAdvisorSuppresionTTL(t *testing.T) {
	testCases := []struct {
		input       string
		expect int
	}{
		{ "-1",-1},
		{ "0:30:0",1800},
		{ "1.1:0:0",90000},
		{ "7.0:0:0",604800},
		{ "0:0:1",1},
		{"0:1:0",60},
		{ "1:0:0",3600},
		{ "3000.0:0:0",259200000},
	}

	for _, test := range testCases {
		es := ParseAdvisorSuppresionTTL(test.input)

		if es!=test.expect {
			t.Fatalf("Expected %q convert to %q", test.input,test.expect)
		}
	}
}