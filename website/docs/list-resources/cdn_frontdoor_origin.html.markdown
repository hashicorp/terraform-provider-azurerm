---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_origin"
description: |-
    Lists Cdn Frontdoor Origin resources.
---

# List resource: azurerm_cdn_frontdoor_origin

Lists Cdn Frontdoor Origin resources.

## Example Usage

### List Cdn Frontdoor Origins in a Cdn Frontdoor Origin Group

```hcl
list "azurerm_cdn_frontdoor_origin" "example" {
  provider = azurerm
  config {
    cdn_frontdoor_origin_group_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cdn_frontdoor_origin_group_id` - (Required) The ID of the CDN FrontDoor Origin Group to query.
