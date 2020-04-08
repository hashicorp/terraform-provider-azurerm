package helper

import "testing"

func TestConvertToAdvisorSuppresionTTL(t *testing.T) {
	testCases := []struct {
		input  int
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
		es := FormatSuppressionTTL(test.input)

		if es != test.expect {
			t.Fatalf("Expected %q convert to %q", test.input, test.expect)
		}
	}
}

func TestParseSuppresionTTL(t *testing.T) {
	testCases := []struct {
		input       string
		expect      int
		shouldError bool
	}{
		{"-1", -1, false},
		{"0:30:0", 1800, false},
		{"1.1:0:0", 90000, false},
		{"7.0:0:0", 604800, false},
		{"0:0:1", 1, false},
		{"0:1:0", 60, false},
		{"1:0:0", 3600, false},
		{"3000.0:0:0", 259200000, false},
		{"3000.24:0:0", 0, true},
		{"3000.0:70:0", 0, true},
		{"3000.0:0:90", 0, true},
		{"", 0, true},
		{"1.2.3.4", 0, true},
	}

	for _, test := range testCases {
		es, err := ParseSuppresionTTL(test.input)
		if err == nil && test.shouldError == true {
			t.Fatalf("Expected %q to raise error", test.input)
		}
		if es != test.expect {
			t.Fatalf("Expected %q convert to %q", test.input, test.expect)
		}
	}
}
