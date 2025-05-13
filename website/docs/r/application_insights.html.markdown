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

## Example Usage - Workspace Mode

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tf-test"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "workspace-test"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_application_insights" "example" {
  name                = "tf-test-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  workspace_id        = azurerm_log_analytics_workspace.example.id
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

* `name` - (Required) Specifies the name of the Application Insights component. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Application Insights component. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `application_type` - (Required) Specifies the type of Application Insights to create. Valid values are `ios` for _iOS_, `java` for _Java web_, `MobileCenter` for _App Center_, `Node.JS` for _Node.js_, `other` for _General_, `phone` for _Windows Phone_, `store` for _Windows Store_ and `web` for _ASP.NET_. Please note these values are case sensitive; unmatched values are treated as _ASP.NET_ by Azure. Changing this forces a new resource to be created.

* `daily_data_cap_in_gb` - (Optional) Specifies the Application Insights component daily data volume cap in GB. Defaults to `100`.

* `daily_data_cap_notifications_disabled` - (Optional) Specifies if a notification email will be sent when the daily data volume cap is met. Defaults to `false`.

* `retention_in_days` - (Optional) Specifies the retention period in days. Possible values are `30`, `60`, `90`, `120`, `180`, `270`, `365`, `550` or `730`. Defaults to `90`.

* `sampling_percentage` - (Optional) Specifies the percentage of the data produced by the monitored application that is sampled for Application Insights telemetry. Defaults to `100`.

* `disable_ip_masking` - (Optional) By default the real client IP is masked as `0.0.0.0` in the logs. Use this argument to disable masking and log the real client IP. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `workspace_id` - (Optional) Specifies the id of a log analytics workspace resource.

~> **Note:** `workspace_id` cannot be removed after set. More details can be found at [Migrate to workspace-based Application Insights resources](https://docs.microsoft.com/azure/azure-monitor/app/convert-classic-resource#migration-process). If `workspace_id` is not specified but you encounter a diff, this might indicate a Microsoft initiated automatic migration from classic resources to workspace-based resources. If this is the case, please update `workspace_id` in the config file to the new value.

* `local_authentication_disabled` - (Optional) Disable Non-Azure AD based Auth. Defaults to `false`.

* `internet_ingestion_enabled` - (Optional) Should the Application Insights component support ingestion over the Public Internet? Defaults to `true`.

* `internet_query_enabled` - (Optional) Should the Application Insights component support querying over the Public Internet? Defaults to `true`.

* `force_customer_storage_for_profiler` - (Optional) Should the Application Insights component force users to create their own storage account for profiling? Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Application Insights component.

* `app_id` - The App ID associated with this Application Insights component.

* `instrumentation_key` - The Instrumentation Key for this Application Insights component. (Sensitive)

* `connection_string` - The Connection String for this Application Insights component. (Sensitive)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Application Insights Component.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Insights Component.
* `update` - (Defaults to 30 minutes) Used when updating the Application Insights Component.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Insights Component.

## Import

Application Insights instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_insights.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Insights/components/instance1
```
