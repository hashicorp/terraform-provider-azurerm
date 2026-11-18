---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_edge_action_version"
description: |-
  Lists Edge Action Version resources.
---

# List resource: azurerm_cdn_frontdoor_edge_action_version

Lists Edge Action Version resources.

## Example Usage

### List all Versions of an Edge Action

```hcl
list "azurerm_cdn_frontdoor_edge_action_version" "example" {
  provider = azurerm
  config {
    cdn_frontdoor_edge_action_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Cdn/edgeActions/exampleedgeaction"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cdn_frontdoor_edge_action_id` - (Required) The ID of the Front Door Edge Action to query.
