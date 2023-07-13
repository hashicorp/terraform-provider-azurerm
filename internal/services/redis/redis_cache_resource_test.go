// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redis_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2023-04-01/redis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RedisCacheResource struct{}

func TestAccRedisCache_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("minimum_tls_version").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				testCheckSSLInConnectionString(data.ResourceName, "primary_connection_string", true),
				testCheckSSLInConnectionString(data.ResourceName, "secondary_connection_string", true),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_withoutSSL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckSSLInConnectionString(data.ResourceName, "primary_connection_string", false),
				testCheckSSLInConnectionString(data.ResourceName, "secondary_connection_string", false),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccRedisCache_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_premiumSharded(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumSharded(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_premiumShardedScaling(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumSharded(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.premiumShardedScaled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRedisCache_BackupDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backupDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			// `redis_configuration.0.aof_storage_connection_string_0` and `redis_configuration.0.aof_storage_connection_string_1` are returned as:
			// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf"
			// TODO: remove this once the Bug's been fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/3037
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccRedisCache_BackupEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backupEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			// `redis_configuration.0.rdb_storage_connection_string` is returned as:
			// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf"
			// TODO: remove this once the Bug's been fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/3037
			ExpectNonEmptyPlan: true,
		},
		data.ImportStep("redis_configuration.0.rdb_storage_connection_string"),
	})
}

func TestAccRedisCache_BackupEnabledDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backupEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			// `redis_configuration.0.rdb_storage_connection_string` is returned as:
			// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf..."
			// TODO: remove this once the Bug's been fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/3037
			ExpectNonEmptyPlan: true,
		},
		{
			Config: r.backupDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			// `redis_configuration.0.aof_storage_connection_string_0` and `redis_configuration.0.aof_storage_connection_string_1` are returned as:
			// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf..."
			// TODO: remove this once the Bug's been fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/3037
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccRedisCache_AOFBackupEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aofBackupEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectNonEmptyPlan: true,
		},
		data.ImportStep("redis_configuration.0.aof_storage_connection_string_0",
			"redis_configuration.0.aof_storage_connection_string_1"),
	})
}

func TestAccRedisCache_AOFBackupEnabledDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aofBackupEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			// `redis_configuration.0.aof_storage_connection_string_0` and `aof_storage_connection_string_1` are returned as:
			// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf..."
			// TODO: remove this once the Bug's been fixed:
			ExpectNonEmptyPlan: true,
		},
		{
			Config: r.aofBackupDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			// `redis_configuration.0.rdb_storage_connection_string` is returned as:
			// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf..."
			// TODO: remove this once the Bug's been fixed:
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccRedisCache_PatchSchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.patchSchedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_PatchScheduleUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.patchSchedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.premium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRedisCache_PublicNetworkAccessDisabledEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicNetworkAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_InternalSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.internalSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_InternalSubnetStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.internalSubnetStaticIP(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_InternalSubnet_withZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.internalSubnet_withZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
				check.That(data.ResourceName).Key("zones.0").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_SubscribeAllEvents(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.subscribeAllEvents(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRedisCache_WithoutAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withoutAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("redis_configuration.0.enable_authentication").HasValue("false"),
			),
		},
	})
}

