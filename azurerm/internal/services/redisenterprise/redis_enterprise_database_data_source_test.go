package redisenterprise_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type RedisEnterpriseDatabaseDataSource struct {
}

func TestAccRedisEnterpriseDatabaseDataSource_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_redis_enterprise_database", "test")
	r := RedisEnterpriseDatabaseDataSource{}

	resourceGroupName := fmt.Sprintf("acctestRG-redisEnterprise-%d", data.RandomInteger)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("default"),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("cluster_id").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
	})
}

func (r RedisEnterpriseDatabaseDataSource) dataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_redis_enterprise_database" "test" {
  depends_on = [azurerm_redis_enterprise_database.test]

  name                = "default"
  resource_group_name = azurerm_resource_group.test.name
  cluster_id          = azurerm_redis_enterprise_cluster.test.id
}
`, RedisenterpriseDatabaseResource{}.basic(data))
}
