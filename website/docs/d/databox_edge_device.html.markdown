---
subcategory: "Databox Edge"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databox_edge_device"
description: |-
  Get information about a Databox Edge Device.
---

# azurerm_databox_edge_device

Get information about a Databox Edge Device.

## Example Usage

```hcl
data "azurerm_databox_edge_device" "example" {
  name                = "example-device"
  resource_group_name = "example-rg"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Databox Edge Device. Changing this forces a new Databox Edge Device to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Databox Edge Device should exist. Changing this forces a new Databox Edge Device to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Databox Edge Device.

* `location` - The Azure Region where the Databox Edge Device should exist.

* `sku_name` - The `sku_name` is comprised of two segments separated by a hyphen (e.g. `TEA_1Node_UPS_Heater-Standard`). The first segment of the `sku_name` defines the `name` of the SKU. The second segment defines the `tier` of the `sku_name`. For more information see the [product documentation]("https://docs.microsoft.com/dotnet/api/microsoft.azure.management.databoxedge.models.sku?view=azure-dotnet"). 

* `device_properties` - A `device_properties` block as defined below.

* `tags` - A mapping of tags which should be assigned to the Databox Edge Device.

---

The `device_properties` block exports the following:

* `configured_role_types` - Type of compute roles configured.

* `culture` - The Data Box Edge/Gateway device culture.

* `hcs_version` - The device software version number of the device (e.g. 1.2.18105.6).

* `capacity` - The Data Box Edge/Gateway device local capacity in MB.

* `model` - The Data Box Edge/Gateway device model.

* `software_version` - The Data Box Edge/Gateway device software version.

* `status` - The status of the Data Box Edge/Gateway device.

* `type` - The type of the Data Box Edge/Gateway device.

* `node_count` - The number of nodes in the cluster.

* `serial_number` - The Serial Number of Data Box Edge/Gateway device.

* `time_zone` - The Data Box Edge/Gateway device timezone.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Databox Edge Device.
