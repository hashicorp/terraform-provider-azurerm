---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_sql_pool_workload_group"
description: |-
  Manages a Synapse SQL Pool Workload Group.
---

# azurerm_synapse_sql_pool_workload_group

Manages a Synapse SQL Pool Workload Group.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "west europe"
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

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_sql_pool" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  sku_name             = "DW100c"
  create_mode          = "Default"
}

resource "azurerm_synapse_sql_pool_workload_group" "example" {
  name                               = "example"
  sql_pool_id                        = azurerm_synapse_sql_pool.example.id
  importance                         = "normal"
  max_resource_percent               = 100
  min_resource_percent               = 0
  max_resource_percent_per_request   = 3
  min_resource_percent_per_request   = 3
  query_execution_timeout_in_seconds = 0
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse SQL Pool Workload Group. Changing this forces a new Synapse SQL Pool Workload Group to be created.

* `sql_pool_id` - (Required) The ID of the Synapse SQL Pool. Changing this forces a new Synapse SQL Pool Workload Group to be created.

* `max_resource_percent` - (Required) The workload group cap percentage resource.

* `min_resource_percent` - (Required) The workload group minimum percentage resource.

---

* `importance` - (Optional) The workload group importance level. Defaults to `normal`.

* `max_resource_percent_per_request` - (Optional) The workload group request maximum grant percentage. Defaults to `3`.

* `min_resource_percent_per_request` - (Optional) The workload group request minimum grant percentage.

* `query_execution_timeout_in_seconds` - (Optional) The workload group query execution timeout.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Synapse SQL Pool Workload Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse SQL Pool Workload Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse SQL Pool Workload Group.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse SQL Pool Workload Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse SQL Pool Workload Group.

## Import

Synapse SQL Pool Workload Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_sql_pool_workload_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/sqlPools/sqlPool1/workloadGroups/workloadGroup1
```
