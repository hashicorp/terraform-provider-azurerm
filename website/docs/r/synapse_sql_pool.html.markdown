---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_sql_pool"
description: |-
  Manages a Synapse Sql Pool.
---

# azurerm_synapse_sql_pool

Manages a Synapse Sql Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "BlobStorage"
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
}

resource "azurerm_synapse_sql_pool" "example" {
  name                 = "examplesqlpool"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  sku_name             = "DW100c"
  create_mode          = "Default"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse Sql Pool. Changing this forces a new synapse SqlPool to be created.

* `synapse_workspace_id` - (Required) The ID of Synapse Workspace within which this Sql Pool should be created. Changing this forces a new Synapse Sql Pool to be created.

* `sku_name` - (Required) Specifies the SKU Name for this Synapse Sql Pool. Possible values are `DW100c`, `DW200c`, `DW300c`, `DW400c`, `DW500c`, `DW1000c`, `DW1500c`, `DW2000c`, `DW2500c`, `DW3000c`, `DW5000c`, `DW6000c`, `DW7500c`, `DW10000c`, `DW15000c` or `DW30000c`.

* `create_mode` - (Optional) Specifies how to create the Sql Pool. Valid values are: `Default`, `Recovery` or `PointInTimeRestore`. Must be `Default` to create a new database. Defaults to `Default`.

* `collation` - (Optional) The name of the collation to use with this pool, only applicable when `create_mode` is set to `Default`. Azure default is `SQL_LATIN1_GENERAL_CP1_CI_AS`. Changing this forces a new resource to be created.

* `recovery_database_id` - (Optional) The ID of the Synapse Sql Pool or Sql Database which is to back up, only applicable when `create_mode` is set to `Recovery`. Changing this forces a new Synapse Sql Pool to be created.

* `restore` - (Optional)  A `restore` block as defined below. only applicable when `create_mode` is set to `PointInTimeRestore`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Synapse Sql Pool.

---

An `restore` block supports the following:

* `source_database_id` - (Optional) The ID of the Synapse Sql Pool or Sql Database which is to restore. Changing this forces a new Synapse Sql Pool to be created.

* `point_in_time` - (Optional) Specifies the Snapshot time to restore. Changing this forces a new Synapse Sql Pool to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Synapse Sql Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Sql Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Sql Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Sql Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Sql Pool.

## Import

Synapse Sql Pool can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_sql_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1/sqlPools/sqlPool1
```
