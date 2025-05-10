---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_consumer_group"
description: |-
  Manages a Event Hubs Consumer Group as a nested resource within an Event Hub.
---

# azurerm_eventhub_consumer_group

Manages a Event Hubs Consumer Group as a nested resource within an Event Hub.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "acceptanceTestEventHubNamespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Basic"
  capacity            = 2

  tags = {
    environment = "Production"
  }
}

resource "azurerm_eventhub" "example" {
  name                = "acceptanceTestEventHub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 2
}

resource "azurerm_eventhub_consumer_group" "example" {
  name                = "acceptanceTestEventHubConsumerGroup"
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name
  user_metadata       = "some-meta-data"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventHub Consumer Group resource. Changing this forces a new resource to be created.

* `namespace_name` - (Required) Specifies the name of the grandparent EventHub Namespace. Changing this forces a new resource to be created.

* `eventhub_name` - (Required) Specifies the name of the EventHub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the EventHub Consumer Group's grandparent Namespace exists. Changing this forces a new resource to be created.

* `user_metadata` - (Optional) Specifies the user metadata.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the EventHub Consumer Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventHub Consumer Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Consumer Group.
* `update` - (Defaults to 30 minutes) Used when updating the EventHub Consumer Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventHub Consumer Group.

## Import

EventHub Consumer Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_consumer_group.consumerGroup1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumerGroups/consumerGroup1
```
