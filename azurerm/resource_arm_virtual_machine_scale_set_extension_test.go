package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	computeSvc "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute"
)

func TestAccAzureRMVirtualMachineScaleSetExtension_basicLinux(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set_extension.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_basicLinux(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_basicWindows(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set_extension.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_basicWindows(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_virtual_machine_scale_set_extension.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_basicLinux(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualMachineScaleSetExtension_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_machine_scale_set_extension"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_autoUpgradeDisabled(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set_extension.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_autoUpgradeDisabled(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_extensionChaining(t *testing.T) {
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_extensionChaining(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists("azurerm_virtual_machine_scale_set_extension.first"),
					testCheckAzureRMVirtualMachineScaleSetExtensionExists("azurerm_virtual_machine_scale_set_extension.second"),
				),
			},
			{
				ResourceName:      "azurerm_virtual_machine_scale_set_extension.first",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "azurerm_virtual_machine_scale_set_extension.second",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_forceUpdateTag(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set_extension.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_forceUpdateTag(ri, location, "first"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_forceUpdateTag(ri, location, "second"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_protectedSettings(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set_extension.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_protectedSettings(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"protected_settings"},
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_protectedSettingsOnly(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set_extension.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_protectedSettingsOnly(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"protected_settings"},
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSetExtension_updateVersion(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set_extension.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetExtensionDestroy,
		Steps: []resource.TestStep{
			{
				// old version
				Config: testAccAzureRMVirtualMachineScaleSetExtension_updateVersion(ri, location, "1.2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualMachineScaleSetExtension_updateVersion(ri, location, "1.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMVirtualMachineScaleSetExtensionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetExtensionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		name := rs.Primary.Attributes["name"]
		virtualMachineScaleSetIdRaw := rs.Primary.Attributes["virtual_machine_scale_set_id"]
		virtualMachineScaleSetId, err := computeSvc.ParseVirtualMachineScaleSetID(virtualMachineScaleSetIdRaw)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, virtualMachineScaleSetId.ResourceGroup, virtualMachineScaleSetId.Name, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on vmScaleSetClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Extension %q (VirtualMachineScaleSet %q / Resource Group: %q) does not exist", name, virtualMachineScaleSetId.Name, virtualMachineScaleSetId.ResourceGroup)
		}

		return err
	}
}

func testCheckAzureRMVirtualMachineScaleSetExtensionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetExtensionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_machine_scale_set_extension" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		virtualMachineScaleSetIdRaw := rs.Primary.Attributes["virtual_machine_scale_set_id"]
		virtualMachineScaleSetId, err := computeSvc.ParseVirtualMachineScaleSetID(virtualMachineScaleSetIdRaw)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, virtualMachineScaleSetId.ResourceGroup, virtualMachineScaleSetId.Name, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual Machine Scale Set Extension still exists:\n%#v", resp.VirtualMachineScaleSetExtensionProperties)
		}
	}

	return nil
}

func testAccAzureRMVirtualMachineScaleSetExtension_basicLinux(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  settings                     = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
}
`, template, rInt)
}

func testAccAzureRMVirtualMachineScaleSetExtension_basicWindows(rInt int, rString, location string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateWindows(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_windows_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  settings                     = jsonencode({
    "commandToExecute" = "Write-Host \"Hello\""
  })
}
`, template, rInt)
}

func testAccAzureRMVirtualMachineScaleSetExtension_autoUpgradeDisabled(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  auto_upgrade_minor_version   = false
  settings                     = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
}
`, template, rInt)
}

func testAccAzureRMVirtualMachineScaleSetExtension_extensionChaining(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(rInt, location)
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
  settings                     = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
  provision_after_extensions = [ azurerm_virtual_machine_scale_set_extension.first.name ]
}
`, template, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSetExtension_forceUpdateTag(rInt int, location, tag string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  force_update_tag             = %q
  settings                     = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
}
`, template, rInt, tag)
}

func testAccAzureRMVirtualMachineScaleSetExtension_updateVersion(rInt int, location, version string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.OSTCExtensions"
  type                         = "CustomScriptForLinux"
  type_handler_version         = %q
  settings                     = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
}
`, template, rInt, version)
}

func testAccAzureRMVirtualMachineScaleSetExtension_protectedSettings(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  settings                     = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
  protected_settings           = jsonencode({
    "secretValue" = "P@55W0rd1234!"
  })
}
`, template, rInt)
}

func testAccAzureRMVirtualMachineScaleSetExtension_protectedSettingsOnly(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "acctestExt-%d"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  protected_settings           = jsonencode({
    "commandToExecute" = "echo $HOSTNAME",
    "secretValue"      = "P@55W0rd1234!"
  })
}
`, template, rInt)
}

func testAccAzureRMVirtualMachineScaleSetExtension_requiresImport(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineScaleSetExtension_basicLinux(rInt, location)
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

func testAccAzureRMVirtualMachineScaleSetExtension_templateLinux(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
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
`, rInt, location, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSetExtension_templateWindows(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
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
`, rInt, location, rInt, rString)
}
