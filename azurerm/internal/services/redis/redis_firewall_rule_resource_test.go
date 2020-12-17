package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redis/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type RedisFirewallRuleResource struct {
}

func TestAccRedisFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_firewall_rule", "test")
	r := RedisFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRedisFirewallRule_multi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_firewall_rule", "test")
	r := RedisFirewallRuleResource{}
	ruleTwo := "azurerm_redis_firewall_rule.double"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multi(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(ruleTwo).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.ImportStepFor(ruleTwo),
	})
}

func TestAccRedisFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_firewall_rule", "test")
	r := RedisFirewallRuleResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccRedisFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_firewall_rule", "test")
	r := RedisFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t RedisFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.FirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Redis.FirewallRulesClient.Get(ctx, id.ResourceGroup, id.RediName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Redis Firewall Rule (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.FirewallRuleProperties != nil), nil
}

func (RedisFirewallRuleResource) basic(data acceptance.TestData) string {
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

func (r RedisFirewallRuleResource) multi(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redis_firewall_rule" "double" {
  name                = "fwruletwo%d"
  redis_cache_name    = azurerm_redis_cache.test.name
  resource_group_name = azurerm_resource_group.test.name
  start_ip            = "4.5.6.7"
  end_ip              = "8.9.0.1"
}
`, r.basic(data), data.RandomInteger)
}

func (r RedisFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redis_firewall_rule" "import" {
  name                = azurerm_redis_firewall_rule.test.name
  redis_cache_name    = azurerm_redis_firewall_rule.test.redis_cache_name
  resource_group_name = azurerm_redis_firewall_rule.test.resource_group_name
  start_ip            = azurerm_redis_firewall_rule.test.start_ip
  end_ip              = azurerm_redis_firewall_rule.test.end_ip
}
`, r.basic(data))
}

func (RedisFirewallRuleResource) update(data acceptance.TestData) string {
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
