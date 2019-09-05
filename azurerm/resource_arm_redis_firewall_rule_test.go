package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
			ErrCount: 0,
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
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisFirewallRule_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRedisFirewallRule_multi(t *testing.T) {
	ruleOne := "azurerm_redis_firewall_rule.test"
	ruleTwo := "azurerm_redis_firewall_rule.double"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisFirewallRule_multi(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(ruleOne),
					testCheckAzureRMRedisFirewallRuleExists(ruleTwo),
				),
			},
			{
				ResourceName:      ruleOne,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      ruleTwo,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRedisFirewallRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_redis_firewall_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRedisFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRedisFirewallRule_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRedisFirewallRuleExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMRedisFirewallRule_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_redis_firewall_rule"),
			},
		},
	})
}

func TestAccAzureRMRedisFirewallRule_update(t *testing.T) {
	resourceName := "azurerm_redis_firewall_rule.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMRedisFirewallRule_basic(ri, testLocation())
	updatedConfig := testAccAzureRMRedisFirewallRule_update(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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

func testCheckAzureRMRedisFirewallRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		cacheName := rs.Primary.Attributes["redis_cache_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).redis.FirewallRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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
	client := testAccProvider.Meta().(*ArmClient).resource.GroupsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMRedisFirewallRule_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

resource "azurerm_redis_firewall_rule" "test" {
  name                = "fwrule%d"
  redis_cache_name    = "${azurerm_redis_cache.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  start_ip            = "1.2.3.4"
  end_ip              = "2.3.4.5"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRedisFirewallRule_multi(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redis_firewall_rule" "double" {
  name                = "fwruletwo%d"
  redis_cache_name    = "${azurerm_redis_cache.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  start_ip            = "4.5.6.7"
  end_ip              = "8.9.0.1"
}
`, testAccAzureRMRedisFirewallRule_basic(rInt, location), rInt)
}

func testAccAzureRMRedisFirewallRule_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redis_firewall_rule" "import" {
  name                = "${azurerm_redis_firewall_rule.test.name}"
  redis_cache_name    = "${azurerm_redis_firewall_rule.test.redis_cache_name}"
  resource_group_name = "${azurerm_redis_firewall_rule.test.resource_group_name}"
  start_ip            = "${azurerm_redis_firewall_rule.test.start_ip}"
  end_ip              = "${azurerm_redis_firewall_rule.test.end_ip}"
}
`, testAccAzureRMRedisFirewallRule_basic(rInt, location))
}

func testAccAzureRMRedisFirewallRule_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

resource "azurerm_redis_firewall_rule" "test" {
  name                = "fwrule%d"
  redis_cache_name    = "${azurerm_redis_cache.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  start_ip            = "2.3.4.5"
  end_ip              = "6.7.8.9"
}
`, rInt, location, rInt, rInt)
}
