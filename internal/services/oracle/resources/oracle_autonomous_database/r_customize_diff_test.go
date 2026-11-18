// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestAutonomousDatabaseCharacterSetsCustomizeDiff(t *testing.T) {
	tests := []struct {
		name         string
		workload     cty.Value
		characterSet cty.Value
		nationalSet  cty.Value
		wantError    string
	}{
		{
			name: "AJD omitted", workload: cty.StringVal("AJD"),
		},
		{
			name: "AJD explicit default character set", workload: cty.StringVal("AJD"),
			characterSet: cty.StringVal("AL32UTF8"), wantError: "`character_set` must be omitted",
		},
		{
			name: "AJD explicit default national character set", workload: cty.StringVal("AJD"),
			nationalSet: cty.StringVal("AL16UTF16"), wantError: "`national_character_set` must be omitted",
		},
		{
			name: "AJD custom character set", workload: cty.StringVal("AJD"),
			characterSet: cty.StringVal("WE8MSWIN1252"), wantError: "`character_set` must be omitted",
		},
		{
			name: "AJD custom national character set", workload: cty.StringVal("AJD"),
			nationalSet: cty.StringVal("UTF8"), wantError: "`national_character_set` must be omitted",
		},
		{
			name: "AJD unknown character set", workload: cty.StringVal("AJD"),
			characterSet: cty.UnknownVal(cty.String),
		},
		{
			name: "AJD unknown national character set", workload: cty.StringVal("AJD"),
			nationalSet: cty.UnknownVal(cty.String),
		},
		{
			name: "unknown workload", workload: cty.UnknownVal(cty.String),
			characterSet: cty.StringVal("AL32UTF8"), nationalSet: cty.StringVal("AL16UTF16"),
		},
		{
			name: "OLTP custom values", workload: cty.StringVal("OLTP"),
			characterSet: cty.StringVal("WE8MSWIN1252"), nationalSet: cty.StringVal("UTF8"),
		},
		{
			name: "DW defaults", workload: cty.StringVal("DW"),
		},
		{
			name: "APEX explicit defaults", workload: cty.StringVal("APEX"),
			characterSet: cty.StringVal("AL32UTF8"), nationalSet: cty.StringVal("AL16UTF16"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wrapper := sdk.NewResourceWrapper(Resource{})
			resource, err := wrapper.Resource()
			if err != nil {
				t.Fatalf("building resource: %v", err)
			}

			config := map[string]any{
				"name": "testdb", "display_name": "testdb", "resource_group_name": "test-rg",
				"location": "eastus", "admin_password": "TestPass#2024#", "db_version": "19c",
				"compute_model": "ECPU", "compute_count": 2.0, "data_storage_size_in_tbs": 1,
				"license_model": "LicenseIncluded", "backup_retention_period_in_days": 12,
				"auto_scaling_enabled": false, "auto_scaling_for_storage_enabled": false,
				"mtls_connection_required": true,
			}
			raw := map[string]cty.Value{
				"db_workload": test.workload, "character_set": test.characterSet, "national_character_set": test.nationalSet,
			}
			for name, value := range raw {
				if value == cty.NilVal {
					raw[name] = cty.NullVal(cty.String)
				} else if !value.IsKnown() {
					continue // Set ComputedKeys below for values that are not yet known.
				} else {
					config[name] = value.AsString()
				}
			}

			// Resource.Diff receives raw configuration through state, matching the
			// SDK's protocol shim. Omitted fields still acquire schema defaults.
			state := &terraform.InstanceState{RawConfig: cty.ObjectVal(raw)}
			resourceConfig := terraform.NewResourceConfigRaw(config)
			for name, value := range raw {
				if !value.IsKnown() {
					resourceConfig.ComputedKeys = append(resourceConfig.ComputedKeys, name)
				}
			}
			diff, err := resource.Diff(context.Background(), state, resourceConfig, &clients.Client{})
			if test.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("expected %q, got %v", test.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("building diff: %v", err)
			}
			for name, want := range map[string]string{"character_set": "AL32UTF8", "national_character_set": "AL16UTF16"} {
				if raw[name].IsNull() {
					if attr := diff.Attributes[name]; attr == nil || attr.New != want {
						t.Errorf("expected omitted %s to default to %q, got %+v", name, want, attr)
					}
				}
			}
		})
	}
}
