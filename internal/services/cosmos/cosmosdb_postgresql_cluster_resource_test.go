// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosDbPostgreSQLClusterResource struct{}

func TestAccCosmosDbPostgreSQLCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_cluster", "test")
	r := CosmosDbPostgreSQLClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("servers.0.fqdn").IsNotEmpty(),
				check.That(data.ResourceName).Key("servers.0.name").IsNotEmpty(),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccCosmosDbPostgreSQLCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_cluster", "test")
	r := CosmosDbPostgreSQLClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCosmosDbPostgreSQLCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_cluster", "test")
	r := CosmosDbPostgreSQLClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccCosmosDbPostgreSQLCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_cluster", "test")
	r := CosmosDbPostgreSQLClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccCosmosDbPostgreSQLCluster_withSourceCluster(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_postgresql_cluster", "test")
	r := CosmosDbPostgreSQLClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			PreConfig: func() { time.Sleep(15 * time.Minute) },
			Config:    r.withSourceCluster(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "source_location", "source_resource_id", "point_in_time_in_utc"),
	})
}

func (r CosmosDbPostgreSQLClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := clusters.ParseServerGroupsv2ID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cosmos.ClustersClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r CosmosDbPostgreSQLClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-pshsc-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CosmosDbPostgreSQLClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_postgresql_cluster" "test" {
  name                            = "acctestcluster%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  administrator_login_password    = "H@Sh1CoR3!"
  coordinator_storage_quota_in_mb = 131072
  coordinator_vcore_count         = 2
  node_count                      = 0
}
`, r.template(data), data.RandomInteger)
}

func (r CosmosDbPostgreSQLClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_postgresql_cluster" "import" {
  name                            = azurerm_cosmosdb_postgresql_cluster.test.name
  resource_group_name             = azurerm_cosmosdb_postgresql_cluster.test.resource_group_name
  location                        = azurerm_cosmosdb_postgresql_cluster.test.location
  administrator_login_password    = azurerm_cosmosdb_postgresql_cluster.test.administrator_login_password
  coordinator_storage_quota_in_mb = azurerm_cosmosdb_postgresql_cluster.test.coordinator_storage_quota_in_mb
  coordinator_vcore_count         = azurerm_cosmosdb_postgresql_cluster.test.coordinator_vcore_count
  node_count                      = azurerm_cosmosdb_postgresql_cluster.test.node_count
}
`, r.basic(data))
}

func (r CosmosDbPostgreSQLClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_postgresql_cluster" "test" {
  name                = "acctestcluster%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  administrator_login_password    = "H@Sh1CoR3!"
  coordinator_storage_quota_in_mb = 131072
  coordinator_vcore_count         = 2
  node_count                      = 0

  citus_version                        = "11.1"
  coordinator_public_ip_access_enabled = true
  ha_enabled                           = false
  coordinator_server_edition           = "GeneralPurpose"

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  node_public_ip_access_enabled = false
  node_server_edition           = "MemoryOptimized"
  sql_version                   = "14"
  preferred_primary_zone        = 1
  node_storage_quota_in_mb      = 524288
  node_vcores                   = 4
  shards_on_coordinator_enabled = true

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CosmosDbPostgreSQLClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_postgresql_cluster" "test" {
  name                = "acctestcluster%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  administrator_login_password    = "H@Sh1CoR4!"
  coordinator_storage_quota_in_mb = 524288
  coordinator_vcore_count         = 4
  node_count                      = 2

  citus_version                        = "12.1"
  coordinator_public_ip_access_enabled = false
  ha_enabled                           = true
  coordinator_server_edition           = "MemoryOptimized"

  maintenance_window {
    day_of_week  = 1
    start_hour   = 9
    start_minute = 1
  }

  node_public_ip_access_enabled = true
  node_server_edition           = "GeneralPurpose"
  sql_version                   = "16"
  preferred_primary_zone        = 2
  node_storage_quota_in_mb      = 524288
  node_vcores                   = 4
  shards_on_coordinator_enabled = false

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CosmosDbPostgreSQLClusterResource) withSourceCluster(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_postgresql_cluster" "test2" {
  name                 = "acctesttcluster%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  source_location      = azurerm_cosmosdb_postgresql_cluster.test.location
  source_resource_id   = azurerm_cosmosdb_postgresql_cluster.test.id
  point_in_time_in_utc = azurerm_cosmosdb_postgresql_cluster.test.earliest_restore_time
  node_count           = 0

  lifecycle {
    ignore_changes = ["coordinator_storage_quota_in_mb", "coordinator_vcore_count"]
  }
}
`, r.basic(data), data.RandomInteger)
}
