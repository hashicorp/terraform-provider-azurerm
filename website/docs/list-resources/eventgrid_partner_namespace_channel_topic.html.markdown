---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_namespace_channel_topic"
description: |-
  Lists Event Grid Partner Namespace Channels for Partner Topics.
---

# List resource: azurerm_eventgrid_partner_namespace_channel_topic

Lists Event Grid Partner Namespace Channels of type `PartnerTopic` within a Partner Namespace.

## Example Usage

```hcl
list "azurerm_eventgrid_partner_namespace_channel_topic" "example" {
  provider = azurerm
  config {
    partner_namespace_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.EventGrid/partnerNamespaces/example-namespace"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `partner_namespace_id` - (Required) The ID of the Event Grid Partner Namespace to query.
