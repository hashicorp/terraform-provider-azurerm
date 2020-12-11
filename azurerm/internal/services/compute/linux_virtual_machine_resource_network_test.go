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

func TestAccLinuxVirtualMachine_networkIPv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkIPv6(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkMultiple(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "0"),
				),
			},
			data.ImportStep(),
			{
				// update the Primary IP
				Config: testLinuxVirtualMachine_networkMultipleUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "0"),
				),
			},
			data.ImportStep(),
			{
				// remove the secondary IP
				Config: testLinuxVirtualMachine_networkMultipleRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkMultiplePublic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkMultiplePublic(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "2"),
				),
			},
			data.ImportStep(),
			{
				// update the Primary IP
				Config: testLinuxVirtualMachine_networkMultiplePublicUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "2"),
				),
			},
			data.ImportStep(),
			{
				// remove the secondary IP
				Config: testLinuxVirtualMachine_networkMultiplePublicRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_addresses.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_addresses.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkPrivateDynamicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkPrivateStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkPrivateUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
				),
			},
			data.ImportStep(),
			{
				Config: testLinuxVirtualMachine_networkPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkPublicDynamicPrivateDynamicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkPublicDynamicPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkPublicDynamicPrivateStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkPublicDynamicPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkPublicDynamicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkPublicDynamicPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(),
			{
				Config: testLinuxVirtualMachine_networkPublicDynamicPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkPublicStaticPrivateDynamicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkPublicStaticPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkPublicStaticPrivateStaticIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkPublicStaticPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccLinuxVirtualMachine_networkPublicStaticPrivateUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_networkPublicStaticPrivateDynamicIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(),
			{
				Config: testLinuxVirtualMachine_networkPublicStaticPrivateStaticIP(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_address"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testLinuxVirtualMachine_networkIPv6(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_templateBase(data)
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

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
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
`, template, data.RandomInteger, data.RandomInteger)
}

func testLinuxVirtualMachine_networkMultipleTemplate(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_templateBase(data)
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

func testLinuxVirtualMachine_networkMultiple(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_networkMultipleTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.first.id,
    azurerm_network_interface.second.id,
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkMultipleUpdated(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_networkMultipleTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.second.id,
    azurerm_network_interface.first.id,
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkMultipleRemoved(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_networkMultipleTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.second.id,
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkMultiplePublicTemplate(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_templateBase(data)
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

func testLinuxVirtualMachine_networkMultiplePublic(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_networkMultiplePublicTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.first.id,
    azurerm_network_interface.second.id,
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkMultiplePublicUpdated(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_networkMultiplePublicTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.second.id,
    azurerm_network_interface.first.id,
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkMultiplePublicRemoved(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_networkMultiplePublicTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.second.id,
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkPrivateDynamicIP(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_template(data)
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkPrivateStaticIP(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_template(data)
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkPublicDynamicPrivateDynamicIP(data acceptance.TestData) string {
	privateIPIsStatic := false
	template := testLinuxVirtualMachine_templatePrivateIP(data, privateIPIsStatic)
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkPublicDynamicPrivateStaticIP(data acceptance.TestData) string {
	privateIPIsStatic := true
	template := testLinuxVirtualMachine_templatePrivateIP(data, privateIPIsStatic)
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkPublicStaticPrivateDynamicIP(data acceptance.TestData) string {
	privateIPIsStatic := false
	publicIPIsStatic := true
	template := testLinuxVirtualMachine_templatePublicIP(data, privateIPIsStatic, publicIPIsStatic)
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_networkPublicStaticPrivateStaticIP(data acceptance.TestData) string {
	privateIPIsStatic := true
	publicIPIsStatic := true
	template := testLinuxVirtualMachine_templatePublicIP(data, privateIPIsStatic, publicIPIsStatic)
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
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_templatePrivateIP(data acceptance.TestData, static bool) string {
	template := testLinuxVirtualMachine_templateBase(data)
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

func testLinuxVirtualMachine_templatePublicIP(data acceptance.TestData, privateStatic, publicStatic bool) string {
	template := testLinuxVirtualMachine_templateBase(data)
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

func allocationType(static bool) string {
	if static {
		return "Static"
	}

	return "Dynamic"
}
