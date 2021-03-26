---
subcategory: "Billing"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_billing_mca_account_scope"
description: |-
  This is a helper Data Source to provide a correctly formatted Billing Scope ID for a Microsoft Customer Agreement Account.
---

# Data Source: azurerm_billing_mca_account_scope

Use this data source to access an ID for your MCA Account billing scope.

## Example Usage

```hcl
data "azurerm_billing_mca_account_scope" "example" {
  billing_account_name = "e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31"
  billing_profile_name = "PE2Q-NOIT-BG7-TGB"
  invoice_section_name = "MTT4-OBS7-PJA-TGB"
}

output "id" {
  value = data.azurerm_billing_mca_account_scope.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `billing_account_name` - (Required) The Billing Account Name of the MCA account.

* `billing_profile_name` - (Required) The Billing Profile Name in the above Billing Account.

* `invoice_section_name` - (Required) The Invoice Section Name in the above Billing Profile.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Billing Scope.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Billing Scope.
