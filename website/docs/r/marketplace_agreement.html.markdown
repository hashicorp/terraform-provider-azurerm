---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_marketplace_agreement"
sidebar_current: "docs-azurerm-resource-compute-marketplace-agreement"
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

* `offer` - (Required) The Offer of the Marketplace Image. Changing this forces a new resource to be created.

* `plan` - (Required) The Plan of the Marketplace Image. Changing this forces a new resource to be created.

* `publisher` - (Required) The Publisher of the Marketplace Image. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the Marketplace Agreement.

## Import

Marketplace Agreement can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_marketplace_agreement.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.MarketplaceOrdering/offerTypes/virtualmachine/publishers/publisher1/offers/offer1/plans/plan1/agreements/current
```
