---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription"
description: |-
  Manages a Subscription.
---

# azurerm_subscription

Manages an Alias for a Subscription - which adds an Alias to an existing Subscription, allowing it to be managed in Terraform - or create a new Subscription with a new Alias.

~> **NOTE:** Destroying a Subscription controlled by this resource will place the Subscription into a cancelled state. It is possible to re-activate a subscription within 90-days of cancellation, after which time the Subscription is irrevocably deleted, and the Subscription ID cannot be re-used. For further information see [here](https://docs.microsoft.com/en-us/azure/cost-management-billing/manage/cancel-azure-subscription#what-happens-after-subscription-cancellation). Users can optionally delete a Subscription once 72 hours have passed, however, this functionality is not suitable for Terraform. A `Deleted` subscription cannot be reactivated.

~> **NOTE:** It is not possible to destroy (cancel) a subscription if it contains resources. If resources are present that are not managed by Terraform then these will need to be removed before the Subscription can be destroyed.

~> **NOTE:** Azure supports Multiple Aliases per Subscription, however, to reliably manage this resource in Terraform only a single Alias is supported.

## Example Usage - creating a new Alias and Subscription for an Enrollment Account

```hcl
data "azurerm_billing_enrollment_account_scope" "example" {
  billing_account_name    = "1234567890"
  enrollment_account_name = "0123456"
}

resource "azurerm_subscription" "example" {
  subscription_name = "My Example EA Subscription"
  billing_scope_id  = data.azurerm_billing_enrollment_account_scope.example.id
}
```

## Example Usage - creating a new Alias and Subscription for a Microsoft Customer Account

```hcl
data "azurerm_billing_mca_account_scope" "example" {
  billing_account_name = "e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31"
  billing_profile_name = "PE2Q-NOIT-BG7-TGB"
  invoice_section_name = "MTT4-OBS7-PJA-TGB"
}

resource "azurerm_subscription" "example" {
  subscription_name = "My Example MCA Subscription"
  billing_scope_id  = data.azurerm_billing_mca_account_scope.example.id
}
```

## Example Usage - adding an Alias to an existing Subscription

```hcl
resource "azurerm_subscription" "example" {
  alias             = "examplesub"
  subscription_name = "My Example Subscription"
  subscription_id   = "12345678-12234-5678-9012-123456789012"
}
```

## Arguments Reference

The following arguments are supported:

* `subscription_name` - (Required) The Name of the Subscription. This is the Display Name in the portal.

---

* `alias` - (Optional) The Alias name for the subscription. Terraform will generate a new GUID if this is not supplied. Changing this forces a new Subscription to be created.

* `billing_scope_id` - (Optional) The Azure Billing Scope ID. Can be either a Microsoft Customer Account Billing Scope ID or an Enrollment Billing Scope ID.

* `subscription_id` - (Optional) The ID of the Subscription. Changing this forces a new Subscription to be created.
 
~> **NOTE:** This value can be specified only for adopting control of an existing Subscription, it cannot be used to provide a custom Subscription ID.

~> **NOTE:** Either `billing_scope_id` or `subscription_id` has to be specified.

* `workload` - (Optional) The workload type of the Subscription.  Possible values are `Production` (default) and `DevTest`. Changing this forces a new Subscription to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The Resource ID of the Alias.

* `tenant_id` - The ID of the Tenant to which the subscription belongs.

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
In this scenario, the `subscription_id` property can be completed and Terraform will assume control of the existing subscription by creating an Alias. See the `adding an Alias to an existing Subscription` above. Terrafom requires an alias to correctly manage Subscription resources due to Azure Subscription API design.
