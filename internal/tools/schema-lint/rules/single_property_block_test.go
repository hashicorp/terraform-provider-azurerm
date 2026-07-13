// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func block(maxItems int, schema map[string]providerschema.SchemaJSON) providerschema.SchemaJSON {
	return providerschema.SchemaJSON{Type: TypeList, MaxItems: maxItems, Elem: &providerschema.ResourceJSON{Schema: schema}}
}

func TestSinglePropertyBlock(t *testing.T) {
	cases := []struct {
		name     string
		schema   providerschema.SchemaJSON
		wantHits int
	}{
		{name: "single enabled property, maxitems 1", schema: block(1, map[string]providerschema.SchemaJSON{"enabled": {Type: TypeBool, Optional: true}}), wantHits: 1},
		{name: "single non-enabled property, maxitems 1", schema: block(1, map[string]providerschema.SchemaJSON{"version": {Type: TypeString, Optional: true}}), wantHits: 1},
		{name: "two properties, maxitems 1", schema: block(1, map[string]providerschema.SchemaJSON{"a": {Type: TypeString}, "b": {Type: TypeString}}), wantHits: 0},
		{name: "single property, maxitems 0", schema: block(0, map[string]providerschema.SchemaJSON{"enabled": {Type: TypeBool}}), wantHits: 0},
		{name: "not a block", schema: providerschema.SchemaJSON{Type: TypeString, MaxItems: 1}, wantHits: 0},
	}

	rule := singlePropertyBlock{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: "feature", Path: "feature", Schema: tc.schema})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
		})
	}
}

func TestSinglePropertyBlock_EnabledFixSuggestion(t *testing.T) {
	rule := singlePropertyBlock{}
	got := rule.CheckProperty(PropertyContext{
		ResourceType: "azurerm_test",
		Kind:         KindResource,
		Name:         "blob_csi_driver",
		Path:         "storage_profile.blob_csi_driver",
		Schema:       block(1, map[string]providerschema.SchemaJSON{"enabled": {Type: TypeBool, Optional: true}}),
	})

	if len(got) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(got))
	}
	if !strings.Contains(got[0].FixSuggestion, "blob_csi_driver_enabled") {
		t.Fatalf("expected fix suggestion to reference blob_csi_driver_enabled, got %q", got[0].FixSuggestion)
	}
}
