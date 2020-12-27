package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type NetAppPoolDataSource struct {
}

func TestAccDataSourceNetAppPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_pool", "test")
	r := NetAppPoolDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("account_name").Exists(),
				check.That(data.ResourceName).Key("service_level").Exists(),
				check.That(data.ResourceName).Key("size_in_tb").Exists(),
			),
		},
	})
}

func (NetAppPoolDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_pool" "test" {
  resource_group_name = azurerm_netapp_pool.test.resource_group_name
  account_name        = azurerm_netapp_pool.test.account_name
  name                = azurerm_netapp_pool.test.name
}
`, NetAppPoolResource{}.basic(data))
}
