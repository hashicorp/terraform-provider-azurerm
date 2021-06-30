package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccLinuxVirtualMachine_otherAllowExtensionOperationsDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherAllowExtensionOperationsDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_extension_operations").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherAllowExtensionOperationsDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherAllowExtensionOperationsDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_extension_operations").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherAllowExtensionOperationsUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherAllowExtensionOperationsDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_extension_operations").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.otherAllowExtensionOperationsDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_extension_operations").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherAllowExtensionOperationsUpdatedWithoutVmAgent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherAllowExtensionOperationsDisabledWithoutVmAgent(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_extension_operations").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.otherAllowExtensionOperationsEnabledWithoutVmAgent(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("allow_extension_operations").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherExtensionsTimeBudget(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherExtensionsTimeBudget(data, "PT30M"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("extensions_time_budget").HasValue("PT30M"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherExtensionsTimeBudgetUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherExtensionsTimeBudget(data, "PT30M"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("extensions_time_budget").HasValue("PT30M"),
			),
		},
		data.ImportStep(),
		{
			Config: r.otherExtensionsTimeBudget(data, "PT50M"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("extensions_time_budget").HasValue("PT50M"),
			),
		},
		data.ImportStep(),
		{
			Config: r.otherExtensionsTimeBudget(data, "PT30M"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("extensions_time_budget").HasValue("PT30M"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherBootDiagnostics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Enabled
			Config: r.otherBootDiagnostics(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Disabled
			Config: r.otherBootDiagnosticsDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Enabled
			Config: r.otherBootDiagnostics(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherBootDiagnosticsManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Enabled
			Config: r.otherBootDiagnosticsManaged(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Disabled
			Config: r.otherBootDiagnosticsDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Enabled
			Config: r.otherBootDiagnosticsManaged(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherComputerNameDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherComputerNameDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("computer_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherComputerNameDefaultInvalid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.otherComputerNameDefaultInvalid(data),
			ExpectError: regexp.MustCompile("unable to assume default computer name"),
		},
	})
}

func TestAccLinuxVirtualMachine_otherComputerNameCustom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherComputerNameCustom(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("computer_name").HasValue("custom-linux-hostname-123"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherCustomData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherCustomData(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("custom_data"),
	})
}

func TestAccLinuxVirtualMachine_otherSkipShutdownAndForceDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherSkipShutdownAndForceDelete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherLicenseType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherLicenseType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("license_type").HasValue("SLES_BYOS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherPrioritySpot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherPrioritySpot(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherPrioritySpotMaxBidPrice(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// expensive, but guarantees this test will pass
			Config: r.otherPrioritySpotMaxBidPrice(data, "0.5000"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// no limit
			Config: r.otherPrioritySpotMaxBidPrice(data, "-1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherProvisionVMAgentDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherProvisionVMAgentDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("provision_vm_agent").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherProvisionVMAgentDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherProvisionVMAgentDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("provision_vm_agent").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSSH(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccLinuxVirtualMachine_otherSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherSecret(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.otherSecretUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.otherSecretRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("secret.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.otherTagsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherUltraSsdDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherUltraSsd(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("additional_capabilities.0.ultra_ssd_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherUltraSsdEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherUltraSsd(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("additional_capabilities.0.ultra_ssd_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherUltraSsdUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherUltraSsd(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("additional_capabilities.0.ultra_ssd_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.otherUltraSsd(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("additional_capabilities.0.ultra_ssd_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherEncryptionAtHostEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherEncryptionAtHostEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherEncryptionAtHostEnabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherEncryptionAtHostEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.otherEncryptionAtHostEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.otherEncryptionAtHostEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherEncryptionAtHostEnabledWithCMK(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherEncryptionAtHostEnabledWithCMK(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxVirtualMachine_otherGracefulShutdownDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherGracefulShutdown(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccLinuxVirtualMachine_otherGracefulShutdownEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherGracefulShutdown(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r LinuxVirtualMachineResource) otherAllowExtensionOperationsDefault(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherAllowExtensionOperationsDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                       = "acctestVM-%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  size                       = "Standard_F2"
  admin_username             = "adminuser"
  allow_extension_operations = false
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherAllowExtensionOperationsDisabledWithoutVmAgent(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                       = "acctestVM-%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  size                       = "Standard_F2"
  admin_username             = "adminuser"
  allow_extension_operations = false
  provision_vm_agent         = false
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherAllowExtensionOperationsEnabledWithoutVmAgent(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                       = "acctestVM-%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  size                       = "Standard_F2"
  admin_username             = "adminuser"
  allow_extension_operations = true
  provision_vm_agent         = false
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherExtensionsTimeBudget(data acceptance.TestData, duration string) string {
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

  extensions_time_budget = "%s"
}
`, r.template(data), data.RandomInteger, duration)
}

func (r LinuxVirtualMachineResource) otherBootDiagnostics(data acceptance.TestData) string {
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

  boot_diagnostics {
    storage_account_uri = azurerm_storage_account.test.primary_blob_endpoint
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
`, r.otherBootDiagnosticsTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherBootDiagnosticsManaged(data acceptance.TestData) string {
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

  boot_diagnostics {}

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
`, r.otherBootDiagnosticsTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherBootDiagnosticsDisabled(data acceptance.TestData) string {
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
`, r.otherBootDiagnosticsTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherBootDiagnosticsTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "accsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, r.template(data), data.RandomString)
}

func (r LinuxVirtualMachineResource) otherComputerNameDefault(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherComputerNameDefaultInvalid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-this-name-too-long-to-be-a-linux-vm-computer-name-1234567890"
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
`, r.template(data))
}

func (r LinuxVirtualMachineResource) otherComputerNameCustom(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  computer_name       = "custom-linux-hostname-123"
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherCustomData(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  custom_data         = base64encode("/bin/bash")
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherSkipShutdownAndForceDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    virtual_machine {
      skip_shutdown_and_force_delete = true
    }
  }
}

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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherLicenseType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  license_type        = "SLES_BYOS"
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherPrioritySpot(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  eviction_policy     = "Deallocate"
  priority            = "Spot"
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherPrioritySpotMaxBidPrice(data acceptance.TestData, maxBidPrice string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  eviction_policy     = "Deallocate"
  priority            = "Spot"
  max_bid_price       = %q
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
`, r.template(data), data.RandomInteger, maxBidPrice)
}

func (r LinuxVirtualMachineResource) otherProvisionVMAgentDefault(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherProvisionVMAgentDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                       = "acctestVM-%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  size                       = "Standard_F2"
  provision_vm_agent         = false
  allow_extension_operations = false
  admin_username             = "adminuser"
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
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "import" {
  name                  = azurerm_linux_virtual_machine.test.name
  resource_group_name   = azurerm_linux_virtual_machine.test.resource_group_name
  location              = azurerm_linux_virtual_machine.test.location
  size                  = azurerm_linux_virtual_machine.test.size
  admin_username        = azurerm_linux_virtual_machine.test.admin_username
  network_interface_ids = azurerm_linux_virtual_machine.test.network_interface_ids

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
`, r.authSSH(data))
}

func (r LinuxVirtualMachineResource) otherSecretTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                        = "standard"
  enabled_for_template_deployment = true
  enabled_for_deployment          = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]

    storage_permissions = [
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-first"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_certificate" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-second"
      validity_in_months = 12
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r LinuxVirtualMachineResource) otherSecret(data acceptance.TestData) string {
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

  secret {
    key_vault_id = azurerm_key_vault.test.id

    certificate {
      url = azurerm_key_vault_certificate.first.secret_id
    }
  }
}
`, r.otherSecretTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherSecretRemoved(data acceptance.TestData) string {
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
`, r.otherSecretTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherSecretUpdated(data acceptance.TestData) string {
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

  secret {
    key_vault_id = azurerm_key_vault.test.id

    certificate {
      url = azurerm_key_vault_certificate.first.secret_id
    }

    certificate {
      url = azurerm_key_vault_certificate.second.secret_id
    }
  }
}
`, r.otherSecretTemplate(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherTags(data acceptance.TestData) string {
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

  tags = {
    Hello = "World"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherTagsUpdated(data acceptance.TestData) string {
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

  tags = {
    Hello     = "World"
    Dimension = "C-137"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineResource) otherUltraSsd(data acceptance.TestData, ultraSsdEnabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3"
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

  additional_capabilities {
    ultra_ssd_enabled = %t
  }
}
`, r.template(data), data.RandomInteger, ultraSsdEnabled)
}

func (r LinuxVirtualMachineResource) otherEncryptionAtHostEnabled(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_DS3_v2"
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

  encryption_at_host_enabled = %t
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r LinuxVirtualMachineResource) otherEncryptionAtHostEnabledWithCMK(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_DS3_v2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching                = "ReadWrite"
    storage_account_type   = "Standard_LRS"
    disk_encryption_set_id = azurerm_disk_encryption_set.test.id
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  encryption_at_host_enabled = %t

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
  ]
}
`, r.diskOSDiskDiskEncryptionSetResource(data), data.RandomInteger, enabled)
}

func (LinuxVirtualMachineResource) otherGracefulShutdown(data acceptance.TestData, gracefulShutdown bool) string {
	return fmt.Sprintf(`
locals {
  first_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
}

provider "azurerm" {
  features {
    virtual_machine {
      graceful_shutdown = %t
    }
  }
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
`, gracefulShutdown, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
