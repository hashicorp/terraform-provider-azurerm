package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type VirtualMachinePowerAction struct{}

func TestAccVirtualMachinePowerAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_power", "test")
	a := VirtualMachinePowerAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.restart(data),
				Check:  nil, // TODO - plugin-testing release?
			},
		},
	})
}

func TestAccVirtualMachinePowerAction_tryTurningItOffAndOnAgain(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_power", "test")
	a := VirtualMachinePowerAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.techSupport(data, "create"),
			},
			{
				Config: a.techSupport(data, "update"),
				Check:  nil, // TODO - plugin-testing release?
			},
		},
	})
}

func (a *VirtualMachinePowerAction) restart(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[2]d"
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
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  lifecycle {
    action_trigger {
      events  = [after_create, after_update] // Restart the vm after create and update
      actions = [action.azurerm_virtual_machine_power.test]
    }
  }
}

data "azurerm_virtual_machine" "test" {
  name                = "acctestVM-%[2]d" // sidestep cyclic reference issue
  resource_group_name = azurerm_resource_group.test.name
}

action "azurerm_virtual_machine_power" "test" {
  config {
    virtual_machine_id = data.azurerm_virtual_machine.test.id
    power_action       = "restart"
  }
}
`, a.templateLinux(data), data.RandomInteger)
}

func (a *VirtualMachinePowerAction) techSupport(data acceptance.TestData, tagVal string) string {
	return fmt.Sprintf(`

%[1]s

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

  tags = {
    triggerme = %[3]s
  }

  lifecycle {
    action_trigger {
      events  = [before_update]
      actions = [action.azurerm_virtual_machine_power.power_off]
    }

    action_trigger {
      events  = [after_update]
      actions = [action.azurerm_virtual_machine_power.power_on]
    }
  }
}

data "azurerm_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
}

resource "terraform_data" "trigger" {
  input = azurerm_linux_virtual_machine.test.tags
  lifecycle {
    action_trigger {
      events  = [before_update]
      actions = [action.azurerm_virtual_machine_power.power_off]
    }
    action_trigger {
      events  = [after_update]
      actions = [action.azurerm_virtual_machine_power.power_on]
    }
  }
}


action "azurerm_virtual_machine_power" "power_off" {
	config {
		virtual_machine_id = azurerm_windows_virtual_machine.test.id
		power_action = "power_off"
	}
}

action "azurerm_virtual_machine_power" "power_on" {
	config {
		virtual_machine_id = azurerm_windows_virtual_machine.test.id
		power_action = "power_on"
	}
}
`, a.templateWindows(data), data.RandomInteger, tagVal)
}

func (a *VirtualMachinePowerAction) templateLinux(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[3]s

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, LinuxVirtualMachineResource{}.templateBasePublicKey())
}

func (a *VirtualMachinePowerAction) templateWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

locals {
  vm_name = "acctestvm%[1]s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[3]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomString, data.RandomInteger, data.Locations.Primary)
}
