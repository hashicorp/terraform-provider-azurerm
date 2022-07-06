---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_profile"
description: |-
  Gets information about an existing CDN FrontDoor Profile.
---

# Data Source: azurerm_cdn_frontdoor_profile

Use this data source to access information about an existing CDN FrontDoor Profile.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_profile" "example" {
  name                = "existing-cdn-profile"
  resource_group_name = "existing-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the FrontDoor Profile.

* `resource_group_name` - (Required) The name of the Resource Group where this FrontDoor Profile exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of this CDN FrontDoor Profile.

* `resource_guid` - The UUID of this CDN FrontDoor Profile.

* `sku_name` - Specifies the SKU for this CDN FrontDoor Profile.

* `response_timeout_seconds` - Specifies the maximum response timeout in seconds.

* `tags` - Specifies a mapping of Tags assigned to this CDN FrontDoor Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the CDN FrontDoor Profile.
