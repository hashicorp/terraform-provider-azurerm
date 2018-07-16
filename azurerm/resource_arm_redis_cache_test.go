package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMRedisCacheFamily_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "C",
			ErrCount: 0,
		},
		{
			Value:    "P",
			ErrCount: 0,
		},
		{
			Value:    "c",
			ErrCount: 0,
		},
		{
			Value:    "p",
			ErrCount: 0,
		},
		{
			Value:    "a",
			ErrCount: 1,
		},
		{
			Value:    "b",
			ErrCount: 1,
		},
		{
			Value:    "D",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateRedisFamily(tc.Value, "azurerm_redis_cache")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Redis Cache Family to trigger a validation error")
		}
	}
}

func TestAccAzureRMRedisCacheMaxMemoryPolicy_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{Value: "noeviction", ErrCount: 0},
		{Value: "allkeys-lru", ErrCount: 0},
		{Value: "volatile-lru", ErrCount: 0},
		{Value: "allkeys-random", ErrCount: 0},
		{Value: "volatile-random", ErrCount: 0},
		{Value: "volatile-ttl", ErrCount: 0},
		{Value: "something-else", ErrCount: 1},
	}

	for _, tc := range cases {
		_, errors := validateRedisMaxMemoryPolicy(tc.Value, "azurerm_redis_cache")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Redis Cache Max Memory Policy to trigger a validation error")
		}
	}
}

func TestAccAzureRMRedisCacheBackupFrequency_validation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{Value: 1, ErrCount: 1},
		{Value: 15, ErrCount: 0},
		{Value: 30, ErrCount: 0},
		{Value: 45, ErrCount: 1},
		{Value: 60, ErrCount: 0},
		{Value: 120, ErrCount: 1},
		{Value: 240, ErrCount: 1},
		{Value: 360, ErrCount: 0},
		{Value: 720, ErrCount: 0},
		{Value: 1440, ErrCount: 0},
	}

	for _, tc := range cases {
		_, errors := validateRedisBackupFrequency(tc.Value, "azurerm_redis_cache")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the AzureRM Redis Cache Backup Frequency to trigger a validation error for '%d'", tc.Value)
		}
	}
}

func TestAccAzureRMRedisCache_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_standard(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_standard(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_premium(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_premium(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_premiumSharded(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_premiumSharded(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_NonStandardCasing(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRedisCacheNonStandardCasing(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},

			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAzureRMRedisCache_BackupDisabled(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRedisCacheBackupDisabled(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_BackupEnabled(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMRedisCacheBackupEnabled(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
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

func TestAccAzureRMRedisCache_BackupEnabledDisabled(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccAzureRMRedisCacheBackupEnabled(ri, rs, location)
	updatedConfig := testAccAzureRMRedisCacheBackupDisabled(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
				// `redis_configuration.0.rdb_storage_connection_string` is returned as:
				// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf"
				// TODO: remove this once the Bug's been fixed:
				// https://github.com/Azure/azure-rest-api-specs/issues/3037
				ExpectNonEmptyPlan: true,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
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

func TestAccAzureRMRedisCache_PatchSchedule(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMRedisCachePatchSchedule(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_PatchScheduleUpdated(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMRedisCachePatchSchedule(ri, location)
	updatedConfig := testAccAzureRMRedisCache_premium(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_InternalSubnet(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_internalSubnet(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRedisCache_InternalSubnetStaticIP(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_internalSubnetStaticIP(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
		},
	})
}

func testCheckAzureRMRedisCacheExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		redisName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Redis Instance: %s", redisName)
		}

		conn := testAccProvider.Meta().(*ArmClient).redisClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, redisName)
		if err != nil {
			return fmt.Errorf("Bad: Get on redisClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Redis Instance %q (resource group: %q) does not exist", redisName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMRedisCacheDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).redisClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func TestAccAzureRMRedisCache_SubscribeAllEvents(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMRedisCacheSubscribeAllEvents(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisCacheExists("azurerm_redis_cache.test"),
				),
			},
		},
	})
}

func testAccAzureRMRedisCache_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_redis_cache" "test" {
    name                = "acctestRedis-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    capacity            = 1
    family              = "C"
    sku_name            = "Basic"
    enable_non_ssl_port = false

    redis_configuration {
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMRedisCache_standard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_redis_cache" "test" {
    name                = "acctestRedis-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    capacity            = 1
    family              = "C"
    sku_name            = "Standard"
    enable_non_ssl_port = false
    redis_configuration {
    }

    tags {
    	environment = "production"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMRedisCache_premium(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_redis_cache" "test" {
    name                = "acctestRedis-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    capacity            = 1
    family              = "P"
    sku_name            = "Premium"
    enable_non_ssl_port = false
    redis_configuration {
      maxmemory_reserved = 2
      maxmemory_delta    = 2
      maxmemory_policy   = "allkeys-lru"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMRedisCache_premiumSharded(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_redis_cache" "test" {
    name                = "acctestRedis-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    capacity            = 1
    family              = "P"
    sku_name            = "Premium"
    enable_non_ssl_port = true
    shard_count         = 3
    redis_configuration {
      maxmemory_reserved = 2
      maxmemory_delta    = 2
      maxmemory_policy   = "allkeys-lru"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMRedisCacheNonStandardCasing(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_redis_cache" "test" {
    name                = "acctestRedis-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    capacity            = 1
    family              = "c"
    sku_name            = "basic"
    enable_non_ssl_port = false
    redis_configuration {
    }
}
`, ri, location, ri)
}

func testAccAzureRMRedisCacheBackupDisabled(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_redis_cache" "test" {
    name                = "acctestRedis-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    capacity            = 3
    family              = "P"
    sku_name            = "Premium"
    enable_non_ssl_port = false
    redis_configuration {
      rdb_backup_enabled = false
    }
}
`, ri, location, ri)
}

func testAccAzureRMRedisCacheBackupEnabled(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags {
    environment = "staging"
  }
}

resource "azurerm_redis_cache" "test" {
    name                = "acctestRedis-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
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
`, rInt, location, rString, rInt)
}

func testAccAzureRMRedisCachePatchSchedule(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}
resource "azurerm_redis_cache" "test" {
    name                = "acctestRedis-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
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
`, rInt, location, rInt)
}

func testAccAzureRMRedisCacheSubscribeAllEvents(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}
resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location     = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
  tags {
    environment = "staging"
  }
}
resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  redis_configuration {
    notify_keyspace_events = "KAE"
  }
}
`, rInt, location, rString, rInt)
}
func testAccAzureRMRedisCache_internalSubnet(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestRedis-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  subnet_id           = "${azurerm_subnet.test.id}"
  redis_configuration {
  }
}
`, ri, location, ri, ri)
}

func testAccAzureRMRedisCache_internalSubnetStaticIP(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_redis_cache" "test" {
  name                      = "acctestRedis-%d"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  capacity                  = 1
  family                    = "P"
  sku_name                  = "Premium"
  enable_non_ssl_port       = false
  subnet_id                 = "${azurerm_subnet.test.id}"
  private_static_ip_address = "10.0.1.20"
  redis_configuration {
  }
}
`, ri, location, ri, ri)
}
