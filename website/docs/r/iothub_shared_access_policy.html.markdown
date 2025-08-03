---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_shared_access_policy"
description: |-
  Manages an IotHub Shared Access Policy
---

# azurerm_iothub_shared_access_policy

Manages an IotHub Shared Access Policy

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_iothub" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_iothub_shared_access_policy" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  iothub_name         = azurerm_iothub.example.name

  registry_read  = true
  registry_write = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the IotHub Shared Access Policy resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the IotHub Shared Access Policy resource has to be created. Changing this forces a new resource to be created.

* `iothub_name` - (Required) The name of the IoTHub to which this Shared Access Policy belongs. Changing this forces a new resource to be created.

* `registry_read` - (Optional) Adds `RegistryRead` permission to this Shared Access Account. It allows read access to the identity registry.

* `registry_write` - (Optional) Adds `RegistryWrite` permission to this Shared Access Account. It allows write access to the identity registry.

-> **Note:** When `registry_write` is set to `true`, `registry_read` must also be set to true. This is a limitation of the Azure REST API

* `service_connect` - (Optional) Adds `ServiceConnect` permission to this Shared Access Account. It allows sending and receiving on the cloud-side endpoints.

* `device_connect` - (Optional) Adds `DeviceConnect` permission to this Shared Access Account. It allows sending and receiving on the device-side endpoints.

-> **Note:** At least one of `registry_read`, `registry_write`, `service_connect`, `device_connect` permissions must be set to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoTHub Shared Access Policy.

* `primary_key` - The primary key used to create the authentication token.

* `primary_connection_string` - The primary connection string of the Shared Access Policy.

* `secondary_key` - The secondary key used to create the authentication token.

* `secondary_connection_string` - The secondary connection string of the Shared Access Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub Shared Access Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub Shared Access Policy.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub Shared Access Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub Shared Access Policy.

## Import

IoTHub Shared Access Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_shared_access_policy.shared_access_policy1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/iotHubs/hub1/iotHubKeys/shared_access_policy1
```
