---
subcategory: "Billing"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_billing_mpa_account_scope"
description: |-
  This is a helper Data Source to provide a correctly formatted Billing Scope ID for a Microsoft Partner Agreement Account.
---

# Data Source: azurerm_billing_mpa_account_scope

Use this data source to access an ID for your MPA Account billing scope.

## Example Usage

```hcl
data "azurerm_billing_mpa_account_scope" "example" {
  billing_account_name = "e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31"
  customer_name        = "2281f543-7321-4cf9-1e23-edb4Oc31a31c"
}

output "id" {
  value = data.azurerm_billing_mpa_account_scope.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `billing_account_name` - (Required) The Billing Account Name of the MPA account.

* `customer_name` - (Required) The Customer Name in the above Billing Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Billing Scope.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Billing Scope.
