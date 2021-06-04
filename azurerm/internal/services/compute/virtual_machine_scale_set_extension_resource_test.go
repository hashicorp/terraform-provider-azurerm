package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type VirtualMachineScaleSetExtensionResource struct {
}

func TestAccVirtualMachineScaleSetExtension_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	r := VirtualMachineScaleSetExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineScaleSetExtension_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	r := VirtualMachineScaleSetExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineScaleSetExtension_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	r := VirtualMachineScaleSetExtensionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccVirtualMachineScaleSetExtension_autoUpgradeDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	r := VirtualMachineScaleSetExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoUpgradeDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineScaleSetExtension_extensionChaining(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "first")
	r := VirtualMachineScaleSetExtensionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionChaining(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			ResourceName:      "azurerm_virtual_machine_scale_set_extension.second",
			ImportState:       true,
			ImportStateVerify: true,
		},
	})
}

func TestAccVirtualMachineScaleSetExtension_forceUpdateTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	r := VirtualMachineScaleSetExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.forceUpdateTag(data, "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.forceUpdateTag(data, "second"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineScaleSetExtension_protectedSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	r := VirtualMachineScaleSetExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.protectedSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_settings"),
	})
}

func TestAccVirtualMachineScaleSetExtension_protectedSettingsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	r := VirtualMachineScaleSetExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.protectedSettingsOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_settings"),
	})
}

func TestAccVirtualMachineScaleSetExtension_updateVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_extension", "test")
	r := VirtualMachineScaleSetExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// old version
			Config: r.updateVersion(data, "1.2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateVersion(data, "1.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VirtualMachineScaleSetExtensionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VirtualMachineScaleSetExtensionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.VMScaleSetExtensionsClient.Get(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.ExtensionName, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Virtual Machine Scale Set Extension %q", id.String())
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r VirtualMachineScaleSetExtensionResource) basicLinux(data acceptance.TestData) string {
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
`, r.templateLinux(data), data.RandomInteger)
}

func (r VirtualMachineScaleSetExtensionResource) basicWindows(data acceptance.TestData) string {
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
`, r.templateWindows(data), data.RandomInteger)
}

func (r VirtualMachineScaleSetExtensionResource) autoUpgradeDisabled(data acceptance.TestData) string {
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
`, r.templateLinux(data), data.RandomInteger)
}

func (r VirtualMachineScaleSetExtensionResource) extensionChaining(data acceptance.TestData) string {
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
`, r.templateLinux(data), data.RandomInteger, data.RandomInteger)
}

func (r VirtualMachineScaleSetExtensionResource) forceUpdateTag(data acceptance.TestData, tag string) string {
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
`, r.templateLinux(data), data.RandomInteger, tag)
}

func (r VirtualMachineScaleSetExtensionResource) updateVersion(data acceptance.TestData, version string) string {
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
`, r.templateLinux(data), data.RandomInteger, version)
}

func (r VirtualMachineScaleSetExtensionResource) protectedSettings(data acceptance.TestData) string {
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
`, r.templateLinux(data), data.RandomInteger)
}

func (r VirtualMachineScaleSetExtensionResource) protectedSettingsOnly(data acceptance.TestData) string {
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
`, r.templateLinux(data), data.RandomInteger)
}

func (r VirtualMachineScaleSetExtensionResource) requiresImport(data acceptance.TestData) string {
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
`, r.templateLinux(data))
}

func (VirtualMachineScaleSetExtensionResource) templateLinux(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineScaleSetExtensionResource) templateWindows(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
