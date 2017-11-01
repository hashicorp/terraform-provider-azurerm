package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMVnet_basic(t *testing.T) {
	dataSourceName := "data.azurerm_vnet.test"
	ri := acctest.RandInt()

	name := fmt.Sprintf("acctestvnet-%d", ri)
	config := testAccDataSourceAzureRMVnet_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "dns_servers.0", "10.0.0.4"),
					resource.TestCheckResourceAttr(dataSourceName, "address_spaces.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0", "subnet1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMVnet_peering(t *testing.T) {
	firstDataSourceName := "data.azurerm_vnet.test"
	ri := acctest.RandInt()

	name_vnet_1 := fmt.Sprintf("acctestvnet-1-%d", ri)
	//name_peer := fmt.Sprintf("acctestpeer-1-%d", ri)
	config := testAccDataSourceAzureRMVnet_peering(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(firstDataSourceName, "name", name_vnet_1),
					resource.TestCheckResourceAttr(firstDataSourceName, "address_spaces.0", "10.0.1.0/24"),
					resource.TestCheckResourceAttr(firstDataSourceName, "vnet_peerings.#", "1"),
					//resource.TestCheckResourceAttr(firstDataSourceName, "vnet_peerings.0.0", "peer-1to2"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMVnet_basic(rInt int, location string) string {
	return fmt.Sprintf(`
		resource "azurerm_resource_group" "test" {
			name     = "acctest%d-rg"
			location = "%s"
		  }
		  
		  resource "azurerm_virtual_network" "test" {
			name                = "acctestvnet-%d"
			address_space       = ["10.0.0.0/16"]
			location            = "${azurerm_resource_group.test.location}"
			resource_group_name = "${azurerm_resource_group.test.name}"
			dns_servers			= ["10.0.0.4"]

			subnet {
				name 			= "subnet1"
				address_prefix	= "10.0.1.0/24"
			}
		  }

		  data "azurerm_vnet" "test" {
			  resource_group_name = "${azurerm_resource_group.test.name}"
			  name = "${azurerm_virtual_network.test.name}"
		  }

	`, rInt, location, rInt)
}

func testAccDataSourceAzureRMVnet_peering(rInt int, location string) string {
	return fmt.Sprintf(`
		resource "azurerm_resource_group" "test" {
			name     = "acctest%d-rg"
			location = "%s"
		  }
		  
		  resource "azurerm_virtual_network" "test1" {
			name                = "acctestvnet-1-%d"
			address_space       = ["10.0.1.0/24"]
			location            = "${azurerm_resource_group.test.location}"
			resource_group_name = "${azurerm_resource_group.test.name}"
		  }

		  resource "azurerm_virtual_network" "test2" {
			name                = "acctestvnet-2-%d"
			address_space       = ["10.0.2.0/24"]
			location            = "${azurerm_resource_group.test.location}"
			resource_group_name = "${azurerm_resource_group.test.name}"
		  }

		  resource "azurerm_virtual_network_peering" "test1" {
			name 						= "peer-1to2"
			resource_group_name 		= "${azurerm_resource_group.test.name}"
			virtual_network_name 		= "${azurerm_virtual_network.test1.name}"
			remote_virtual_network_id 	= "${azurerm_virtual_network.test2.id}"
		  }

		  data "azurerm_vnet" "test" {
			  resource_group_name = "${azurerm_resource_group.test.name}"
			  name = "${azurerm_virtual_network.test1.name}"
		  }
	`, rInt, location, rInt, rInt)
}
