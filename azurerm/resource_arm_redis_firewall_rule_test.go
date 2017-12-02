package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAzureRMRedisFirewallRuleName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "ab",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "webapp1",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 1,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloworld21!",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateRedisFirewallRuleName(tc.Value, "azurerm_redis_firewall_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Redis Firewall Rule Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAccAzureRMRedisFirewallRule_basic(t *testing.T) {
	resourceName := "azurerm_redis_firewall_rule.test"
	ri := acctest.RandInt()
	config := testAccAzureRMRedisFirewallRule_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMRedisFirewallRule_update(t *testing.T) {
	resourceName := "azurerm_redis_firewall_rule.test"
	ri := acctest.RandInt()
	config := testAccAzureRMRedisFirewallRule_basic(ri, testLocation())
	updatedConfig := testAccAzureRMRedisFirewallRule_update(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMRedisFirewallRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		name := rs.Primary.Attributes["name"]
		cacheName := rs.Primary.Attributes["redis_cache_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).redisFirewallClient
		resp, err := client.Get(resourceGroup, cacheName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Firewall Rule %q (cache %q resource group: %q) does not exist", name, cacheName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on redisFirewallClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRedisFirewallRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).resourceGroupClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_redis_firewall_rule" {
			continue
		}

		resourceGroup := rs.Primary.ID

		resp, err := client.Get(resourceGroup)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Firewall Rule still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMRedisFirewallRule_basic(rInt int, location string) string {
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
      maxclients         = 256,
      maxmemory_reserved = 2,
      maxmemory_delta    = 2
      maxmemory_policy   = "allkeys-lru"
    }
}

resource "azurerm_redis_firewall_rule" "test" {
  name                = "fwrule%d"
  redis_cache_name    = "${azurerm_redis_cache.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  start_ip            = "1.2.3.4"
  end_ip              = "2.3.4.5"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRedisFirewallRule_update(rInt int, location string) string {
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
      maxclients         = 256,
      maxmemory_reserved = 2,
      maxmemory_delta    = 2
      maxmemory_policy   = "allkeys-lru"
    }
}

resource "azurerm_redis_firewall_rule" "test" {
  name                = "fwrule%d"
  redis_cache_name    = "${azurerm_redis_cache.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  start_ip            = "2.3.4.5"
  end_ip              = "6.7.8.9"
}
`, rInt, location, rInt, rInt)
}
