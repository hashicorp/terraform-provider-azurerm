---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_device_security_group"
description: |-
  Manages a Device Security Group.
---

# azurerm_device_security_group

Manages a Device Security Group.

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

resource "azurerm_device_security_group" "example" {
  name = "example-device-security-group"
  target_resource_id = azurerm_iothub.example.id

  allow_list_rule {
    type = "LocalUserNotAllowed"
    values = ["user0"]
  }

  time_window_rule {
    type = "ActiveConnectionsNotInAllowedRange"
    min_threshold = 0
    max_threshold = 30
    time_window_size = "PT5M"
  }

  depends_on = [azurerm_iot_security_solution.example]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Device Security Group. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the Azure Resource which to create the Device Security Group. Changing this forces a new resource to be created.

* `allow_list_rule` - (Optional) One or more `allow_list_rule` blocks as defined below.

* `time_window_rule` - (Optional) One or more `time_window_rule` blocks as defined below.

---

An `allow_list_rule` block supports the following:

* `type` - (Required) The type of supported rule type. Possible Values are `ConnectionToIpNotAllowed`, `LocalUserNotAllowed` and `LocalUserNotAllowed`.

* `values` - (Required) The values to allow.

---

An `time_window_rule` block supports the following:

* `type` - (Required) The type of supported rule type. Possible Values are `ActiveConnectionsNotInAllowedRange`, `AmqpC2DMessagesNotInAllowedRange`, `MqttC2DMessagesNotInAllowedRange`, `HttpC2DMessagesNotInAllowedRange`, `AmqpC2DRejectedMessagesNotInAllowedRange`, `MqttC2DRejectedMessagesNotInAllowedRange`, `HttpC2DRejectedMessagesNotInAllowedRange`, `AmqpD2CMessagesNotInAllowedRange`, `MqttD2CMessagesNotInAllowedRange`, `HttpD2CMessagesNotInAllowedRange`, `DirectMethodInvokesNotInAllowedRange`, `FailedLocalLoginsNotInAllowedRange`, `FileUploadsNotInAllowedRange`, `QueuePurgesNotInAllowedRange`, `TwinUpdatesNotInAllowedRange` and `UnauthorizedOperationsNotInAllowedRange`.

* `max_threshold` - (Required) The maximum threshold in the given time window.

* `min_threshold` - (Required) The minimum threshold in the given time window.

* `time_window_size` - (Required) Specifies the time range. represented in ISO 8601 duration format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Device Security Group resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Device Security Group.
* `update` - (Defaults to 30 minutes) Used when updating the Device Security Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Device Security Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Device Security Group.

## Import

Device Security Group can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_device_security_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/cwz/providers/Microsoft.Security/DeviceSecurityGroups/group1
```
