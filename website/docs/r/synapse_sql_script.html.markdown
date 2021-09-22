---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_sql_script"
description: |-
  Manages a Synapse SQL Script.
---

# azurerm_synapse_sql_script

Manages a Synapse SQL Script.

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

resource "azurerm_synapse_sql_script" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  description          = "test"
  language             = "sql"
  query                = "SELECT TOP 100 * FROM example_table_name;"
  sql_connection {
    name = "master"
    type = "SqlOnDemand"
  }

  depends_on = [
    azurerm_synapse_firewall_rule.example,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse SQL Script. Changing this forces a new Synapse SQL Script to be created.

* `synapse_workspace_id` - (Required) The ID of the Synapse Workspace. Changing this forces a new Synapse SQL Script to be created.

---

* `description` - (Optional) The description for the Synapse SQL Script.

* `language` - (Optional) The language of the Synapse SQL script.

* `query` - (Optional) SQL query to execute.

* `sql_connection` - (Optional) A `sql_connection` block as defined below.

* `type` - (Optional) The type of the Synapse SQL script. Possible values include: `SQLQuery`. Defaults to `SQLQuery`.

---

A `sql_connection` block supports the following:

* `name` - (Required) The name which should be used for this connection.

* `type` - (Required) The type of the connection. Possible values include: `SQLConnectionTypeSQLOnDemand`, `SQLConnectionTypeSQLPool`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Synapse SQL Script.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse SQL Script.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse SQL Script.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse SQL Script.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse SQL Script.

## Import

Synapse SQL Scripts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_sql_script.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/sqlscripts/sqlscript1
```
