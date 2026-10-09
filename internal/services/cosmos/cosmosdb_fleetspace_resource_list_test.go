// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccCosmosdbFleetspace_list_basic(t *testing.T) {
	r := CosmosdbFleetspaceResource{}
	listResourceAddress := "azurerm_cosmosdb_fleetspace.list"

	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_fleetspace", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r CosmosdbFleetspaceResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_fleetspace" "test" {
  count = 3

  name                = "acctest-cosfleetspace-${count.index}-%d"
  resource_group_name = azurerm_resource_group.test.name
  fleet_name          = azurerm_cosmosdb_fleet.test.name
  service_tier        = "GeneralPurpose"
  data_regions = [
    azurerm_resource_group.test.location
  ]
}
`, r.template(data), data.RandomInteger)
}

func (CosmosdbFleetspaceResource) basicQuery() string {
	return `
list "azurerm_cosmosdb_fleetspace" "list" {
  provider = azurerm
  config {
    fleet_id = azurerm_cosmosdb_fleet.test.id
  }
}
`
}
