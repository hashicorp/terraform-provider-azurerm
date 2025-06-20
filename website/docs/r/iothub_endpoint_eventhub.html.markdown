---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_endpoint_eventhub"
description: |-
  Manages an IotHub EventHub Endpoint
---

# azurerm_iothub_endpoint_eventhub

Manages an IotHub EventHub Endpoint

~> **Note:** Endpoints can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_endpoint_*` resources - but the two ways of defining the endpoints cannot be used together. If both are used against the same IoTHub, spurious changes will occur. Also, defining a `azurerm_iothub_endpoint_*` resource and another endpoint of a different type directly on the `azurerm_iothub` resource is not supported.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "exampleEventHubNamespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "example" {
  name                = "exampleEventHub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "example" {
  name                = "exampleRule"
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name

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

resource "azurerm_iothub_endpoint_eventhub" "example" {
  resource_group_name = azurerm_resource_group.example.name
  iothub_id           = azurerm_iothub.example.id
  name                = "example"

  connection_string = azurerm_eventhub_authorization_rule.example.primary_connection_string
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the endpoint. The name must be unique across endpoint types. The following names are reserved: `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the Event Hub has been created. Changing this forces a new resource to be created.

* `authentication_type` - (Optional) Type used to authenticate against the Event Hub endpoint. Possible values are `keyBased` and `identityBased`. Defaults to `keyBased`.

* `identity_id` - (Optional) ID of the User Managed Identity used to authenticate against the Event Hub endpoint.

-> **Note:** `identity_id` can only be specified when `authentication_type` is `identityBased`. It must be one of the `identity_ids` of the Iot Hub. If not specified when `authentication_type` is `identityBased`, System Assigned Managed Identity of the Iot Hub will be used.

* `endpoint_uri` - (Optional) URI of the Event Hubs Namespace endpoint. This attribute can only be specified and is mandatory when `authentication_type` is `identityBased`.

* `entity_path` - (Optional) Name of the Event Hub. This attribute can only be specified and is mandatory when `authentication_type` is `identityBased`.

* `connection_string` - (Optional) The connection string for the endpoint. This attribute can only be specified and is mandatory when `authentication_type` is `keyBased`.

* `iothub_id` - (Required) The IoTHub ID for the endpoint. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoTHub EventHub Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub EventHub Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub EventHub Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub EventHub Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub EventHub Endpoint.

## Import

IoTHub EventHub Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_endpoint_eventhub.eventhub1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/iotHubs/hub1/endpoints/eventhub_endpoint1
```
