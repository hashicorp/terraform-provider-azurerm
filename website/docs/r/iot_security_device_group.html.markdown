---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iot_security_device_group"
description: |-
  Manages a Iot Security Device Group.
---

# azurerm_iot_security_device_group

Manages a Iot Security Device Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_iothub" "example" {
  name                = "example-IoTHub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_iot_security_solution" "example" {
  name                = "example-Iot-Security-Solution"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  display_name        = "Iot Security Solution"
  iothub_ids          = [azurerm_iothub.example.id]
}

resource "azurerm_iot_security_device_group" "example" {
  name      = "example-device-security-group"
  iothub_id = azurerm_iothub.example.id

  allow_rule {
    connection_to_ips_not_allowed = ["10.0.0.0/24"]
  }

  range_rule {
    type     = "ActiveConnectionsNotInAllowedRange"
    min      = 0
    max      = 30
    duration = "PT5M"
  }

  depends_on = [azurerm_iot_security_solution.example]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Device Security Group. Changing this forces a new resource to be created.

* `iothub_id` - (Required) The ID of the IoT Hub which to link the Security Device Group to. Changing this forces a new resource to be created.

* `allow_rule` - (Optional) an `allow_rule` blocks as defined below.

* `range_rule` - (Optional) One or more `range_rule` blocks as defined below.

---

An `allow_rule` block supports the following:

* `connection_from_ips_not_allowed` - (Optional) Specifies which IP is not allowed to be connected to in current device group for inbound connection.

* `connection_to_ips_not_allowed` - (Optional) Specifies which IP is not allowed to be connected to in current device group for outbound connection.

* `local_users_not_allowed` - (Optional) Specifies which local user is not allowed to login in current device group.

* `processes_not_allowed` - (Optional) Specifies which process is not allowed to be executed in current device group.

---

An `range_rule` block supports the following:

* `duration` - (Required) Specifies the time range. represented in ISO 8601 duration format.

* `max` - (Required) The maximum threshold in the given time window.

* `min` - (Required) The minimum threshold in the given time window.

* `type` - (Required) The type of supported rule type. Possible Values are `ActiveConnectionsNotInAllowedRange`, `AmqpC2DMessagesNotInAllowedRange`, `MqttC2DMessagesNotInAllowedRange`, `HttpC2DMessagesNotInAllowedRange`, `AmqpC2DRejectedMessagesNotInAllowedRange`, `MqttC2DRejectedMessagesNotInAllowedRange`, `HttpC2DRejectedMessagesNotInAllowedRange`, `AmqpD2CMessagesNotInAllowedRange`, `MqttD2CMessagesNotInAllowedRange`, `HttpD2CMessagesNotInAllowedRange`, `DirectMethodInvokesNotInAllowedRange`, `FailedLocalLoginsNotInAllowedRange`, `FileUploadsNotInAllowedRange`, `QueuePurgesNotInAllowedRange`, `TwinUpdatesNotInAllowedRange` and `UnauthorizedOperationsNotInAllowedRange`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Iot Security Device Group resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Iot Security Device Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Iot Security Device Group.
* `update` - (Defaults to 30 minutes) Used when updating the Iot Security Device Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Iot Security Device Group.

## Import

Iot Security Device Group can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_security_device_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/providers/Microsoft.Security/deviceSecurityGroups/group1
```
