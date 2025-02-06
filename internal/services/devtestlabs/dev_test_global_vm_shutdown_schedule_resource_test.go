// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devtestlabs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/globalschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DevTestGlobalVMShutdownScheduleResource struct{}

func TestAccDevTestGlobalVMShutdownSchedule_autoShutdownBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_global_vm_shutdown_schedule", "test")
	r := DevTestGlobalVMShutdownScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoShutdownBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("timezone").HasValue("Pacific Standard Time"),
				check.That(data.ResourceName).Key("notification_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification_settings.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("notification_settings.0.time_in_minutes").HasValue("30"),
				check.That(data.ResourceName).Key("notification_settings.0.webhook_url").HasValue(""),
				check.That(data.ResourceName).Key("daily_recurrence_time").HasValue("0100"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDevTestGlobalVMShutdownSchedule_autoShutdownComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_global_vm_shutdown_schedule", "test")
	r := DevTestGlobalVMShutdownScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoShutdownComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
				check.That(data.ResourceName).Key("timezone").HasValue("Central Standard Time"),
				check.That(data.ResourceName).Key("notification_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("notification_settings.0.time_in_minutes").HasValue("15"),
				check.That(data.ResourceName).Key("notification_settings.0.webhook_url").HasValue("https://www.bing.com/2/4"),
				check.That(data.ResourceName).Key("notification_settings.0.email").HasValue("alerts@devtest.com"),
				check.That(data.ResourceName).Key("daily_recurrence_time").HasValue("1100"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("Production"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDevTestGlobalVMShutdownSchedule_autoShutdownUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_global_vm_shutdown_schedule", "test")
	r := DevTestGlobalVMShutdownScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoShutdownBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("timezone").HasValue("Pacific Standard Time"),
				check.That(data.ResourceName).Key("notification_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification_settings.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("notification_settings.0.time_in_minutes").HasValue("30"),
				check.That(data.ResourceName).Key("notification_settings.0.webhook_url").HasValue(""),
				check.That(data.ResourceName).Key("daily_recurrence_time").HasValue("0100"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoShutdownComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
				check.That(data.ResourceName).Key("timezone").HasValue("Central Standard Time"),
				check.That(data.ResourceName).Key("notification_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("notification_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("notification_settings.0.time_in_minutes").HasValue("15"),
				check.That(data.ResourceName).Key("notification_settings.0.webhook_url").HasValue("https://www.bing.com/2/4"),
				check.That(data.ResourceName).Key("daily_recurrence_time").HasValue("1100"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("Production"),
			),
		},
	})
}

func (DevTestGlobalVMShutdownScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := globalschedules.ParseScheduleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevTestLabs.GlobalLabSchedulesClient.Get(ctx, *id, globalschedules.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (DevTestGlobalVMShutdownScheduleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dtl-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVN-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestSN-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctestNIC-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                  = "acctestVM-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  size                  = "Standard_B2s"

  admin_username                  = "testadmin"
  admin_password                  = "Password1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    name                 = "myosdisk-%d"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r DevTestGlobalVMShutdownScheduleResource) autoShutdownBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_global_vm_shutdown_schedule" "test" {
  location              = azurerm_resource_group.test.location
  virtual_machine_id    = azurerm_linux_virtual_machine.test.id
  daily_recurrence_time = "0100"
  timezone              = "Pacific Standard Time"

  notification_settings {
    enabled = false
  }

  tags = {
    environment = "Production"
  }
}
`, r.template(data))
}

func (r DevTestGlobalVMShutdownScheduleResource) autoShutdownComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_global_vm_shutdown_schedule" "test" {
  location           = azurerm_resource_group.test.location
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  enabled            = false

  daily_recurrence_time = "1100"
  timezone              = "Central Standard Time"

  notification_settings {
    time_in_minutes = 15
    webhook_url     = "https://www.bing.com/2/4"
    email           = "alerts@devtest.com"
    enabled         = true
  }

  tags = {
    Environment = "Production"
  }
}
`, r.template(data))
}
