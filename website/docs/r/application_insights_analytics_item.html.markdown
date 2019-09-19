---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_analytics_item"
sidebar_current: "docs-azurerm-resource-application-insights-x"
description: |-
  Manages an Application Insights Analytics Item component.
---

# azurerm_application_insights_analytics_item

Manages an Application Insights Analytics Item component.

## Example Usage

```terraform
resource "azurerm_resource_group" "test" {
  name     = "tf-test"
  location = "West Europe"
}

resource "azurerm_application_insights" "test" {
  name                = "tf-test-appinsights"
  location            = "West Europe"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_application_insights_analytics_item" "test" {
  name                    = "testquery"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  application_insights_id = "${azurerm_application_insights.test.id}"
  content                 = "requests //simple example query"
  scope                   = "shared"
  type                    = "query"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Application Insights Analytics Item. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Application Insights component.

* `application_insights_id` - (Required) The ID of the Application Insights component on which the Analytics Item exists. Changing this forces a new resource to be created.

* `type` - (Required) The type of Analytics Item to create. Can be one of `query`, `function`, `folder`, `recent`. Changing this forces a new resource to be created.

* `scope` - (Required) The scope for the Analytics Item. Can be `shared` or `user`. Changing this forces a new resource to be created. Must be `shared` for for functions.

* `content` - (Required) The content for the Analytics Item, for example the query text if `type` is `query`.

* `function_alias` - (Optional) The alias to use for the function. Required when `type` is `function`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Application Insights Analytics Item.

* `time_created` - A string containing the time the Analytics Item was created.

* `time_modified` - A string containing the time the Analytics Item was last modified.
