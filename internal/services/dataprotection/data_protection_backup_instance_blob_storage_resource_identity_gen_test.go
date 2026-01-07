// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccDataProtectionBackupInstanceBlobStorage_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_blob_storage", "test")
	r := DataProtectionBackupInstanceBlobStorageResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"name":                {},
		"backup_vault_name":   {},
		"resource_group_name": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_data_protection_backup_instance_blob_storage.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_data_protection_backup_instance_blob_storage.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_data_protection_backup_instance_blob_storage.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_protection_backup_instance_blob_storage.test", tfjsonpath.New("backup_vault_name"), tfjsonpath.New("vault_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_protection_backup_instance_blob_storage.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("vault_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
