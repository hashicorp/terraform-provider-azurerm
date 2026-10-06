---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall_policy_rule_collection_group"
description: |-
    Lists Firewall Policy Rule Collection Group resources.
---

# List resource: azurerm_firewall_policy_rule_collection_group

Lists Firewall Policy Rule Collection Group resources.

## Example Usage

### List Firewall Policy Rule Collection Groups in a Firewall Policy

```hcl
list "azurerm_firewall_policy_rule_collection_group" "example" {
  provider = azurerm
  config {
    firewall_policy_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/firewallPolicies/policy1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `firewall_policy_id` - (Required) The ID of the Firewall Policy to query.
