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

func TestAccRouteTable_list_basic(t *testing.T) {
	r := RouteTableResource{}

	data := acceptance.BuildTestData(t, "azurerm_route_table", "test1")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:             true,
				Config:            r.basicQuery(),
				ConfigQueryChecks: []querycheck.QueryCheck{}, // TODO
			},
			{
				Query:             true,
				Config:            r.basicQueryByResourceGroupName(data),
				ConfigQueryChecks: []querycheck.QueryCheck{}, // TODO
			},
		},
	})
}

func (r RouteTableResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_route_table" "test1" {
  name                = "acctestrt1-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_route_table" "test2" {
  name                = "acctestrt2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_route_table" "test3" {
  name                = "acctestrt3-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r RouteTableResource) basicQuery() string {
	return `
list "azurerm_route_table" "list" {
  provider = azurerm
  config {}
}
`
}

func (r RouteTableResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_route_table" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-%[1]d"
  }
}
`, data.RandomInteger)
}
