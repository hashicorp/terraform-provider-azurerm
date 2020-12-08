package locks

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicatesInStringArray(t *testing.T) {
	cases := []struct {
		Name   string
		Input  []string
		Result []string
	}{
		{
			Name:   "contain duplicates",
			Input:  []string{"string1", "string2", "string1", "string3", ""},
			Result: []string{"string1", "string2", "string3", ""},
		},
		{
			Name:   "does not contain duplicates",
			Input:  []string{"string1", "string2", "string3", ""},
			Result: []string{"string1", "string2", "string3", ""},
		},
		{
			Name:   "empty array",
			Input:  []string{},
			Result: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if !reflect.DeepEqual(removeDuplicatesFromStringArray(tc.Input), tc.Result) {
				t.Fatalf("Expected TestRemoveDuplicatesInStringArray to return %v", tc.Result)
			}
		})
	}
}
