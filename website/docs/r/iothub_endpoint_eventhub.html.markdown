---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_endpoint_eventhub"
sidebar_current: "docs-azurerm-resource-messaging-iothub-endpoint-eventhub"
description: |-
  Manages an IotHub EventHub Endpoint
---

# azurerm_iothub_endpoint_eventhub

Manages an IotHub EventHub Endpoint

~> **NOTE:** Endpoints can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_endpoint_*` resources - but the two ways of defining the endpoints cannot be used together. If both are used against the same IoTHub, spurious changes will occur. Also, defining a `azurerm_iothub_endpoint_*` resource and another endpoint of a different type directly on the `azurerm_iothub` resource is not supported.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "exampleEventHubNamespace"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "example" {
  name                = "exampleEventHub"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "example" {
  name                = "exampleRule"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  eventhub_name       = "${azurerm_eventhub.example.name}"
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

resource "azurerm_iothub_endpoint_eventhub" "example" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  iothub_name         = "${azurerm_iothub.example.name}"
  name                = "example"
  
  connection_string = "${azurerm_eventhub_authorization_rule.example.primary_connection_string}"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the endpoint. The name must be unique across endpoint types. The following names are reserved:  `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`.

* `connection_string` - (Required) The connection string for the endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTHub EventHub Endpoint.

## Import

IoTHub EventHub Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_endpoint_eventhub.eventhub1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/IotHubs/hub1/Endpoints/eventhub_endpoint1
```
