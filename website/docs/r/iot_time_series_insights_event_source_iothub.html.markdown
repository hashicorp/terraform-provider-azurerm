---
subcategory: "Time Series Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iot_time_series_insights_event_source_iothub"
description: |-
  Manages an Azure IoT Time Series Insights Event Source for IoTHub.
---

# azurerm_time_series_insights_event_source_iothub

Manages an Azure IoT Time Series Insights Event Source for IoTHub.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}
resource "azurerm_iot_time_series_insights_standard_environment" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "S1_1"
  data_retention_time = "P30D"
}
resource "azurerm_iot_time_series_insights_event_source_iothub" "test" {
  name                                = "accTEst_tsiap%d"
  time_series_insights_environment_id = azurerm_iot_time_series_insights_standard_environment.test.id
  location                            = azurerm_resource_group.test.location

  event_source_resource_id = azurerm_resource_group.test.id
  iothub_name              = azurerm_iothub.test.name
  key_name                 = azurerm_iothub.test.shared_access_policy.0.key_name
  shared_access_key        = azurerm_iothub.test.shared_access_policy.0.primary_key
  consumer_group_name      = "tsiquickstart"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure IoT Time Series Insights Event Source. Changing this forces a new resource to be created. Must be globally unique.

* `time_series_insights_environment_id` - (Required) The resource ID of the Azure IoT Time Series Insights Environment in which to create the Azure IoT Time Series Insights Reference Data Set. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `event_source_resource_id` - (Required) Specifies the resource id of the Event Source.

* `iothub_name` - (Required) The name of the iot hub that will be used with the Event Source.

* `consumer_group_name` - (Required) The name of the iot hub's consumer group that holds the partitions from which events will be read.

* `key_name` - (Required) The name of the Shared Access Policy key that grants the Time Series Insights service access to the iot hub. This shared access policy key must grant 'service connect' permissions to the iot hub.

* `shared_access_key` - (Required) The value of the Shared Access Policy key that grants the Time Series Insights service read access to the iot hub.

* `timestamp_property_name` - (Optional) The event property that will be used as the event source's timestamp. If a value isn't specified for timestampPropertyName, the event creation time will be used.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoT Time Series Insights Event Source.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Time Series Insights Event Source.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Time Series Insights Event Source.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Time Series Insights Event Source.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Time Series Insights Event Source.

## Import

Azure IoT Time Series Insights Event Source can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_time_series_event_source_iothub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.TimeSeriesInsights/environments/environment1/eventSources/example
```
