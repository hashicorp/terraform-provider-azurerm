// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccManagedDisk_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk", "test")
	r := ManagedDiskResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.empty(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValue("azurerm_managed_disk.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_managed_disk.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_managed_disk.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
