// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestStorageAccountCustomerManagedKeyIdentityValidation(t *testing.T) {
	const userAssignedIdentityID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test"

	tests := []struct {
		name                   string
		identityType           string
		userAssignedIdentityID cty.Value
		customerManagedKey     bool
		expectError            bool
	}{
		{
			name:                   "system assigned identity",
			identityType:           "SystemAssigned",
			userAssignedIdentityID: cty.NullVal(cty.String),
			customerManagedKey:     true,
		},
		{
			name:                   "user assigned identity",
			identityType:           "UserAssigned",
			userAssignedIdentityID: cty.StringVal(userAssignedIdentityID),
			customerManagedKey:     true,
		},
		{
			name:                   "user assigned identity unknown until apply",
			identityType:           "UserAssigned",
			userAssignedIdentityID: cty.UnknownVal(cty.String),
			customerManagedKey:     true,
		},
		{
			name:                   "user assigned identity missing",
			identityType:           "UserAssigned",
			userAssignedIdentityID: cty.NullVal(cty.String),
			customerManagedKey:     true,
			expectError:            true,
		},
		{
			name:                   "combined identity missing user assigned identity",
			identityType:           "SystemAssigned, UserAssigned",
			userAssignedIdentityID: cty.NullVal(cty.String),
			customerManagedKey:     true,
			expectError:            true,
		},
		{
			name:         "no inline customer managed key",
			identityType: "UserAssigned",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resource := resourceStorageAccount()
			identity := map[string]cty.Value{"type": cty.StringVal(test.identityType)}
			if strings.Contains(test.identityType, "UserAssigned") {
				identity["identity_ids"] = cty.SetVal([]cty.Value{cty.StringVal(userAssignedIdentityID)})
			}
			config := map[string]cty.Value{
				"name":                     cty.StringVal("accteststorage"),
				"resource_group_name":      cty.StringVal("acctestRG"),
				"location":                 cty.StringVal("westeurope"),
				"account_tier":             cty.StringVal("Standard"),
				"account_replication_type": cty.StringVal("LRS"),
				"account_kind":             cty.StringVal("StorageV2"),
				"identity":                 cty.ListVal([]cty.Value{cty.ObjectVal(identity)}),
			}
			if test.customerManagedKey {
				config["customer_managed_key"] = cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
					"key_vault_key_id":          cty.StringVal("https://acctest.vault.azure.net/keys/test/00000000000000000000000000000000"),
					"user_assigned_identity_id": test.userAssignedIdentityID,
				})})
			}

			rawConfig, err := resource.CoreConfigSchema().CoerceValue(cty.ObjectVal(config))
			if err != nil {
				t.Fatalf("building configuration: %v", err)
			}
			resourceConfig := terraform.NewResourceConfigShimmed(rawConfig, resource.CoreConfigSchema())
			if diags := resource.Validate(resourceConfig); diags.HasError() {
				t.Fatalf("validating configuration: %v", diags)
			}

			_, err = resource.Diff(context.Background(), &terraform.InstanceState{RawConfig: rawConfig}, resourceConfig, nil)
			if test.expectError {
				if err == nil || !strings.Contains(err.Error(), "user_assigned_identity_id` must be specified") {
					t.Fatalf("expected missing user assigned identity error, got: %v", err)
				}
			} else if err != nil {
				t.Fatalf("unexpected plan error: %v", err)
			}
		})
	}
}
