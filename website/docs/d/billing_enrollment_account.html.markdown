---
subcategory: "Billing"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_billing_enrollment_account"
description: |-
  Use this data source to access information about an Enterprise Account Enrollment.
---

# Data Source: azurerm_billing_enrollment_account

Use this data source to access information about an existing Enterprise Billing Enrollment Account.

!> **Note:** In order to use the `azurerm_billing_enrollment_account` data source your service principal may need an additional Enterprise Agreement role assigned. Please see the Microsoft documentation on [`Managing Azure Enterprise Agreement roles`](https://learn.microsoft.com/en-us/azure/cost-management-billing/manage/understand-ea-roles) for more information.

## Example Usage

```hcl
data "azurerm_billing_account" "example" {
  name = "12345678"
}

data "azurerm_billing_enrollment_account" "example" {
  billing_account_name    = data.azurerm_billing_account.example.name
  enrollment_account_name = "123456"
}

output "id" {
  value = data.azurerm_billing_enrollment_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `billing_account_name` - (Required) The Billing Account Name of the Enterprise Account. Note that in the Azure Portal this is actually referred to as the "Billing Account ID".

* `enrollment_account_name` - (Required) The Enrollment Account Name in the above Enterprise Account. Note that in the Azure Portal this is actually referred to as the "Enrollment Account ID".

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Enrollment Account.

* `account_name` - The name of the enrollment account.

* `account_owner` - The owner of the enrollment account.

* `cost_center` - The cost center associated with the enrollment account.

* `end_date` - The end date of the enrollment account.

* `start_date` - The start date of the enrollment account.

* `status` - The status of the enrollment account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Enrollment Account.
