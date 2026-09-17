---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_origin_group"
description: |-
    Lists CDN FrontDoor Origin Group resources.
---

# List resource: azurerm_cdn_frontdoor_origin_group

Lists CDN FrontDoor Origin Group resources.

## Example Usage

### List all CDN FrontDoor Origin Groups in a CDN FrontDoor Profile

```hcl
list "azurerm_cdn_frontdoor_origin_group" "example" {
  provider = azurerm
  config {
    cdn_frontdoor_profile_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Cdn/profiles/example-profile"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cdn_frontdoor_profile_id` - (Required) The ID of the CDN FrontDoor Profile to query.
