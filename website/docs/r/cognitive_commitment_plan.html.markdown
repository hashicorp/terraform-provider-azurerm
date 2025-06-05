---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_commitment_plan"
description: |-
  Manages a Cognitive Commitment Plan associated with the cognitive services account.
---

# azurerm_cognitive_commitment_plan

Manages a Cognitive Commitment Plan associated with the cognitive services account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                = "example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}

resource "azurerm_cognitive_commitment_plan" "example" {
  name                          = "example-ccp"
  cognitive_services_account_id = azurerm_cognitive_account.example.id
  auto_renew_enabled            = false
  hosting_model                 = "Web"
  plan_type                     = "STT"
  current_tier                  = "T1"
  renewal_tier                  = "T2"

  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Cognitive Commitment Plan. Changing this forces a new resource to be created.

* `cognitive_account_id` - (Required) Specifies the ID of the Cognitive Commitment Plan. Changing this forces a new resource to be created.

* `current_tier` - (Required) Specifies the current commitment period commitment tier of the Cognitive Commitment Plan. Changing this forces a new resource to be created.

* `hosting_model` - (Required) Specifies the account hosting model of the Cognitive Commitment Plan. Possible values are `ConnectedContainer`, `DisconnectedContainer`, `ProvisionedWeb` and `Web`. Changing this forces a new resource to be created.

* `plan_type` - (Required) Specifies the type of the Cognitive Commitment Plan. Changing this forces a new resource to be created.

* `auto_renew_enabled` - (Optional) Whether auto-renewal commitment plan is enabled. Defaults to `false`.

* `renewal_tier` - (Optional) Specifies to renew or change the current tier starting with the next billing cycle for the Cognitive Commitment Plan.

* `tags` - (Optional) A mapping of tags to assign to the Cognitive Commitment Plan.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Commitment Plan.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Commitment Plan.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Commitment Plan.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Commitment Plan.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Commitment Plan.

## Import

Cognitive Commitment Plan can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_commitment_plan.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.CognitiveServices/accounts/account1/commitmentPlans/commitmentPlan1
```