---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_consumer_group"
description: |-
  Manages a Consumer Group within an IotHub
---

# azurerm_iothub_consumer_group

Manages a Consumer Group within an IotHub

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_iothub" "example" {
  name                = "test"
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

resource "azurerm_iothub_consumer_group" "example" {
  name                   = "terraform"
  iothub_name            = azurerm_iothub.example.name
  eventhub_endpoint_name = "events"
  resource_group_name    = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of this Consumer Group. Changing this forces a new resource to be created.

* `iothub_name` - (Required) The name of the IoT Hub. Changing this forces a new resource to be created.

* `eventhub_endpoint_name` - (Required) The name of the Event Hub-compatible endpoint in the IoT hub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group that contains the IoT hub. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTHub Consumer Group.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub Consumer Group.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub Consumer Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub Consumer Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub Consumer Group.

## Import

IoTHub Consumer Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_consumer_group.group1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/IotHubs/hub1/eventHubEndpoints/events/ConsumerGroups/group1
```
