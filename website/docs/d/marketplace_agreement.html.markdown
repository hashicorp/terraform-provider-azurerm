---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_marketplace_agreement"
description: |-
  Gets information about an existing Marketplace Agreement.
---

# Data Source: azurerm_marketplace_agreement

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

output "azurerm_marketplace_agreement_accepted" {
  value = data.azurerm_marketplace_agreement.accepted
}
```

## Arguments Reference

The following arguments are supported:

* `offer` - The Offer of the Marketplace Image.

* `plan` - The Plan of the Marketplace Image.

* `publisher` - The Publisher of the Marketplace Image.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Marketplace Agreement.

* `accepted` - Whether the Marketplace Agreement has been accepted.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Marketplace Agreement.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.MarketplaceOrdering` - 2015-06-01
