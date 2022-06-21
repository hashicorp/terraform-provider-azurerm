---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_sql_pool_workload_classifier"
description: |-
  Manages a Synapse SQL Pool Workload Classifier.
---

# azurerm_synapse_sql_pool_workload_classifier

Manages a Synapse SQL Pool Workload Classifier.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
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

resource "azurerm_synapse_sql_pool_workload_classifier" "example" {
  name              = "example"
  workload_group_id = azurerm_synapse_sql_pool_workload_group.example.id

  context     = "example_context"
  end_time    = "14:00"
  importance  = "high"
  label       = "example_label"
  member_name = "dbo"
  start_time  = "12:00"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse SQL Pool Workload Classifier. Changing this forces a new Synapse SQL Pool Workload Classifier to be created.

* `workload_group_id` - (Required) The ID of the Synapse SQL Pool Workload Group. Changing this forces a new Synapse SQL Pool Workload Classifier to be created.

* `member_name` - (Required) The workload classifier member name used to classified against.

---

* `context` - (Optional) Specifies the session context value that a request can be classified against.

* `end_time` - (Optional) The workload classifier end time for classification. It's of the `HH:MM` format in UTC time zone.

* `importance` - (Optional) The workload classifier importance. The allowed values are `low`, `below_normal`, `normal`, `above_normal` and `high`.

* `label` - (Optional) Specifies the label value that a request can be classified against.

* `start_time` - (Optional) The workload classifier start time for classification. It's of the `HH:MM` format in UTC time zone.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Synapse SQL Pool Workload Classifier.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse SQL Pool Workload Classifier.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse SQL Pool Workload Classifier.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse SQL Pool Workload Classifier.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse SQL Pool Workload Classifier.

## Import

Synapse SQL Pool Workload Classifiers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_sql_pool_workload_classifier.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/sqlPools/sqlPool1/workloadGroups/workloadGroup1/workloadClassifiers/workloadClassifier1
```
