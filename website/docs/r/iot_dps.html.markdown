---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iot_dps"
sidebar_current: "docs-azurerm-resource-messaging-iot-dps-x"
description: |-
  Manages an IoT Device Provisioning Service.
---

# azurerm_iot_dps

Manages an IoT Device Provisioning Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_iot_dps" "example" {
  name                = "example"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
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

* `name` - (Required) The name of the sku. Possible values are `B1`, `B2`, `B3`, `F1`, `S1`, `S2`, and `S3`.

* `tier` - (Required) The billing tier for the IoT Device Provisioning Service. Possible values are `Basic`, `Free` or `Standard`.

* `capacity` - (Required) The number of provisioned IoT Device Provisioning Service units.

---

A `linked_hub` block supports the following:

* `connection_string` - (Required) The connection string to connect to the IoT Hub. Changing this forces a new resource.

* `location` - (Required) The location of the IoT hub. Changing this forces a new resource.

* `apply_application_policy` - (Optional) Determines whether to apply application policies to the IoT Hub. Defaults to false.

* `allocation_weight` - (Optional) The weight applied to the IoT Hub. Defaults to 0.

* `hostname` - (Computed) The IoT Hub hostname.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoT Device Provisioning Service.

## Import

IoT Device Provisioning Service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_dps.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/provisioningServices/example
```
