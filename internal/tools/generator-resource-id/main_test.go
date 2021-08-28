package main

import (
	"testing"
)

func Test(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{
			"",
			"",
		},
		{
			"a",
			"a",
		},
		{
			"A",
			"a",
		},
		{
			"alllower",
			"alllower",
		},
		{
			"ALLUPPER",
			"allupper",
		},
		{
			"lowerUPPER",
			"lower_upper",
		},
		{
			"UPPERLower",
			"upper_lower",
		},
		{
			"aAa",
			"a_aa",
		},
		{
			"AaA",
			"aa_a",
		},
	}

	for idx, c := range cases {
		out := convertToSnakeCase(c.in)
		if c.out != out {
			t.Fatalf("%d. %q (expect) != %q (actual)", idx, c.out, out)
		}
	}
}
