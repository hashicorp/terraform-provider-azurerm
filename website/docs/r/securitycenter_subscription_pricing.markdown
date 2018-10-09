---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_securitycenter_subscription_pricing"
sidebar_current: "docs-azurerm-securitycenter-subscription-pricing"
description: |-
    Manages the subscription's Security Center pricing tier.
---

# azurerm_securitycenter_subscription_pricing

Manages the subscription's Security Center pricing tier.

## Example Usage

```hcl
resource "azurerm_securitycenter_subscription_pricing" "example" {
    tier = "Standard"
}
```

## Argument Reference

The following arguments are supported:

* `tier` - (Required) The pricing tier to use. Must be one of `Free` or `Standard`

~> **NOTE:** Changing the pricing tier to `Standard` affects all resources in the subscription and could be quite costly.

## Attributes Reference

The following attributes are exported:

* `id` - The subscription pricing ID.


## Import

Resource Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_securitycenter_subscription_pricing.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/pricings/default
```
