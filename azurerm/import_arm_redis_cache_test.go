package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMRedisCache_importBasic(t *testing.T) {
	resourceName := "azurerm_redis_cache.test"

	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRedisCache_importStandard(t *testing.T) {
	resourceName := "azurerm_redis_cache.test"

	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_standard(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRedisCache_importPremium(t *testing.T) {
	resourceName := "azurerm_redis_cache.test"

	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_premium(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRedisCache_importPremiumSharded(t *testing.T) {
	resourceName := "azurerm_redis_cache.test"

	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_premiumSharded(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRedisCache_importNonStandardCasing(t *testing.T) {
	resourceName := "azurerm_redis_cache.test"

	ri := acctest.RandInt()
	config := testAccAzureRMRedisCacheNonStandardCasing(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRedisCache_importBackupEnabled(t *testing.T) {
	resourceName := "azurerm_redis_cache.test"

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
				// `redis_configuration.0.rdb_storage_connection_string` is returned as:
				// "...;AccountKey=[key hidden]" rather than "...;AccountKey=fsjfvjnfnf"
				// TODO: remove this once the Bug's been fixed:
				// https://github.com/Azure/azure-rest-api-specs/issues/3037
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"redis_configuration.0.rdb_storage_connection_string"},
			},
		},
	})
}

func TestAccAzureRMRedisCache_importPatchSchedule(t *testing.T) {
	resourceName := "azurerm_redis_cache.test"
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccAzureRMRedisCache_importInternalSubnet(t *testing.T) {
	resourceName := "azurerm_redis_cache.test"
	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_internalSubnet(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRedisCache_importInternalSubnetStaticIP(t *testing.T) {
	resourceName := "azurerm_redis_cache.test"
	ri := acctest.RandInt()
	config := testAccAzureRMRedisCache_internalSubnetStaticIP(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
