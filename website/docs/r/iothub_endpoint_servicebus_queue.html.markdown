---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_endpoint_servicebus_queue"
sidebar_current: "docs-azurerm-resource-messaging-iothub-servicebus-queue"
description: |-
  Manages an IotHub ServiceBus Queue Endpoint
---

# azurerm_iothub_endpoint_eventhub

Manages an IotHub ServiceBus Queue Endpoint

~> **NOTE:** Endpoints can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_endpoint_*` resources - but the two ways of defining the endpoints cannot be used together. If both are used against the same IoTHub, spurious changes will occur. Also, defining a `azurerm_iothub_endpoint_*` resource and another endpoint of a different type directly on the `azurerm_iothub` resource is not supported.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "exampleNamespace"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "example" {
  name                = "exampleQueue"
  resource_group_name = "${azurerm_resource_group.example.name}"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"

  enable_partitioning = true
}

resource "azurerm_servicebus_queue_authorization_rule" "example" {
  name                = "exampleRule"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"
  queue_name          = "${azurerm_servicebus_queue.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  listen = false
  send   = true
  manage = false
}

resource "azurerm_iothub" "example" {
  name                = "exampleIothub"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  sku {
    name     = "B1"
    tier     = "Basic"
    capacity = "1"
  }

  tags = {
    purpose = "example"
  }
}

resource "azurerm_iothub_endpoint_servicebus_queue" "example" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  iothub_name         = "${azurerm_iothub.example.name}"
  name                = "example"

  connection_string = "${azurerm_servicebus_queue_authorization_rule.example.primary_connection_string}"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the endpoint. The name must be unique across endpoint types. The following names are reserved:  `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`.

* `connection_string` - (Required) The connection string for the endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTHub ServiceBus Queue Endpoint.

## Import

IoTHub ServiceBus Queue Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_endpoint_servicebus_queue.servicebus_queue1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/IotHubs/hub1/Endpoints/servicebusqueue_endpoint1
```
