---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_forwarding_rule"
description: |-
    Lists Private DNS Resolver Forwarding Rule resources.
---

# List resource: azurerm_private_dns_resolver_forwarding_rule

Lists Private DNS Resolver Forwarding Rule resources.

## Example Usage

### List all Private DNS Resolver Forwarding Rules in a DNS Forwarding Ruleset

```hcl
list "azurerm_private_dns_resolver_forwarding_rule" "example" {
  provider = azurerm
  config {
    dns_forwarding_ruleset_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/dnsForwardingRulesets/example-ruleset"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `dns_forwarding_ruleset_id` - (Required) The ID of the Private DNS Resolver DNS Forwarding Ruleset to query.
