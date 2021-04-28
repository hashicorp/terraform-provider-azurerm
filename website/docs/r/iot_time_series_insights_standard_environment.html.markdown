---
subcategory: "Time Series Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iot_time_series_insights_standard_environment"
description: |-
  Manages an Azure IoT Time Series Insights Standard Environment.
---

# azurerm_time_series_insights_standard_environment

Manages an Azure IoT Time Series Insights Standard Environment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_iot_time_series_insights_standard_environment" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "S1_1"
  data_retention_time = "P30D"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure IoT Time Series Insights Standard Environment. Changing this forces a new resource to be created. Must be globally unique.

* `resource_group_name` - (Required) The name of the resource group in which to create the Azure IoT Time Series Insights Standard Environment.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this IoT Time Series Insights Standard Environment. It is string consisting of two parts separated by an underscore(\_).The fist part is the `name`, valid values include: `S1` and `S2`. The second part is the `capacity` (e.g. the number of deployed units of the `sku`), which must be a positive `integer` (e.g. `S1_1`). Changing this forces a new resource to be created.

* `data_retention_time` - (Required) Specifies the ISO8601 timespan specifying the minimum number of days the environment's events will be available for query. Changing this forces a new resource to be created.

* `storage_limit_exceeded_behavior` - (Optional) Specifies the behaviour the IoT Time Series Insights service should take when the environment's capacity has been exceeded. Valid values include `PauseIngress` and `PurgeOldData`. Defaults to `PurgeOldData`.

* `partition_key` - (Optional) The name of the event property which will be used to partition data. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoT Time Series Insights Standard Environment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Time Series Insights Standard Environment.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Time Series Insights Standard Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Time Series Insights Standard Environment.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Time Series Insights Standard Environment.

## Import

Azure IoT Time Series Insights Standard Environment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_time_series_environment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.TimeSeriesInsights/environments/example
```
