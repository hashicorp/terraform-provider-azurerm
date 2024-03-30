---
subcategory: "Billing"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_billing_account"
description: |-
  Use this data source to access information about a Billing Account.
---

# Data Source: azurerm_billing_account

Use this data source to access information about an existing Billing Account.

## Example Usage

```hcl
data "azurerm_billing_account" "example" {
  name    = "12345678"
}

output "id" {
  value = data.azurerm_billing_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Billing Account Name. Note that in the Azure Portal this is actually referred to as the "Billing Account ID".

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Enrollment Account.

* `account_status` - The status of the billing account (`Active`, `Deleted`, `Disabled`, `Expired`, `Extended`, `Terminated`, `Transferred`).

* `account_type` - The type of customer (`Enterprise`, `Individual`, `Partner`).

* `agreement_type` - The type of agreement (`EnterpriseAgreement`, `MicrosoftCustomerAgreement`, `MicrosoftOnlineServicesProgram`, `MicrosoftPartnerAgreement`).

* `display_name` - The billing account's display name.

* `has_read_access` - Indicates whether user has read access to the billing account.

* `sold_to` - A `sold_to` block as defined below.

---

A `sold_to` block exports the following. Note that most of these fields are optional and may be empty:

* `address_line_1`

* `address_line_2`

* `address_line_3`

* `city`

* `company_name`

* `country` - Country code uses ISO2, 2-digit format.

* `district`

* `email`

* `first_name`

* `last_name`

* `middle_name`

* `phone_number`

* `postal_code`

* `region`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Enrollment Account.
