---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_query_pack_query"
description: |-
  Manages a Log Analytics Query Pack Query.
---

# azurerm_log_analytics_query_pack_query

Manages a Log Analytics Query Pack Query.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_query_pack" "example" {
  name                = "example-laqp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_log_analytics_query_pack_query" "example" {
  query_pack_id = azurerm_log_analytics_query_pack.example.id
  body          = "let newExceptionsTimeRange = 1d;\nlet timeRangeToCheckBefore = 7d;\nexceptions\n| where timestamp < ago(timeRangeToCheckBefore)\n| summarize count() by problemId\n| join kind= rightanti (\nexceptions\n| where timestamp >= ago(newExceptionsTimeRange)\n| extend stack = tostring(details[0].rawStack)\n| summarize count(), dcount(user_AuthenticatedId), min(timestamp), max(timestamp), any(stack) by problemId  \n) on problemId \n| order by  count_ desc\n"
  display_name  = "Exceptions - New in the last 24 hours"
}
```

## Arguments Reference

The following arguments are supported:

* `query_pack_id` - (Required) The ID of the Log Analytics Query Pack. Changing this forces a new resource to be created.

* `body` - (Required) The body of the Log Analytics Query Pack Query.

* `display_name` - (Required) The unique display name for the query within the Log Analytics Query Pack.

* `name` - (Optional) A unique UUID/GUID which identifies this lighthouse assignment- one will be generated if not specified. Changing this forces a new resource to be created.

* `description` - (Optional) The description of the Log Analytics Query Pack Query.

* `properties_json` - (Optional) The additional properties that can be set for the Log Analytics Query Pack Query.

* `related` - (Optional) A `related` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Log Analytics Query Pack.

---

A `related` block supports the following:

* `categories` - (Optional) A list of the related categories for the function.

* `resource_types` - (Optional) A list of the related resource types for the function.

* `solutions` - (Optional) A list of the related Log Analytics solutions for the function.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Log Analytics Query Pack Query.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Query Pack Query.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Query Pack Query.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Query Pack Query.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Query Pack Query.

## Import

Log Analytics Query Pack Queries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_query_pack.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.OperationalInsights/queryPacks/queryPack1/queries/15b49e87-8555-4d92-8a7b-2014b469a9df
```
