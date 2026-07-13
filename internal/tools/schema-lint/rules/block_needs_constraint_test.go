// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func optionalBlock(schema map[string]providerschema.SchemaJSON) providerschema.SchemaJSON {
	return providerschema.SchemaJSON{Type: TypeList, Optional: true, MaxItems: 1, Elem: &providerschema.ResourceJSON{Schema: schema}}
}

func TestBlockNeedsConstraint(t *testing.T) {
	cases := []struct {
		name     string
		schema   providerschema.SchemaJSON
		wantHits int
	}{
		{name: "all optional, no constraint", schema: optionalBlock(map[string]providerschema.SchemaJSON{"a": {Type: TypeString, Optional: true}, "b": {Type: TypeString, Optional: true}}), wantHits: 1},
		{name: "has a required field", schema: optionalBlock(map[string]providerschema.SchemaJSON{"a": {Type: TypeString, Required: true}}), wantHits: 0},
		{name: "has AtLeastOneOf", schema: optionalBlock(map[string]providerschema.SchemaJSON{"a": {Type: TypeString, Optional: true, AtLeastOneOf: []string{"x.0.a", "x.0.b"}}}), wantHits: 0},
		{name: "computed-only block", schema: providerschema.SchemaJSON{Type: TypeList, Computed: true, Elem: &providerschema.ResourceJSON{Schema: map[string]providerschema.SchemaJSON{"a": {Type: TypeString}}}}, wantHits: 0},
		{name: "not a block", schema: providerschema.SchemaJSON{Type: TypeString, Optional: true}, wantHits: 0},
	}

	rule := blockNeedsConstraint{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: "setting", Path: "setting", Schema: tc.schema})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
		})
	}
}
