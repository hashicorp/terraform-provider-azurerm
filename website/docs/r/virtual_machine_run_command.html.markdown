---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_run_command"
description: |-
  Manages a Virtual Machine Run Command.
---

# azurerm_virtual_machine_run_command

Manages a Virtual Machine Run Command.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-uai"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_linux_virtual_machine" "example" {
  name                            = "example-VM"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  size                            = "Standard_B2s"
  admin_username                  = "adminuser"
  admin_password                  = "P@$$w0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.example.id,
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
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_storage_account" "example" {
  name                     = "exampleaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "example-sc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "example1" {
  name                   = "script1"
  storage_account_name   = azurerm_storage_account.example.name
  storage_container_name = azurerm_storage_container.example.name
  type                   = "Block"
  source_content         = "echo 'hello world'"
}

resource "azurerm_storage_blob" "example2" {
  name                   = "output"
  storage_account_name   = azurerm_storage_account.example.name
  storage_container_name = azurerm_storage_container.example.name
  type                   = "Append"
}

resource "azurerm_storage_blob" "example3" {
  name                   = "error"
  storage_account_name   = azurerm_storage_account.example.name
  storage_container_name = azurerm_storage_container.example.name
  type                   = "Append"
}

data "azurerm_storage_account_sas" "example" {
  connection_string = azurerm_storage_account.example.primary_connection_string
  https_only        = true
  signed_version    = "2019-10-10"
  start             = "2023-04-01T00:00:00Z"
  expiry            = "2024-04-01T00:00:00Z"

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

# basic example
resource "azurerm_virtual_machine_run_command" "example" {
  name               = "example-vmrc"
  location           = azurerm_resource_group.example.location
  virtual_machine_id = azurerm_linux_virtual_machine.example.id
  source {
    script = "echo 'hello world'"
  }
}

# authorize to storage blob using user assigned identity
resource "azurerm_virtual_machine_run_command" "example2" {
  location           = azurerm_resource_group.example.location
  name               = "example2-vmrc"
  virtual_machine_id = azurerm_linux_virtual_machine.example.id
  output_blob_uri    = azurerm_storage_blob.example2.id
  error_blob_uri     = azurerm_storage_blob.example3.id
  run_as_password    = "P@$$w0rd1234!"
  run_as_user        = "adminuser"

  source {
    script_uri = azurerm_storage_blob.example1.id
    script_uri_managed_identity {
      client_id = azurerm_user_assigned_identity.example.client_id
    }
  }

  error_blob_managed_identity {
    client_id = azurerm_user_assigned_identity.example.client_id
  }

  output_blob_managed_identity {
    client_id = azurerm_user_assigned_identity.example.client_id
  }

  parameter {
    name  = "examplev1"
    value = "val1"
  }

  protected_parameter {
    name  = "examplev2"
    value = "val2"
  }

  tags = {
    environment = "terraform-examples"
    some_key    = "some-value"
  }

  depends_on = [
    azurerm_role_assignment.example,
  ]
}

# authorize to storage blob using SAS token
resource "azurerm_virtual_machine_run_command" "example3" {
  location           = azurerm_resource_group.example.location
  name               = "example3-vmrc"
  virtual_machine_id = azurerm_linux_virtual_machine.example.id
  run_as_password    = "P@$$w0rd1234!"
  run_as_user        = "adminuser"
  error_blob_uri     = "${azurerm_storage_blob.example3.id}${data.azurerm_storage_account_sas.example.sas}"
  output_blob_uri    = "${azurerm_storage_blob.example2.id}${data.azurerm_storage_account_sas.example.sas}"

  source {
    script_uri = "${azurerm_storage_blob.example1.id}${data.azurerm_storage_account_sas.example.sas}"
  }

  parameter {
    name  = "example-vm1"
    value = "val1"
  }

  tags = {
    environment = "terraform-example-s"
    some_key    = "some-value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Virtual Machine Run Command should exist. Changing this forces a new Virtual Machine Run Command to be created.

* `name` - (Required) Specifies the name of this Virtual Machine Run Command. Changing this forces a new Virtual Machine Run Command to be created.

* `virtual_machine_id` - (Required) Specifies the Virtual Machine ID within which this Virtual Machine Run Command should exist. Changing this forces a new Virtual Machine Run Command to be created.

* `source` - (Required) A `source` block as defined below. The source of the run command script.

* `error_blob_managed_identity` - (Optional) An `error_blob_managed_identity` block as defined below. User-assigned managed Identity that has access to errorBlobUri storage blob.

* `error_blob_uri` - (Optional) Specifies the Azure storage blob where script error stream will be uploaded.

* `output_blob_managed_identity` - (Optional) An `output_blob_managed_identity` block as defined below. User-assigned managed Identity that has access to outputBlobUri storage blob.

* `output_blob_uri` - (Optional) Specifies the Azure storage blob where script output stream will be uploaded. It can be basic blob URI with SAS token.

* `parameter` - (Optional) A list of `parameter` blocks as defined below. The parameters used by the script.

* `protected_parameter` - (Optional) A list of `protected_parameter` blocks as defined below. The protected parameters used by the script.

* `run_as_password` - (Optional) Specifies the user account password on the VM when executing the Virtual Machine Run Command.

* `run_as_user` - (Optional) Specifies the user account on the VM when executing the Virtual Machine Run Command.

* `tags` - (Optional) A mapping of tags which should be assigned to the Virtual Machine Run Command.

---

An `error_blob_managed_identity` block supports the following arguments:

* `client_id` - (Optional) The client ID of the managed identity.
* `object_id` - (Optional) The object ID of the managed identity.

---

An `output_blob_managed_identity` block supports the following arguments:

* `client_id` - (Optional) The client ID of the managed identity.
* `object_id` - (Optional) The object ID of the managed identity.

---

A `parameter` block supports the following arguments:

* `name` - (Required) The run parameter name.
* `value` - (Required) The run parameter value.

---

A `protected_parameter` block supports the following arguments:

* `name` - (Required) The run parameter name.
* `value` - (Required) The run parameter value.

---

A `script_uri_managed_identity` block supports the following arguments:

* `client_id` - (Optional) The client ID of the managed identity.
* `object_id` - (Optional) The object ID of the managed identity.

---

A `source` block supports the following arguments:

* `command_id` - (Optional) 
* `script` - (Optional) 
* `script_uri` - (Optional) 
* `script_uri_managed_identity` - (Optional) A `script_uri_managed_identity` block as defined above.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Machine Run Command.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Machine Run Command.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Run Command.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Machine Run Command.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Machine Run Command.

## Import

An existing Virtual Machine Run Command can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_run_command.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/vm1/runCommands/rc1
```
