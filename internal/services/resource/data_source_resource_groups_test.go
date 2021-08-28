package resource_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ResourceGroupsDataSource struct{}

const NumberOfResourceGroups = 3

func TestAccDataSourceResourceGroups_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resource_groups", "test")
	r := ResourceGroupsDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.template(data, NumberOfResourceGroups),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_groups.0.id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.name").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.type").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.location").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.subscription_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.name").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.type").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.location").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.subscription_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.tenant_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.name").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.type").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.location").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.subscription_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.tenant_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceResourceGroups_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resource_groups", "test")
	r := ResourceGroupsDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.template(data, NumberOfResourceGroups),
		},
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_groups.0.id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.name").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.type").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.location").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.subscription_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.name").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.type").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.location").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.subscription_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.1.tenant_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.name").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.type").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.location").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.subscription_id").Exists(),
				check.That(data.ResourceName).Key("resource_groups.2.tenant_id").Exists(),
			),
		},
	})
}

func (d ResourceGroupsDataSource) basic(data acceptance.TestData) string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_resource_groups" "test" {}
`
}

func (d ResourceGroupsDataSource) complete(data acceptance.TestData) string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

data "azurerm_resource_groups" "test" {
  filter_by_subscription_id = [data.azurerm_client_config.current.subscription_id]
}
`
}

func (ResourceGroupsDataSource) template(data acceptance.TestData, count int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  count    = "%d"
  name     = "acctestRG-${count.index}"
  location = "%s"
  lifecycle {
    ignore_changes = [tags]
  }
}

`, count, data.Locations.Primary)
}
