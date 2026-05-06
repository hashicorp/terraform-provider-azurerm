---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription_zone_peers"
description: |-
  Get information about availability zone mappings between subscriptions.
---

# Data Source: azurerm_subscription_zone_peers

Use this data source to check the availability zone peering between the current subscription and a peer subscription. This is useful for determining how logical availability zones map between subscriptions, which is critical for ensuring resources in different subscriptions are co-located in the same physical zones.

## Example Usage

```hcl
data "azurerm_subscription_zone_peers" "example" {
  location             = "eastus"
  peer_subscription_id = "00000000-0000-0000-0000-000000000000"
}

output "zone_peers" {
  value = data.azurerm_subscription_zone_peers.example.availability_zone_peers
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure location for which to check zone peering.

* `peer_subscription_id` - (Required) The subscription ID of the peer subscription to check zone peering against.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this data source.

* `subscription_id` - The ID of the current subscription.

* `availability_zone_peers` - A list of `availability_zone_peers` blocks as defined below.

---

An `availability_zone_peers` block exports the following:

* `availability_zone` - The availability zone in the current subscription.

* `peers` - A list of `peers` blocks as defined below.

---

A `peers` block exports the following:

* `subscription_id` - The subscription ID of the peer subscription.

* `availability_zone` - The corresponding availability zone in the peer subscription.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when checking zone peers.
