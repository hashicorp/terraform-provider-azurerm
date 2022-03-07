package redisenterprise_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type RedisEnterpriseGeoDatabaseDataSource struct{}

func TestAccRedisEnterpriseDatabaseDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_redis_enterprise_geo_database", "test")
	r := RedisEnterpriseGeoDatabaseDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("default"),
				check.That(data.ResourceName).Key("cluster_id").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
	})
}

func (r RedisEnterpriseGeoDatabaseDataSource) dataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_redis_enterprise_geo_database" "test" {
  depends_on = [azurerm_redis_enterprise_geo_database.test]
  name       = "default"
  cluster_id = azurerm_redis_enterprise_cluster.test.id
}


`, RedisenterpriseGeoDatabaseResource{}.basic(data))
}
