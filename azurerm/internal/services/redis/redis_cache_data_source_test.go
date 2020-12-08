package redis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMRedisCache_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_redis_cache", "test")

	name := fmt.Sprintf("acctestRedis-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRedisCacheDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRedisCache_standardWithDataSource(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", name),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_port", "6380"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMRedisCache_standardWithDataSource(data acceptance.TestData) string {
	config := testAccAzureRMRedisCache_standard(data)
	return fmt.Sprintf(`
%s

data "azurerm_redis_cache" "test" {
  name                = azurerm_redis_cache.test.name
  resource_group_name = azurerm_redis_cache.test.resource_group_name
}
`, config)
}
