---
subcategory: "Digital Twins"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_digital_twins_endpoint_servicebus"
description: |-
  Manages a Digital Twins Service Bus Endpoint.
---

# azurerm_digital_twins_endpoint_servicebus

Manages a Digital Twins Service Bus Endpoint.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example_resources"
  location = "West Europe"
}

resource "azurerm_digital_twins_instance" "example" {
  name                = "example-DT"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "exampleservicebusnamespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "example" {
  name         = "exampleservicebustopic"
  namespace_id = azurerm_servicebus_namespace.example.id
}

resource "azurerm_servicebus_topic_authorization_rule" "example" {
  name     = "example-rule"
  topic_id = azurerm_servicebus_topic.example.id

  listen = false
  send   = true
  manage = false
}

resource "azurerm_digital_twins_endpoint_servicebus" "example" {
  name                                   = "example-EndpointSB"
  digital_twins_id                       = azurerm_digital_twins_instance.example.id
  servicebus_primary_connection_string   = azurerm_servicebus_topic_authorization_rule.example.primary_connection_string
  servicebus_secondary_connection_string = azurerm_servicebus_topic_authorization_rule.example.secondary_connection_string
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Digital Twins Service Bus Endpoint. Changing this forces a new Digital Twins Service Bus Endpoint to be created.

* `digital_twins_id` - (Required) The ID of the Digital Twins Instance. Changing this forces a new Digital Twins Service Bus Endpoint to be created.

* `servicebus_primary_connection_string` - (Required) The primary connection string of the Service Bus Topic Authorization Rule with a minimum of `send` permission. .

* `servicebus_secondary_connection_string` - (Required) The secondary connection string of the Service Bus Topic Authorization Rule with a minimum of `send` permission.

* `dead_letter_storage_secret` - (Optional) The storage secret of the dead-lettering, whose format is `https://<storageAccountname>.blob.core.windows.net/<containerName>?<SASToken>`. When an endpoint can't deliver an event within a certain time period or after trying to deliver the event a certain number of times, it can send the undelivered event to a storage account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Digital Twins Service Bus Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Digital Twins Servicebus Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Digital Twins Servicebus Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Digital Twins Servicebus Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Digital Twins Servicebus Endpoint.

## Import

Digital Twins Service Bus Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_digital_twins_endpoint_servicebus.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/dt1/endpoints/ep1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DigitalTwins`: 2023-01-31
