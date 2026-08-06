// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestDatabaseName(t *testing.T) {
	cases := []struct {
		Name  string
		Input string
		Valid bool
	}{
		{Name: "simple", Input: "my-database", Valid: true},
		{Name: "short", Input: "ab", Valid: true},
		{Name: "at max length", Input: strings.Repeat("a", 260), Valid: true},
		{Name: "too long", Input: strings.Repeat("a", 261), Valid: false},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, errs := DatabaseName(tc.Input, "name")
			if (len(errs) == 0) != tc.Valid {
				t.Fatalf("expected valid=%t for %q but got errors: %v", tc.Valid, tc.Input, errs)
			}
			for _, err := range errs {
				if strings.Contains(err.Error(), "between 4 and 22") {
					t.Fatalf("error message cites the wrong length bound: %v", err)
				}
			}
		})
	}
}

func TestDatabasePrincipalAssignmentName(t *testing.T) {
	_, errs := DatabasePrincipalAssignmentName(strings.Repeat("a", 261), "name")
	if len(errs) == 0 {
		t.Fatal("expected an error for a name longer than 260 characters")
	}
	for _, err := range errs {
		if strings.Contains(err.Error(), "between 4 and 22") {
			t.Fatalf("error message cites the wrong length bound: %v", err)
		}
	}
}

func TestClusterPrincipalAssignmentName(t *testing.T) {
	_, errs := ClusterPrincipalAssignmentName(strings.Repeat("a", 261), "name")
	if len(errs) == 0 {
		t.Fatal("expected an error for a name longer than 260 characters")
	}
	for _, err := range errs {
		if strings.Contains(err.Error(), "between 4 and 22") {
			t.Fatalf("error message cites the wrong length bound: %v", err)
		}
	}
}
