// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccStorageMoverSourceEndpoint_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_source_endpoint", "test")
	r := StorageMoverSourceEndpointResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"resource_group_name": {},
		"storage_mover_name":  {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_storage_mover_source_endpoint.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_storage_mover_source_endpoint.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_mover_source_endpoint.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("storage_mover_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_mover_source_endpoint.test", tfjsonpath.New("storage_mover_name"), tfjsonpath.New("storage_mover_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_mover_source_endpoint.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("storage_mover_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
