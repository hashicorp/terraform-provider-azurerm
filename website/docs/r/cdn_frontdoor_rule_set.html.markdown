---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_rule_set"
description: |-
  Manages a Front Door (standard/premium) Rule Set.
---

# azurerm_cdn_frontdoor_rule_set

Manages a Front Door (standard/premium) Rule Set.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_rule_set" "example" {
  name                     = "ExampleRuleSet"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Front Door Rule Set. Changing this forces a new resource to be created.

* `cdn_frontdoor_profile_id` - (Required) The ID of the Front Door Profile. Changing this forces a new resource to be created.

~> **Note:** This resource manages the non-batch Front Door Standard/Premium Rule Set path. To create a batch-mode Rule Set with inline `rules`, use `azurerm_cdn_frontdoor_batch_rule_set` instead.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Rule Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 4 hours) Used when creating the Front Door Rule Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Rule Set.
* `delete` - (Defaults to 6 hours) Used when deleting the Front Door Rule Set.

## Import

A Front Door Rule Set can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_rule_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cdn` - 2025-12-01
