---
subcategory: "Billing"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_billing_invoice_section"
description: |-
  Manages an Invoice Section within a Billing Profile.
---

# azurerm_billing_invoice_section

Manages an Invoice Section within a Billing Profile.

~> **Note:** This resource is only supported for Billing Accounts with an agreement type of `Microsoft Customer Agreement`. The Billing Account and the Billing Profile must already exist, as neither can be created by Terraform.

## Example Usage

```hcl
resource "azurerm_billing_invoice_section" "example" {
  name               = "invoice-section-1"
  billing_profile_id = "/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000:00000000-0000-0000-0000-000000000000_2019-05-31/billingProfiles/AAAA-BBBB-CCC-DDD"
  display_name       = "Invoice Section 1"

  tags = {
    costCategory = "Support"
  }
}
```

## Example Usage - creating a Subscription within the Invoice Section

The ID of an Invoice Section is a Microsoft Customer Account Billing Scope ID, so it can be passed directly to [`azurerm_subscription`](subscription.html):

```hcl
resource "azurerm_billing_invoice_section" "example" {
  name               = "invoice-section-1"
  billing_profile_id = "/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000:00000000-0000-0000-0000-000000000000_2019-05-31/billingProfiles/AAAA-BBBB-CCC-DDD"
  display_name       = "Invoice Section 1"
}

resource "azurerm_subscription" "example" {
  subscription_name = "My Example MCA Subscription"
  billing_scope_id  = azurerm_billing_invoice_section.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Invoice Section. Must be between 1 and 128 characters in length and may only contain letters, numbers, hyphens and underscores. Changing this forces a new Invoice Section to be created.

* `billing_profile_id` - (Required) The ID of the Billing Profile in which this Invoice Section should exist, in the format `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}`. Changing this forces a new Invoice Section to be created.

* `display_name` - (Required) The display name of the Invoice Section.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Invoice Section. Keys and values are limited to 256 characters and keys may not contain any of `<`, `>`, `%`, `&`, `\`, `?` or `/`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Invoice Section.

* `system_id` - The system generated unique identifier for the Invoice Section.

* `state` - The status of the Invoice Section. Possible values are `Other`, `Active`, `Deleted`, `Disabled`, `UnderReview`, `Warned` and `Restricted`.

* `reason_code` - The reason for the current `state` of the Invoice Section. Possible values are `Other`, `PastDue`, `UnusualActivity`, `SpendingLimitReached` and `SpendingLimitExpired`.

* `target_cloud` - The cloud environments which are associated with the Invoice Section. This is a system managed field which is updated as the Invoice Section becomes associated with accounts in various clouds.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Invoice Section.
* `read` - (Defaults to 5 minutes) Used when retrieving the Invoice Section.
* `update` - (Defaults to 30 minutes) Used when updating the Invoice Section.
* `delete` - (Defaults to 30 minutes) Used when deleting the Invoice Section.

~> **Note:** An Invoice Section cannot be deleted whilst Subscriptions or Products are still associated with it. Terraform checks this before attempting the deletion and returns the reason reported by the Billing API.

## Import

Invoice Sections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_billing_invoice_section.example /providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000:00000000-0000-0000-0000-000000000000_2019-05-31/billingProfiles/AAAA-BBBB-CCC-DDD/invoiceSections/invoice-section-1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Billing` - 2024-04-01
