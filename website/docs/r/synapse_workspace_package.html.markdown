---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_workspace_package"
description: |-
  Manages a Synapse Workspace Package.
---

# azurerm_synapse_workspace_package

Manages a Synapse Workspace Package.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "example"
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  managed_virtual_network_enabled      = true
}

resource "azurerm_synapse_firewall_rule" "example" {
  name                 = "allowAll"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}

resource "azurerm_synapse_workspace_package" "example" {
  name                 = "example.jar"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  source               = "/home/example.jar"
  source_md5           = filemd5("/home/example.jar")
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse Workspace Package. Changing this forces a new Synapse Workspace Package to be created.

* `synapse_workspace_id` - (Required) The ID of the Synapse Workspace. Changing this forces a new Synapse Workspace Package to be created.
  
* `source` - (Required) The filepath of the package. Changing this forces a new Synapse Workspace Package to be created.

* `source_md5` - (Required) The file md5 of the package. Changing this forces a new Synapse Workspace Package to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Synapse Workspace Package.

* `container_name` - Container name of the Synapse Workspace Package.

* `path` - Location of Synapse Workspace Package in storage account.

* `type` - Type of the Synapse Workspace Package.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Workspace Package.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Workspace Package.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Workspace Package.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Workspace Package.

## Import

Synapse Workspace Packages can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_workspace_package.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/libraries/library1
```
