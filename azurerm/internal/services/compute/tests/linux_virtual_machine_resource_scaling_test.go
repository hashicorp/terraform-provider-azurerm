package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccLinuxVirtualMachine_scalingAdditionalCapabilitiesUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// NOTE: this requires a large-ish machine to provision
			Config: r.scalingAdditionalCapabilitiesUltraSSD(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_scalingAvailabilitySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.scalingAvailabilitySet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_scalingDedicatedHost(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.scalingDedicatedHost(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_scalingDedicatedHostUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.scalingDedicatedHostInitial(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.scalingDedicatedHost(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.scalingDedicatedHostUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.scalingDedicatedHostRemoved(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_scalingProximityPlacementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.scalingProximityPlacementGroup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_scalingMachineSizeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.scalingMachineSize(data, "Standard_F2"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.scalingMachineSize(data, "Standard_F4"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.scalingMachineSize(data, "Standard_F4s_v2"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_scalingZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.scalingZone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LinuxVirtualMachineResource) scalingAdditionalCapabilitiesUltraSSD(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2S_V3"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]
  zone = 1

  additional_capabilities {
    ultra_ssd_enabled = true
  }

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) scalingAvailabilitySet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  managed             = true
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  availability_set_id = azurerm_availability_set.test.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r LinuxVirtualMachineResource) scalingDedicatedHostInitial(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3" # NOTE: SKU's are limited by the Dedicated Host
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r LinuxVirtualMachineResource) scalingDedicatedHost(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}

resource "azurerm_dedicated_host" "test" {
  name                    = "acctestDH-%d"
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  location                = azurerm_resource_group.test.location
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3" # NOTE: SKU's are limited by the Dedicated Host
  admin_username      = "adminuser"
  dedicated_host_id   = azurerm_dedicated_host.test.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LinuxVirtualMachineResource) scalingDedicatedHostUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}

resource "azurerm_dedicated_host" "test" {
  name                    = "acctestDH-%d"
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  location                = azurerm_resource_group.test.location
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
}

resource "azurerm_dedicated_host" "second" {
  name                    = "acctestDH2-%d"
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  location                = azurerm_resource_group.test.location
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3" # NOTE: SKU's are limited by the Dedicated Host
  admin_username      = "adminuser"
  dedicated_host_id   = azurerm_dedicated_host.second.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LinuxVirtualMachineResource) scalingDedicatedHostRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}

resource "azurerm_dedicated_host" "second" {
  name                    = "acctestDH2-%d"
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  location                = azurerm_resource_group.test.location
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3" # NOTE: SKU's are limited by the Dedicated Host
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LinuxVirtualMachineResource) scalingProximityPlacementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_linux_virtual_machine" "test" {
  name                         = "acctestVM-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  size                         = "Standard_F2"
  admin_username               = "adminuser"
  proximity_placement_group_id = azurerm_proximity_placement_group.test.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r LinuxVirtualMachineResource) scalingMachineSize(data acceptance.TestData, size string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = %q
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger, size)
}

func (r LinuxVirtualMachineResource) scalingZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]
  zone = 1

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger)
}
