---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_firewall_policy"
sidebar_current: "docs-azurerm-resource-front-door-firewall-policy"
description: |-
  Manages an Azure Front Door Web Application Firewall Policy instance.
---

# azurerm_web_application_firewall_policy

Manages an Azure Front Door Web Application Firewall Policy instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West US 2"
}

resource "azurerm_frontdoor_firewall_policy" "example" {
  name                              = "example-fdwafpolicy"
  resource_group_name               = "${azurerm_resource_group.example.name}"
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  custom_rule {
    name                           = "Rule1"
    enabled                        = true
    priority                       = 1
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    rule_type                      = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_value        = ["192.168.1.0/24", "10.0.0.0/24"]
      transforms         = ["Lowercase", "Trim"]
    }
  }

  custom_rules {
    name                           = "Rule2"
    enabled                        = true
    priority                       = 2
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    rule_type                      = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_value        = ["192.168.1.0/24"]
    }

    match_condition {
      variable_name      = "RequestHeaders"
      selector           = "UserAgent"
      operator           = "Contains"
      negation_condition = false
      match_value        = ["windows"]
      transforms         = ["Lowercase", "Trim"]
    }
  }

  managed_rule {
    type    = "DefaultRuleSet"
    version = "preview-0.1"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933111"
        enable  = false
        action  = "Block"
      }
    }
  }

  managed_rule {
    type      = "BotProtection"
    version   = "preview-0.1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `enabled` - (Optional) Describes if the policy is in enabled state or disabled state Defaults to `Enabled`.

* `mode` - (Optional) Describes if it is in detection mode  or prevention mode at the policy level Defaults to `Prevention`.

* `redirect_url` - (Optional) If action type is redirect, this field represents redirect URL for the client.

* `custom_block_response_status_code` - (Optional) If the `action` type is block, customer can override the response status code. Valid values for the `custom_block_response_status_code` are `200`, `403`, `405`, `406`, or `429`.

* `custom_block_response_body` - (Optional) If the `action` type is block, customer can override the response body. The body must be specified in base64 encoding.

* `custom_rule` - (Optional) One or more `custom_rule` blocks as defined below.

* `managed_rule` - (Optional) One or more `managed_rule` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Web Application Firewall Policy.

---

The `custom_rule` block supports the following:

* `name` - (Required) Gets name of the resource that is unique within a policy. This name can be used to access the resource.

* `priority` - (Required) Describes priority of the rule. Rules with a lower value will be evaluated before rules with a higher value

* `rule_type` - (Required) Describes the type of rule

* `match_condition` - (Required) One or more `match_condition` block defined below.

* `action` - (Required) Type of Actions

---

The `managed_rule` block supports the following:

* `type` - (Required) The name of the managed rule to use with this resource.

* `version` - (Required) The version on the managed rule to use with this resource.

* `override` - (Optional) One or more `override` blocks as defined below.

---

The `match_condition` block supports the following:

* `variable_name` - (Required) The name of the Match Variable

* `selector` - (Optional) Match against a specific key from the `QueryString`, `PostArgs`, `RequestHeader` or `Cookies` variables.

* `match_variable` - (Required) The request variable to compare with.

* `operator` - (Required) Comparison type to use for matching with the variable value.

* `negation_condition` - (Optional) Describes if the result of this condition should be negated.

* `match_values` - (Required) A list of possible values to match.

---

The `override` block supports the following:

* `rule_group_name` - (Required) Describes the managed rule group to override.

* `rule` - (Optional) One or more `rule` blocks as defined below. If none are specified, all of the rules in the group will be disabled.

---

The `rule` block supports the following:

* `rule_id` - (Required) Identifier for the managed rule.

* `action` - (Required) Describes the override action to be applied when the rule matches.

* `enabled` - (Optional) Describes if the managed rule is in enabled or disabled state. Defaults to Disabled if not specified.

## Attributes Reference

The following attributes are exported:

* `id` - Resource ID.

* `location` - Resource location.

* `frontend_endpoint_ids` - Describes Frontend Endpoints associated with this Front Door Web Application Firewall policy.

## Import

Front Door Web Application Firewall Policy can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_frontdoor_firewall_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/frontdoorwebapplicationfirewallpolicies/example-fdwafpolicy
```
