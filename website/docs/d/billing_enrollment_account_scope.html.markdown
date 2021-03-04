---
subcategory: "Billing"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_billing_enrollment_account_scope"
description: |-
  This is a helper Data Source to provide a correctly formatted Billing Scope ID for an Enterprise Account Enrollment.
---

# Data Source: azurerm_billing_enrollment_account_scope

Use this data source to access information about an existing Enrollment Account Billing Scope.

## Example Usage

```hcl
data "azurerm_billing_enrollment_account_scope" "example" {
  billing_account_name    = "existing"
  enrollment_account_name = "existing"
}

output "id" {
  value = data.azurerm_billing_enrollment_account_scope.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `billing_account_name` - (Required) The Billing Account Name of the Enterprise Account.

* `enrollment_account_name` - (Required) The Enrollment Account Name in the above Enterprise Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Enrollment Account Billing Scope.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Enrollment Account Billing Scope.
