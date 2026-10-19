// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestNodePoolSecuritySchemaValidation(t *testing.T) {
	testCases := []struct {
		name     string
		security []any
		valid    bool
	}{
		{name: "omitted", valid: true},
		{name: "empty block", security: []any{map[string]any{}}},
		{name: "null flags", security: []any{map[string]any{"secure_boot_enabled": nil, "vtpm_enabled": nil}}},
		{name: "secure boot disabled", security: []any{map[string]any{"secure_boot_enabled": false}}, valid: true},
		{name: "vtpm enabled", security: []any{map[string]any{"vtpm_enabled": true}}, valid: true},
	}
	for _, parent := range []struct {
		name   string
		path   string
		nested bool
	}{
		{name: "standalone", path: "security"},
		{name: "default node pool", path: "default_node_pool.0.security", nested: true},
	} {
		t.Run(parent.name, func(t *testing.T) {
			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					resource := &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"security": schemaNodePoolSecurity(parent.path),
					}}
					config := map[string]any{}
					if parent.nested {
						resource.Schema = map[string]*pluginsdk.Schema{
							"default_node_pool": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
									"name":     {Type: pluginsdk.TypeString, Required: true},
									"security": schemaNodePoolSecurity(parent.path),
								}},
							},
						}
						pool := map[string]any{"name": "test"}
						if tc.security != nil {
							pool["security"] = tc.security
						}
						config["default_node_pool"] = []any{pool}
					} else if tc.security != nil {
						config["security"] = tc.security
					}
					diagnostics := resource.Validate(terraform.NewResourceConfigRaw(config))
					if diagnostics.HasError() == tc.valid {
						t.Fatalf("expected valid %t, got diagnostics: %v", tc.valid, diagnostics)
					}
				})
			}
		})
	}
}

func TestNodePoolSecurityFlagsChanged(t *testing.T) {
	disabled := []any{map[string]any{"secure_boot_enabled": false, "vtpm_enabled": false}}
	enabled := []any{map[string]any{"secure_boot_enabled": true, "vtpm_enabled": true}}

	for _, tc := range []struct {
		name string
		old  []any
		new  []any
		want bool
	}{
		{name: "omitted to disabled", old: []any{}, new: disabled},
		{name: "disabled to omitted", old: disabled, new: []any{}},
		{name: "omitted to enabled", old: []any{}, new: enabled, want: true},
		{name: "enabled to omitted", old: enabled, new: []any{}, want: true},
		{name: "enabled to disabled", old: enabled, new: disabled, want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := nodePoolSecurityFlagsChanged(tc.old, tc.new); got != tc.want {
				t.Fatalf("nodePoolSecurityFlagsChanged() = %t, want %t", got, tc.want)
			}
		})
	}
}
