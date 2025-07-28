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

type ManagedRedisDatabaseResource struct{}

func TestAccManagedRedisDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
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

func TestAccManagedRedisDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
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

func TestAccManagedRedisDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisDatabase_geoDatabase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.threeLinkedDatabases(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisDatabase_geoDatabaseOtherEvictionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoDatabaseOtherEvictionPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisDatabase_geoDatabaseModule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoDatabasewithModuleEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisDatabase_geoDatabaseWithRedisJsonModule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoDatabasewithRedisJsonModuleEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisDatabase_unlinkDatabase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.threeLinkedDatabases(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.twoLinkedDatabases(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisDatabase_linkDatabase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.twoLinkedDatabases(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.threeLinkedDatabases(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisDatabase_customizeDiffDatabaseIdError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.customizeDiffDatabaseIdError(data),
			ExpectError: regexp.MustCompile("linked database list must include the current database ID"),
		},
	})
}

func TestAccManagedRedisDatabase_customizeDiffModuleError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_database", "test")
	r := ManagedRedisDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.customizeDiffModuleError(),
			ExpectError: regexp.MustCompile("Only `RediSearch` and `RedisJSON` modules are allowed with geo-replication"),
		},
	})
}

func (r ManagedRedisDatabaseResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r ManagedRedisDatabaseResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_redis_cluster" "test" {
  name                = "acctest-rec1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Balanced_B3"
}`, data.RandomInteger, data.Locations.Primary)
}

func (r ManagedRedisDatabaseResource) templateThreeClusters(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
resource "azurerm_managed_redis_cluster" "test1" {
  name                = "acctest-rec2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Balanced_B3"
}
resource "azurerm_managed_redis_cluster" "test2" {
  name                = "acctest-rec3-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Balanced_B3"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ManagedRedisDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_redis_database" "test" {
  name       = "default"
  cluster_id = azurerm_managed_redis_cluster.test.id
}
`, r.template(data))
}

func (r ManagedRedisDatabaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_redis_database" "import" {
  name       = azurerm_managed_redis_database.test.name
  cluster_id = azurerm_managed_redis_database.test.cluster_id
}
`, r.basic(data))
}

func (r ManagedRedisDatabaseResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_redis_database" "test" {
  cluster_id = azurerm_managed_redis_cluster.test.id

  name                               = "default"
  access_keys_authentication_enabled = true
  client_protocol                    = "Encrypted"
  clustering_policy                  = "EnterpriseCluster"
  eviction_policy                    = "NoEviction"

  module {
    name = "RediSearch"
    args = ""
  }

  module {
    name = "RedisBloom"
    args = ""
  }

  module {
    name = "RedisTimeSeries"
    args = ""
  }

  module {
    name = "RedisJSON"
    args = ""
  }

  port = 10000
}
`, r.template(data))
}

func (r ManagedRedisDatabaseResource) threeLinkedDatabases(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_redis_database" "test" {
  cluster_id = azurerm_managed_redis_cluster.test.id

  client_protocol   = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy   = "NoEviction"

  linked_database_id = [
    "${azurerm_managed_redis_cluster.test.id}/databases/default",
    "${azurerm_managed_redis_cluster.test1.id}/databases/default",
    "${azurerm_managed_redis_cluster.test2.id}/databases/default"
  ]

  linked_database_group_nickname = "tftestGeoGroup"
}
`, r.templateThreeClusters(data))
}

func (r ManagedRedisDatabaseResource) geoDatabaseOtherEvictionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_redis_database" "test" {
  cluster_id = azurerm_managed_redis_cluster.test.id

  client_protocol   = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy   = "AllKeysLRU"

  linked_database_id = [
    "${azurerm_managed_redis_cluster.test.id}/databases/default",
    "${azurerm_managed_redis_cluster.test1.id}/databases/default",
    "${azurerm_managed_redis_cluster.test2.id}/databases/default"
  ]

  linked_database_group_nickname = "tftestGeoGroup"
}
`, r.templateThreeClusters(data))
}

func (r ManagedRedisDatabaseResource) twoLinkedDatabases(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_redis_database" "test" {
  cluster_id = azurerm_managed_redis_cluster.test.id

  client_protocol   = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy   = "NoEviction"

  linked_database_id = [
    "${azurerm_managed_redis_cluster.test.id}/databases/default",
    "${azurerm_managed_redis_cluster.test1.id}/databases/default",
  ]

  linked_database_group_nickname = "tftestGeoGroup"
}
`, r.templateThreeClusters(data))
}

func (r ManagedRedisDatabaseResource) geoDatabasewithModuleEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_redis_database" "test" {
  cluster_id = azurerm_managed_redis_cluster.test.id

  client_protocol   = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy   = "NoEviction"
  module {
    name = "RediSearch"
    args = ""
  }
  linked_database_id = [
    "${azurerm_managed_redis_cluster.test.id}/databases/default",
    "${azurerm_managed_redis_cluster.test1.id}/databases/default",
    "${azurerm_managed_redis_cluster.test2.id}/databases/default"
  ]

  linked_database_group_nickname = "tftestGeoGroup"
}
`, r.templateThreeClusters(data))
}

func (r ManagedRedisDatabaseResource) geoDatabasewithRedisJsonModuleEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_managed_redis_database" "test" {
  cluster_id = azurerm_managed_redis_cluster.test.id

  client_protocol   = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy   = "NoEviction"
  module {
    name = "RedisJSON"
    args = ""
  }
  linked_database_id = [
    "${azurerm_managed_redis_cluster.test.id}/databases/default",
    "${azurerm_managed_redis_cluster.test1.id}/databases/default",
    "${azurerm_managed_redis_cluster.test2.id}/databases/default"
  ]

  linked_database_group_nickname = "tftestGeoGroup"
}
`, r.templateThreeClusters(data))
}

func (r ManagedRedisDatabaseResource) customizeDiffDatabaseIdError(data acceptance.TestData) string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_managed_redis_database" "test" {
  name       = "default"
  cluster_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Cache/redisEnterprise/redisEnterprise1"

  linked_database_id = [
    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Cache/redisEnterprise/redisEnterprise2/databases/default",
  ]
}
`
}

func (r ManagedRedisDatabaseResource) customizeDiffModuleError() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_managed_redis_database" "test" {
  name       = "default"
  cluster_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Cache/redisEnterprise/redisEnterprise1"

	module {
    name = "RedisTimeSeries"
    args = ""
  }

  linked_database_id = [
    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Cache/redisEnterprise/redisEnterprise1/databases/default",
  ]
}
`
}
