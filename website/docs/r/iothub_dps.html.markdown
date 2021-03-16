---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_dps"
description: |-
  Manages an IoT Device Provisioning Service.
---

# azurerm_iothub_dps

Manages an IotHub Device Provisioning Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_iothub_dps" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Iot Device Provisioning Service resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the Iot Device Provisioning Service resource has to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource has to be createc. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `linked_hub` - (Optional) A `linked_hub` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block supports the following:

* `name` - (Required) The name of the sku. Currently can only be set to `S1`.

* `capacity` - (Required) The number of provisioned IoT Device Provisioning Service units.

---

A `linked_hub` block supports the following:

* `connection_string` - (Required) The connection string to connect to the IoT Hub. Changing this forces a new resource.

* `location` - (Required) The location of the IoT hub. Changing this forces a new resource.

* `apply_allocation_policy` - (Optional) Determines whether to apply allocation policies to the IoT Hub. Defaults to false.

* `allocation_weight` - (Optional) The weight applied to the IoT Hub. Defaults to 0.

* `hostname` - (Computed) The IoT Hub hostname.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoT Device Provisioning Service.

* `allocation_policy` - The allocation policy of the IoT Device Provisioning Service.

* `device_provisioning_host_name` - The device endpoint of the IoT Device Provisioning Service.

* `id_scope` - The unique identifier of the IoT Device Provisioning Service.

* `service_operations_host_name` - The service endpoint of the IoT Device Provisioning Service.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub Device Provisioning Service.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub Device Provisioning Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub Device Provisioning Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub Device Provisioning Service.

## Import

IoT Device Provisioning Service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_dps.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/provisioningServices/example
```
