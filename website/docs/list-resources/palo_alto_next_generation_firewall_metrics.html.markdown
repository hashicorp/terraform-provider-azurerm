---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_next_generation_firewall_metrics"
description: |-
  Lists Palo Alto Next Generation Firewall Metrics resources.
---

# List resource: azurerm_palo_alto_next_generation_firewall_metrics

Lists Palo Alto Next Generation Firewall Metrics resources.

## Example Usage

### List the Metrics configuration of a specific Palo Alto Next Generation Firewall

```hcl
list "azurerm_palo_alto_next_generation_firewall_metrics" "example" {
  provider = azurerm
  config {
    firewall_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/PaloAltoNetworks.Cloudngfw/firewalls/firewall1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `firewall_id` - (Required) The ID of the Palo Alto Next Generation Firewall to query.
