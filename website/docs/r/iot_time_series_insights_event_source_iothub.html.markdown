---
subcategory: "Time Series Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iot_time_series_insights_event_source_iothub"
description: |-
  Manages an Azure IoT Time Series Insights IoTHub Event Source.
---

# azurerm_time_series_insights_event_source_iothub

Manages an Azure IoT Time Series Insights IoTHub Event Source.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Eurpoe"
}

resource "azurerm_iothub" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "B1"
    capacity = "1"
  }
}

resource "azurerm_iothub_consumer_group" "example" {
  name                   = "example"
  iothub_name            = azurerm_iothub.example.name
  eventhub_endpoint_name = "events"
  resource_group_name    = azurerm_resource_group.example.name
}

resource "azurerm_storage_account" "storage" {
  name                     = "example"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iot_time_series_insights_gen2_environment" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "L1"
  id_properties       = ["id"]

  storage {
    name = azurerm_storage_account.storage.name
    key  = azurerm_storage_account.storage.primary_access_key
  }
}

resource "azurerm_iot_time_series_insights_event_source_iothub" "example" {
  name                     = "example"
  location                 = azurerm_resource_group.example.location
  environment_id           = azurerm_iot_time_series_insights_gen2_environment.example.id
  iothub_name              = azurerm_iothub.example.name
  shared_access_key        = azurerm_iothub.example.shared_access_policy.0.primary_key
  shared_access_key_name   = azurerm_iothub.example.shared_access_policy.0.key_name
  consumer_group_name      = azurerm_iothub_consumer_group.example.name
  event_source_resource_id = azurerm_iothub.example.id
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure IoT Time Series Insights IoTHub Event Source. Changing this forces a new resource to be created. Must be globally unique.

* `environment_id` - (Required) Specifies the id of the IoT Time Series Insights Environment that the Event Source should be associated with. Changing this forces a new resource to created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `iothub_name` - (Required) Specifies the name of the IotHub which will be associated with this resource.

* `consumer_group_name` - (Required) Specifies the name of the IotHub Consumer Group that holds the partitions from which events will be read.

* `event_source_resource_id` - (Required) Specifies the resource id where events will be coming from.

* `shared_access_key_name` - (Required) Specifies the name of the Shared Access key that grants the Event Source access to the IotHub.

* `shared_access_key` - (Required) Specifies the value of the Shared Access Policy key that grants the Time Series Insights service read access to the IotHub.

* `timestamp_property_name` - (Optional) Specifies the value that will be used as the event source's timestamp. This value defaults to the event creation time.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoT Time Series Insights IoTHub Event Source.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Time Series Insights IoTHub Event Source.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Time Series Insights IoTHub Event Source.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Time Series Insights IoTHub Event Source.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Time Series Insights IoTHub Event Source.

## Import

Azure IoT Time Series Insights IoTHub Event Source can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_time_series_insights_event_source_iothub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.TimeSeriesInsights/environments/environment1/eventSources/example
```
