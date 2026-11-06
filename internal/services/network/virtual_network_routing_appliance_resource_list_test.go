// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

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

func testAccVirtualNetworkRoutingAppliance_list_basic(t *testing.T) {
	r := VirtualNetworkRoutingApplianceResource{}
	listResourceAddress := "azurerm_virtual_network_routing_appliance.list"

	data := acceptance.BuildTestData(t, "azurerm_virtual_network_routing_appliance", "test")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 2),
				},
			},
		},
	})
}

func (VirtualNetworkRoutingApplianceResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  count = 2

  name                = "acctestVNet-${count.index}-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.${count.index}.0.0/16"]
}

resource "azurerm_subnet" "test" {
  count = 2

  name                 = "VirtualNetworkApplianceSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test[count.index].name
  address_prefixes     = ["10.${count.index}.0.0/24"]
}

resource "azurerm_virtual_network_routing_appliance" "test" {
  count = 2

  name                = "acctestVNRA-${count.index}-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  bandwidth_in_gbps   = 10
  subnet_id           = azurerm_subnet.test[count.index].id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (VirtualNetworkRoutingApplianceResource) basicQuery() string {
	return `
list "azurerm_virtual_network_routing_appliance" "list" {
  provider         = azurerm
  include_resource = true

  config {}
}
`
}

func (VirtualNetworkRoutingApplianceResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_virtual_network_routing_appliance" "list" {
  provider         = azurerm
  include_resource = true

  config {
    resource_group_name = "acctestRG-%d"
  }
}
`, data.RandomInteger)
}
