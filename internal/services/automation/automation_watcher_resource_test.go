// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WatcherResource struct{}

func TestAccWatcher_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.WatcherResource{}.ResourceType(), "test")
	r := WatcherResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("execution_frequency_in_seconds").HasValue("2"),
			),
		},
		data.ImportStep("tags", "etag", "location"),
	})
}

func TestAccWatcher_update(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.WatcherResource{}.ResourceType(), "test")
	r := WatcherResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("execution_frequency_in_seconds").HasValue("2"),
			),
		},
		data.ImportStep("tags", "etag", "location"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("execution_frequency_in_seconds").HasValue("20"),
			),
		},
		data.ImportStep("tags", "etag", "location"),
	})
}

func (a WatcherResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := watcher.ParseWatcherID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automation.WatcherClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (a WatcherResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_watcher" "test" {
  automation_account_id = azurerm_automation_account.test.id
  name                  = "acctest-watcher-%[2]d"
  location              = "%[3]s"

  tags = {
    foo = "bar"
  }

  script_parameters = {
    param_foo = "arg_bar"
  }

  script_run_on                  = azurerm_automation_hybrid_runbook_worker_group.test.name
  description                    = "example-watcher desc"
  etag                           = "etag example"
  script_name                    = azurerm_automation_runbook.test.name
  execution_frequency_in_seconds = 2
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a WatcherResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_watcher" "test" {
  automation_account_id = azurerm_automation_account.test.id
  name                  = "acctest-watcher-%[2]d"
  location              = "%[3]s"

  tags = {
    "foo" = "bar"
  }

  script_parameters = {
    foo = "bar"
  }

  etag                           = "etag example"
  execution_frequency_in_seconds = 20
  script_name                    = azurerm_automation_runbook.test.name
  script_run_on                  = azurerm_automation_hybrid_runbook_worker_group.test.name
  description                    = "example-watcher desc"
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a WatcherResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_credential" "test" {
  name                    = "acctest-%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  username                = "test_user"
  password                = "test_pwd"
}

resource "azurerm_automation_hybrid_runbook_worker_group" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  name                    = "acctest-%[1]d"
  credential_name         = azurerm_automation_credential.test.name
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
  name                = "acctni-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"

  disable_password_authentication = false

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

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

  tags = {
    azsecpack                                                                  = "nonprod"
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true"
  }
}

resource "azurerm_automation_runbook" "test" {
  name                    = "acc-runbook-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  content = <<CONTENT
# Some test content
# for Terraform acceptance test
CONTENT
  tags = {
    ENV = "runbook_test"
  }
}
`, data.RandomInteger, data.Locations.Primary, uuid.New().String())
}
