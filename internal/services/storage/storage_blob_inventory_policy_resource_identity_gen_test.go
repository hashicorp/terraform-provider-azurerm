// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccStorageBlobInventoryPolicy_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_blob_inventory_policy", "test")
	r := StorageBlobInventoryPolicyResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_blob_inventory_policy.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("storage_account_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_blob_inventory_policy.test", tfjsonpath.New("storage_account_name"), tfjsonpath.New("storage_account_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_blob_inventory_policy.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("storage_account_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
