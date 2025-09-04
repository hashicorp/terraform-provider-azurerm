---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_route"
description: |-
  Manages an IotHub Route
---

# azurerm_iothub_route

Manages an IotHub Route

~> **Note:** Routes can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_route` resource - but the two cannot be used together. If both are used against the same IoTHub, spurious changes will occur.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "example"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_iothub" "example" {
  name                = "exampleIothub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "example" {
  resource_group_name = azurerm_resource_group.example.name
  iothub_id           = azurerm_iothub.example.id
  name                = "example"

  connection_string          = azurerm_storage_account.example.primary_blob_connection_string
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  container_name             = azurerm_storage_container.example.name
  encoding                   = "Avro"
  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
}

resource "azurerm_iothub_route" "example" {
  resource_group_name = azurerm_resource_group.example.name
  iothub_name         = azurerm_iothub.example.name
  name                = "example"

  source         = "DeviceMessages"
  condition      = "true"
  endpoint_names = [azurerm_iothub_endpoint_storage_container.example.name]
  enabled        = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the route. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the IotHub Route resource has to be created. Changing this forces a new resource to be created.

* `iothub_name` - (Required) The name of the IoTHub to which this Route belongs. Changing this forces a new resource to be created.

* `source` - (Required) The source that the routing rule is to be applied to. Possible values include: `DeviceConnectionStateEvents`, `DeviceJobLifecycleEvents`, `DeviceLifecycleEvents`, `DeviceMessages`, `DigitalTwinChangeEvents`, `Invalid`, `TwinChangeEvents`.

* `condition` - (Optional) The condition that is evaluated to apply the routing rule. For grammar, see: <https://docs.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language>. Defaults to `true`.

* `endpoint_names` - (Required) The list of endpoints to which messages that satisfy the condition are routed. Currently only one endpoint is allowed.

* `enabled` - (Required) Specifies whether a route is enabled.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoTHub Route.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub Route.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub Route.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub Route.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub Route.

## Import

IoTHub Route can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_route.route1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/iotHubs/hub1/routes/route1
```
