// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package providerjson

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func TestValidatorAcceptFunc(t *testing.T) {
	cases := []struct {
		name    string
		s       *schema.Schema
		wantNil bool
		accepts map[string]bool
	}{
		{
			name:    "classic enum",
			s:       &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringInSlice([]string{"None", "Auto"}, false)},
			accepts: map[string]bool{"None": true, "Auto": true, "Nope": false},
		},
		{
			name:    "diag enum",
			s:       &schema.Schema{Type: schema.TypeString, ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"None"}, false))},
			accepts: map[string]bool{"None": true, "Nope": false},
		},
		{
			name:    "length validator",
			s:       &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringLenBetween(1, 4)},
			accepts: map[string]bool{"None": true, "toolong": false},
		},
		{name: "no validator", s: &schema.Schema{Type: schema.TypeString}, wantNil: true},
		{name: "non-string is skipped", s: &schema.Schema{Type: schema.TypeInt}, wantNil: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			accept := validatorAcceptFunc(tc.s)
			if tc.wantNil {
				if accept != nil {
					t.Fatalf("expected nil accept func, got non-nil")
				}
				return
			}
			if accept == nil {
				t.Fatalf("expected non-nil accept func")
			}
			for value, want := range tc.accepts {
				if got := accept(value); got != want {
					t.Fatalf("accept(%q) = %v, want %v", value, got, want)
				}
			}
		})
	}
}
