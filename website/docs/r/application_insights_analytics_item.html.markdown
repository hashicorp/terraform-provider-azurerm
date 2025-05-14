---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_analytics_item"
description: |-
  Manages an Application Insights Analytics Item component.
---

# azurerm_application_insights_analytics_item

Manages an Application Insights Analytics Item component.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tf-test"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "tf-test-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_application_insights_analytics_item" "example" {
  name                    = "testquery"
  application_insights_id = azurerm_application_insights.example.id
  content                 = "requests //simple example query"
  scope                   = "shared"
  type                    = "query"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Application Insights Analytics Item. Changing this forces a new resource to be created.

* `application_insights_id` - (Required) The ID of the Application Insights component on which the Analytics Item exists. Changing this forces a new resource to be created.

* `type` - (Required) The type of Analytics Item to create. Can be one of `query`, `function`, `folder`, `recent`. Changing this forces a new resource to be created.

* `scope` - (Required) The scope for the Analytics Item. Can be `shared` or `user`. Changing this forces a new resource to be created. Must be `shared` for functions.

* `content` - (Required) The content for the Analytics Item, for example the query text if `type` is `query`.

* `function_alias` - (Optional) The alias to use for the function. Required when `type` is `function`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Application Insights Analytics Item.

* `time_created` - A string containing the time the Analytics Item was created.

* `time_modified` - A string containing the time the Analytics Item was last modified.

* `version` - A string indicating the version of the query format

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Insights Analytics Item.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Insights Analytics Item.
* `update` - (Defaults to 30 minutes) Used when updating the Application Insights Analytics Item.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Insights Analytics Item.

## Import

Application Insights Analytics Items can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_insights_analytics_item.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Insights/components/mycomponent1/analyticsItems/11111111-1111-1111-1111-111111111111
```

-> **Note:** This is a Terraform Unique ID matching the format: `{appInsightsID}/analyticsItems/{itemId}` for items with `scope` set to `shared`, or  `{appInsightsID}/myAnalyticsItems/{itemId}` for items with `scope` set to `user`

To find the Analytics Item ID you can query the REST API using the [`az rest` CLI command](https://docs.microsoft.com/cli/azure/reference-index?view=azure-cli-latest#az-rest), e.g.

```shell
az rest --method GET --uri "https://management.azure.com/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.insights/components/appinsightstest/analyticsItems?api-version=2015-05-01"
```
