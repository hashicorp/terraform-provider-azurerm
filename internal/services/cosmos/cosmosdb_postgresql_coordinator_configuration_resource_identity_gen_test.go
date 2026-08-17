// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccCosmosdbPostgresqlCoordinatorConfiguration_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_coordinator_configuration", "test")
	r := CosmosdbPostgresqlCoordinatorConfigurationResource{}

	checkedFields := map[string]struct{}{
		"name":                  {},
		"resource_group_name":   {},
		"server_groups_v2_name": {},
		"subscription_id":       {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data, "on"),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_cosmosdb_postgresql_coordinator_configuration.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_cosmosdb_postgresql_coordinator_configuration.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cosmosdb_postgresql_coordinator_configuration.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cosmosdb_postgresql_coordinator_configuration.test", tfjsonpath.New("server_groups_v2_name"), tfjsonpath.New("cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_cosmosdb_postgresql_coordinator_configuration.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("cluster_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
