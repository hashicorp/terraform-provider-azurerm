// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"regexp"
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

func TestAccVirtualMachineRunCommand_recreate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_run_command", "test")
	r := VirtualMachineRunCommandTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basicWithScriptError(data),
			ExpectError: regexp.MustCompile("running the command"),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineRunCommand_sourceCommandId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_run_command", "test")
	r := VirtualMachineRunCommandTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sourceCommandId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineRunCommand_storageBlobSAS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_run_command", "test")
	r := VirtualMachineRunCommandTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageBlobSAS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("error_blob_managed_identity", "output_blob_managed_identity", "protected_parameter", "run_as_password", "source.0.script_uri_managed_identity", "source.0.script_uri", "error_blob_uri", "output_blob_uri"),
	})
}

func TestAccVirtualMachineRunCommand_storageBlobSystemIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_run_command", "test")
	r := VirtualMachineRunCommandTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageBlobSystemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("error_blob_managed_identity", "output_blob_managed_identity", "protected_parameter", "run_as_password", "source.0.script_uri_managed_identity"),
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

func TestAccVirtualMachineRunCommand_updateParameters(t *testing.T) {
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
			Config: r.basicWithParameters(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_parameter"),
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

func (r VirtualMachineRunCommandTestResource) basicWithParameters(data acceptance.TestData) string {
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

  parameter {
    name  = "acctestvmrc-${var.random_string}"
    value = "val-${var.random_string}"
  }

  protected_parameter {
    name  = "acctestvmrc2-${var.random_string}"
    value = "val-${var.random_string}"
  }
}
`, r.template(data))
}

func (r VirtualMachineRunCommandTestResource) basicWithScriptError(data acceptance.TestData) string {
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
    script = "echo1 'hello world'"
  }
}
`, r.template(data))
}

func (r VirtualMachineRunCommandTestResource) sourceCommandId(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_virtual_machine_run_command" "pretest" {
  name               = "acctestvmrc-pre-${var.random_string}"
  location           = azurerm_resource_group.test.location
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  source {
    script = "sudo apt update && sudo apt install -y net-tools"
  }
}

resource "azurerm_virtual_machine_run_command" "test" {
  name               = "acctestvmrc-${var.random_string}"
  location           = azurerm_resource_group.test.location
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  source {
    command_id = "ifconfig"
  }
  depends_on = [
    azurerm_virtual_machine_run_command.pretest,
  ]
}
`, r.template(data))
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
  location           = azurerm_resource_group.test.location
  name               = "acctestvmrc-${var.random_string}"
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  run_as_password    = "Pa-${var.random_string}"
  run_as_user        = "adminuser"
  error_blob_uri     = azurerm_storage_blob.test3.id
  output_blob_uri    = azurerm_storage_blob.test2.id

  error_blob_managed_identity {
    client_id = azurerm_user_assigned_identity.test.client_id
  }

  output_blob_managed_identity {
    client_id = azurerm_user_assigned_identity.test.client_id
  }

  source {
    script_uri = azurerm_storage_blob.test1.id
    script_uri_managed_identity {
      client_id = azurerm_user_assigned_identity.test.client_id
    }
  }

  parameter {
    name  = "acctestvmrc-${var.random_string}"
    value = "val-${var.random_string}"
  }

  protected_parameter {
    name  = "acctestvmrc-${var.random_string}"
    value = "val-${var.random_string}"
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

func (r VirtualMachineRunCommandTestResource) storageBlobSystemIdentity(data acceptance.TestData) string {
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
  principal_id         = azurerm_linux_virtual_machine.test.identity[0].principal_id
}

resource "azurerm_virtual_machine_run_command" "test" {
  location           = azurerm_resource_group.test.location
  name               = "acctestvmrc-${var.random_string}"
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  run_as_password    = "Pa-${var.random_string}"
  run_as_user        = "adminuser"
  error_blob_uri     = azurerm_storage_blob.test3.id
  output_blob_uri    = azurerm_storage_blob.test2.id

  source {
    script_uri = azurerm_storage_blob.test1.id
    script_uri_managed_identity {
      client_id = azurerm_linux_virtual_machine.test.identity[0].principal_id
    }
  }

  parameter {
    name  = "acctestvmrc-${var.random_string}"
    value = "val-${var.random_string}"
  }

  protected_parameter {
    name  = "acctestvmrc-${var.random_string}"
    value = "val-${var.random_string}"
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

func (r VirtualMachineRunCommandTestResource) storageBlobSAS(data acceptance.TestData) string {
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

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true
  signed_version    = "2019-10-10"
  start             = "2023-04-01T00:00:00Z"
  expiry            = "2123-04-01T00:00:00Z"

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  permissions {
    read    = true
    write   = true
    delete  = false
    list    = false
    add     = true
    create  = true
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

resource "azurerm_virtual_machine_run_command" "test" {
  location           = azurerm_resource_group.test.location
  name               = "acctestvmrc-${var.random_string}"
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  run_as_password    = "Pa-${var.random_string}"
  run_as_user        = "adminuser"
  error_blob_uri     = "${azurerm_storage_blob.test3.id}${data.azurerm_storage_account_sas.test.sas}"
  output_blob_uri    = "${azurerm_storage_blob.test2.id}${data.azurerm_storage_account_sas.test.sas}"

  source {
    script_uri = "${azurerm_storage_blob.test1.id}${data.azurerm_storage_account_sas.test.sas}"
  }

  parameter {
    name  = "acctestvmrc-${var.random_string}"
    value = "val-${var.random_string}"
  }

  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
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
  size                            = "Standard_B2s"
  admin_username                  = "adminuser"
  admin_password                  = "Pa-${var.random_string}"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Premium_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
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
