package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccWindowsVirtualMachineScaleSet_extensionDoNotRunOnOverProvisionedMachines(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionDoNotRunOnOverProvisionedMachines(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachineScaleSet_extensionsDoNotRunOnOverProvisionedMachinesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionDoNotRunOnOverProvisionedMachines(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.extensionDoNotRunOnOverProvisionedMachines(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.extensionDoNotRunOnOverProvisionedMachines(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccWindowsVirtualMachineScaleSet_extensionBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_extensionForceUpdateTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionForceUpdateTag(data, "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
		{
			Config: r.extensionForceUpdateTag(data, "second"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_extensionMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionMultiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_extensionOnlySettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionOnlySettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_extensionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
		{
			Config: r.extensionUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
		{
			Config: r.extensionBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_extensionWithTimeBudget(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionWithTimeBudget(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_extensionWithTimeBudgetUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionWithTimeBudget(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
		{
			Config: r.extensionWithTimeBudget(data, "PT1H"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
		{
			Config: r.extensionWithTimeBudget(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_extensionTimeBudgetWithoutExtensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionTimeBudgetWithoutExtensions(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_extensionTimeBudgetWithoutExtensionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionTimeBudgetWithoutExtensions(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
		{
			Config: r.extensionTimeBudgetWithoutExtensions(data, "PT1H"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
		{
			Config: r.extensionTimeBudgetWithoutExtensions(data, "PT30M"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings"),
	})
}

func (r WindowsVirtualMachineScaleSetResource) extensionDoNotRunOnOverProvisionedMachines(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  overprovision       = true

  do_not_run_extensions_on_overprovisioned_machines = %t

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
`, r.template(data), enabled)
}

func (r WindowsVirtualMachineScaleSetResource) extensionBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Compute"
    type                       = "CustomScriptExtension"
    type_handler_version       = "1.10"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "powershell.exe -c \"Get-Content env:computername\""
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) extensionForceUpdateTag(data acceptance.TestData, updateTag string) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Compute"
    type                       = "CustomScriptExtension"
    type_handler_version       = "1.10"
    auto_upgrade_minor_version = true
    force_update_tag           = %q

    settings = jsonencode({
      "commandToExecute" = "powershell.exe -c \"Get-Content env:computername\""
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })
  }
}
`, r.template(data), updateTag)
}

func (r WindowsVirtualMachineScaleSetResource) extensionMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Compute"
    type                       = "CustomScriptExtension"
    type_handler_version       = "1.10"
    auto_upgrade_minor_version = true

    provision_after_extensions = ["AADLoginForWindows"]

    settings = jsonencode({
      "commandToExecute" = "powershell.exe -c \"Get-Content env:computername\""
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })
  }

  extension {
    name                       = "AADLoginForWindows"
    publisher                  = "Microsoft.Azure.ActiveDirectory"
    type                       = "AADLoginForWindows"
    type_handler_version       = "1.0"
    auto_upgrade_minor_version = true
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) extensionOnlySettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Compute"
    type                       = "CustomScriptExtension"
    type_handler_version       = "1.10"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "powershell.exe -c \"Get-Content env:computername\""
    })
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) extensionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Compute"
    type                       = "CustomScriptExtension"
    type_handler_version       = "1.10"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "powershell.exe -c \"Get-Process | Where-Object { $_.CPU -gt 10000 }\""
    })
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) extensionWithTimeBudget(data acceptance.TestData, duration string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Compute"
    type                       = "CustomScriptExtension"
    type_handler_version       = "1.10"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "powershell.exe -c \"Get-Content env:computername\""
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })
  }

  extensions_time_budget = "%s"
}
`, template, duration)
}

func (r WindowsVirtualMachineScaleSetResource) extensionTimeBudgetWithoutExtensions(data acceptance.TestData, duration string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

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

  extensions_time_budget = "%s"
}
`, template, duration)
}
