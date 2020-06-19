---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_firewall_policy_rule_collection_group"
description: |-
  Gets information about an existing Firewall Policy Rule Collection Group.
---

# Data Source: azurerm_firewall_policy_rule_collection_group

Use this data source to access information about an existing Firewall Policy Rule Collection Group.

## Example Usage

```hcl
data "azurerm_firewall_policy_rule_collection_group" "example" {
  name               = "existing"
  firewall_policy_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/firewallPolicies/policy1"
}

output "id" {
  value = data.azurerm_firewall_policy_rule_collection_group.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Firewall Policy Rule Collection Group.

* `firewall_policy_id` - (Required) The ID of the Firewall Policy where the Firewall Policy Rule Collection Group resides in.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Firewall Policy Rule Collection Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall Policy Rule Collection Group.
