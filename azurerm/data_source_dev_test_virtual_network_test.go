package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceArmDevTestVirtualNetwork_basic(t *testing.T) {
	dataSourceName := "data.azurerm_dev_test_virtual_network.test"
	ri := tf.AccRandTimeInt()

	name := fmt.Sprintf("acctestdtvn%d", ri)
	labName := fmt.Sprintf("acctestdtl%d", ri)
	resGroup := fmt.Sprintf("acctestRG-%d", ri)
	subnetName := name + "Subnet"
	subnetResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, name, subnetName)

	config := testAccDataSourceArmDevTestVirtualNetwork_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "lab_name", labName),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resGroup),
					resource.TestCheckResourceAttr(dataSourceName, "allowed_subnets.0.allow_public_ip", "Allow"),
					resource.TestCheckResourceAttr(dataSourceName, "allowed_subnets.0.lab_subnet_name", subnetName),
					resource.TestCheckResourceAttr(dataSourceName, "allowed_subnets.0.resource_id", subnetResourceID),
					resource.TestCheckResourceAttr(dataSourceName, "subnet_overrides.0.lab_subnet_name", subnetName),
					resource.TestCheckResourceAttr(dataSourceName, "subnet_overrides.0.resource_id", subnetResourceID),
					resource.TestCheckResourceAttr(dataSourceName, "subnet_overrides.0.use_in_vm_creation_permission", "Allow"),
					resource.TestCheckResourceAttr(dataSourceName, "subnet_overrides.0.use_public_ip_address_permission", "Allow"),
					resource.TestCheckResourceAttr(dataSourceName, "subnet_overrides.0.virtual_network_pool_name", ""),
				),
			},
		},
	})
}

func testAccDataSourceArmDevTestVirtualNetwork_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
	name                = "acctestdtl%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dev_test_virtual_network" "test" {
	name                = "acctestdtvn%d"
	lab_name            = "${azurerm_dev_test_lab.test.name}"
	resource_group_name = "${azurerm_resource_group.test.name}"
  
	subnet {
	  use_public_ip_address           = "Allow"
	  use_in_virtual_machine_creation = "Allow"
	}
}

data "azurerm_dev_test_virtual_network" "test" {
	name                = "${azurerm_dev_test_virtual_network.test.name}"
	lab_name            = "${azurerm_dev_test_lab.test.name}"
	resource_group_name = "${azurerm_resource_group.test.name}"
}


`, rInt, location, rInt, rInt)
}
