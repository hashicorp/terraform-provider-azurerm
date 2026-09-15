---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_namespace_topic_id_association"
description: |-
  Manages an Event Grid Namespace Topic ID Association.
---

# azurerm_eventgrid_namespace_topic_id_association

Manages an Event Grid Namespace Topic ID Association.

~> **Note:** `azurerm_eventgrid_namespace_topic_id_association` should not be created at the same time with `topic_spaces_configuration.0.route_topic_id` of `azurerm_eventgrid_namespace` resource being set to avoid perpetual change due to different topic ID set via the two methods.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "southeastasia"
}

resource "azurerm_eventgrid_namespace" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_eventgrid_namespace_topic" "example" {
  name                   = "example"
  eventgrid_namespace_id = azurerm_eventgrid_namespace.example.id
}

resource "azurerm_eventgrid_namespace_topic_id_association" example {
  eventgrid_namespace_id       = azurerm_eventgrid_namespace.example.id
  eventgrid_namespace_topic_id = azurerm_eventgrid_namespace_topic.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `eventgrid_namespace_id` - (Required) The Event Grid Namespace ID. Changing this forces a new resource to be created.

* `eventgrid_namespace_topic_id` - (Required) The ID of Event Grid Namespace Topic which is associated with the Event Grid Namespace. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Event Grid Namespace ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Event Grid Namespace Topic ID Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Event Grid Namespace Topic ID Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Event Grid Namespace Topic ID Association.

## Import

Event Grid Namespace Topic ID Association can be imported using the `resource id` of Event Grid Namespace, e.g.

```shell
terraform import azurerm_eventgrid_namespace_topic_id_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.EventGrid/namespaces/namespace1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.EventGrid` - 2025-02-15
