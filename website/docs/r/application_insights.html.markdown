---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights"
sidebar_current: "docs-azurerm-resource-application-insights"
description: |-
  Create an Application Insights component.
---

# azurerm_application_insights

Create an Application Insights component.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "tf-test"
  location = "West Europe"
}

resource "azurerm_application_insights" "test" {
  name                = "tf-test-appinsights"
  location            = "West Europe"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "Web"
}

output "instrumentation_key" {
  value = "${azurerm_application_insights.test.instrumentation_key}"
}

output "app_id" {
  value = "${azurerm_application_insights.test.app_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Application Insights component. Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the Application Insights component.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `application_type` - (Required) Specifies the type of Application Insights to create. Valid values are `Web` and `Other`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Application Insights component.

* `app_id` - The App ID associated with this Application Insights component.

* `instrumentation_key` - The Instrumentation Key for this Application Insights component.


## Import

Application Insights instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_insights.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.insights/components/instance1
```
