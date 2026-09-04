---
subcategory: "Billing"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_billing_invoice_section"
description: |-
  Gets information about an existing Invoice Section.
---

# Data Source: azurerm_billing_invoice_section

Use this data source to access information about an existing Invoice Section within a Billing Profile.

~> **Note:** This data source is only supported for Billing Accounts with an agreement type of `Microsoft Customer Agreement`.

## Example Usage

```hcl
data "azurerm_billing_invoice_section" "example" {
  name               = "invoice-section-1"
  billing_profile_id = "/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000:00000000-0000-0000-0000-000000000000_2019-05-31/billingProfiles/AAAA-BBBB-CCC-DDD"
}

output "invoice_section_display_name" {
  value = data.azurerm_billing_invoice_section.example.display_name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Invoice Section.

* `billing_profile_id` - (Required) The ID of the Billing Profile in which this Invoice Section exists, in the format `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Invoice Section. This is a Microsoft Customer Account Billing Scope ID, so it can be passed directly to the `billing_scope_id` property of `azurerm_subscription`.

* `display_name` - The display name of the Invoice Section.

* `system_id` - The system generated unique identifier for the Invoice Section.

* `state` - The status of the Invoice Section.

* `reason_code` - The reason for the current `state` of the Invoice Section.

* `target_cloud` - The cloud environments which are associated with the Invoice Section.

* `tags` - A mapping of tags assigned to the Invoice Section.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Invoice Section.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Billing` - 2024-04-01
