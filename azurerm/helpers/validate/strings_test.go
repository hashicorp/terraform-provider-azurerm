package validate

import (
	"testing"
)

func TestNoEmptyStrings(t *testing.T) {
	cases := []struct {
		Value    string
		TestName string
		ErrCount int
	}{
		{
			Value:    "!",
			TestName: "Exclamation",
			ErrCount: 0,
		},
		{
			Value:    ".",
			TestName: "Period",
			ErrCount: 0,
		},
		{
			Value:    "-",
			TestName: "Hyphen",
			ErrCount: 0,
		},
		{
			Value:    "_",
			TestName: "Underscore",
			ErrCount: 0,
		},
		{
			Value:    "10.1.0.0/16",
			TestName: "IP",
			ErrCount: 0,
		},
		{
			Value:    "",
			TestName: "Empty",
			ErrCount: 1,
		},
		{
			Value:    " ",
			TestName: "Space",
			ErrCount: 1,
		},
		{
			Value:    "     ",
			TestName: "FiveSpaces",
			ErrCount: 1,
		},
		{
			Value:    "  1",
			TestName: "DoubleSpaceOne",
			ErrCount: 0,
		},
		{
			Value:    "1 ",
			TestName: "OneSpace",
			ErrCount: 0,
		},
		{
			Value:    "\r",
			TestName: "CarriageReturn",
			ErrCount: 1,
		},
		{
			Value:    "\n",
			TestName: "NewLine",
			ErrCount: 1,
		},
		{
			Value:    "\t",
			TestName: "HorizontalTab",
			ErrCount: 1,
		},
		{
			Value:    "\f",
			TestName: "FormFeed",
			ErrCount: 1,
		},
		{
			Value:    "\v",
			TestName: "VerticalTab",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			_, errors := NoEmptyStrings(tc.Value, tc.TestName)

			if len(errors) != tc.ErrCount {
				t.Fatalf("Expected NoEmptyStrings to have %d not %d errors for %q", tc.ErrCount, len(errors), tc.TestName)
			}
		})
	}
}
