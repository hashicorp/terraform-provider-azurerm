package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccSubnet_list_basic(t *testing.T) {
	r := SubnetResource{}

	data := acceptance.BuildTestData(t, "azurerm_subnet", "test1")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			//{
			//	Query:  true,
			//	Config: r.basicQuery(data),
			//	QueryResultChecks: []querycheck.QueryResultCheck{
			//		querycheck.ExpectLength("list.azurerm_virtual_network.test", 2),
			//		querycheck.ExpectLength("list.azurerm_subnet.test", 3),
			//	},
			//},
			{
				Query:  true,
				Config: r.multipleParentsQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_virtual_network.test", 2),
					querycheck.ExpectLength("azurerm_subnet.test[*]", 5),
					querycheck.ExpectLength("azurerm_subnet.test-single", 3),
					querycheck.ExpectLength("azurerm_subnet.test-count[*]", 5), // if for whatever reason you wanted to use count, that also works.
					// ExpectLength expects no `list.` prefix
					// querycheck.ExpectLength("list.azurerm_virtual_network.test", 2),
					// querycheck.ExpectLength("list.azurerm_subnet.test", 5),
				},
			},
		},
	})
}

func (r SubnetResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvnet1-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvnet2-%[1]d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  count = 3

  name                 = "acctestsubnet1-${count.index}-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test1.name

  address_prefixes = ["10.0.${count.index}.0/24"]
}

resource "azurerm_subnet" "test2" {
  count = 2

  name                 = "acctestsubnet2-${count.index}-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test2.name

  address_prefixes = ["10.1.${count.index}.0/24"]
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r SubnetResource) basicQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_virtual_network" "test" {
  provider = azurerm

  include_resource = true

  config {
    resource_group_name = "acctestRG-%d"
  }
}

list "azurerm_subnet" "test" {
  provider = azurerm

  config {
    virtual_network_id = list.azurerm_virtual_network.test.data[0].state.id
  }
}
`, data.RandomInteger)
}

func (r SubnetResource) multipleParentsQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_virtual_network" "test" {
  provider = azurerm

  include_resource = true

  config {
    resource_group_name = "acctestRG-%[1]d"
  }
}

list "azurerm_subnet" "test" {
  for_each = toset([for vnet in list.azurerm_virtual_network.test.data : vnet.state.id])

  provider = azurerm

  config {
    virtual_network_id = each.key
  }
}

list "azurerm_subnet" "test-single" {
  provider = azurerm

  config {
    virtual_network_id = list.azurerm_virtual_network.test.data[0].state.id
  }
}

list "azurerm_subnet" "test-count" {
  count = 2

  provider = azurerm

  config {
    virtual_network_id = "/subscriptions/%[2]s/resourceGroups/acctestRG-%[1]d/providers/Microsoft.Network/virtualNetworks/acctestvnet${count.index + 1}-%[1]d"
  }
}
`, data.RandomInteger, data.Subscriptions.Primary)
}
