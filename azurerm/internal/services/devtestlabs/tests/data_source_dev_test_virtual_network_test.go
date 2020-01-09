package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceArmDevTestVirtualNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_test_virtual_network", "test")

	name := fmt.Sprintf("acctestdtvn%d", data.RandomInteger)
	labName := fmt.Sprintf("acctestdtl%d", data.RandomInteger)
	resGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	subnetName := name + "Subnet"
	subnetResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, name, subnetName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmDevTestVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", name),
					resource.TestCheckResourceAttr(data.ResourceName, "lab_name", labName),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resGroup),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_subnets.0.allow_public_ip", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_subnets.0.lab_subnet_name", subnetName),
					resource.TestCheckResourceAttr(data.ResourceName, "allowed_subnets.0.resource_id", subnetResourceID),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet_overrides.0.lab_subnet_name", subnetName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet_overrides.0.resource_id", subnetResourceID),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet_overrides.0.use_in_vm_creation_permission", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet_overrides.0.use_public_ip_address_permission", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet_overrides.0.virtual_network_pool_name", ""),
				),
			},
		},
	})
}

func testAccDataSourceArmDevTestVirtualNetwork_basic(data acceptance.TestData) string {
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


`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
