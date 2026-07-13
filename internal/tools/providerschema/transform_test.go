// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package providerschema

import (
	"encoding/json"
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

// TestNestedBlocksRoundTrip guards against a regression where blocks nested two
// or more levels deep lost their Elem (and therefore all of their child paths)
// when a schema was loaded from JSON. Level-1 properties are decoded by
// SchemaJSON.UnmarshalJSON, but deeper ones go through SchemaFromMap, which must
// decode a nested elem map the same way.
func TestNestedBlocksRoundTrip(t *testing.T) {
	source := &ProviderWrapper{
		ProviderName: "azurerm",
		ProviderSchema: &ProviderSchemaJSON{
			ResourcesMap: map[string]ResourceJSON{
				"azurerm_test": {
					Schema: map[string]SchemaJSON{
						"level1": {
							Type: SchemaTypeList,
							Elem: &ResourceJSON{Schema: map[string]SchemaJSON{
								"level2": {
									Type: SchemaTypeList,
									Elem: &ResourceJSON{Schema: map[string]SchemaJSON{
										"leaf": {Type: "TypeString", Optional: true},
									}},
								},
							}},
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got ProviderWrapper
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	level1 := got.ProviderSchema.ResourcesMap["azurerm_test"].Schema["level1"]
	l1Elem, ok := level1.Elem.(ResourceJSON)
	if !ok {
		t.Fatalf("level1.Elem = %T, want ResourceJSON", level1.Elem)
	}

	level2, ok := l1Elem.Schema["level2"]
	if !ok {
		t.Fatal("level2 missing from the level1 block after round-trip")
	}

	// Before the fix this Elem was nil: SchemaFromMap decoded a nested elem map
	// via decodeElem, which has no case for map[string]interface{}.
	l2Elem, ok := level2.Elem.(ResourceJSON)
	if !ok {
		t.Fatalf("level2.Elem = %T, want ResourceJSON (deep nested block lost its Elem)", level2.Elem)
	}
	if _, ok := l2Elem.Schema["leaf"]; !ok {
		t.Fatal("leaf missing from the level2 block after round-trip")
	}
}
