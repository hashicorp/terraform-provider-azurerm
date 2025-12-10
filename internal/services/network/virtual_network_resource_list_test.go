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

	data := acceptance.BuildTestData(t, "azurerm_virtual_network", "test1")

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
				Config:            r.basicList_query(data), // TODO - Testing not currently functional
				QueryResultChecks: []querycheck.QueryResultCheck{
					// querycheck.ExpectIdentityValue("azurerm_virtual_network.test1", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
					// querycheck.ExpectIdentityValue("azurerm_virtual_network.test1", tfjsonpath.New("name"), knownvalue.StringExact(fmt.Sprintf("acctestvnet1%d", data.RandomInteger))),
					// querycheck.ExpectIdentityValue("azurerm_virtual_network.test1", tfjsonpath.New("resource_group_name"), knownvalue.StringExact(fmt.Sprintf("acctestRG-%d", data.RandomInteger))),
					// querycheck.ExpectIdentityValue("azurerm_virtual_network.test2", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
					// querycheck.ExpectIdentityValue("azurerm_virtual_network.test2", tfjsonpath.New("name"), knownvalue.StringExact(fmt.Sprintf("acctestvnet2%d", data.RandomInteger))),
					// querycheck.ExpectIdentityValue("azurerm_virtual_network.test2", tfjsonpath.New("resource_group_name"), knownvalue.StringExact(fmt.Sprintf("acctestRG-%d", data.RandomInteger))),
					// querycheck.ExpectIdentityValue("azurerm_virtual_network.test3", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
					// querycheck.ExpectIdentityValue("azurerm_virtual_network.test3", tfjsonpath.New("name"), knownvalue.StringExact(fmt.Sprintf("acctestvnet2%d", data.RandomInteger))),
					// querycheck.ExpectIdentityValue("azurerm_virtual_network.test3", tfjsonpath.New("resource_group_name"), knownvalue.StringExact(fmt.Sprintf("acctestRG-%d", data.RandomInteger))),
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

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvnet1%[1]d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  subnet {
    name             = "subnet1"
    address_prefixes = ["10.1.1.0/24"]
  }
  tags = {
    environment = "Production"
  }
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvnet2%[1]d"
  address_space       = ["10.2.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  subnet {
    name             = "subnet1"
    address_prefixes = ["10.2.1.0/24"]
  }
  tags = {
    environment = "Production"
  }
}

resource "azurerm_virtual_network" "test3" {
  name                = "acctestvnet3%[1]d"
  address_space       = ["10.3.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  subnet {
    name             = "subnet1"
    address_prefixes = ["10.3.1.0/24"]
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
