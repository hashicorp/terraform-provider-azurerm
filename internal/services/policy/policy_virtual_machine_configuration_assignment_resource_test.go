// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2020-06-25/guestconfigurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PolicyVirtualMachineConfigurationAssignmentResource struct{}

func TestAccPolicyVirtualMachineConfigurationAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_virtual_machine_configuration_assignment", "test")
	r := PolicyVirtualMachineConfigurationAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPolicyVirtualMachineConfigurationAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_virtual_machine_configuration_assignment", "test")
	r := PolicyVirtualMachineConfigurationAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPolicyVirtualMachineConfigurationAssignment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_virtual_machine_configuration_assignment", "test")
	r := PolicyVirtualMachineConfigurationAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateGuestConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PolicyVirtualMachineConfigurationAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := guestconfigurationassignments.ParseProviders2GuestConfigurationAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Policy.GuestConfigurationAssignmentsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r PolicyVirtualMachineConfigurationAssignmentResource) templateBase(data acceptance.TestData) string {
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

func (r PolicyVirtualMachineConfigurationAssignmentResource) template(data acceptance.TestData) string {
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

func (r PolicyVirtualMachineConfigurationAssignmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_virtual_machine_configuration_assignment" "test" {
  name     = "WhitelistedApplication"
  location = azurerm_windows_virtual_machine.test.location

  virtual_machine_id = azurerm_windows_virtual_machine.test.id

  configuration {
    version = "1.*"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
`, r.template(data))
}

func (r PolicyVirtualMachineConfigurationAssignmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_virtual_machine_configuration_assignment" "import" {
  name               = azurerm_policy_virtual_machine_configuration_assignment.test.name
  location           = azurerm_policy_virtual_machine_configuration_assignment.test.location
  virtual_machine_id = azurerm_policy_virtual_machine_configuration_assignment.test.virtual_machine_id

  configuration {
    version = "1.*"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
`, r.basic(data))
}

func (r PolicyVirtualMachineConfigurationAssignmentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_virtual_machine_configuration_assignment" "test" {
  name               = "WhitelistedApplication"
  location           = azurerm_windows_virtual_machine.test.location
  virtual_machine_id = azurerm_windows_virtual_machine.test.id

  configuration {
    version         = "1.1.1.1"
    assignment_type = "ApplyAndAutoCorrect"
    content_hash    = "testcontenthash"
    content_uri     = "https://testcontenturi/package"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
`, r.template(data))
}

func (r PolicyVirtualMachineConfigurationAssignmentResource) updateGuestConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_virtual_machine_configuration_assignment" "test" {
  name               = "WhitelistedApplication"
  location           = azurerm_windows_virtual_machine.test.location
  virtual_machine_id = azurerm_windows_virtual_machine.test.id

  configuration {
    version         = "1.1.1.1"
    assignment_type = "Audit"
    content_hash    = "testcontenthash2"
    content_uri     = "https://testcontenturi/package2"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
`, r.template(data))
}
