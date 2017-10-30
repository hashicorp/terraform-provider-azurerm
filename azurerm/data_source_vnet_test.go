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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "dns_servers.0", "10.0.0.10"),
					resource.TestCheckResourceAttr(dataSourceName, "address_spaces.0", "10.0.0.0/16"),
				),
			},
		},
	})
}

// func TestAccDataSourceAzureRMVnet_Subnets(t *testing.T) {
// 	dataSourceName := "data.azurerm_vnet.test"
// 	ri := acctest.RandInt()

// 	name := fmt.Sprintf("acctestvnet-%d", ri)
// 	//resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)

// 	config := testAccDataSourceAzureRMVnet_basic(ri, testLocation())

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testCheckAzureRMPublicIpDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: config,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(dataSourceName, "name", name),
// 					resource.TestCheckResourceAttr(dataSourceName, "dns_servers", `["10.0.0.10","10.0.0.11"]`),
// 				),
// 			},
// 		},
// 	})
// }

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
			dns_servers			= ["10.0.0.10"]
		  }

		  data "azurerm_vnet" "test" {
			  resource_group_name = "${azurerm_resource_group.test.name}"
			  name = "${azurerm_virtual_network.test.name}"
		  }

	`, rInt, location, rInt)
}

// func testAccDataSourceAzureRMVnet_Subnets(rInt int, location string) string {
// 	return fmt.Sprintf(`
// 		resource "azurerm_resource_group" "test" {
// 			name     = "acctest%d-rg"
// 			location = "%s"
// 		  }

// 		  resource "azurerm_virtual_network" "test" {
// 			name                = "acctestvnet-%d"
// 			address_space       = ["10.0.0.0/16"]
// 			location            = "${azurerm_resource_group.test.location}"
// 			resource_group_name = "${azurerm_resource_group.test.name}"
// 			dns_servers			= ["10.0.0.10","10.0.0.11"]
// 		  }

// 		  data "azurerm_vnet" "test" {
// 			  resource_group_name = "${azurerm_resource_group.test.name}"
// 			  name = "${azurerm_virtual_network.test.name}"
// 			  dns_servers = "${azurerm_virtual_network.test.dns_servers}"
// 		  }

// 	`, rInt, location, rInt)
// }
