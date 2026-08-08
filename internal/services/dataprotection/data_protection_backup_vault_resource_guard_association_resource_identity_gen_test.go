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

func TestAccDataProtectionBackupVaultResourceGuardAssociation_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault_resource_guard_association", "test")
	r := DataProtectionBackupVaultResourceGuardAssociationResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"backup_vault_name":   {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_data_protection_backup_vault_resource_guard_association.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_data_protection_backup_vault_resource_guard_association.test", tfjsonpath.New("name"), knownvalue.StringExact("DppResourceGuardProxy")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_protection_backup_vault_resource_guard_association.test", tfjsonpath.New("backup_vault_name"), tfjsonpath.New("data_protection_backup_vault_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_protection_backup_vault_resource_guard_association.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("data_protection_backup_vault_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_protection_backup_vault_resource_guard_association.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("data_protection_backup_vault_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
