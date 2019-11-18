---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_consumer_group"
sidebar_current: "docs-azurerm-resource-messaging-iothub-consumer-group"
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
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_consumer_group" "example" {
  name                   = "terraform"
  iothub_name            = "${azurerm_iothub.example.name}"
  eventhub_endpoint_name = "events"
  resource_group_name    = "${azurerm_resource_group.foo.name}"
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

## Import

IoTHub Consumer Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_consumer_group.group1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/IotHubs/hub1/eventHubEndpoints/events/ConsumerGroups/group1
```
