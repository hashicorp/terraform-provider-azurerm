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

Creates an asset with 2 events and 1 dataset containing 2 data points:
```hcl
resource "azurerm_device_registry_asset" "example" {
  asset_endpoint_profile_ref = "myAssetEndpointProfileRef"
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
  default_topic_path   = "/path/defaultTopic"
  default_topic_retain = "Keep"
  description          = "This is my asset"
  discovered_asset_refs = [
    "foo",
    "bar",
    "baz",
  ]
  display_name           = "My Asset"
  documentation_uri      = "https://example.com/about"
  enabled                = false
  extended_location_name = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group/providers/Microsoft.ExtendedLocation/customLocations/my-custom-location"
  extended_location_type = "CustomLocation"
  external_asset_id      = "foobar"
  hardware_revision      = "1.0"
  location               = "westus2"
  manufacturer           = "Contoso"
  manufacturer_uri       = "https://example.com"
  model                  = "Model123"
  name                   = "myassetbasic"
  product_code           = "3E1YZ7"
  resource_group_name    = "my-resource-group"
  serial_number          = "1234"
  software_revision      = "2.0"
  tags = {
    "sensor" = "temperature,humidity"
  }

  datasets {
    name = "dataset1"
    dataset_configuration = jsonencode(
      {
        publishingInterval = 7
        queueSize          = 8
        samplingInterval   = 1000
      }
    )
    topic_path   = "/path/dataset1"
    topic_retain = "Keep"

    data_points {
      name               = "dataPoint1"
      data_source        = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt1"
      observability_mode = "Log"
      data_point_configuration = jsonencode(
        {
          publishingInterval = 7
          queueSize          = 8
          samplingInterval   = 1000
        }
      )
    }

    data_points {
      name               = "dataPoint2"
      data_source        = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt2"
      observability_mode = "Counter"
      data_point_configuration = jsonencode(
        {
          publishingInterval = 7
          queueSize          = 8
          samplingInterval   = 1000
        }
      )
    }
  }

  events {
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
    topic_path         = "/path/event1"
    topic_retain       = "Never"
  }
  events {
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
    topic_path         = "/path/event2"
    topic_retain       = "Keep"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `asset_endpoint_profile_ref` - (Required) A reference to the asset endpoint profile (connection information) used by brokers to connect to an endpoint that provides data points for this asset. Must provide asset endpoint profile name.

* `extended_location_name` - (Required) The extended location name.

* `extended_location_type` - (Required) The extended location type.

* `location` - (Required) The Azure Region where the Asset should exist. Changing this forces a new Asset to be created.

* `name` - (Required) The name which should be used for this Asset.

* `resource_group_name` - (Required) The name of the Resource Group where the Asset should exist. Changing this forces a new Asset to be created.

---

* `attributes` - (Optional) A set of key-value pairs that contain custom attributes set by the customer.

* `datasets` - (Optional) Array of datasets that are part of the asset. Each dataset describes the data points that make up the set.

* `default_datasets_configuration` - (Optional) Stringified JSON that contains connector-specific default configuration for all datasets. Each dataset can have its own configuration that overrides the default settings here.

* `default_events_configuration` - (Optional) Stringified JSON that contains connector-specific default configuration for all events. Each event can have its own configuration that overrides the default settings here.

* `default_topic` - (Optional) Object that describes the default topic information for the asset.

* `description` - (Optional) Human-readable description of the asset.

* `discovered_asset_refs` - (Optional) Reference to a list of discovered assets. Populated only if the asset has been created from discovery flow. Discovered asset names must be provided.

* `display_name` - (Optional) Human-readable display name.

* `documentation_uri` - (Optional) Reference to the documentation.

* `enabled` - (Optional) Enabled/Disabled status of the asset.

* `events` - (Optional) Array of events that are part of the asset. Each event can have per-event configuration.

* `external_asset_id` - (Optional) Asset id provided by the customer.

* `hardware_revision` - (Optional) Revision number of the hardware.

* `manufacturer` - (Optional) Asset manufacturer name.

* `manufacturer_uri` - (Optional) Asset manufacturer URI.

* `model` - (Optional) Asset model name.

* `product_code` - (Optional) Asset product code.

* `serial_number` - (Optional) Asset serial number.

* `software_revision` - (Optional) Revision number of the software.

* `tags` - (Optional) A mapping of tags which should be assigned to the Asset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Asset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Asset.
* `update` - (Defaults to 30 minutes) Used when updating the Asset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Asset.

## Import

Assets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_device_registry_asset.example C:/Program Files/Git/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/adr-terraform-rg/providers/Microsoft.DeviceRegistry/assets/test-asset
```
