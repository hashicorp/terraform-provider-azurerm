---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_application_firewall_policy"
description: |-
  Manages a Azure Web Application Firewall Policy instance.
---

# azurerm_web_application_firewall_policy

Manages a Azure Web Application Firewall Policy instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West US 2"
}

resource "azurerm_web_application_firewall_policy" "example" {
  name                = "example-wafpolicy"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  custom_rules {
    name      = "Rule1"
    priority  = 1
    rule_type = "MatchRule"

    match_conditions {
      match_variables {
        variable_name = "RemoteAddr"
      }

      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24", "10.0.0.0/24"]
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

      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24"]
    }

    match_conditions {
      match_variables {
        variable_name = "RequestHeaders"
        selector      = "UserAgent"
      }

      operator           = "Contains"
      negation_condition = false
      match_values       = ["Windows"]
    }

    action = "Block"
  }

  policy_settings {
    enabled                     = true
    mode                        = "Prevention"
    request_body_check          = true
    file_upload_limit_in_mb     = 100
    max_request_body_size_in_kb = 128
  }

  managed_rules {
    exclusion {
      match_variable          = "RequestHeaderNames"
      selector                = "x-company-secret-header"
      selector_match_operator = "Equals"
    }
    exclusion {
      match_variable          = "RequestCookieNames"
      selector                = "too-tasty"
      selector_match_operator = "EndsWith"
    }

    managed_rule_set {
      type    = "OWASP"
      version = "3.1"
      rule_group_override {
        rule_group_name = "REQUEST-920-PROTOCOL-ENFORCEMENT"
        disabled_rules = [
          "920300",
          "920440"
        ]
      }
    }
  }

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `custom_rules` - (Optional) One or more `custom_rules` blocks as defined below.

* `policy_settings` - (Optional) A `policy_settings` block as defined below.

* `managed_rules` - (Required) A `managed_rules` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Web Application Firewall Policy.

---

The `custom_rules` block supports the following:

* `name` - (Optional) Gets name of the resource that is unique within a policy. This name can be used to access the resource.

* `priority` - (Required) Describes priority of the rule. Rules with a lower value will be evaluated before rules with a higher value.

* `rule_type` - (Required) Describes the type of rule.

* `match_conditions` - (Required) One or more `match_conditions` blocks as defined below.

* `action` - (Required) Type of action.

---

The `match_conditions` block supports the following:

* `match_variables` - (Required) One or more `match_variables` blocks as defined below.

* `match_values` - (Required) A list of match values.

* `operator` - (Required) Describes operator to be matched.

* `negation_condition` - (Optional) Describes if this is negate condition or not

* `transforms` - (Optional) A list of transformations to do before the match is attempted.

---

The `match_variables` block supports the following:

* `variable_name` - (Required) The name of the Match Variable

* `selector` - (Optional) Describes field of the matchVariable collection

---

The `policy_settings` block supports the following:

* `enabled` - (Optional) Describes if the policy is in enabled state or disabled state. Defaults to `Enabled`.

* `mode` - (Optional) Describes if it is in detection mode or prevention mode at the policy level. Defaults to `Prevention`.

* `file_upload_limit_mb` - (Optional) The File Upload Limit in MB. Accepted values are in the range `1` to `750`. Defaults to `100`.

* `request_body_check` - (Optional) Is Request Body Inspection enabled? Defaults to `true`.

* `max_request_body_size_kb` - (Optional) The Maximum Request Body Size in KB.  Accepted values are in the range `8` to `128`. Defaults to `128`.

---

The `managed_rules` block supports the following:

* `exclusion` - (Optional) One or more `exclusion` block defined below.

* `managed_rule_set` - (Optional) One or more `managed_rule_set` block defined below.

---

The `exclusion` block supports the following:

* `match_variable` - (Required) The name of the Match Variable. Possible values: `RequestArgNames`, `RequestCookieNames`, `RequestHeaderNames`.

* `selector` - (Optional) Describes field of the matchVariable collection.

* `selector_match_operator` - (Required) Describes operator to be matched. Possible values: `Contains`, `EndsWith`, `Equals`, `EqualsAny`, `StartsWith`.

---

The `managed_rule_set` block supports the following:

* `type` - (Optional) The rule set type. Possible values: `Microsoft_BotManagerRuleSet` and `OWASP`.

* `version` - (Required) The rule set version. Possible values: `0.1`, `1.0`, `2.2.9`, `3.0` and `3.1`.

* `rule_group_override` - (Optional) One or more `rule_group_override` block defined below.

---

The `rule_group_override` block supports the following:

* `rule_group_name` - (Required) The name of the Rule Group

* `disabled_rules` - (Optional) One or more Rule ID's

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Web Application Firewall Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web Application Firewall Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Web Application Firewall Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web Application Firewall Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Web Application Firewall Policy.

## Import

Web Application Firewall Policy can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_web_application_firewall_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/ApplicationGatewayWebApplicationFirewallPolicies/example-wafpolicy
```
