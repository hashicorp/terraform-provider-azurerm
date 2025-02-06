---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_endpoint_servicebus_topic"
description: |-
  Manages an IotHub ServiceBus Topic Endpoint
---

# azurerm_iothub_endpoint_servicebus_topic

Manages an IotHub ServiceBus Topic Endpoint

~> **NOTE:** Endpoints can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_endpoint_*` resources - but the two ways of defining the endpoints cannot be used together. If both are used against the same IoTHub, spurious changes will occur. Also, defining a `azurerm_iothub_endpoint_*` resource and another endpoint of a different type directly on the `azurerm_iothub` resource is not supported.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "exampleNamespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "example" {
  name         = "exampleTopic"
  namespace_id = azurerm_servicebus_namespace.example.id
}

resource "azurerm_servicebus_topic_authorization_rule" "example" {
  name     = "exampleRule"
  topic_id = azurerm_servicebus_topic.example.id

  listen = false
  send   = true
  manage = false
}

resource "azurerm_iothub" "example" {
  name                = "exampleIothub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "example"
  }
}

resource "azurerm_iothub_endpoint_servicebus_topic" "example" {
  resource_group_name = azurerm_resource_group.example.name
  iothub_id           = azurerm_iothub.example.id
  name                = "example"

  connection_string = azurerm_servicebus_topic_authorization_rule.example.primary_connection_string
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the endpoint. The name must be unique across endpoint types. The following names are reserved: `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the Service Bus Topic has been created. Changing this forces a new resource to be created.

* `authentication_type` - (Optional) Type used to authenticate against the Service Bus Topic endpoint. Possible values are `keyBased` and `identityBased`. Defaults to `keyBased`.

* `identity_id` - (Optional) ID of the User Managed Identity used to authenticate against the Service Bus Topic endpoint.

-> **NOTE:** `identity_id` can only be specified when `authentication_type` is `identityBased`. It must be one of the `identity_ids` of the Iot Hub. If not specified when `authentication_type` is `identityBased`, System Assigned Managed Identity of the Iot Hub will be used.

* `endpoint_uri` - (Optional) URI of the Service Bus endpoint. This attribute can only be specified and is mandatory when `authentication_type` is `identityBased`.

* `entity_path` - (Optional) Name of the Service Bus Topic. This attribute can only be specified and is mandatory when `authentication_type` is `identityBased`.

* `connection_string` - (Optional) The connection string for the endpoint. This attribute can only be specified and is mandatory when `authentication_type` is `keyBased`.

* `iothub_id` - (Required) The IoTHub ID for the endpoint. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoTHub ServiceBus Topic Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub ServiceBus Topic Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub ServiceBus Topic Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub ServiceBus Topic Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub ServiceBus Topic Endpoint.

## Import

IoTHub ServiceBus Topic Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_endpoint_servicebus_topic.servicebus_topic1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/iotHubs/hub1/endpoints/servicebustopic_endpoint1
```
