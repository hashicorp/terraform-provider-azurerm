---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_namespace_channel"
description: |-
  Gets information about an existing Event Grid Partner Namespace Channel
---

# Data Source: azurerm_eventgrid_partner_namespace_channel

Use this data source to access information about an existing Event Grid Partner Namespace Channel.

## Example Usage

```hcl
data "azurerm_eventgrid_partner_namespace_channel" "example" {
  name                   = "my-eventgrid-partner-namespace-channel"
  partner_namespace_name = "example-partner-namespace"
  resource_group_name    = "example-resources"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Event Grid Partner Namespace Channel.

* `partner_namespace_name` - (Required) The name of the Event Grid Partner Namespace in which the Event Grid Partner Namespace exists.

* `resource_group_name` - (Required) The name of the Resource Group in which the Event Grid Partner Namespace Channel exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Event Grid Partner Namespace Channel.

* `channel_type` -  The type of the channel which represents the direction flow of events.

* `expiration_time_if_not_activated_in_utc` - The expiration time of the channel if not activated.

* `partner_topic` - A `partner_topic` block as defined below.

* `readiness_state` - The readiness state of the corresponding partner topic.

---

A `partner_topic` block exports the following:

* `name` - The name of the partner topic.

* `subscription_id` - The subscription ID of the subscriber in which the partner topic associated with the channel exists.

* `resource_group_name` - The resource group name of the subscriber in which the partner topic associated with the channel exists.

* `source` - The source information to determine the scope or context from which events are originating.

* `event_type_definitions` - An `event_type_definitions` block as defined below.

---

An `event_type_definitions` block exports the following:

* `inline_event_type` - One or more `inline_event_type` blocks as defined below.

* `kind` - The kind of event type definition.

---

An `inline_event_type` block exports the following:

* `name` - The name of the inline event type.

* `display_name` - The display name of the inline event type.

* `data_schema_url` - The data schema URL of the inline event type.

* `description` - The description of the inline event type.

* `documentation_url` - The documentation URL of the inline event type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Event Grid Partner Namespace Channel.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.EventGrid` - 2025-02-15
