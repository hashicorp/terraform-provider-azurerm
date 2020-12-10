package redis_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMRedisCache_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCache_basic(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "minimum_tls_version"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
					testCheckSSLInConnectionString(data.ResourceName, "primary_connection_string", true),
					testCheckSSLInConnectionString(data.ResourceName, "secondary_connection_string", true),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisCache_withoutSSL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCache_basic(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
					testCheckSSLInConnectionString(data.ResourceName, "primary_connection_string", false),
					testCheckSSLInConnectionString(data.ResourceName, "secondary_connection_string", false),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisCache_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCache_basic(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMRedisCache_requiresImport),
		},
	})
}

func TestAccAzureRMRedisCache_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCache_standard(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisCache_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCache_premium(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisCache_premiumSharded(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCache_premiumSharded(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisCache_premiumShardedScaling(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCache_premiumSharded(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMRedisCache_premiumShardedScaled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_NonStandardCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCacheNonStandardCasing(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
			{
				Config:             testAccAzureRMRedisCacheNonStandardCasing(data),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAzureRMRedisCache_BackupDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCacheBackupDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_BackupEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCacheBackupEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
				// `redis_configuration.0.rdb_storage_connection_string` is returned as:
				// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf"
				// TODO: remove this once the Bug's been fixed:
				// https://github.com/Azure/azure-rest-api-specs/issues/3037
				ExpectNonEmptyPlan: true,
			},
			data.ImportStep("redis_configuration.0.rdb_storage_connection_string"),
		},
	})
}

func TestAccAzureRMRedisCache_BackupEnabledDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCacheBackupEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
				// `redis_configuration.0.rdb_storage_connection_string` is returned as:
				// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf"
				// TODO: remove this once the Bug's been fixed:
				// https://github.com/Azure/azure-rest-api-specs/issues/3037
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAzureRMRedisCacheBackupDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
				// `redis_configuration.0.rdb_storage_connection_string` is returned as:
				// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf"
				// TODO: remove this once the Bug's been fixed:
				// https://github.com/Azure/azure-rest-api-specs/issues/3037
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMRedisCache_AOFBackupEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCacheAOFBackupEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
			data.ImportStep("redis_configuration.0.aof_storage_connection_string_0",
				"redis_configuration.0.aof_storage_connection_string_1"),
		},
	})
}

func TestAccAzureRMRedisCache_AOFBackupEnabledDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCacheAOFBackupEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAzureRMRedisCacheAOFBackupDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMRedisCache_PatchSchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCachePatchSchedule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisCache_PatchScheduleUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCachePatchSchedule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMRedisCache_premium(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_InternalSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCache_internalSubnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisCache_InternalSubnetStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCache_internalSubnetStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisCache_InternalSubnet_withZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCache_internalSubnet_withZone(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.0", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisCache_SubscribeAllEvents(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCacheSubscribeAllEvents(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_WithoutAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_cache", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisCacheWithoutAuth(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "redis_configuration.0.enable_authentication", "false"),
				),
			},
		},
	})
}

func testCheckAzureRMRedisCacheExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Redis.Client
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		redisName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Redis Instance: %s", redisName)
		}

		resp, err := conn.Get(ctx, resourceGroup, redisName)
		if err != nil {
			return fmt.Errorf("Bad: Get on redis.Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Redis Instance %q (resource group: %q) does not exist", redisName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMRedisCacheDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Redis.Client
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_redis_cache" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Redis Instance still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMRedisCache_basic(data acceptance.TestData, requireSSL bool) string {
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

func testAccAzureRMRedisCache_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMRedisCache_basic(data, true)
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

func testAccAzureRMRedisCache_standard(data acceptance.TestData) string {
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

func testAccAzureRMRedisCache_premium(data acceptance.TestData) string {
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
    maxmemory_reserved              = 2
    maxfragmentationmemory_reserved = 2
    maxmemory_delta                 = 2
    maxmemory_policy                = "allkeys-lru"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRedisCache_premiumSharded(data acceptance.TestData) string {
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
    maxmemory_reserved              = 2
    maxfragmentationmemory_reserved = 2
    maxmemory_delta                 = 2
    maxmemory_policy                = "allkeys-lru"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRedisCache_premiumShardedScaled(data acceptance.TestData) string {
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
    maxmemory_reserved              = 2
    maxfragmentationmemory_reserved = 2
    maxmemory_delta                 = 2
    maxmemory_policy                = "allkeys-lru"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRedisCacheNonStandardCasing(data acceptance.TestData) string {
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
  family              = "c"
  sku_name            = "basic"
  enable_non_ssl_port = false
  redis_configuration {
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRedisCacheBackupDisabled(data acceptance.TestData) string {
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
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    rdb_backup_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRedisCacheBackupEnabled(data acceptance.TestData) string {
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
    rdb_storage_connection_string = "DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.test.primary_blob_endpoint};AccountName=${azurerm_storage_account.test.name};AccountKey=${azurerm_storage_account.test.primary_access_key}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMRedisCacheAOFBackupDisabled(data acceptance.TestData) string {
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
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    aof_backup_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRedisCacheAOFBackupEnabled(data acceptance.TestData) string {
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
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    aof_backup_enabled              = true
    aof_storage_connection_string_0 = "DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.test.primary_blob_endpoint};AccountName=${azurerm_storage_account.test.name};AccountKey=${azurerm_storage_account.test.primary_access_key}"
    aof_storage_connection_string_1 = "DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.test.primary_blob_endpoint};AccountName=${azurerm_storage_account.test.name};AccountKey=${azurerm_storage_account.test.secondary_access_key}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMRedisCachePatchSchedule(data acceptance.TestData) string {
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
    maxmemory_reserved = 2
    maxmemory_delta    = 2
    maxmemory_policy   = "allkeys-lru"
  }

  patch_schedule {
    day_of_week    = "Tuesday"
    start_hour_utc = 8
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRedisCacheSubscribeAllEvents(data acceptance.TestData) string {
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

func testAccAzureRMRedisCache_internalSubnet(data acceptance.TestData) string {
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
  address_prefix       = "10.0.1.0/24"
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

func testAccAzureRMRedisCache_internalSubnetStaticIP(data acceptance.TestData) string {
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
  address_prefix       = "10.0.1.0/24"
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

func testAccAzureRMRedisCache_internalSubnet_withZone(data acceptance.TestData) string {
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
  address_prefix       = "10.0.1.0/24"
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

func testAccAzureRMRedisCacheWithoutAuth(data acceptance.TestData) string {
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
  address_prefix       = "10.0.1.0/24"
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

func testCheckSSLInConnectionString(resourceName string, propertyName string, requireSSL bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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
