---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall_policy_rule_collection_group"
description: |-
    Lists firewall policy rule collection group resources.
---

# List resource: azurerm_firewall_policy_rule_collection_group

Lists firewall policy rule collection group resources.

## Example Usage

### List firewall policy rule collection groups in a firewall firewall policy

```hcl
list "azurerm_firewall_policy_rule_collection_group" "example" {
  provider = azurerm
  config {
    firewall_policy_id = "example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `firewall_policy_id` - (Required) The ID of the firewall policy to query.
