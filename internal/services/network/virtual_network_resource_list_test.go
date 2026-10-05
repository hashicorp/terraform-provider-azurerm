// Copyright IBM Corp. 2014, 2025
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

func TestAccVirtualNetwork_list_basic(t *testing.T) {
	r := VirtualNetworkResource{}

	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test")

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
				Query:  true,
				Config: r.basicList_query(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_virtual_network.test", 3),
				},
			},
		},
	})
}

func (r VirtualNetworkResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  count = 3

  name                = "acctestvnet${count.index}-%[1]d"
  address_space       = ["10.${count.index + 1}.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  subnet {
    name             = "subnet1"
    address_prefixes = ["10.${count.index + 1}.1.0/24"]
  }
  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r VirtualNetworkResource) basicList_query(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_virtual_network" "test" {
  provider = azurerm

  config {
    resource_group_name = "acctestRG-%d"
  }
}
`, data.RandomInteger)
}
