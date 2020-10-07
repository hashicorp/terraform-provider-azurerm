package validate

import (
	"testing"
)

func TestLowerCasedStrings(t *testing.T) {
	cases := []struct {
		Value    string
		TestName string
		ErrCount int
	}{
		{
			Value:    "",
			TestName: "Empty",
			ErrCount: 1,
		},
		{
			Value:    " ",
			TestName: "Whitespace",
			ErrCount: 1,
		},
		{
			Value:    "Hello",
			TestName: "TitleCaseSingleWord",
			ErrCount: 1,
		},
		{
			Value:    "HELLO",
			TestName: "TitleCaseSingleWord",
			ErrCount: 1,
		},
		{
			Value:    "hello",
			TestName: "LowerCaseSingleWord",
			ErrCount: 0,
		},
		{
			Value:    "hello-there.com",
			TestName: "LowerCaseMultipleWords",
			ErrCount: 0,
		},
		{
			Value:    "hello there.com",
			TestName: "LowerCaseMultipleWordsWhitespace",
			ErrCount: 1,
		},
		{
			Value:    "Hello There.com",
			TestName: "TitleCaseMultipleWordsWhitespace",
			ErrCount: 1,
		},
		{
			Value:    "Hello-There.com",
			TestName: "TitleCaseMultipleWordsDash",
			ErrCount: 1,
		},
		{
			Value:    "HELLO-THERE.COM",
			TestName: "UpperCaseMultipleWordsDash",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			_, errors := LowerCasedString(tc.Value, tc.TestName)

			if len(errors) != tc.ErrCount {
				t.Fatalf("Expected NoEmptyStrings to have %d not %d errors for %q", tc.ErrCount, len(errors), tc.TestName)
			}
		})
	}
}
