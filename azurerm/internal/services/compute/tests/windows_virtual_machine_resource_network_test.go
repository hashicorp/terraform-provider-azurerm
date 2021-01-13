package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

/*
Load Balancer + updating the backend pool
App Gateway + updating the backend pool
FrontDoor?
*/

func TestAccWindowsVirtualMachine_networkIPv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkIPv6(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").HasValue(""),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkMultiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("private_ip_addresses.#").HasValue("2"),
				check.That(data.ResourceName).Key("public_ip_address").HasValue(""),
				check.That(data.ResourceName).Key("public_ip_addresses.#").HasValue("0"),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// update the Primary IP
			Config: r.networkMultipleUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("private_ip_addresses.#").HasValue("2"),
				check.That(data.ResourceName).Key("public_ip_address").HasValue(""),
				check.That(data.ResourceName).Key("public_ip_addresses.#").HasValue("0"),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// remove the secondary IP
			Config: r.networkMultipleRemoved(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("private_ip_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("public_ip_address").HasValue(""),
				check.That(data.ResourceName).Key("public_ip_addresses.#").HasValue("0"),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkMultiplePublic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkMultiplePublic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("private_ip_addresses.#").HasValue("2"),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_addresses.#").HasValue("2"),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// update the Primary IP
			Config: r.networkMultiplePublicUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("private_ip_addresses.#").HasValue("2"),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_addresses.#").HasValue("2"),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// remove the secondary IP
			Config: r.networkMultiplePublicRemoved(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("private_ip_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_addresses.#").HasValue("1"),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkPrivateDynamicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkPrivateDynamicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").HasValue(""),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkPrivateStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkPrivateStaticIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").HasValue(""),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkPrivateUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkPrivateDynamicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").HasValue(""),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.networkPrivateStaticIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").HasValue(""),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkPublicDynamicPrivateDynamicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkPublicDynamicPrivateDynamicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkPublicDynamicPrivateStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkPublicDynamicPrivateStaticIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkPublicDynamicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkPublicDynamicPrivateDynamicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.networkPublicDynamicPrivateStaticIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkPublicStaticPrivateDynamicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkPublicStaticPrivateDynamicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkPublicStaticPrivateStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkPublicStaticPrivateStaticIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachine_networkPublicStaticPrivateUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkPublicStaticPrivateDynamicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.networkPublicStaticPrivateStaticIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_ip_address").Exists(),
				check.That(data.ResourceName).Key("public_ip_address").Exists(),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func (r WindowsVirtualMachineResource) networkIPv6(data acceptance.TestData) string {
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
`, r.templateBase(data), data.RandomInteger)
}

func (r WindowsVirtualMachineResource) networkMultipleTemplate(data acceptance.TestData) string {
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
`, r.templateBase(data), data.RandomInteger, data.RandomInteger)
}

func (r WindowsVirtualMachineResource) networkMultiple(data acceptance.TestData) string {
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
`, r.networkMultipleTemplate(data))
}

func (r WindowsVirtualMachineResource) networkMultipleUpdated(data acceptance.TestData) string {
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
`, r.networkMultipleTemplate(data))
}

func (r WindowsVirtualMachineResource) networkMultipleRemoved(data acceptance.TestData) string {
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
`, r.networkMultipleTemplate(data))
}

func (r WindowsVirtualMachineResource) networkMultiplePublicTemplate(data acceptance.TestData) string {
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
`, r.templateBase(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r WindowsVirtualMachineResource) networkMultiplePublic(data acceptance.TestData) string {
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
`, r.networkMultiplePublicTemplate(data))
}

func (r WindowsVirtualMachineResource) networkMultiplePublicUpdated(data acceptance.TestData) string {
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
`, r.networkMultiplePublicTemplate(data))
}

func (r WindowsVirtualMachineResource) networkMultiplePublicRemoved(data acceptance.TestData) string {
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
`, r.networkMultiplePublicTemplate(data))
}

func (r WindowsVirtualMachineResource) networkPrivateDynamicIP(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r WindowsVirtualMachineResource) networkPrivateStaticIP(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r WindowsVirtualMachineResource) networkPublicDynamicPrivateDynamicIP(data acceptance.TestData) string {
	privateIPIsStatic := false
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
`, r.templatePrivateIP(data, privateIPIsStatic))
}

func (r WindowsVirtualMachineResource) networkPublicDynamicPrivateStaticIP(data acceptance.TestData) string {
	privateIPIsStatic := true
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
`, r.templatePrivateIP(data, privateIPIsStatic))
}

func (r WindowsVirtualMachineResource) networkPublicStaticPrivateDynamicIP(data acceptance.TestData) string {
	privateIPIsStatic := false
	publicIPIsStatic := true
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
`, r.templatePublicIP(data, privateIPIsStatic, publicIPIsStatic))
}

func (r WindowsVirtualMachineResource) networkPublicStaticPrivateStaticIP(data acceptance.TestData) string {
	privateIPIsStatic := true
	publicIPIsStatic := true
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
`, r.templatePublicIP(data, privateIPIsStatic, publicIPIsStatic))
}

func (r WindowsVirtualMachineResource) templatePrivateIP(data acceptance.TestData, static bool) string {
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
`, r.templateBase(data), data.RandomInteger)
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
`, r.templateBase(data), data.RandomInteger)
}

func (r WindowsVirtualMachineResource) templatePublicIP(data acceptance.TestData, privateStatic, publicStatic bool) string {
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
`, r.templateBase(data), data.RandomInteger, publicAllocationType, data.RandomInteger)
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
`, r.templateBase(data), data.RandomInteger, publicAllocationType, data.RandomInteger)
}
