// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PolicyVirtualMachineConfigurationAssignmentDataSource struct{}

func TestAccPolicyVirtualMachineConfigurationAssignmentDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_virtual_machine_configuration_assignment", "test")
	r := PolicyVirtualMachineConfigurationAssignmentDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("compliance_status").Exists(),
			),
		},
	})
}

func (r PolicyVirtualMachineConfigurationAssignmentDataSource) templateBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

locals {
  vm_name = "acctestvm%s"
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
  address_prefixes     = ["10.0.2.0/24"]
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
`, data.RandomString, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r PolicyVirtualMachineConfigurationAssignmentDataSource) template(data acceptance.TestData) string {
	tags := ""
	if strings.HasPrefix(strings.ToLower(data.Client().SubscriptionID), "85b3dbca") {
		tags = `
  tags = {
    "azsecpack"                                                                = "nonprod"
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true"
  }
`
	}
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

  identity {
    type = "SystemAssigned"
  }

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

%s
}
`, r.templateBase(data), tags)
}

func (r PolicyVirtualMachineConfigurationAssignmentDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_virtual_machine_configuration_assignment" "test" {
  name               = "AzureWindowsBaseline"
  location           = azurerm_windows_virtual_machine.test.location
  virtual_machine_id = azurerm_windows_virtual_machine.test.id

  configuration {
    assignment_type = "ApplyAndMonitor"
    version         = "1.*"

    parameter {
      name  = "Minimum Password Length;ExpectedValue"
      value = "16"
    }
    parameter {
      name  = "Minimum Password Age;ExpectedValue"
      value = "0"
    }
    parameter {
      name  = "Maximum Password Age;ExpectedValue"
      value = "30,45"
    }
    parameter {
      name  = "Enforce Password History;ExpectedValue"
      value = "10"
    }
    parameter {
      name  = "Password Must Meet Complexity Requirements;ExpectedValue"
      value = "1"
    }
  }
}

data "azurerm_policy_virtual_machine_configuration_assignment" "test" {
  name                 = azurerm_policy_virtual_machine_configuration_assignment.test.name
  resource_group_name  = azurerm_resource_group.test.name
  virtual_machine_name = azurerm_windows_virtual_machine.test.name
}
`, r.template(data))
}
