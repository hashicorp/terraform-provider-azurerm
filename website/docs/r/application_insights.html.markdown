---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights"
description: |-
  Manages an Application Insights component.
---

# azurerm_application_insights

Manages an Application Insights component.

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

output "instrumentation_key" {
  value = azurerm_application_insights.example.instrumentation_key
}

output "app_id" {
  value = azurerm_application_insights.example.app_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Application Insights component. Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the Application Insights component.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `application_type` - (Required) Specifies the type of Application Insights to create. Valid values are `ios` for _iOS_, `java` for _Java web_, `MobileCenter` for _App Center_, `Node.JS` for _Node.js_, `other` for _General_, `phone` for _Windows Phone_, `store` for _Windows Store_ and `web` for _ASP.NET_. Please note these values are case sensitive; unmatched values are treated as _ASP.NET_ by Azure. Changing this forces a new resource to be created.

* `daily_data_cap_in_gb` - (Optional) Specifies the Application Insights component daily data volume cap in GB.

* `daily_data_cap_notifications_disabled` - (Optional) Specifies if a notification email will be send when the daily data volume cap is met.

* `retention_in_days` - (Optional) Specifies the retention period in days. Possible values are `30`, `60`, `90`, `120`, `180`, `270`, `365`, `550` or `730`. Defaults to `90`.

* `sampling_percentage` - (Optional) Specifies the percentage of the data produced by the monitored application that is sampled for Application Insights telemetry.

* `disable_ip_masking` - (Optional) By default the real client ip is masked as `0.0.0.0` in the logs. Use this argument to disable masking and log the real client ip. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Application Insights component.

* `app_id` - The App ID associated with this Application Insights component.

* `instrumentation_key` - The Instrumentation Key for this Application Insights component.

* `connection_string` - The Connection String for this Application Insights component. (Sensitive)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Insights Component.
* `update` - (Defaults to 30 minutes) Used when updating the Application Insights Component.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Insights Component.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Insights Component.

## Import

Application Insights instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_insights.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.insights/components/instance1
```
