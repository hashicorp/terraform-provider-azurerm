---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_saved_search"
description: |-
  Manages a Log Analytics (formally Operational Insights) Saved Search.
---

# azurerm_log_analytics_saved_search

Manages a Log Analytics (formally Operational Insights) Saved Search.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "acctest-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_saved_search" "example" {
  name                       = "exampleSavedSearch"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  category     = "exampleCategory"
  display_name = "exampleDisplayName"
  query        = "exampleQuery"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Log Analytics Saved Search. Changing this forces a new resource to be created.

* `log_analytics_workspace_id` - (Required) Specifies the ID of the Log Analytics Workspace that the Saved Search will be associated with. Changing this forces a new resource to be created.

* `display_name` - (Required) The name that Saved Search will be displayed as. Changing this forces a new resource to be created.

* `category` - (Required) The category that the Saved Search will be listed under. Changing this forces a new resource to be created.

* `query` - (Required) The query expression for the saved search. Changing this forces a new resource to be created.

* `function_alias` - (Optional) The function alias if the query serves as a function. Changing this forces a new resource to be created.

* `function_parameters` - (Optional) The function parameters if the query serves as a function. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Log Analytics Saved Search ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Saved Search.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Saved Search.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Saved Search.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Saved Search.

## Import

Log Analytics Saved Searches can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_saved_search.search1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/savedSearches/search1
```
