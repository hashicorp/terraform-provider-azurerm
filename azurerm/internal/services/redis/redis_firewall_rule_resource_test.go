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

func TestAccAzureRMRedisFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_firewall_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisFirewallRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRedisFirewallRule_multi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_firewall_rule", "test")
	ruleTwo := "azurerm_redis_firewall_rule.double"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisFirewallRule_multi(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(data.ResourceName),
					testCheckAzureRMRedisFirewallRuleExists(ruleTwo),
				),
			},
			data.ImportStep(),
			{
				ResourceName:      ruleTwo,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRedisFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_firewall_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisFirewallRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMRedisFirewallRule_requiresImport),
		},
	})
}

func TestAccAzureRMRedisFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_firewall_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisFirewallRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMRedisFirewallRule_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(data.ResourceName),
				),
			},
		},
	})
}

func testCheckAzureRMRedisFirewallRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Redis.FirewallRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		cacheName := rs.Primary.Attributes["redis_cache_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, cacheName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Firewall Rule %q (cache %q resource group: %q) does not exist", name, cacheName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on redis.FirewallRulesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRedisFirewallRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.GroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_redis_firewall_rule" {
			continue
		}

		resourceGroup := rs.Primary.ID

		resp, err := client.Get(ctx, resourceGroup)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Firewall Rule still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMRedisFirewallRule_basic(data acceptance.TestData) string {
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
}

resource "azurerm_redis_firewall_rule" "test" {
  name                = "fwrule%d"
  redis_cache_name    = azurerm_redis_cache.test.name
  resource_group_name = azurerm_resource_group.test.name
  start_ip            = "1.2.3.4"
  end_ip              = "2.3.4.5"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMRedisFirewallRule_multi(data acceptance.TestData) string {
	template := testAccAzureRMRedisFirewallRule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_firewall_rule" "double" {
  name                = "fwruletwo%d"
  redis_cache_name    = azurerm_redis_cache.test.name
  resource_group_name = azurerm_resource_group.test.name
  start_ip            = "4.5.6.7"
  end_ip              = "8.9.0.1"
}
`, template, data.RandomInteger)
}

func testAccAzureRMRedisFirewallRule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMRedisFirewallRule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_firewall_rule" "import" {
  name                = azurerm_redis_firewall_rule.test.name
  redis_cache_name    = azurerm_redis_firewall_rule.test.redis_cache_name
  resource_group_name = azurerm_redis_firewall_rule.test.resource_group_name
  start_ip            = azurerm_redis_firewall_rule.test.start_ip
  end_ip              = azurerm_redis_firewall_rule.test.end_ip
}
`, template)
}

func testAccAzureRMRedisFirewallRule_update(data acceptance.TestData) string {
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
}

resource "azurerm_redis_firewall_rule" "test" {
  name                = "fwrule%d"
  redis_cache_name    = azurerm_redis_cache.test.name
  resource_group_name = azurerm_resource_group.test.name
  start_ip            = "2.3.4.5"
  end_ip              = "6.7.8.9"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
