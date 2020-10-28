package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualMachineScaleSetExtension_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_basicLinux(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_basicWindows(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_basicLinux(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMVirtualMachineScaleSetExtension_requiresImport),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_autoUpgradeDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_autoUpgradeDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_extensionChaining(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "first")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_extensionChaining(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists("azurerm_virtual_machine_scale_set_extension.first"),
					testCheckAzureRMVirtualMachineScaleSetExtensionExists("azurerm_virtual_machine_scale_set_extension.second"),
				),
			},
			data.ImportStep(),
			{
				ResourceName:      "azurerm_virtual_machine_scale_set_extension.second",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_forceUpdateTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_forceUpdateTag(data, "first"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_forceUpdateTag(data, "second"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_protectedSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_protectedSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(data.ResourceName),
				),
			},
			data.ImportStep("protected_settings"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_protectedSettingsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_protectedSettingsOnly(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(data.ResourceName),
				),
			},
			data.ImportStep("protected_settings"),
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_updateVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				// old version
				Config: testAccAzureRMVirtualMachineScaleSetExtension_updateVersion(data, "1.2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_updateVersion(data, "1.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetExtensionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.VirtualMachineScaleSetExtensionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.Name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Extension %q (VirtualMachineScaleSet %q / Resource Group: %q) does not exist", id.Name, id.VirtualMachineScaleSetName, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on vmScaleSetClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetExtensionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetExtensionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_machine_scale_set_extension" {
			continue
		}

		id, err := parse.VirtualMachineScaleSetExtensionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.Name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Compute.VMScaleSetExtensionsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualMachineScaleSetExtension_basicLinux(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  settings = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualMachineScaleSetExtension_basicWindows(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateWindows(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_windows_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  settings = jsonencode({
    "commandToExecute" = "Write-Host \"Hello\""
  })
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualMachineScaleSetExtension_autoUpgradeDisabled(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  auto_upgrade_minor_version   = false
  settings = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualMachineScaleSetExtension_extensionChaining(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "first" {
  name                         = "acctestExt1-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "DockerExtension"
  type_handler_version         = "1.0"
}

resource "azurerm_virtual_machine_scale_set_extension" "second" {
  name                         = "acctestExt2-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  settings = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
  provision_after_extensions = [azurerm_virtual_machine_scale_set_extension.first.name]
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualMachineScaleSetExtension_forceUpdateTag(data acceptance.TestData, tag string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  force_update_tag             = %q
  settings = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
}
`, template, data.RandomInteger, tag)
}

func testAccAzureRMVirtualMachineScaleSetExtension_updateVersion(data acceptance.TestData, version string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.OSTCExtensions"
  type                         = "CustomScriptForLinux"
  type_handler_version         = %q
  settings = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
}
`, template, data.RandomInteger, version)
}

func testAccAzureRMVirtualMachineScaleSetExtension_protectedSettings(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  settings = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
  protected_settings = jsonencode({
    "secretValue" = "P@55W0rd1234!"
  })
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualMachineScaleSetExtension_protectedSettingsOnly(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  protected_settings = jsonencode({
    "commandToExecute" = "echo $HOSTNAME",
    "secretValue"      = "P@55W0rd1234!"
  })
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualMachineScaleSetExtension_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_basicLinux(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "import" {
  name                         = azurerm_virtual_machine_scale_set_extension.test.name
  virtual_machine_scale_set_id = azurerm_virtual_machine_scale_set_extension.test.virtual_machine_scale_set_id
  publisher                    = azurerm_virtual_machine_scale_set_extension.test.publisher
  type                         = azurerm_virtual_machine_scale_set_extension.test.type
  type_handler_version         = azurerm_virtual_machine_scale_set_extension.test.type_handler_version
  settings                     = azurerm_virtual_machine_scale_set_extension.test.settings
}
`, template)
}

func testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name               = "internal"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"
}


resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualMachineScaleSetExtension_templateWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name               = "internal"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"
}

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                 = "acctestvm%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  sku                  = "Standard_F2"
  instances            = 1
  admin_username       = "adminuser"
  admin_password       = "P@ssword1234!"
  computer_name_prefix = "acctestvm"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
