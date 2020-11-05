---
subcategory: "DigitalTwins"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_digital_twins_endpoint_eventgrid"
description: |-
  Manages a Digital Twins Eventgrid Endpoint.
---

# azurerm_digital_twins_endpoint_eventgrid

Manages a Digital Twins Eventgrid Endpoint.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example_resources"
  location = "West Europe"
}

resource "azurerm_digital_twins" "example" {
  name                = "example-DT"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_eventgrid_topic" "example" {
  name                = "example-topic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_digital_twins_endpoint_eventgrid" "example" {
  name                                 = "example-EG"
  digital_twins_id                     = azurerm_digital_twins.example.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.example.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.example.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.example.secondary_access_key
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Digital Twins Eventgrid Endpoint. Changing this forces a new Digital Twins Eventgrid Endpoint to be created.

* `digital_twins_id` - (Required) The resource ID of the Digital Twins. Changing this forces a new Digital Twins Eventgrid Endpoint to be created.

* `eventgrid_topic_endpoint` - (Required) The endpoint of the EventGrid Topic.

* `eventgrid_topic_primary_access_key` - (Required) The primary accesskey of the EventGrid Topic.

* `eventgrid_topic_secondary_access_key` - (Required) The secondary accesskey of the EventGrid Topic.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Digital Twins Eventgrid Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Digital Twins Eventgrid Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Digital Twins Eventgrid Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Digital Twins Eventgrid Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Digital Twins Eventgrid Endpoint.

## Import

Digital Twins Eventgrid Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_digital_twins_endpoint_eventgrid.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/dt1/endpoints/ep1
```