func TestAccRedisCache_ReplicasPerMaster(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.replicasPerMaster(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRedisCache_ReplicasPerPrimary(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.replicasPerPrimary(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRedisCache_RedisVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.redisVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRedisCache_TenantSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tenantSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRedisCache_redisConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.redisConfiguration(data, "volatile-lru"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.redisConfiguration(data, "allkeys-lru"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisCache_SkuDowngrade(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")
	r := RedisCacheResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.SkuDowngrade(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t RedisCacheResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := redis.ParseRediID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Redis.Redis.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (RedisCacheResource) basic(data acceptance.TestData, requireSSL bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "C"
  sku_name            = "Basic"
  enable_non_ssl_port = %t
  minimum_tls_version = "1.2"

  redis_configuration {
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, !requireSSL)
}

func (RedisCacheResource) requiresImport(data acceptance.TestData) string {
	template := RedisCacheResource{}.basic(data, true)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_cache" "import" {
  name                = azurerm_redis_cache.test.name
  location            = azurerm_redis_cache.test.location
  resource_group_name = azurerm_redis_cache.test.resource_group_name
  capacity            = azurerm_redis_cache.test.capacity
  family              = azurerm_redis_cache.test.family
  sku_name            = azurerm_redis_cache.test.sku_name
  enable_non_ssl_port = azurerm_redis_cache.test.enable_non_ssl_port

  redis_configuration {
  }
}
`, template)
}

func (RedisCacheResource) standard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "C"
  sku_name            = "Standard"
  enable_non_ssl_port = false
  redis_configuration {
  }

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) premium(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    maxmemory_reserved              = 642
    maxfragmentationmemory_reserved = 642
    maxmemory_delta                 = 642
    maxmemory_policy                = "allkeys-lru"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) premiumSharded(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = true
  shard_count         = 3

  redis_configuration {
    maxmemory_reserved              = 642
    maxfragmentationmemory_reserved = 642
    maxmemory_delta                 = 642
    maxmemory_policy                = "allkeys-lru"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) premiumShardedScaled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 2
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = true
  shard_count         = 3

  redis_configuration {
    maxmemory_reserved              = 1328
    maxfragmentationmemory_reserved = 1328
    maxmemory_delta                 = 1328
    maxmemory_policy                = "allkeys-lru"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) backupDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account" "test3" {
  name                     = "acctestsa3%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    rdb_backup_enabled              = false
    aof_backup_enabled              = true
    aof_storage_connection_string_0 = azurerm_storage_account.test2.primary_connection_string
    aof_storage_connection_string_1 = azurerm_storage_account.test3.primary_connection_string
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString, data.RandomInteger)
}

func (RedisCacheResource) backupEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    rdb_backup_enabled            = true
    rdb_backup_frequency          = 60
    rdb_backup_max_snapshot_count = 1
    rdb_storage_connection_string = azurerm_storage_account.test.primary_connection_string
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (RedisCacheResource) aofBackupDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account" "test3" {
  name                     = "acctestsa3%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    aof_backup_enabled            = false
    rdb_backup_enabled            = true
    rdb_backup_frequency          = 60
    rdb_backup_max_snapshot_count = 1
    rdb_storage_connection_string = azurerm_storage_account.test3.primary_connection_string
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString, data.RandomInteger)
}

func (RedisCacheResource) aofBackupEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    aof_backup_enabled              = true
    aof_storage_connection_string_0 = azurerm_storage_account.test.primary_connection_string
    aof_storage_connection_string_1 = azurerm_storage_account.test2.primary_connection_string
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger)
}

func (RedisCacheResource) patchSchedule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    maxmemory_reserved = 642
    maxmemory_delta    = 642
    maxmemory_policy   = "allkeys-lru"
  }

  patch_schedule {
    day_of_week        = "Tuesday"
    start_hour_utc     = 8
    maintenance_window = "PT7H"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) publicNetworkAccessDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                          = "acctestRedis-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  capacity                      = 1
  family                        = "C"
  sku_name                      = "Basic"
  minimum_tls_version           = "1.2"
  enable_non_ssl_port           = false
  public_network_access_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) subscribeAllEvents(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    notify_keyspace_events = "KAE"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (RedisCacheResource) internalSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  subnet_id           = azurerm_subnet.test.id
  redis_configuration {
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (RedisCacheResource) internalSubnetStaticIP(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_redis_cache" "test" {
  name                      = "acctestRedis-%d"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  capacity                  = 1
  family                    = "P"
  sku_name                  = "Premium"
  enable_non_ssl_port       = false
  subnet_id                 = azurerm_subnet.test.id
  private_static_ip_address = "10.0.1.20"
  redis_configuration {
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (RedisCacheResource) internalSubnet_withZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  subnet_id           = azurerm_subnet.test.id
  redis_configuration {
  }
  zones = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (RedisCacheResource) withoutAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  subnet_id           = azurerm_subnet.test.id
  redis_configuration {
    enable_authentication = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (RedisCacheResource) replicasPerMaster(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-redis-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  replicas_per_master = 3
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) replicasPerPrimary(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-redis-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                 = "acctestRedis-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  capacity             = 3
  family               = "P"
  sku_name             = "Premium"
  enable_non_ssl_port  = false
  replicas_per_primary = 3
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) redisVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-redis-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  redis_version       = "6"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) tenantSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-redis-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  tenant_settings = {
    config = "config"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) redisConfiguration(data acceptance.TestData, maxMemoryPolicy string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-redis-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 2
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  minimum_tls_version = "1.2"

  redis_configuration {
    maxmemory_policy = "%s"
  }

  patch_schedule {
    day_of_week        = "Tuesday"
    start_hour_utc     = 4
    maintenance_window = "PT5H"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, maxMemoryPolicy)
}

func (RedisCacheResource) systemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "C"
  sku_name            = "Standard"
  enable_non_ssl_port = false
  redis_configuration {
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RedisCacheResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "C"
  sku_name            = "Standard"
  enable_non_ssl_port = false
  redis_configuration {
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (RedisCacheResource) SkuDowngrade(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 1
  family              = "C"
  sku_name            = "Basic"
  enable_non_ssl_port = false
  redis_configuration {
  }

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testCheckSSLInConnectionString(resourceName string, propertyName string, requireSSL bool) acceptance.TestCheckFunc {
	return func(s *acceptance.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		connectionString := rs.Primary.Attributes[propertyName]
		if strings.Contains(connectionString, fmt.Sprintf("ssl=%t", requireSSL)) {
			return nil
		}
		if strings.Contains(connectionString, fmt.Sprintf("ssl=%t", !requireSSL)) {
			return fmt.Errorf("Bad: wrong SSL setting in connection string: %s", propertyName)
		}

		return fmt.Errorf("Bad: missing SSL setting in connection string: %s", propertyName)
	}
}
