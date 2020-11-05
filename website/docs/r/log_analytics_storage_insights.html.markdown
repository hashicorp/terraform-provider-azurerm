---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_storage_insights"
description: |-
  Manages a Log Analytics Storage Insights resource.
---

# azurerm_log_analytics_storage_insights

Manages a Log Analytics Storage Insights resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "exampleworkspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_log_analytics_storage_insights" "example" {
  name                = "example-storageinsightconfig"
  resource_group_name = azurerm_resource_group.example.name
  workspace_id        = azurerm_log_analytics_workspace.example.id

  storage_account_id  = azurerm_storage_account.example.id
  storage_account_key = azurerm_storage_account.example.primary_access_key
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Log Analytics Storage Insights. Changing this forces a new Log Analytics Storage Insights to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Log Analytics Storage Insights should exist. Changing this forces a new Log Analytics Storage Insights to be created.

* `workspace_id` - (Required) The resource ID of the workspace to create the Log Analytics Storage Insights within. Changing this forces a new Log Analytics Storage Insights to be created.

* `storage_account_id` - (Required) The resource ID of the storage account to be used by this Log Analytics Storage Insights.

* `storage_account_key` - (Required) The storage access key to be used to connect to the storage account.

* `blob_container_names` - (Optional) The names of the blob containers that the workspace should read.

* `table_names` - (Optional) The names of the Azure tables that the workspace should read.

* `tags` - (Optional) A mapping of tags which should be assigned to the Log Analytics Storage Insights.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Storage Insights.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Storage Insights.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Storage Insights.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Storage Insights.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Storage Insights.

## Import

Log Analytics Storage Insight Configs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_storage_insights.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/storageInsightConfigs/storageInsight1
```