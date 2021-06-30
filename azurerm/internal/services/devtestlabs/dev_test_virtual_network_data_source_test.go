package devtestlabs_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ArmDevTestVirtualNetworkDataSource struct {
}

func TestAccArmDevTestVirtualNetworkDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_test_virtual_network", "test")
	r := ArmDevTestVirtualNetworkDataSource{}

	name := fmt.Sprintf("acctestdtvn%d", data.RandomInteger)
	labName := fmt.Sprintf("acctestdtl%d", data.RandomInteger)
	resGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	subnetName := name + "Subnet"
	subnetResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s", os.Getenv("ARM_SUBSCRIPTION_ID"), resGroup, name, subnetName)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("lab_name").HasValue(labName),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resGroup),
				check.That(data.ResourceName).Key("allowed_subnets.0.allow_public_ip").HasValue("Allow"),
				check.That(data.ResourceName).Key("allowed_subnets.0.lab_subnet_name").HasValue(subnetName),
				check.That(data.ResourceName).Key("allowed_subnets.0.resource_id").HasValue(subnetResourceID),
				check.That(data.ResourceName).Key("subnet_overrides.0.lab_subnet_name").HasValue(subnetName),
				check.That(data.ResourceName).Key("subnet_overrides.0.resource_id").HasValue(subnetResourceID),
				check.That(data.ResourceName).Key("subnet_overrides.0.use_in_vm_creation_permission").HasValue("Allow"),
				check.That(data.ResourceName).Key("subnet_overrides.0.use_public_ip_address_permission").HasValue("Allow"),
				check.That(data.ResourceName).Key("subnet_overrides.0.virtual_network_pool_name").HasValue(""),
			),
		},
	})
}

func (ArmDevTestVirtualNetworkDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_virtual_network" "test" {
  name                = "acctestdtvn%d"
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name

  subnet {
    use_public_ip_address           = "Allow"
    use_in_virtual_machine_creation = "Allow"
  }
}

data "azurerm_dev_test_virtual_network" "test" {
  name                = azurerm_dev_test_virtual_network.test.name
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
