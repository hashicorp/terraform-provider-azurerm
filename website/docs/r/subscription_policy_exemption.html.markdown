---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription_policy_exemption"
description: |-
  Manages a Subscription Policy Exemption.
---

# azurerm_subscription_policy_exemption

Manages a Subscription Policy Exemption.

## Example Usage

```hcl
data "azurerm_subscription" "example" {}

data "azurerm_policy_set_definition" "example" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_subscription_policy_assignment" "example" {
  name                 = "exampleAssignment"
  subscription_id      = data.azurerm_subscription.example.id
  policy_definition_id = data.azurerm_policy_set_definition.example.id
  location             = "westus"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_subscription_policy_exemption" "example" {
  name                 = "exampleExemption"
  subscription_id      = data.azurerm_subscription.example.id
  policy_assignment_id = azurerm_subscription_policy_assignment.example.id
  exemption_category   = "Mitigated"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Policy Exemption. Changing this forces a new resource to be created.

* `subscription_id` - (Required) The Subscription ID where the Policy Exemption should be applied. Changing this forces a new resource to be created.

* `exemption_category` - (Required) The category of this policy exemption. Possible values are `Waiver` and `Mitigated`.

* `policy_assignment_id` - (Required) The ID of the Policy Assignment to be exempted at the specified Scope. Changing this forces a new resource to be created.

* `description` - (Optional) A description to use for this Policy Exemption.

* `display_name` - (Optional) A friendly display name to use for this Policy Exemption.

* `expires_on` - (Optional) The expiration date and time in UTC ISO 8601 format of this policy exemption.

* `policy_definition_reference_ids` - (Optional) The policy definition reference ID list when the associated policy assignment is an assignment of a policy set definition.

* `metadata` - (Optional) The metadata for this policy exemption. This is a JSON string representing additional metadata that should be stored with the policy exemption.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Policy Exemption id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Exemption.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Exemption.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Exemption.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Exemption.

## Import

Policy Exemptions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subscription_policy_exemption.exemption1 /subscriptions/00000000-0000-0000-000000000000/providers/Microsoft.Authorization/policyExemptions/exemption1
```
