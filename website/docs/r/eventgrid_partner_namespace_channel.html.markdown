---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_namespace_channel"
description: |-
  Manages an EventGrid Partner Namespace Channel.
---

# azurerm_eventgrid_partner_namespace_channel

Manages an EventGrid Partner Namespace Channel.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_client_config" "example" {}

resource "azurerm_eventgrid_partner_registration" "example" {
  name                = "example-partner-registration"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_eventgrid_partner_namespace" "example" {
  name                    = "example-partner-namespace"
  resource_group_name     = azurerm_resource_group.example.name
  location                = azurerm_resource_group.example.location
  partner_registration_id = azurerm_eventgrid_partner_registration.example.id
}

resource "azurerm_eventgrid_partner_namespace_channel" "example" {
  name                 = "example-partner-namespace-channel"
  partner_namespace_id = azurerm_eventgrid_partner_namespace.example.id
  channel_type         = "PartnerTopic"
  partner_topic {
    name                = "example-partner-topic"
    subscription_id     = data.azurerm_client_config.example.subscription_id
    resource_group_name = azurerm_resource_group.example.name
    source              = "example.source"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Event Grid Partner Namespace Channel. Changing this forces a new Event Grid Partner Namespace Channel to be created.

* `partner_namespace_id` - (Required) The ID of the Event Grid Partner Namespace where the channel should exist. Changing this forces a new Event Grid Partner Namespace Channel to be created.

* `partner_topic` - (Required) A `partner_topic` block as defined below.

---

* `channel_type` - (Optional) The type of the channel which represents the direction flow of events. The only possible value is `PartnerTopic`. Defaults to `PartnerTopic`. Changing this forces a new Event Grid Partner Namespace Channel to be created.

* `expiration_time_if_not_activated_in_utc` - (Optional) The expiration time of the channel if not activated (Datetime Format `RFC 3339`).

~> **Note:** Once `readiness_state` is `Activated`, this field can no longer be updated.

~> **Note:** If this timer expires while the corresponding partner topic is never activated, the channel and corresponding partner topic are deleted.

---

A `partner_topic` block supports the following:

* `name` - (Required) The name of the partner topic. Changing this forces a new Event Grid Partner Namespace Channel to be created.

* `subscription_id` - (Required) The subscription ID of the subscriber in which the partner topic associated with the channel will be created under. Changing this forces a new Event Grid Partner Namespace Channel to be created.

* `resource_group_name` - (Required) The resource group name of the subscriber in which the partner topic associated with the channel will be created under. Changing this forces a new Event Grid Partner Namespace Channel to be created.

* `source` - (Required) The source information to determine the scope or context from which events are originating. Changing this forces a new Event Grid Partner Namespace Channel to be created.

* `event_type_definitions` - (Optional) An `event_type_definitions` block as defined below.

---

An `event_type_definitions` block supports the following:

* `inline_event_type` - (Required) One or more `inline_event_type` blocks as defined below.

* `kind` - (Optional) The kind of event type definition. The only possible value is `Inline`. Defaults to `Inline`.

---

An `inline_event_type` block supports the following:

* `name` - (Required) The name of the inline event type.

* `display_name` - (Required) The display name of the inline event type.

* `data_schema_url` - (Optional) The data schema URL of the inline event type.

* `description` - (Optional) The description of the inline event type.

* `documentation_url` - (Optional) The documentation URL of the inline event type.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Event Grid Partner Namespace Channel.

* `readiness_state` - The readiness state of the corresponding partner topic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Event Grid Partner Namespace Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Event Grid Partner Namespace Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Event Grid Partner Namespace Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Event Grid Partner Namespace Channel.

## Import

Event Grid Partner Namespace Channels can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_partner_namespace_channel.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.EventGrid/partnerNamespaces/example/channels/example
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.EventGrid` - 2025-02-15
