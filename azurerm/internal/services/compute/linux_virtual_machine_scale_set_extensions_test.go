package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccLinuxVirtualMachineScaleSet_extensionDoNotRunExtensionsOnOverProvisionedMachines(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionDoNotRunExtensionsOnOverProvisionedMachines(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_extensionDoNotRunExtensionsOnOverProvisionedMachinesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionDoNotRunExtensionsOnOverProvisionedMachines(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.extensionDoNotRunExtensionsOnOverProvisionedMachines(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.extensionDoNotRunExtensionsOnOverProvisionedMachines(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_extensionBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

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

func TestAccLinuxVirtualMachineScaleSet_extensionOnlySettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

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

func TestAccLinuxVirtualMachineScaleSet_extensionForceUpdateTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

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

func TestAccLinuxVirtualMachineScaleSet_extensionsMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionMultiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "extension.0.protected_settings", "extension.1.protected_settings"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_extensionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

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

func TestAccLinuxVirtualMachineScaleSet_extensionWithTimeBudget(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

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

func TestAccLinuxVirtualMachineScaleSet_extensionWithTimeBudgetUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

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

func TestAccLinuxVirtualMachineScaleSet_extensionTimeBudgetWithoutExtensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

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

func TestAccLinuxVirtualMachineScaleSet_extensionTimeBudgetWithoutExtensionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

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

func (r LinuxVirtualMachineScaleSetResource) extensionDoNotRunExtensionsOnOverProvisionedMachines(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  overprovision       = true

  disable_password_authentication                   = false
  do_not_run_extensions_on_overprovisioned_machines = %t

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
`, r.template(data), data.RandomInteger, enabled)
}

func (r LinuxVirtualMachineScaleSetResource) extensionOnlySettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

  }

  tags = {
    accTest = "true"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) extensionBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })

  }

  tags = {
    accTest = "true"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) extensionForceUpdateTag(data acceptance.TestData, updateTag string) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true
    force_update_tag           = %q

    settings = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })
  }

  tags = {
    accTest = "true"
  }
}
`, r.template(data), data.RandomInteger, updateTag)
}

func (r LinuxVirtualMachineScaleSetResource) extensionMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    provision_after_extensions = ["VMAccessForLinux"]

    settings = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })
  }

  extension {
    name                       = "VMAccessForLinux"
    publisher                  = "Microsoft.OSTCExtensions"
    type                       = "VMAccessForLinux"
    type_handler_version       = "1.5"
    auto_upgrade_minor_version = true

    protected_settings = jsonencode({
      "reset_ssh" = "True"
    })

  }

  tags = {
    accTest = "true"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) extensionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "echo $(date)"
    })
  }

  tags = {
    accTest = "true"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LinuxVirtualMachineScaleSetResource) extensionWithTimeBudget(data acceptance.TestData, duration string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

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

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = jsonencode({
      "commandToExecute" = "echo $HOSTNAME"
    })

    protected_settings = jsonencode({
      "managedIdentity" = {}
    })
  }

  extensions_time_budget = "%s"
}
`, template, data.RandomInteger, duration)
}

func (r LinuxVirtualMachineScaleSetResource) extensionTimeBudgetWithoutExtensions(data acceptance.TestData, duration string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

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

  extensions_time_budget = "%s"
}
`, template, data.RandomInteger, duration)
}
