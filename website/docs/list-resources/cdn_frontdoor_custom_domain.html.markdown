---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_custom_domain"
description: |-
    Lists CDN FrontDoor Custom Domain resources.
---

# List resource: azurerm_cdn_frontdoor_custom_domain

Lists CDN FrontDoor Custom Domain resources.

## Example Usage

### List all CDN FrontDoor Custom Domains in a CDN FrontDoor Profile

```hcl
list "azurerm_cdn_frontdoor_custom_domain" "example" {
  provider = azurerm
  config {
    cdn_frontdoor_profile_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/myprofile1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cdn_frontdoor_profile_id` - (Required) The ID of the CDN FrontDoor Profile to query.
