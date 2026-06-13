// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccCosmosdbFleetspace_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_fleetspace", "test")
	r := CosmosdbFleetspaceResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"fleet_name":          {},
		"name":                {},
		"resource_group_name": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_cosmosdb_fleetspace.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_cosmosdb_fleetspace.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_cosmosdb_fleetspace.test", tfjsonpath.New("fleet_name"), tfjsonpath.New("fleet_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_cosmosdb_fleetspace.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_cosmosdb_fleetspace.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
