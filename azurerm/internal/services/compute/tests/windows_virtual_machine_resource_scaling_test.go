package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccWindowsVirtualMachine_scalingAdditionalCapabilitiesUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				// NOTE: this requires a large-ish machine to provision
				Config: testWindowsVirtualMachine_scalingAdditionalCapabilitiesUltraSSD(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_scalingAvailabilitySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_scalingAvailabilitySet(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_scalingDedicatedHost(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_scalingDedicatedHost(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_scalingDedicatedHostUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_scalingDedicatedHostInitial(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_scalingDedicatedHost(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_scalingDedicatedHostUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_scalingDedicatedHostRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_scalingProximityPlacementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_scalingProximityPlacementGroup(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_scalingMachineSizeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_scalingMachineSize(data, "Standard_F2"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_scalingMachineSize(data, "Standard_F4"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_scalingMachineSize(data, "Standard_F4s_v2"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_scalingZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_scalingZone(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func testWindowsVirtualMachine_scalingAdditionalCapabilitiesUltraSSD(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2S_V3"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]
  zone = 1

  additional_capabilities {
    ultra_ssd_enabled = true
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_scalingAvailabilitySet(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  managed             = true
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  availability_set_id = azurerm_availability_set.test.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func testWindowsVirtualMachine_scalingDedicatedHostInitial(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3" # NOTE: SKU's are limited by the Dedicated Host
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func testWindowsVirtualMachine_scalingDedicatedHost(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
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

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3" # NOTE: SKU's are limited by the Dedicated Host
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  dedicated_host_id   = azurerm_dedicated_host.test.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testWindowsVirtualMachine_scalingDedicatedHostUpdate(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
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

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3" # NOTE: SKU's are limited by the Dedicated Host
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  dedicated_host_id   = azurerm_dedicated_host.second.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testWindowsVirtualMachine_scalingDedicatedHostRemoved(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
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

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3" # NOTE: SKU's are limited by the Dedicated Host
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  dedicated_host_id   = azurerm_dedicated_host.second.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testWindowsVirtualMachine_scalingProximityPlacementGroup(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_windows_virtual_machine" "test" {
  name                         = local.vm_name
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  size                         = "Standard_F2"
  admin_username               = "adminuser"
  admin_password               = "P@$$w0rd1234!"
  proximity_placement_group_id = azurerm_proximity_placement_group.test.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func testWindowsVirtualMachine_scalingMachineSize(data acceptance.TestData, size string) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = %q
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, size)
}

func testWindowsVirtualMachine_scalingZone(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]
  zone = 1

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}
