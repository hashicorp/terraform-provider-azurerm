package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

/*
Load Balancer + updating the backend pool
App Gateway + updating the backend pool
FrontDoor?
*/

func TestAccWindowsVirtualMachine_networkIPv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkIPv6(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkMultiple(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "0"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// update the Primary IP
				Config: testWindowsVirtualMachine_networkMultipleUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "0"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// remove the secondary IP
				Config: testWindowsVirtualMachine_networkMultipleRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "0"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkMultiplePublic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkMultiplePublic(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "2"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// update the Primary IP
				Config: testWindowsVirtualMachine_networkMultiplePublicUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "2"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// remove the secondary IP
				Config: testWindowsVirtualMachine_networkMultiplePublicRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "1"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkPrivateDynamicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkPrivateStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkPrivateUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_networkPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkPublicDynamicPrivateDynamicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkPublicDynamicPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkPublicDynamicPrivateStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkPublicDynamicPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkPublicDynamicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkPublicDynamicPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_networkPublicDynamicPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkPublicStaticPrivateDynamicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkPublicStaticPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkPublicStaticPrivateStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkPublicStaticPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_networkPublicStaticPrivateUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_networkPublicStaticPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_networkPublicStaticPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func testWindowsVirtualMachine_networkIPv6(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_templateBase(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "secondary"
    private_ip_address_version    = "IPv6"
    private_ip_address_allocation = "dynamic"
  }
}

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

func testWindowsVirtualMachine_networkMultipleTemplate(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_templateBase(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "first" {
  name                = "acctestnic1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface" "second" {
  name                = "acctestnic2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testWindowsVirtualMachine_networkMultiple(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_networkMultipleTemplate(data)
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
    azurerm_network_interface.first.id,
    azurerm_network_interface.second.id,
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
`, template)
}

func testWindowsVirtualMachine_networkMultipleUpdated(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_networkMultipleTemplate(data)
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
    azurerm_network_interface.second.id,
    azurerm_network_interface.first.id,
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
`, template)
}

func testWindowsVirtualMachine_networkMultipleRemoved(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_networkMultipleTemplate(data)
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
    azurerm_network_interface.second.id,
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
`, template)
}

func testWindowsVirtualMachine_networkMultiplePublicTemplate(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_templateBase(data)
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "first" {
  name                = "acctpip1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_interface" "first" {
  name                = "acctestnic1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.first.id
  }
}

resource "azurerm_public_ip" "second" {
  name                = "acctpip2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_interface" "second" {
  name                = "acctestnic2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.second.id
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testWindowsVirtualMachine_networkMultiplePublic(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_networkMultiplePublicTemplate(data)
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
    azurerm_network_interface.first.id,
    azurerm_network_interface.second.id,
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
`, template)
}

func testWindowsVirtualMachine_networkMultiplePublicUpdated(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_networkMultiplePublicTemplate(data)
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
    azurerm_network_interface.second.id,
    azurerm_network_interface.first.id,
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
`, template)
}

func testWindowsVirtualMachine_networkMultiplePublicRemoved(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_networkMultiplePublicTemplate(data)
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
    azurerm_network_interface.second.id,
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
`, template)
}

func testWindowsVirtualMachine_networkPrivateDynamicIP(data acceptance.TestData) string {
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

func testWindowsVirtualMachine_networkPrivateStaticIP(data acceptance.TestData) string {
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

func testWindowsVirtualMachine_networkPublicDynamicPrivateDynamicIP(data acceptance.TestData) string {
	privateIPIsStatic := false
	template := testWindowsVirtualMachine_templatePrivateIP(data, privateIPIsStatic)
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

func testWindowsVirtualMachine_networkPublicDynamicPrivateStaticIP(data acceptance.TestData) string {
	privateIPIsStatic := true
	template := testWindowsVirtualMachine_templatePrivateIP(data, privateIPIsStatic)
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

func testWindowsVirtualMachine_networkPublicStaticPrivateDynamicIP(data acceptance.TestData) string {
	privateIPIsStatic := false
	publicIPIsStatic := true
	template := testWindowsVirtualMachine_templatePublicIP(data, privateIPIsStatic, publicIPIsStatic)
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

func testWindowsVirtualMachine_networkPublicStaticPrivateStaticIP(data acceptance.TestData) string {
	privateIPIsStatic := true
	publicIPIsStatic := true
	template := testWindowsVirtualMachine_templatePublicIP(data, privateIPIsStatic, publicIPIsStatic)
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

func testWindowsVirtualMachine_templatePrivateIP(data acceptance.TestData, static bool) string {
	template := testWindowsVirtualMachine_templateBase(data)
	if static {
		return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.2.30"
  }
}
`, template, data.RandomInteger)
	}

	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, template, data.RandomInteger)
}

func testWindowsVirtualMachine_templatePublicIP(data acceptance.TestData, privateStatic, publicStatic bool) string {
	template := testWindowsVirtualMachine_templateBase(data)
	publicAllocationType := allocationType(publicStatic)

	if privateStatic {
		return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = %q
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.2.30"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}
`, template, data.RandomInteger, publicAllocationType, data.RandomInteger)
	}

	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = %q
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}
`, template, data.RandomInteger, publicAllocationType, data.RandomInteger)
}
