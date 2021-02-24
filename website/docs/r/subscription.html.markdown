---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription"
description: |-
  Manages a Subscription.
---

# azurerm_subscription

Manages a Subscription.

~> **NOTE:** Destroying a Subscription controlled by this resource will place the Subscription into a cancelled state. It is possible to re-activate a subscription within 90-days of cancellation, after which time the Subscription is irrevocably deleted and the Subscription ID cannot be re-used. For further information see [here](https://docs.microsoft.com/en-us/azure/cost-management-billing/manage/cancel-azure-subscription#what-happens-after-subscription-cancellation)

## Example Usage

```hcl
resource "azurerm_subscription" "example" {
  alias              = "examplesub"
  subscription_name  = "My Example Subscription"
  billing_account    = "1234567890"
  enrollment_account = "0123456"
}
```

## Arguments Reference

The following arguments are supported:

* `subscription_name` - (Required) The Name of the Subscription. This is the Display Name in the portal.

---

* `alias` - (Optional) The Alias name for the subscription. Terraform will generate a new GUID if this is not supplied. Changing this forces a new Subscription to be created.

* `billing_account` - (Optional) The Azure Billing Account Name. Changing this forces a new Subscription to be created.

* `billing_profile` - (Optional) The Billing Profile within the Billing Account. Conflicts with `enrollment_account`,`subscription_id`. Changing this forces a new Subscription to be created.

~> **NOTE:** This value is only used for MCA and Partner Account types and must be used with `invoice_section`.

* `enrollment_account` - (Optional) The Enrollment Account Name. Conflicts with `invoice_section`,`billing_profile`, and `subscription_id`. Changing this forces a new Subscription to be created.

~> **NOTE:** This value is only valid for Enterprise Agreements.

* `invoice_section` - (Optional) The Invoice Section name. Conflicts with `enrollment_account`, and `subscription_id`. Changing this forces a new Subscription to be created.

~> **NOTE:** This value is only valid for MCA and Partner Account types and must be used with `billing_profile`.

* `state` - (Optional) The target state of the subscription. Possible values are `Active` (default) and `Cancelled`.

* `subscription_id` - (Optional) The ID of the Subscription. Cannot be specified with `billing_account`, `billing_profile`, `enrollment_account`, or `invoice_section` Changing this forces a new Subscription to be created.

~> **NOTE:** This value can be specified only for adopting control of an existing Subscription, it cannot be used to provide a custom Subscription ID.

* `workload` - (Optional) The workload type of the Subscription.  Possible values are `Production` (default) and `DevTest`. Changing this forces a new Subscription to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Subscription.

* `tenant` - The ID of the Tenant to which the subscription belongs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the Subscription.
* `update` - (Defaults to 30 minutes) Used when updating the Subscription.
* `delete` - (Defaults to 30 minutes) Used when deleting the Subscription.

## Import

Subscriptions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subscription.example "/providers/Microsoft.Subscription/aliases/subscription1"
```

!> **NOTE:** When importing a Subscription that was not created programmatically, either by this Terraform resource or using the Alias API, it will have no Alias ID to import via `terraform import`.  
In this scenario, the `subscription_id` property can be completed and Terraform will assume control of the existing subscription by creating an Alias. e.g.

```hcl
// importing existing Manually Created Subscription with no Alias
resource "azurerm_subscription" "example" {
  alias             = "examplesub"
  subscription_name = "My Example Subscription"
  subscription_id   = "12345678-12234-5678-9012-123456789012"
}
```
