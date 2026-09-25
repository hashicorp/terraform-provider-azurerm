// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestWindowsVirtualMachineWriteOnlyPasswordVersion(t *testing.T) {
	for _, test := range []struct {
		name      string
		version   cty.Value
		writeOnly bool
		wantError bool
	}{
		{"ordinary password without version", cty.NullVal(cty.Number), false, false},
		{"negative version", cty.NumberIntVal(-1), true, true},
		{"zero version", cty.NumberIntVal(0), true, true},
		{"initial version", cty.NumberIntVal(1), true, false},
		{"incremented version", cty.NumberIntVal(2), true, false},
		{"unknown version", cty.UnknownVal(cty.Number), true, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := resourceWindowsVirtualMachine()
			raw := windowsVMPasswordVersionConfig(t, r, test.writeOnly, test.version)
			cfg := terraform.NewResourceConfigShimmed(raw, r.CoreConfigSchema())
			diags := r.Validate(cfg)
			if diags.HasError() != test.wantError {
				t.Fatalf("configuration error status: got %t, want %t", diags.HasError(), test.wantError)
			}
			if test.wantError {
				for _, diagnostic := range diags {
					if strings.Contains(diagnostic.Summary+diagnostic.Detail, "admin_password_wo_version") {
						return
					}
				}
				t.Fatal("expected an error for admin_password_wo_version")
			}
		})
	}
}

func TestWindowsVirtualMachineImportedWriteOnlyPasswordVersion(t *testing.T) {
	for _, version := range []int64{1, 2} {
		t.Run(strconv.FormatInt(version, 10), func(t *testing.T) {
			r := resourceWindowsVirtualMachine()
			ordinary := windowsVMPasswordVersionConfig(t, r, false, cty.NullVal(cty.Number))
			imported := ordinary.AsValueMap()
			imported["id"] = cty.StringVal("imported-vm")
			imported["admin_password"] = cty.StringVal("ignored-as-imported")
			imported["admin_password_wo_version"] = cty.NumberIntVal(0)
			state, err := r.ShimInstanceStateFromValue(cty.ObjectVal(imported))
			if err != nil {
				t.Fatal(err)
			}

			raw := windowsVMPasswordVersionConfig(t, r, true, cty.NumberIntVal(version))
			cfg := terraform.NewResourceConfigShimmed(raw, r.CoreConfigSchema())
			if r.Validate(cfg).HasError() {
				t.Fatal("valid write-only configuration was rejected")
			}
			state.RawConfig = raw
			diff, err := r.SimpleDiff(context.Background(), state, cfg, nil)
			if err != nil {
				t.Fatal(err)
			}
			if diff == nil {
				t.Fatal("expected a replacement diff for the imported VM")
			}
			trigger := diff.Attributes["admin_password_wo_version"]
			if trigger == nil || !trigger.RequiresNew {
				t.Fatal("the write-only version must replace an imported VM even when its ordinary password diff is suppressed")
			}
		})
	}
}

func windowsVMPasswordVersionConfig(t *testing.T, r *pluginsdk.Resource, writeOnly bool, version cty.Value) cty.Value {
	t.Helper()
	const prefix = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/"
	values := map[string]cty.Value{
		"name":                      cty.StringVal("example-vm"),
		"resource_group_name":       cty.StringVal("example"),
		"location":                  cty.StringVal("westeurope"),
		"size":                      cty.StringVal("Standard_F2"),
		"admin_username":            cty.StringVal("adminuser"),
		"admin_password_wo_version": version,
		"network_interface_ids":     cty.ListVal([]cty.Value{cty.StringVal(prefix + "Microsoft.Network/networkInterfaces/example")}),
		"source_image_id":           cty.StringVal(prefix + "Microsoft.Compute/images/example"),
		"os_disk": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"caching":              cty.StringVal("ReadWrite"),
			"storage_account_type": cty.StringVal("Standard_LRS"),
		})}),
	}
	key := "admin_password"
	if writeOnly {
		key = "admin_password_wo"
	}
	values[key] = cty.StringVal("P@ss-test-only-123")
	raw, err := r.CoreConfigSchema().CoerceValue(cty.ObjectVal(values))
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
