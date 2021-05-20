---
subcategory: "Databox Edge"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databox_edge_device"
description: |-
  Manages a Databox Edge Device.
---

# azurerm_databox_edge_device

Manages a Databox Edge Device.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-databoxedge"
  location = "West Europe"
}

resource "azurerm_databox_edge_device" "example" {
  name                = "example-device"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "EdgeP_Base-Standard"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Databox Edge Device. Changing this forces a new Databox Edge Device to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Databox Edge Device should exist. Changing this forces a new Databox Edge Device to be created.

* `location` - (Required) The Azure Region where the Databox Edge Device should exist. Changing this forces a new Databox Edge Device to be created.

* `sku_name` - (Required)  The `sku_name` is comprised of two segments separated by a hyphen (e.g. `TEA_1Node_UPS_Heater-Standard`). The first segment of the `sku_name` defines the `name` of the sku, possible values are `Gateway`, `EdgeMR_Mini`, `EdgeP_Base`, `EdgeP_High`, `EdgePR_Base`, `EdgePR_Base_UPS`, `GPU`, `RCA_Large`, `RCA_Small`, `RDC`, `TCA_Large`, `TCA_Small`, `TDC`, `TEA_1Node`, `TEA_1Node_UPS`, `TEA_1Node_Heater`, `TEA_1Node_UPS_Heater`, `TEA_4Node_Heater`, `TEA_4Node_UPS_Heater` or `TMA`. The second segment defines the `tier` of the `sku_name`, possible values are `Standard`. For more information see the [product documentation]("https://docs.microsoft.com/en-us/dotnet/api/microsoft.azure.management.databoxedge.models.sku?view=azure-dotnet"). Changing this forces a new Databox Edge Device to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Databox Edge Device.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Databox Edge Device.

* `device_properties` - A `device_properties` block as defined below.

---

The `device_properties` block exports the following:

* `configured_role_types` - Type of compute roles configured.

* `culture` - The Data Box Edge/Gateway device culture.

* `hcs_version` - The device software version number of the device (eg: 1.2.18105.6).

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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Databox Edge Device.
* `read` - (Defaults to 5 minutes) Used when retrieving the Databox Edge Device.
* `update` - (Defaults to 30 minutes) Used when updating the Databox Edge Device.
* `delete` - (Defaults to 30 minutes) Used when deleting the Databox Edge Device.

## Import

Databox Edge Devices can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databox_edge_device.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/device1
```
