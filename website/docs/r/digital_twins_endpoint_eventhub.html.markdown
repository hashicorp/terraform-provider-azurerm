e---
subcategory: "Digital Twins"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_digital_twins_endpoint_eventhub"
description: |-
  Manages a Digital Twins Event Hub Endpoint.
---

# azurerm_digital_twins_endpoint_eventhub

Manages a Digital Twins Event Hub Endpoint.

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

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-eh-ns"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku = "Standard"
}

resource "azurerm_eventhub" "example" {
  name                = "example-eh"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name

  partition_count   = 2
  message_retention = 1
}

resource "azurerm_eventhub_authorization_rule" "example" {
  name                = "example-ar"
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name

  listen = false
  send   = true
  manage = false
}

resource "azurerm_digital_twins_endpoint_eventhub" "example" {
  name                                 = "example-EH"
  digital_twins_id                     = azurerm_digital_twins_instance.example.id
  eventhub_primary_connection_string   = azurerm_eventhub_authorization_rule.example.primary_connection_string
  eventhub_secondary_connection_string = azurerm_eventhub_authorization_rule.example.secondary_connection_string
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Digital Twins Event Hub Endpoint. Changing this forces a new Digital Twins Event Hub Endpoint to be created.

* `digital_twins_id` - (Required) The resource ID of the Digital Twins Instance. Changing this forces a new Digital Twins Event Hub Endpoint to be created.

* `eventhub_primary_connection_string` - (Required) The primary connection string of the Event Hub Authorization Rule with a minimum of `send` permission. 

* `eventhub_secondary_connection_string` - (Required) The secondary connection string of the Event Hub Authorization Rule with a minimum of `send` permission.

* `dead_letter_storage_secret` - (Optional) The storage secret of the dead-lettering, whose format is `https://<storageAccountname>.blob.core.windows.net/<containerName>?<SASToken>`. When an endpoint can't deliver an event within a certain time period or after trying to deliver the event a certain number of times, it can send the undelivered event to a storage account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Digital Twins Event Hub Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Digital Twins Eventhub Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Digital Twins Eventhub Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Digital Twins Eventhub Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Digital Twins Eventhub Endpoint.

## Import

Digital Twins Eventhub Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_digital_twins_endpoint_eventhub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/dt1/endpoints/ep1
```
