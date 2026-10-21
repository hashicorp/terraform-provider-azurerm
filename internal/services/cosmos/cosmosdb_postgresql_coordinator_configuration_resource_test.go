// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CosmosdbPostgresqlCoordinatorConfigurationResource struct{}

func TestAccCosmosDbPostgreSQLCoordinatorConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_coordinator_configuration", "test")
	r := CosmosdbPostgresqlCoordinatorConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "on"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbPostgreSQLCoordinatorConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_coordinator_configuration", "test")
	r := CosmosdbPostgresqlCoordinatorConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "on"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "off"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CosmosdbPostgresqlCoordinatorConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := configurations.ParseCoordinatorConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cosmos.ConfigurationsClient
	resp, err := client.GetCoordinator(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r CosmosdbPostgresqlCoordinatorConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-pshsccn-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_postgresql_cluster" "test" {
  name                            = "acctestcluster%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  administrator_login_password    = "H@Sh1CoR3!"
  coordinator_storage_quota_in_mb = 131072
  coordinator_vcore_count         = 2
  node_count                      = 0
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CosmosdbPostgresqlCoordinatorConfigurationResource) basic(data acceptance.TestData, value string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_postgresql_coordinator_configuration" "test" {
  name       = "array_nulls"
  cluster_id = azurerm_cosmosdb_postgresql_cluster.test.id
  value      = "%s"
}
`, r.template(data), value)
}
