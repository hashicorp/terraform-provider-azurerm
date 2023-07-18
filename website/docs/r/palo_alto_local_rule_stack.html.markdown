
---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_local_rule_stack"
description: |-
  Manages a Palo Alto Local Rulestack.
---

# azurerm_palo_alto_local_rule_stack

Manages a Palo Alto Local Rulestack.

## Example Usage

```hcl
resource "azurerm_palo_alto_local_rule_stack" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Palo Alto Local Rulestack. Changing this forces a new Palo Alto Local Rulestack to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Palo Alto Local Rulestack should exist. Changing this forces a new Palo Alto Local Rulestack to be created.

* `location` - (Required) The Azure Region where the Palo Alto Local Rulestack should exist. Changing this forces a new Palo Alto Local Rulestack to be created.

---

* `description` - (Optional) The description for the Palo Alto Local Rulestack.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Local Rulestack.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Palo Alto Local Rulestack.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Local Rulestack.
* `update` - (Defaults to 30 minutes) Used when updating the Palo Alto Local Rulestack.
* `delete` - (Defaults to 30 minutes) Used when deleting the Palo Alto Local Rulestack.

## Import

Palo Alto Local Rulestacks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_local_rule_stack.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/myLocalRulestack
```