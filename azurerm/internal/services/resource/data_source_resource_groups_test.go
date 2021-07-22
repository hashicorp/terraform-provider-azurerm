package resource_test

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ResourceGroupsDataSource struct{}

func TestAccDataSourceResourceGroups_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resource_groups", "test")

	r := ResourceGroupsDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_groups.0.id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.name").HasValue("acctestRG-1"),
				check.That(data.ResourceName).Key("resource_groups.0.type").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.location").HasValue("westeurope"),
				check.That(data.ResourceName).Key("resource_groups.0.subscription_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.tenant_id").Exists(),
			),
		},
	})
}

func (d ResourceGroupsDataSource) template(data acceptance.TestData) string {
	return `
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-1"
  location = "westeurope"
  lifecycle {
    ignore_changes = [tags]
  }
}
  `
}

func (d ResourceGroupsDataSource) basic(data acceptance.TestData) string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_resource_groups" "test" {}
`
}
