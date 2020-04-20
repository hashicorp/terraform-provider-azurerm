---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_consumer_group"
description: |-
  Gets information about an Event Hubs Consumer Group within an Event Hub.
---

# Data Source: azurerm_eventhub_consumer_group

Use this data source to access information about an existing Event Hubs Consumer Group within an Event Hub.

## Example Usage

```hcl
data "azurerm_eventhub_consumer_group" "test" {
  name                = "${azurerm_eventhub_consumer_group.test.name}"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  eventhub_name       = "${azurerm_eventhub.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
```

## Argument Reference

* `name` - Specifies the name of the EventHub Consumer Group resource.

* `namespace_name` - Specifies the name of the grandparent EventHub Namespace.

* `eventhub_name` - Specifies the name of the EventHub.

* `resource_group_name` - The name of the resource group in which the EventHub Consumer Group's grandparent Namespace exists.

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub Consumer Group ID.

* `user_metadata` - Specifies the user metadata.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Consumer Group.
