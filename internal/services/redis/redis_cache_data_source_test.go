// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package redis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type RedisCacheDataSource struct{}

func TestAccRedisCacheDataSource_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_redis_cache", "test")
	r := RedisCacheDataSource{}

	name := fmt.Sprintf("acctestRedis-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.standardWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("ssl_port").HasValue("6380"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("access_keys_authentication_enabled").Exists(),
				check.That(data.ResourceName).Key("identity.#").HasValue("0"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("redis_version").Exists(),
				check.That(data.ResourceName).Key("replicas_per_master").HasValue("0"),
				check.That(data.ResourceName).Key("replicas_per_primary").HasValue("0"),
				check.That(data.ResourceName).Key("tenant_settings.%").HasValue("0"),
			),
		},
	})
}

func TestAccRedisCacheDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_redis_cache", "test")
	r := RedisCacheDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.completeWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("redis_version").Exists(),
				check.That(data.ResourceName).Key("replicas_per_master").HasValue("3"),
				check.That(data.ResourceName).Key("replicas_per_primary").HasValue("3"),
				check.That(data.ResourceName).Key("tenant_settings.%").HasValue("1"),
				check.That(data.ResourceName).Key("tenant_settings.config").HasValue("config"),
			),
		},
	})
}

func (r RedisCacheDataSource) standardWithDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_redis_cache" "test" {
  name                = azurerm_redis_cache.test.name
  resource_group_name = azurerm_redis_cache.test.resource_group_name
}
`, RedisCacheResource{}.standard(data))
}

func (r RedisCacheDataSource) completeWithDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-redis-%d"
  location = "%s"
}

resource "azurerm_redis_cache" "test" {
  name                          = "acctestRedis-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  capacity                      = 3
  family                        = "P"
  sku_name                      = "Premium"
  non_ssl_port_enabled          = false
  public_network_access_enabled = false
  replicas_per_master           = 3
  replicas_per_primary          = 3
  redis_version                 = "6"
  tenant_settings = {
    config = "config"
  }

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_redis_cache" "test" {
  name                = azurerm_redis_cache.test.name
  resource_group_name = azurerm_redis_cache.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
