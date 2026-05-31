---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_batch_rule"
description: |-
  Lists Front Door (standard/premium) Batch Rule collections for a Rule Set.
---

# List resource: azurerm_cdn_frontdoor_batch_rule

Lists Front Door (standard/premium) Batch Rule collections for a Rule Set.

## Example Usage

### List the Batch Rule Collection for a Front Door Rule Set

```hcl
list "azurerm_cdn_frontdoor_batch_rule" "example" {
  provider = azurerm
  config {
    cdn_frontdoor_rule_set_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cdn_frontdoor_rule_set_id` - (Required) The resource ID of the Front Door Rule Set to query.

The list result returns the ordered batch rule collection configured for the specified Rule Set when `batch_mode_enabled` is `true`.
