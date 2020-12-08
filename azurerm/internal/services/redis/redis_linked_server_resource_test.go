package redis_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRedisLinkedServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_linked_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisLinkedServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisLinkedServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisLinkedServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisLinkedServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_linked_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisLinkedServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisLinkedServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisLinkedServerExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMRedisLinkedServer_requiresImport),
		},
	})
}

func testCheckAzureRMRedisLinkedServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Redis.LinkedServerClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		cacheName := rs.Primary.Attributes["target_redis_cache_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, cacheName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Linked Server %q (cache %q resource group: %q) does not exist", name, cacheName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on redis.LinkedServersClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRedisLinkedServerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Redis.LinkedServerClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_redis_linked_server" {
			continue
		}

		redisCacheName := rs.Primary.Attributes["target_redis_cache_name"]
		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, redisCacheName, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Linked Server still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMRedisLinkedServer_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "pri" {
  name     = "acctestRG-redis-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "pri" {
  name                = "acctestRedispri%d"
  location            = azurerm_resource_group.pri.location
  resource_group_name = azurerm_resource_group.pri.name
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

resource "azurerm_resource_group" "sec" {
  name     = "accsecRG-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "sec" {
  name                = "acctestRedissec%d"
  location            = azurerm_resource_group.sec.location
  resource_group_name = azurerm_resource_group.sec.name
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

resource "azurerm_redis_linked_server" "test" {
  target_redis_cache_name     = azurerm_redis_cache.pri.name
  resource_group_name         = azurerm_redis_cache.pri.resource_group_name
  linked_redis_cache_id       = azurerm_redis_cache.sec.id
  linked_redis_cache_location = azurerm_redis_cache.sec.location
  server_role                 = "Secondary"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger,
		data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

func testAccAzureRMRedisLinkedServer_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMRedisLinkedServer_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_linked_server" "import" {
  target_redis_cache_name     = azurerm_redis_linked_server.test.target_redis_cache_name
  resource_group_name         = azurerm_redis_linked_server.test.resource_group_name
  linked_redis_cache_id       = azurerm_redis_linked_server.test.linked_redis_cache_id
  linked_redis_cache_location = azurerm_redis_linked_server.test.linked_redis_cache_location
  server_role                 = azurerm_redis_linked_server.test.server_role
}
`, template)
}
