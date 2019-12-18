---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_subscription_pricing"
sidebar_current: "docs-azurerm-security-center-subscription-pricing"
description: |-
    Manages the Pricing Tier for Azure Security Center in the current subscription.
---

# azurerm_security_center_subscription_pricing

Manages the Pricing Tier for Azure Security Center in the current subscription.

~> **NOTE:** This resource requires the `Owner` permission on the Subscription.

~> **NOTE:** Deletion of this resource does not change or reset the pricing tier to `Free`

## Example Usage

```hcl
resource "azurerm_security_center_subscription_pricing" "example" {
  tier = "Standard"
}
```

## Argument Reference

The following arguments are supported:

* `tier` - (Required) The pricing tier to use. Possible values are `Free` and `Standard`.

~> **NOTE:** Changing the pricing tier to `Standard` affects all resources in the subscription and could be quite costly.

## Attributes Reference

The following attributes are exported:

* `id` - The subscription pricing ID.


## Import

The pricing tier can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_subscription_pricing.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/pricings/default
```
