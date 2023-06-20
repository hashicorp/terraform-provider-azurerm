---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_rule_set"
description: |-
  Gets information about an existing Front Door (standard/premium) Rule Set.
---

# Data Source: azurerm_cdn_frontdoor_rule_set

Use this data source to access information about an existing Front Door (standard/premium) Rule Set.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_rule_set" "example" {
  name                = "existing-rule-set"
  profile_name        = "existing-profile"
  resource_group_name = "existing-resources"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Front Door Rule Set to retrieve.

* `profile_name` - (Required) Specifies the name of the Front Door Profile where this Front Door Rule Set exists.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Front Door Profile exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Rule Set.

* `cdn_frontdoor_profile_id` - The ID of the Front Door Profile within which this Front Door Rule Set exists.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Rule Set.
