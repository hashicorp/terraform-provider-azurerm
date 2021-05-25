package netapp_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type NetAppAccountDataSource struct {
}

func testAccDataSourceNetAppAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_account", "test")
	r := NetAppAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (r NetAppAccountDataSource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_account" "test" {
  resource_group_name = azurerm_netapp_account.test.resource_group_name
  name                = azurerm_netapp_account.test.name
}
`, NetAppAccountResource{}.basicConfig(data))
}
