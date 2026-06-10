// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

func TestValidateFrontDoorConditionBlocksRequireMatchValues(t *testing.T) {
	tests := []struct {
		name        string
		condition   cty.Value
		errContains string
	}{
		{
			name: "missing match values invalid",
			condition: cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"match_values": cty.ListValEmpty(cty.String),
				}),
			}),
			errContains: "requires `match_values`",
		},
		{
			name:      "unknown condition value ignored",
			condition: cty.UnknownVal(cty.List(cty.Object(map[string]cty.Type{"match_values": cty.List(cty.String)}))),
		},
		{
			name: "unknown match values ignored",
			condition: cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"match_values": cty.UnknownVal(cty.List(cty.String)),
				}),
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateFrontDoorConditionBlocksRequireMatchValues(map[string]cty.Value{
				"request_scheme_condition": test.condition,
			}, []string{"request_scheme_condition"})

			if test.errContains == "" {
				if err != nil {
					t.Fatalf("expected no error but got %q", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error containing %q but got nil", test.errContains)
			}

			if err.Error() != "the `request_scheme_condition` block requires `match_values`" {
				t.Fatalf("expected canonical error but got %q", err)
			}
		})
	}
}
