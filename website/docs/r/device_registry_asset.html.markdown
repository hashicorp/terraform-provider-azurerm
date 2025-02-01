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

```hcl
resource "azurerm_device_registry_asset" "example" {
  name = "example"
  resource_group_name = "example"
  location = "West Europe"
  extended_location_name = "example"
  extended_location_type = "TODO"
  asset_endpoint_profile_ref = "TODO"
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

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Asset.

* `provisioning_state` - Provisioning state of the resource.

* `status` - Read only object to reflect changes that have occurred on the Edge. Similar to Kubernetes status property for custom resources.

* `type` - Azure resource type. Defaults to `Microsoft.DeviceRegistry/Assets`.

* `uuid` - Globally unique, immutable, non-reusable id."

* `version` - An integer that is incremented each time the resource is modified.

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