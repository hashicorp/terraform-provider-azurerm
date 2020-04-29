---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_marketplace_agreement"
description: |-
  Allows accepting the Legal Terms for a Marketplace Image.
---

# azurerm_marketplace_agreement

Allows accepting the Legal Terms for a Marketplace Image.

## Example Usage

```hcl
resource "azurerm_marketplace_agreement" "barracuda" {
  publisher = "barracudanetworks"
  offer     = "waf"
  plan      = "hourly"
}
```

## Argument Reference

The following arguments are supported:

* `publisher` - (Required) The Publisher of the Marketplace Image. Changing this forces a new resource to be created.

* `offer` - (Optional) The Offer of the Marketplace Image, also known as `product`. Changing this forces a new resource to be created.

* `product` - (Optional) The Product of the Marketplace Image, also known as `offer`. Changing this forces a new resource to be created.

~> **Note:** One of `offer` or `product` must be specified.

* `plan` - (Optional) The Plan of the Marketplace Image, also known as `name`. Changing this forces a new resource to be created.

* `name` - (Optional) The Plan of the Marketplace Image,also known as `plan`. Changing this forces a new resource to be created.

~> **Note:** One of `plan` or `name` must be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the Marketplace Agreement.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Marketplace Agreement.
* `update` - (Defaults to 30 minutes) Used when updating the Marketplace Agreement.
* `read` - (Defaults to 5 minutes) Used when retrieving the Marketplace Agreement.
* `delete` - (Defaults to 30 minutes) Used when deleting the Marketplace Agreement.

## Import

Marketplace Agreement can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_marketplace_agreement.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.MarketplaceOrdering/offerTypes/virtualmachine/publishers/publisher1/offers/offer1/plans/plan1/agreements/current
```
