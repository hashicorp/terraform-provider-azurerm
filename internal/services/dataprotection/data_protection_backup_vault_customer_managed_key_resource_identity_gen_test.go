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

func TestAccDataProtectionBackupVaultCustomerManagedKey_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault_customer_managed_key", "test")
	r := DataProtectionBackupVaultCustomerManagedKeyResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValue("azurerm_data_protection_backup_vault_customer_managed_key.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_protection_backup_vault_customer_managed_key.test", tfjsonpath.New("name"), tfjsonpath.New("data_protection_backup_vault_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_protection_backup_vault_customer_managed_key.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("data_protection_backup_vault_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
