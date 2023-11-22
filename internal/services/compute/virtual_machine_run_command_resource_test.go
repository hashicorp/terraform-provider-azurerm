package compute_test

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachineruncommands"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualMachineRunCommandTestResource struct{}

func TestAccVirtualMachineRunCommand_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_run_command", "test")
	r := VirtualMachineRunCommandTestResource{}

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

func TestAccVirtualMachineRunCommand_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_run_command", "test")
	r := VirtualMachineRunCommandTestResource{}

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

func TestAccVirtualMachineRunCommand_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_run_command", "test")
	r := VirtualMachineRunCommandTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("error_blob_managed_identity", "output_blob_managed_identity", "protected_parameter", "run_as_password", "source.0.script_uri_managed_identity"),
	})
}

func TestAccVirtualMachineRunCommand_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_run_command", "test")
	r := VirtualMachineRunCommandTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("error_blob_managed_identity", "output_blob_managed_identity", "protected_parameter", "run_as_password", "source.0.script_uri_managed_identity"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r VirtualMachineRunCommandTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualmachineruncommands.ParseVirtualMachineRunCommandID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.VirtualMachineRunCommandsClient.GetByVirtualMachine(ctx, *id, virtualmachineruncommands.DefaultGetByVirtualMachineOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r VirtualMachineRunCommandTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_virtual_machine_run_command" "test" {
  name               = "acctestvmrc-${var.random_string}"
  location           = azurerm_resource_group.test.location
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  source {
    script = "echo 'hello world'"
  }
}
`, r.template(data))
}

func (r VirtualMachineRunCommandTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_run_command" "import" {
  location           = azurerm_virtual_machine_run_command.test.location
  name               = azurerm_virtual_machine_run_command.test.name
  virtual_machine_id = azurerm_virtual_machine_run_command.test.virtual_machine_id
  source {
    script = azurerm_virtual_machine_run_command.test.source.0.script
  }
}
`, r.basic(data))
}

func (r VirtualMachineRunCommandTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc${var.random_string}"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc${var.random_integer}"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "test1" {
  name                   = "script1"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_content         = "echo 'hello world'"
}

resource "azurerm_storage_blob" "test2" {
  name                   = "output"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Append"
}

resource "azurerm_storage_blob" "test3" {
  name                   = "error"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Append"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_virtual_machine_run_command" "test" {
  location                = azurerm_resource_group.test.location
  name                    = "acctestvmrc-${var.random_string}"
  virtual_machine_id      = azurerm_linux_virtual_machine.test.id
  async_execution_enabled = false
  error_blob_uri          = azurerm_storage_blob.test3.id
  output_blob_uri         = azurerm_storage_blob.test2.id
  run_as_password         = "val-${var.random_string}"
  run_as_user             = "val-${var.random_string}"
  timeout_in_seconds      = 21

  error_blob_managed_identity {
    client_id = azurerm_user_assigned_identity.test.client_id
  }

  output_blob_managed_identity {
    client_id = azurerm_user_assigned_identity.test.client_id
  }

  parameter {
    name  = "acctestvmrc-${var.random_string}"
    value = "val-${var.random_string}"
  }

  protected_parameter {
    name  = "acctestvmrc-${var.random_string}"
    value = "val-${var.random_string}"
  }

  source {
    script_uri = azurerm_storage_blob.test1.id
    script_uri_managed_identity {
      client_id = azurerm_user_assigned_identity.test.client_id
    }
  }

  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }

  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, r.template(data))
}

func (r VirtualMachineRunCommandTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}
variable "random_integer" {
  default = %d
}
variable "random_string" {
  default = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-vmcmd-${var.random_integer}"
  location = var.primary_location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-${var.random_integer}"
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
  name                = "acctestnic-${var.random_integer}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-${var.random_integer}"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_B1ls"
  admin_username                  = "adminuser"
  admin_password                  = "P@$$w0rd1234!"
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

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-${var.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
