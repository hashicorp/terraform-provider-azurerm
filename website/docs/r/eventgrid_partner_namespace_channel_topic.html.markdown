---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_namespace_channel_topic"
description: |-
  Manages an Event Grid Partner Namespace Channel for a Partner Topic.
---

# azurerm_eventgrid_partner_namespace_channel_topic

Manages an Event Grid Partner Namespace Channel for a Partner Topic. Creating the channel also creates its associated Partner Topic in the subscriber's subscription and resource group.

~> **Note:** The Partner Namespace must use the `ChannelNameHeader` partner topic routing mode. The subscriber must authorize the Partner Registration through an Event Grid Partner Configuration before the channel is created. The subscriber is responsible for activating the associated Partner Topic. See [onboarding an Event Grid partner](https://learn.microsoft.com/azure/event-grid/onboard-partner) for details.

~> **Note:** By default, the Channel and its associated Partner Topic expire seven days after creation if the Partner Topic is not activated.

## Example Usage

```hcl
# Partner-managed resources
resource "azurerm_resource_group" "partner" {
  name     = "example-partner-resources"
  location = "East US 2"
}

resource "azurerm_eventgrid_partner_registration" "partner" {
  name                = "example-registration"
  resource_group_name = azurerm_resource_group.partner.name
}

resource "azurerm_eventgrid_partner_namespace" "partner" {
  name                       = "example-namespace"
  resource_group_name        = azurerm_resource_group.partner.name
  location                   = azurerm_resource_group.partner.location
  partner_registration_id    = azurerm_eventgrid_partner_registration.partner.id
  partner_topic_routing_mode = "ChannelNameHeader"
}

# Subscriber-managed resources: destination Resource Group and authorization of the partner.
resource "azurerm_resource_group" "subscriber" {
  name     = "example-subscriber-resources"
  location = "East US 2"
}

resource "azurerm_eventgrid_partner_configuration" "subscriber" {
  resource_group_name = azurerm_resource_group.subscriber.name

  partner_authorization {
    partner_registration_id              = azurerm_eventgrid_partner_registration.partner.partner_registration_id
    partner_name                         = azurerm_eventgrid_partner_registration.partner.name
    authorization_expiration_time_in_utc = "2026-11-09T00:00:00Z"
  }
}

# Partner-managed Channel: creates a Partner Topic in the subscriber's destination Resource Group.
resource "azurerm_eventgrid_partner_namespace_channel_topic" "example" {
  name                 = "example-channel"
  partner_namespace_id = azurerm_eventgrid_partner_namespace.partner.id

  partner_topic {
    # Destination details supplied by the subscriber.
    name                = "example-partner-topic"
    resource_group_name = azurerm_resource_group.subscriber.name
    source              = "Contoso.Accounts.User1"
    subscription_id     = data.azurerm_client_config.subscriber.subscription_id
  }

  event_type {
    name              = "Contoso.OrderCreated"
    description       = "An order was created."
    display_name      = "Order created"
    documentation_url = "https://example.com/events/order-created"
    data_schema_url   = "https://example.com/schemas/order-created.json"
  }

  depends_on = [azurerm_eventgrid_partner_configuration.subscriber]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Channel. Must be between 3 and 50 characters and contain only letters, numbers and hyphens. The name must be unique within the Azure region. Changing this forces a new resource to be created.

* `partner_namespace_id` - (Required) The ID of the Event Grid Partner Namespace in which to create the Channel. Changing this forces a new resource to be created.

* `partner_topic` - (Required) A `partner_topic` block as defined below.

---

* `event_type` - (Optional) One or more `event_type` blocks as defined below.

---

A `partner_topic` block supports the following:

* `name` - (Required) The name of the associated Partner Topic. Must be between 3 and 50 characters and contain only letters, numbers and hyphens. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the subscriber's existing Resource Group in which the Partner Topic will be created. Changing this forces a new resource to be created.

* `source` - (Required) Non-empty contextual information about the source of events, displayed to the subscriber during activation. Changing this forces a new resource to be created.

* `subscription_id` - (Required) The GUID of the subscriber's Azure subscription in which the Partner Topic will be created. Changing this forces a new resource to be created.

---

An `event_type` block supports the following:

* `name` - (Required) The name of the event type, for example `Contoso.OrderCreated`. Must be between 1 and 128 characters. Names must be unique within the Channel.

* `display_name` - (Required) A non-empty display name for the event type.

* `data_schema_url` - (Optional) An HTTP or HTTPS URL describing the event's data schema.

* `description` - (Optional) A description of the event type. Must be between 1 and 256 characters.

* `documentation_url` - (Optional) An HTTP or HTTPS URL for the event type's documentation.

## Attributes Reference

In addition to the Arguments listed above, the following Attributes are exported:

* `id` - The ID of the Event Grid Partner Namespace Channel.

* `activation_message` - The Azure-generated message providing context to the subscriber when activating the associated Partner Topic.

* `expiration_time_if_not_activated` - The RFC3339 expiration timestamp for activation, the channel will be deleted if not activated by this time.

* `readiness_state` - The readiness state of the associated Partner Topic, such as `NeverActivated` or `Activated`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Channel.

## Import

Event Grid Partner Namespace Channels for Partner Topics can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_partner_namespace_channel_topic.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.EventGrid/partnerNamespaces/example-namespace/channels/example-channel
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.EventGrid` - 2025-02-15
