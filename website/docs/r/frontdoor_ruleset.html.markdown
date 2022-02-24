---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_rule_set"
description: |-
  Manages a Frontdoor Rule Set.
---

# azurerm_frontdoor_rule_set

Manages a Frontdoor Rule Set.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-cdn"
  location = "West Europe"
}

resource "azurerm_frontdoor_profile" "test" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_frontdoor_rule_set" "test" {
  name                 = "exampleruleset"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Rule Set. Changing this forces a new Frontdoor Rule Set to be created.

* `frontdoor_profile_id` - (Required) The ID of the Frontdoor Profile. Changing this forces a new Frontdoor Rule Set to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Rule Set.

* `frontdoor_profile_name` - The name of the Frontdoor Profile containing this Frontdoor Rule Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Rule Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Rule Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Rule Set.

## Import

Frontdoor Rule Sets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_rule_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1
```
