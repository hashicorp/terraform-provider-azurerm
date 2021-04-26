package redis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type RedisCacheDataSource struct {
}

func TestAccRedisCacheDataSource_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_redis_cache", "test")
	r := RedisCacheDataSource{}

	name := fmt.Sprintf("acctestRedis-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.standardWithDataSource(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("ssl_port").HasValue("6380"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
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
