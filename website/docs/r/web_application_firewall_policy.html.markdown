---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_application_firewall_policy"
sidebar_current: "docs-azurerm-resource-web-application-firewall-policy"
description: |-
  Manage a Azure WebApplicationFirewallPolicy instance.
---

# azurerm_web_application_firewall_policy

Manage a Azure WebApplicationFirewallPolicy instance.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West US 2"
}

resource "azurerm_web_application_firewall_policy" "example" {
  name                = "example-wafpolicy"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  custom_rules {
    name      = "Rule1"
    priority  = 1
    rule_type = "MatchRule"

    match_conditions {
      match_variables {
        variable_name = "RemoteAddr"
      }

      operator          = "IPMatch"
      negation_conditon = false
      match_values      = ["192.168.1.0/24", "10.0.0.0/24"]
    }

    action = "Block"
  }

  custom_rules {
    name      = "Rule2"
    priority  = 2
    rule_type = "MatchRule"

    match_conditions {
      match_variables {
        variable_name = "RemoteAddr"
      }

      operator          = "IPMatch"
      negation_conditon = false
      match_values      = ["192.168.1.0/24"]
    }

    match_conditions {
      match_variables {
        variable_name = "RequestHeaders"
        selector      = "UserAgent"
      }

      operator          = "Contains"
      negation_conditon = false
      match_values      = ["Windows"]
    }

    action = "Block"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `custom_rules` - (Optional) One or more `custom_rule` block defined below.

* `policy_settings` - (Optional) One `policy_setting` block defined below.

* `tags` - (Optional) Resource tags. Changing this forces a new resource to be created.

---

The `custom_rule` block supports the following:

* `name` - (Optional) Gets name of the resource that is unique within a policy. This name can be used to access the resource.

* `priority` - (Required) Describes priority of the rule. Rules with a lower value will be evaluated before rules with a higher value

* `rule_type` - (Required) Describes type of rule

* `match_conditions` - (Required) One or more `match_condition` block defined below.

* `action` - (Required) Type of Actions


---

The `match_condition` block supports the following:

* `match_variables` - (Required) One or more `match_variable` block defined below.

* `operator` - (Required) Describes operator to be matched

* `negation_conditon` - (Optional) Describes if this is negate condition or not

* `match_values` - (Required) Match value


---

The `match_variable` block supports the following:

* `variable_name` - (Required) Match Variable

* `selector` - (Optional) Describes field of the matchVariable collection

---

The `policy_setting` block supports the following:

* `enabled` - (Optional) Describes if the policy is in enabled state or disabled state Defaults to `Enabled`.

* `mode` - (Optional) Describes if it is in detection mode  or prevention mode at policy level Defaults to `Prevention`.

## Attributes Reference

The following attributes are exported:

* `id` - Resource ID.


## Import

Web Application Firewall Policy can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_web_application_firewall_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/ApplicationGatewayWebApplicationFirewallPolicies/example-wafpolicy
```
