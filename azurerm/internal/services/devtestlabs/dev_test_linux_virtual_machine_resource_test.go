package devtestlabs_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccDevTestLinuxVirtualMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLinuxVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLinuxVirtualMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gallery_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(
				// not returned from the API
				"lab_subnet_name",
				"lab_virtual_network_id",
				"password",
			),
		},
	})
}

func TestAccDevTestLinuxVirtualMachine_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLinuxVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLinuxVirtualMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			{
				Config:      testAccDevTestLinuxVirtualMachine_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dev_test_lab_linux_virtual_machine"),
			},
		},
	})
}

func TestAccDevTestLinuxVirtualMachine_basicSSH(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLinuxVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLinuxVirtualMachine_basicSSH(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gallery_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(
				// not returned from the API
				"lab_subnet_name",
				"lab_virtual_network_id",
				"password",
				"ssh_key",
			),
		},
	})
}

func TestAccDevTestLinuxVirtualMachine_inboundNatRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLinuxVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLinuxVirtualMachine_inboundNatRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "disallow_public_ip_address", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "gallery_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Acceptance", "Test"),
				),
			},
			data.ImportStep(
				// not returned from the API
				"inbound_nat_rule",
				"lab_subnet_name",
				"lab_virtual_network_id",
				"password",
			),
		},
	})
}

func TestAccDevTestLinuxVirtualMachine_updateStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLinuxVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLinuxVirtualMachine_storage(data, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gallery_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_type", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccDevTestLinuxVirtualMachine_storage(data, "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLinuxVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gallery_image_reference.0.publisher", "Canonical"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_type", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testCheckDevTestLinuxVirtualMachineExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.VirtualMachinesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		virtualMachineName := rs.Primary.Attributes["name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, virtualMachineName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get devTestVirtualMachinesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DevTest Linux Virtual Machine %q (Lab %q / Resource Group: %q) does not exist", virtualMachineName, labName, resourceGroup)
		}

		return nil
	}
}

func testCheckDevTestLinuxVirtualMachineDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.VirtualMachinesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_linux_virtual_machine" {
			continue
		}

		virtualMachineName := rs.Primary.Attributes["name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, virtualMachineName, "")
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DevTest Linux Virtual Machine still exists:\n%#v", resp)
	}

	return nil
}

func testAccDevTestLinuxVirtualMachine_basic(data acceptance.TestData) string {
	template := testAccDevTestLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_linux_virtual_machine" "test" {
  name                   = "acctestvm-vm%d"
  lab_name               = azurerm_dev_test_lab.test.name
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  size                   = "Standard_F2"
  username               = "acct5stU5er"
  password               = "Pa$w0rd1234!"
  lab_virtual_network_id = azurerm_dev_test_virtual_network.test.id
  lab_subnet_name        = azurerm_dev_test_virtual_network.test.subnet[0].name
  storage_type           = "Standard"

  gallery_image_reference {
    offer     = "UbuntuServer"
    publisher = "Canonical"
    sku       = "18.04-LTS"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func testAccDevTestLinuxVirtualMachine_requiresImport(data acceptance.TestData) string {
	template := testAccDevTestLinuxVirtualMachine_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_linux_virtual_machine" "import" {
  name                   = azurerm_dev_test_linux_virtual_machine.test.name
  lab_name               = azurerm_dev_test_linux_virtual_machine.test.lab_name
  resource_group_name    = azurerm_dev_test_linux_virtual_machine.test.resource_group_name
  location               = azurerm_dev_test_linux_virtual_machine.test.location
  size                   = azurerm_dev_test_linux_virtual_machine.test.size
  username               = "acct5stU5er"
  password               = "Pa$w0rd1234!"
  lab_virtual_network_id = azurerm_dev_test_virtual_network.test.id
  lab_subnet_name        = azurerm_dev_test_virtual_network.test.subnet[0].name
  storage_type           = "Standard"

  gallery_image_reference {
    offer     = "UbuntuServer"
    publisher = "Canonical"
    sku       = "18.04-LTS"
    version   = "latest"
  }
}
`, template)
}

func testAccDevTestLinuxVirtualMachine_basicSSH(data acceptance.TestData) string {
	template := testAccDevTestLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_linux_virtual_machine" "test" {
  name                   = "acctestvm-vm%d"
  lab_name               = azurerm_dev_test_lab.test.name
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  size                   = "Standard_F2"
  username               = "acct5stU5er"
  ssh_key                = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
  lab_virtual_network_id = azurerm_dev_test_virtual_network.test.id
  lab_subnet_name        = azurerm_dev_test_virtual_network.test.subnet[0].name
  storage_type           = "Standard"

  gallery_image_reference {
    offer     = "UbuntuServer"
    publisher = "Canonical"
    sku       = "18.04-LTS"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func testAccDevTestLinuxVirtualMachine_inboundNatRules(data acceptance.TestData) string {
	template := testAccDevTestLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_linux_virtual_machine" "test" {
  name                       = "acctestvm-vm%d"
  lab_name                   = azurerm_dev_test_lab.test.name
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  size                       = "Standard_F2"
  username                   = "acct5stU5er"
  password                   = "Pa$w0rd1234!"
  disallow_public_ip_address = true
  lab_virtual_network_id     = azurerm_dev_test_virtual_network.test.id
  lab_subnet_name            = azurerm_dev_test_virtual_network.test.subnet[0].name
  storage_type               = "Standard"

  gallery_image_reference {
    offer     = "UbuntuServer"
    publisher = "Canonical"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  inbound_nat_rule {
    protocol     = "Tcp"
    backend_port = 22
  }

  inbound_nat_rule {
    protocol     = "Tcp"
    backend_port = 3389
  }

  tags = {
    "Acceptance" = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccDevTestLinuxVirtualMachine_storage(data acceptance.TestData, storageType string) string {
	template := testAccDevTestLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_linux_virtual_machine" "test" {
  name                   = "acctestvm-vm%d"
  lab_name               = azurerm_dev_test_lab.test.name
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  size                   = "Standard_B1ms"
  username               = "acct5stU5er"
  password               = "Pa$w0rd1234!"
  lab_virtual_network_id = azurerm_dev_test_virtual_network.test.id
  lab_subnet_name        = azurerm_dev_test_virtual_network.test.subnet[0].name
  storage_type           = "%s"

  gallery_image_reference {
    offer     = "UbuntuServer"
    publisher = "Canonical"
    sku       = "18.04-LTS"
    version   = "latest"
  }
}
`, template, data.RandomInteger, storageType)
}

func testAccDevTestLinuxVirtualMachine_template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
