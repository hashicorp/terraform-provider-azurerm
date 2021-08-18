package validate

import (
	"testing"
)

func TestBase64EncodedString(t *testing.T) {
	cases := []struct {
		Input  string
		Errors int
	}{
		{
			Input:  "",
			Errors: 1,
		},
		{
			Input:  "aGVsbG8td29ybGQ=",
			Errors: 0,
		},
		{
			Input:  "hello-world",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			if _, errors := Base64EncodedString(tc.Input, "base64"); len(errors) != tc.Errors {
				t.Fatalf("Expected Base64 string to have %d not %d errors for %q: %v", tc.Errors, len(errors), tc.Input, errors)
			}
		})
	}
}

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
