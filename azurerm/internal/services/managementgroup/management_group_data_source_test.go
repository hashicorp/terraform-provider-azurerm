package managementgroup

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ManagementGroupDataSource struct {
}

func TestAccManagementGroupDataSource_basicByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_management_group", "test")
	r := ManagementGroupDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicByName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestmg-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("subscription_ids.#").HasValue("0"),
			),
		},
	})
}

func TestAccManagementGroupDataSource_basicByDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_management_group", "test")
	r := ManagementGroupDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicByDisplayName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest Management Group %d", data.RandomInteger)),
				check.That(data.ResourceName).Key("subscription_ids.#").HasValue("0"),
			),
		},
	})
}

func (ManagementGroupDataSource) basicByName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%d"
}

data "azurerm_management_group" "test" {
  name = azurerm_management_group.test.name
}
`, data.RandomInteger)
}

func (ManagementGroupDataSource) basicByDisplayName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  display_name = "acctest Management Group %d"
}

data "azurerm_management_group" "test" {
  display_name = azurerm_management_group.test.display_name
}
`, data.RandomInteger)
}
