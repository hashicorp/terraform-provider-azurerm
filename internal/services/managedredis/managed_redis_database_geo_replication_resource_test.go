// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedRedisDatabaseGeoReplicationResource struct{}

func TestAccManagedRedisDatabaseGeoReplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database_geo_replication", "test")
	r := ManagedRedisDatabaseGeoReplicationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisDatabaseGeoReplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database_geo_replication", "test")
	r := ManagedRedisDatabaseGeoReplicationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.addThirdCluster(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.unlinkAllDbs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisDatabaseGeoReplication_linkedDbDoesNotContainSelf(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database_geo_replication", "test")
	r := ManagedRedisDatabaseGeoReplicationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.doesNotContainSelf(),
			ExpectError: regexp.MustCompile("`linked_database_ids` must include `database_id`"),
		},
	})
}

func (r ManagedRedisDatabaseGeoReplicationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := databases.ParseDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ManagedRedis.DatabaseClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagedRedisDatabaseGeoReplicationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_redis_cluster" "c1" {
  name                = "acctest-c1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  sku_name            = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "d1" {
  cluster_id = azurerm_managed_redis_cluster.c1.id

  geo_replication_group_name = "acctest-georep-%[1]d"
}

resource "azurerm_managed_redis_cluster" "c2" {
  name                = "acctest-c2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  sku_name            = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "d2" {
  cluster_id = azurerm_managed_redis_cluster.c2.id

  geo_replication_group_name = "acctest-georep-%[1]d"
}

resource "azurerm_managed_redis_database_geo_replication" "test" {
  database_id = azurerm_managed_redis_database.d1.id
  linked_database_ids = [
    azurerm_managed_redis_database.d1.id,
    azurerm_managed_redis_database.d2.id,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ManagedRedisDatabaseGeoReplicationResource) addThirdCluster(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_redis_cluster" "c1" {
  name                = "acctest-c1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  sku_name            = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "d1" {
  cluster_id = azurerm_managed_redis_cluster.c1.id

  geo_replication_group_name = "acctest-georep-%[1]d"
}

resource "azurerm_managed_redis_cluster" "c2" {
  name                = "acctest-c2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  sku_name            = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "d2" {
  cluster_id = azurerm_managed_redis_cluster.c2.id

  geo_replication_group_name = "acctest-georep-%[1]d"
}

resource "azurerm_managed_redis_cluster" "c3" {
  name                = "acctest-c3-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"
  sku_name            = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "d3" {
  cluster_id = azurerm_managed_redis_cluster.c3.id

  geo_replication_group_name = "acctest-georep-%[1]d"
}

resource "azurerm_managed_redis_database_geo_replication" "test" {
  database_id = azurerm_managed_redis_database.d1.id
  
  linked_database_ids = [
    azurerm_managed_redis_database.d1.id,
    azurerm_managed_redis_database.d2.id,
    azurerm_managed_redis_database.d3.id,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, data.Locations.Ternary)
}

func (r ManagedRedisDatabaseGeoReplicationResource) unlinkAllDbs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_redis_cluster" "c1" {
  name                = "acctest-c1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  sku_name            = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "d1" {
  cluster_id = azurerm_managed_redis_cluster.c1.id

  geo_replication_group_name = "acctest-georep-%[1]d"
}

resource "azurerm_managed_redis_cluster" "c2" {
  name                = "acctest-c2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  sku_name            = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "d2" {
  cluster_id = azurerm_managed_redis_cluster.c2.id

  geo_replication_group_name = "acctest-georep-%[1]d"
}

resource "azurerm_managed_redis_cluster" "c3" {
  name                = "acctest-c3-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"
  sku_name            = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "d3" {
  cluster_id = azurerm_managed_redis_cluster.c3.id

  geo_replication_group_name = "acctest-georep-%[1]d"
}

resource "azurerm_managed_redis_database_geo_replication" "test" {
  database_id = azurerm_managed_redis_database.d1.id
  
	linked_database_ids = [
    azurerm_managed_redis_database.d1.id,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, data.Locations.Ternary)
}

func (r ManagedRedisDatabaseGeoReplicationResource) doesNotContainSelf() string {
	return `
provider "azurerm" {
	features {}
}

resource "azurerm_managed_redis_database_geo_replication" "test" {
  database_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/my-amr-1/databases/default"
	linked_database_ids = [
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/my-amr-2/databases/default",
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Cache/redisEnterprise/my-amr-3/databases/default",
	]
}
`
}
