// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagedRedisDatabaseDataSource struct{}

func TestAccManagedRedisDatabaseDataSource_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cluster_id").Exists(),
			),
		},
	})
}

func TestAccManagedRedisDatabaseDataSource_geoReplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSourceGeoReplication(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cluster_id").Exists(),
				check.That(data.ResourceName).Key("linked_database_id.#").Exists(),
				check.That(data.ResourceName).Key("geo_replication_group_name").Exists(),
			),
		},
	})
}

func TestAccManagedRedisDatabaseDataSource_withAccessKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSourceAccessKeysEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cluster_id").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
	})
}

func (r ManagedRedisDatabaseDataSource) dataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_managed_redis_database" "test" {
  cluster_id = azurerm_managed_redis_cluster.test.id
  depends_on = [azurerm_managed_redis_database.test]
}
`, ManagedRedisDatabaseResource{}.basic(data))
}

func (r ManagedRedisDatabaseDataSource) dataSourceGeoReplication(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_redis_cluster" "test" {
  name                = "acctest-rec1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "test" {
  cluster_id = azurerm_managed_redis_cluster.test.id

  geo_replication_group_name = "acctest-georepg-%[1]d"
}

data "azurerm_managed_redis_database" "test" {
  cluster_id = azurerm_managed_redis_cluster.test.id
  depends_on = [azurerm_managed_redis_database.test]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ManagedRedisDatabaseDataSource) dataSourceAccessKeysEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_redis_cluster" "test" {
  name                = "acctest-rec1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "test" {
  cluster_id                         = azurerm_managed_redis_cluster.test.id
  access_keys_authentication_enabled = true
}

data "azurerm_managed_redis_database" "test" {
  cluster_id = azurerm_managed_redis_cluster.test.id
  depends_on = [azurerm_managed_redis_database.test]
}
`, data.RandomInteger, data.Locations.Primary)
}
