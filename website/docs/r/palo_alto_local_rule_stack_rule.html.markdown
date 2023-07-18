---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_local_rule_stack_rule"
description: |-
  Manages a Palo Alto Local Rulestack Rule.
---

# azurerm_palo_alto_local_rule_stack_rule

Manages a Palo Alto Local Rulestack Rule.

## Example Usage

```hcl
resource "azurerm_palo_alto_local_rule_stack_rule" "example" {
  name = "example"

  source {
    
  }
  priority = 42
  rule_stack_id = "TODO"

  destination {
    
  }
  applications = [ "example" ]
}
```

## Arguments Reference

The following arguments are supported:

* `applications` - (Required) Specifies a list of TODO.

* `destination` - (Required) One or more `destination` blocks as defined below.

* `name` - (Required) The name which should be used for this Palo Alto Local Rulestack Rule.

* `priority` - (Required) TODO. Changing this forces a new Palo Alto Local Rulestack Rule to be created.

* `rule_stack_id` - (Required) The ID of the TODO. Changing this forces a new Palo Alto Local Rulestack Rule to be created.

* `source` - (Required) One or more `source` blocks as defined below.

---

* `action` - (Optional) TODO. Defaults to `Allow`.

* `audit_comment` - (Optional) TODO.

* `category` - (Optional) A `category` block as defined below.

* `decryption_rule_type` - (Optional) TODO. Defaults to `None`.

* `description` - (Optional) TODO.

* `enabled` - (Optional) Should the TODO be enabled? Defaults to `true`.

* `inspection_certificate_id` - (Optional) The ID of the TODO.

* `logging_enabled` - (Optional) Should the TODO be enabled? Defaults to `false`.

* `negate_destination` - (Optional) TODO. Defaults to `false`.

* `negate_source` - (Optional) TODO. Defaults to `false`.

* `protocol` - (Optional) TODO.Conflicts with `protocol_ports`. Defaults to `application-default`.

* `protocol_ports` - (Optional) Specifies a list of TODO.Conflicts with `protocol`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Palo Alto Local Rulestack Rule.

---

A `category` block supports the following:

* `feeds` - (Required) Specifies a list of TODO.

* `custom_urls` - (Optional) Specifies a list of TODO.

---

A `destination` block supports the following:

* `cidrs` - (Optional) Specifies a list of TODO.

* `countries` - (Optional) Specifies a list of TODO.

* `feeds` - (Optional) Specifies a list of TODO.

* `fqdn_lists` - (Optional) Specifies a list of TODO.

* `prefix_lists` - (Optional) Specifies a list of TODO.

---

A `source` block supports the following:

* `cidrs` - (Optional) Specifies a list of TODO.

* `countries` - (Optional) Specifies a list of TODO.

* `feeds` - (Optional) Specifies a list of TODO.

* `prefix_lists` - (Optional) Specifies a list of TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Local Rulestack Rule.

* `etag` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Palo Alto Local Rulestack Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Local Rulestack Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Palo Alto Local Rulestack Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Palo Alto Local Rulestack Rule.

## Import

Palo Alto Local Rulestack Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_local_rule_stack_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/myLocalRulestack/localRules/myRule1
```