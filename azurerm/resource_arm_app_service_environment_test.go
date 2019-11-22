package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceEnvironment_basicWindows(t *testing.T) {
	resourceName := "azurerm_app_service_environment.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_basicWindows(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "number_of_ip_addresses", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccAzureRMAppServiceEnvironment_basicWindows(rInt int, location string) string {
	return fmt.Sprintf(`
	resource "azurerm_resource_group" "test_rg" {
		name     = "aseTest-%d"
		location = "%s"
	  }
	  
	  resource "azurerm_virtual_network" "test_vnet" {
		name                = "asevnettest-%d"
		location            = "${azurerm_resource_group.test_rg.location}"
		resource_group_name = "${azurerm_resource_group.test_rg.name}"
		address_space       = ["10.0.0.0/16"]
	  
	  
		subnet {
		  name           = "asesubnet"
		  address_prefix = "10.0.1.0/24"
		}
	  
		subnet {
		  name           = "gatewaysubnet"
		  address_prefix = "10.0.2.0/24"
		}
	  }
	  
	  resource "azurerm_app_service_environment" "test" {
		name                = "asetest-%d"
		location            = "${azurerm_resource_group.test_rg.location}"
		resource_group_name = "${azurerm_resource_group.test_rg.name}"
		number_of_ip_addresses = 1
		
		virtual_network {
				virtual_network_id = "${azurerm_virtual_network.test_vnet.id}"
				subnet_name = "asesubnet"
			}
	  
		frontend_pool {
				vm_size = "Small"
				number_of_workers = 2
			}
		
	  
		worker_pool {
				worker_size_id = 0
				worker_size = "Small"
				worker_count = 2
			}
	  
	  }
`, rInt, location, rInt, rInt)
}

func testCheckAzureRMAppServiceEnvironmentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).Web.AppServiceEnvironmentsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_environment" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return nil
	}

	return nil
}
