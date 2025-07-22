---
subcategory: "Device Registry"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_device_registry_asset"
description: |-
  Manages a Device Registry Asset.
---

# azurerm_device_registry_asset

Manages a Device Registry Asset.

## Example Usage

Creates an asset with 2 events and 1 dataset containing 2 datapoints:
```hcl
resource "azurerm_device_registry_asset" "example" {
  name                             = "example"
  location                         = "West US 2"
  resource_group_id                = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group"
  asset_endpoint_profile_reference = "myAssetEndpointProfileRef"
  extended_location_id             = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group/providers/Microsoft.ExtendedLocation/customLocations/my-custom-location"
  display_name                     = "my asset"
  enabled                          = true
  external_asset_id                = "8ZBA6LRHU0A458969"
  attributes = {
    "foo" = "bar"
    "x"   = "y"
  }
  default_datasets_configuration = jsonencode(
    {
      defaultPublishingInterval = 200
      defaultQueueSize          = 10
      defaultSamplingInterval   = 500
    }
  )
  default_events_configuration = jsonencode(
    {
      defaultPublishingInterval = 200
      defaultQueueSize          = 10
      defaultSamplingInterval   = 500
    }
  )
  default_topic {
    path   = "/path/defaultTopic"
    retain = "Keep"
  }
  description = "this is my asset"
  discovered_asset_references = [
    "foo",
    "bar",
    "baz",
  ]
  documentation_uri = "https://example.com/about"
  hardware_revision = "1.0"
  manufacturer      = "Contoso"
  manufacturer_uri  = "https://www.contoso.com/manufacturerUri"
  model             = "ContosoModel"
  product_code      = "SA34VDG"
  serial_number     = "64-103816-519918-8"
  software_revision = "2.0"
  tags = {
    "site" = "building-1"
  }

  dataset {
    dataset_configuration = jsonencode(
      {
        publishingInterval = 7
        queueSize          = 8
        samplingInterval   = 1000
      }
    )
    name = "dataset1"
    topic {
      path   = "/path/dataset1"
      retain = "Keep"
    }

    data_point {
      data_point_configuration = jsonencode(
        {
          publishingInterval = 7
          queueSize          = 8
          samplingInterval   = 1000
        }
      )
      data_source        = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt1"
      name               = "datapoint1"
      observability_mode = "Counter"
    }
    data_point {
      data_point_configuration = jsonencode(
        {
          publishingInterval = 7
          queueSize          = 8
          samplingInterval   = 1000
        }
      )
      data_source        = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt2"
      name               = "datapoint2"
      observability_mode = "None"
    }
  }

  event {
    event_configuration = jsonencode(
      {
        publishingInterval = 7
        queueSize          = 8
        samplingInterval   = 1000
      }
    )
    event_notifier     = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt3"
    name               = "event1"
    observability_mode = "Log"
    topic {
      path   = "/path/event1"
      retain = "Never"
    }
  }
  event {
    event_configuration = jsonencode(
      {
        publishingInterval = 7
        queueSize          = 8
        samplingInterval   = 1000
      }
    )
    event_notifier     = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt4"
    name               = "event2"
    observability_mode = "None"
    topic {
      path   = "/path/event2"
      retain = "Keep"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `asset_endpoint_profile_reference` - (Required) A reference to the asset endpoint profile (connection information) used by brokers to connect to an endpoint that provides data points for this asset. Must provide asset endpoint profile name.

* `extended_location_id` - (Required) The ID of the extended location. Must provide a custom location ID.

* `location` - (Required) The Azure Region where the Device Registry Asset should exist. Changing this forces a new Device Registry Asset to be created.

* `name` - (Required) The name which should be used for this Device Registry Asset.

* `resource_group_id` - (Required) The ID of the Resource Group where the Asset should exist.


* `attributes` - (Optional) A set of key-value pairs that contain custom attributes set by the customer.

* `dataset` - (Optional) One or more `dataset` blocks as defined below. Each dataset describes the data points that make up the set.

* `default_datasets_configuration` - (Optional) Stringified JSON that contains connector-specific default configuration for all datasets. Each dataset can have its own configuration that overrides the default settings here..

* `default_events_configuration` - (Optional) Stringified JSON that contains connector-specific default configuration for all events. Each event can have its own configuration that overrides the default settings here..

* `default_topic` - (Optional) A `default_topic` block as defined below. Describes the default topic information for the asset.

* `description` - (Optional) Human-readable description of the asset.

* `discovered_asset_references` - (Optional) Specifies a list of discovered assets references. Populated only if the asset has been created from discovery flow. Discovered asset names must be provided.

* `display_name` - (Optional) Human-readable display name.

* `documentation_uri` - (Optional) Reference to the documentation.

* `enabled` - (Optional) Enabled/Disabled status of the asset.

* `event` - (Optional) One or more `event` blocks as defined below. Each event can have per-event configuration.

* `external_asset_id` - (Optional) Asset ID provided by the customer.

* `hardware_revision` - (Optional) Revision number of the hardware.

* `manufacturer` - (Optional) Asset manufacturer name.

* `manufacturer_uri` - (Optional) Asset manufacturer URI.

* `model` - (Optional) Asset model name.

* `product_code` - (Optional) Asset product code.

* `serial_number` - (Optional) Asset serial number.

* `software_revision` - (Optional) Revision number of the software.

* `tags` - (Optional) A mapping of tags which should be assigned to the Device Registry Asset.

---

A `data_point` block supports the following:

* `data_source` - (Required) The address of the source of the data in the asset (e.g. URL) so that a client can access the data source on the asset.

* `name` - (Required) The name of the data point.

* `data_point_configuration` - (Optional) Stringified JSON that contains connector-specific configuration for the data point. For OPC UA, this could include configuration like, publishingInterval, samplingInterval, and queueSize.

* `observability_mode` - (Optional) An indication of how the data point should be mapped to OpenTelemetry. Possible values are `Counter`, `Gauge`, `Histogram`, `Log` and `None`. Defaults to `None`.

---

A `dataset` block supports the following:

* `name` - (Required) The name of the dataset.

* `data_point` - (Optional) One or more `data_point` blocks as defined above.

* `dataset_configuration` - (Optional) Stringified JSON that contains connector-specific JSON string that describes configuration for the specific dataset.

* `topic` - (Optional) A `topic` block as defined below. Describes the topic information for the specific dataset.

---

A `default_topic` block supports the following:

* `path` - (Required) The topic path for messages published to an MQTT broker.

* `retain` - (Optional) When set to `Keep`, messages published to an MQTT broker will have the retain flag set. Defaults to `Never`.

---

A `event` block supports the following:

* `event_notifier` - (Required) The address of the notifier of the event in the asset (e.g. URL) so that a client can access the event on the asset.

* `name` - (Required) The name of the event.

* `event_configuration` - (Optional) Stringified JSON that contains connector-specific configuration for the event. For OPC UA, this could include configuration like, publishingInterval, samplingInterval, and queueSize.

* `observability_mode` - (Optional) An indication of how the event should be mapped to OpenTelemetry. Possible values are `Log` and `None`. Defaults to `None`.

* `topic` - (Optional) A `topic` block as defined below. Describes the topic information for the specific event.

---

A `topic` block supports the following:

* `path` - (Required) The topic path for messages published to an MQTT broker.

* `retain` - (Optional) When set to `Keep`, messages published to an MQTT broker will have the retain flag set. Defaults to `Never`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Device Registry Asset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Device Registry Asset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Device Registry Asset.
* `update` - (Defaults to 30 minutes) Used when updating the Device Registry Asset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Device Registry Asset.

## Import

Device Registry Assets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_device_registry_asset.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName/providers/Microsoft.DeviceRegistry/assets/assetName
```
