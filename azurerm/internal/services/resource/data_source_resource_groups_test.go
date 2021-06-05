package resource_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ResourceGroupsDataSource struct{}

func TestAccDataSourceResourceGroups_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resource_groups", "current")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: ResourceGroupsDataSource{}.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_groups.0.id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.name").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.type").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.location").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.subscription_id").Exists(),
			),
		},
	})
}

func (d ResourceGroupsDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_resource_groups" "current" {}
`
}
