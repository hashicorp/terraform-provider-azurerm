---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_spark_pool"
description: |-
  Manages a Synapse Spark Pool.
---

# azurerm_synapse_spark_pool

Manages a Synapse Spark Pool.

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
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
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

resource "azurerm_synapse_spark_pool" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"

  auto_scale {
    max_node_count = 50
    min_node_count = 3
  }

  auto_pause {
    delay_in_minutes = 15
  }

  tags = {
    ENV = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse Spark Pool. Changing this forces a new Synapse Spark Pool to be created.

* `synapse_workspace_id` - (Required) The ID of the Synapse Workspace where the Synapse Spark Pool should exist. Changing this forces a new Synapse Spark Pool to be created.

* `node_size_family` - (Required) The kind of nodes that the Spark Pool provides. Possible value is `MemoryOptimized`.

* `node_size` - (Required) The level of node in the Spark Pool. Possible value is `Small`, `Medium` and `Large`.

* `node_count` - (Optional) The number of nodes in the Spark Pool. Exactly one of `node_count` or `auto_scale` must be specified.

* `auto_scale` - (Optional)  An `auto_scale` block as defined below. Exactly one of `node_count` or `auto_scale` must be specified.

* `auto_pause` - (Optional)  An `auto_pause` block as defined below.

* `library_requirement` - (Optional)  A `library_requirement` block as defined below.

* `spark_log_folder` - (Optional) The default folder where Spark logs will be written. Defaults to `/logs`.

* `spark_events_folder` - (Optional) The Spark events folder. Defaults to `/events`.

* `spark_version` - (Optional) The Apache Spark version. Possible values are `2.4` and `3.0`. Defaults to `2.4`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Synapse Spark Pool.

---

An `auto_pause` block supports the following:

* `delay_in_minutes` - (Required) Number of minutes of idle time before the Spark Pool is automatically paused. Must be between `5` and `10080`.

---

An `auto_scale` block supports the following:

* `max_node_count` - (Required) The maximum number of nodes the Spark Pool can support. Must be between `3` and `200`.

* `min_node_count` - (Required) The minimum number of nodes the Spark Pool can support. Must be between `3` and `200`.

---

An `library_requirement` block supports the following:

* `content` - (Required) The content of library requirements.

* `filename` - (Required) The name of the library requirements file.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Synapse Spark Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Spark Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Spark Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Spark Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Spark Pool.

## Import

Synapse Spark Pool can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_spark_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1/bigDataPools/sparkPool1
```
