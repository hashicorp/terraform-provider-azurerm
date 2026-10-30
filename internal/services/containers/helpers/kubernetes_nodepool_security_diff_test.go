// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestNodePoolSecurityCustomizeDiff(t *testing.T) {
	securityType := cty.List(cty.Object(map[string]cty.Type{
		"secure_boot_enabled": cty.Bool,
		"vtpm_enabled":        cty.Bool,
	}))
	security := func(boot, vtpm cty.Value) cty.Value {
		return cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"secure_boot_enabled": boot,
			"vtpm_enabled":        vtpm,
		})})
	}
	testCases := []struct {
		name     string
		osSKU    cty.Value
		security cty.Value
		errorKey string
	}{
		{"omitted", cty.StringVal("AzureContainerLinux"), cty.NullVal(securityType), "secure_boot_enabled"},
		{"empty block list", cty.StringVal("AzureContainerLinux"), cty.ListValEmpty(securityType.ElementType()), "secure_boot_enabled"},
		{"null secure boot", cty.StringVal("AzureContainerLinux"), security(cty.NullVal(cty.Bool), cty.True), "secure_boot_enabled"},
		{"null vtpm", cty.StringVal("AzureContainerLinux"), security(cty.True, cty.NullVal(cty.Bool)), "vtpm_enabled"},
		{"enabled", cty.StringVal("AzureContainerLinux"), security(cty.True, cty.True), ""},
		{"disabled", cty.StringVal("AzureContainerLinux"), security(cty.False, cty.False), "secure_boot_enabled"},
		{"secure boot only", cty.StringVal("AzureContainerLinux"), security(cty.True, cty.False), "vtpm_enabled"},
		{"vtpm only", cty.StringVal("AzureContainerLinux"), security(cty.False, cty.True), "secure_boot_enabled"},
		{"unknown secure boot", cty.StringVal("AzureContainerLinux"), security(cty.UnknownVal(cty.Bool), cty.True), ""},
		{"unknown vtpm", cty.StringVal("AzureContainerLinux"), security(cty.True, cty.UnknownVal(cty.Bool)), ""},
		{"unknown block", cty.StringVal("AzureContainerLinux"), cty.UnknownVal(securityType), ""},
		{"unknown os sku", cty.UnknownVal(cty.String), cty.NullVal(securityType), ""},
		{"known disabled secure boot", cty.StringVal("AzureContainerLinux"), security(cty.False, cty.UnknownVal(cty.Bool)), "secure_boot_enabled"},
		{"known disabled vtpm", cty.StringVal("AzureContainerLinux"), security(cty.UnknownVal(cty.Bool), cty.False), "vtpm_enabled"},
		{"ubuntu omitted", cty.StringVal("Ubuntu"), cty.NullVal(securityType), ""},
		{"ubuntu disabled", cty.StringVal("Ubuntu"), security(cty.False, cty.False), ""},
	}
	for _, nested := range []bool{false, true} {
		name, osSKUPath, securityPath := "standalone", "os_sku", "security"
		if nested {
			name, osSKUPath, securityPath = "default", "default_node_pool.0.os_sku", "default_node_pool.0.security"
		}
		t.Run(name, func(t *testing.T) {
			fields := map[string]*schema.Schema{
				"os_sku":   SchemaDefaultNodePool().Elem.(*schema.Resource).Schema["os_sku"],
				"security": SchemaNodePoolSecurity(securityPath),
			}
			if nested {
				fields = map[string]*schema.Schema{
					"default_node_pool": {Type: schema.TypeList, Required: true, MaxItems: 1, Elem: &schema.Resource{Schema: fields}},
				}
			}
			r := &schema.Resource{
				Schema:        fields,
				CustomizeDiff: pluginsdk.CustomDiffInSequence(NodePoolSecurityCustomizeDiff(osSKUPath, securityPath)),
			}
			config := func(osSKU, security cty.Value) *terraform.ResourceConfig {
				value := cty.ObjectVal(map[string]cty.Value{"os_sku": osSKU, "security": security})
				if nested {
					value = cty.ObjectVal(map[string]cty.Value{"default_node_pool": cty.ListVal([]cty.Value{value})})
				}
				return terraform.NewResourceConfigShimmed(value, r.CoreConfigSchema())
			}
			if nested {
				t.Run("unknown pool", func(t *testing.T) {
					value := cty.ObjectVal(map[string]cty.Value{
						"default_node_pool": cty.UnknownVal(cty.List(cty.Object(map[string]cty.Type{"os_sku": cty.String, "security": securityType}))),
					})
					if _, err := r.Diff(context.Background(), nil, terraform.NewResourceConfigShimmed(value, r.CoreConfigSchema()), nil); err != nil {
						t.Fatal(err)
					}
				})
			}
			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					_, err := r.Diff(context.Background(), nil, config(tc.osSKU, tc.security), nil)
					if tc.errorKey != "" {
						if err == nil || !strings.Contains(err.Error(), securityPath+".0."+tc.errorKey) {
							t.Fatalf("expected an error for %s, got %v", tc.errorKey, err)
						}
					} else if err != nil {
						t.Fatal(err)
					}
				})
			}
			for _, osSKU := range []string{"AzureContainerLinux", "Ubuntu"} {
				t.Run(osSKU+" refreshed state", func(t *testing.T) {
					attrs := map[string]string{osSKUPath: osSKU, securityPath + ".#": "1", securityPath + ".0.secure_boot_enabled": "true", securityPath + ".0.vtpm_enabled": "true"}
					if nested {
						attrs["default_node_pool.#"] = "1"
					}
					state := &terraform.InstanceState{ID: "test", Attributes: attrs}
					diff, err := r.Diff(context.Background(), state, config(cty.StringVal(osSKU), security(cty.True, cty.True)), nil)
					if err != nil || (diff != nil && len(diff.Attributes) != 0) {
						t.Fatalf("expected unchanged enabled security, got diff %v, error %v", diff, err)
					}
					diff, err = r.Diff(context.Background(), state, config(cty.StringVal(osSKU), cty.NullVal(securityType)), nil)
					if osSKU == "AzureContainerLinux" {
						if err == nil {
							t.Fatal("expected removing required ACL security to fail")
						}
					} else if err != nil || diff == nil || diff.Attributes[securityPath+".#"] == nil || diff.Attributes[securityPath+".#"].New != "0" {
						t.Fatalf("expected Ubuntu security removal, got diff %v, error %v", diff, err)
					}
				})
			}
		})
	}
}
