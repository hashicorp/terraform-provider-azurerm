---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_marketplace_agreement"
description: |-
  Gets information about an existing Marketplace Agreement.
---

# azurerm_marketplace_agreement

Uses this data source to access information about an existing Marketplace Agreement.

## Example Usage

```hcl
data "azurerm_marketplace_agreement" "barracuda" {
  publisher = "barracudanetworks"
  offer     = "waf"
  plan      = "hourly"
}

output "azurerm_marketplace_agreement_id" {
  value = data.azurerm_marketplace_agreement.id
}
```

## Argument Reference

The following arguments are supported:

* `offer` - The Offer of the Marketplace Image.

* `plan` - The Plan of the Marketplace Image.

* `publisher` - The Publisher of the Marketplace Image.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Marketplace Agreement.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Marketplace Agreement.
