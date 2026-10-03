---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_batch_rule_set"
description: |-
  Lists Front Door (standard/premium) Batch Rule Set resources.
---

# List resource: azurerm_cdn_frontdoor_batch_rule_set

Lists Front Door (standard/premium) Batch Rule Set resources.

## Example Usage

### List Front Door Batch Rule Sets for a Profile

```hcl
list "azurerm_cdn_frontdoor_batch_rule_set" "example" {
  provider = azurerm
  config {
    cdn_frontdoor_profile_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cdn_frontdoor_profile_id` - (Required) The resource ID of the Front Door Profile to query.
