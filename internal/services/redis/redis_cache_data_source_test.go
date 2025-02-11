// Copyright (c) HashiCorp, Inc.
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
