---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_firewall_policy"
description: |-
  Manages an Azure Front Door Web Application Firewall Policy instance.
---

# azurerm_frontdoor_firewall_policy

Manages an Azure Front Door Web Application Firewall Policy instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West US 2"
}

resource "azurerm_frontdoor_firewall_policy" "example" {
  name                              = "example-fdwafpolicy"
  resource_group_name               = azurerm_resource_group.example.name
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
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24", "10.0.0.0/24"]
    }
  }

  custom_rule {
    name                           = "Rule2"
    enabled                        = true
    priority                       = 2
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24"]
    }

    match_condition {
      match_variable     = "RequestHeader"
      selector           = "UserAgent"
      operator           = "Contains"
      negation_condition = false
      match_values       = ["windows"]
      transforms         = ["Lowercase", "Trim"]
    }
  }

  managed_rule {
    type    = "DefaultRuleSet"
    version = "1.0"

    exclusion {
      match_variable = "QueryStringArgNames"
      operator       = "Equals"
      selector       = "not_suspicious"
    }

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "Block"
      }
    }

    override {
      rule_group_name = "SQLI"

      exclusion {
        match_variable = "QueryStringArgNames"
        operator       = "Equals"
        selector       = "really_not_suspicious"
      }

      rule {
        rule_id = "942200"
        action  = "Block"

        exclusion {
          match_variable = "QueryStringArgNames"
          operator       = "Equals"
          selector       = "innocent"
        }
      }
    }
  }

  managed_rule {
    type    = "Microsoft_BotManagerRuleSet"
    version = "1.0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `enabled` - (Optional) Is the policy a enabled state or disabled state. Defaults to `true`.

* `mode` - (Optional) The firewall policy mode. Possible values are `Detection`, `Prevention` and defaults to `Prevention`.

* `redirect_url` - (Optional) If action type is redirect, this field represents redirect URL for the client.

* `custom_rule` - (Optional) One or more `custom_rule` blocks as defined below.

* `custom_block_response_status_code` - (Optional) If a `custom_rule` block's action type is `block`, this is the response status code. Possible values are `200`, `403`, `405`, `406`, or `429`.

* `custom_block_response_body` - (Optional) If a `custom_rule` block's action type is `block`, this is the response body. The body must be specified in base64 encoding.

* `managed_rule` - (Optional) One or more `managed_rule` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Web Application Firewall Policy.

---

The `custom_rule` block supports the following:

* `name` - (Required) Gets name of the resource that is unique within a policy. This name can be used to access the resource.

* `action` - (Required) The action to perform when the rule is matched. Possible values are `Allow`, `Block`, `Log`, or `Redirect`.

* `enabled` - (Optional) Is the rule is enabled or disabled? Defaults to `true`.

* `priority` - (Required) The priority of the rule. Rules with a lower value will be evaluated before rules with a higher value. Defaults to `1`.

* `type` - (Required) The type of rule. Possible values are `MatchRule` or `RateLimitRule`.

* `match_condition` - (Required) One or more `match_condition` block defined below.

* `rate_limit_duration_in_minutes` - (Optional) The rate limit duration in minutes. Defaults to `1`.

* `rate_limit_threshold` - (Optional) The rate limit threshold. Defaults to `10`.

---

The `match_condition` block supports the following:

* `match_variable` - (Required) The request variable to compare with. Possible values are `Cookies`, `PostArgs`, `QueryString`, `RemoteAddr`, `RequestBody`, `RequestHeader`, `RequestMethod`, `RequestUri`, or `SocketAddr`.

* `match_values` - (Required) Up to `100` possible values to match.

* `operator` - (Required) Comparison type to use for matching with the variable value. Possible values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GeoMatch`, `GreaterThan`, `GreaterThanOrEqual`, `IPMatch`, `LessThan`, `LessThanOrEqual` or `RegEx`.

* `selector` - (Optional) Match against a specific key if the `match_variable` is `QueryString`, `PostArgs`, `RequestHeader` or `Cookies`.

* `negation_condition` - (Optional) Should the result of the condition be negated.

* `transforms` - (Optional) Up to `5` transforms to apply. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `URLDecode` or`URLEncode`.

---

The `managed_rule` block supports the following:

* `type` - (Required) The name of the managed rule to use with this resource.

* `version` - (Required) The version on the managed rule to use with this resource.

* `exclusion` - (Optional) One or more `exclusion` blocks as defined below.

* `override` - (Optional) One or more `override` blocks as defined below.

---

The `override` block supports the following:

* `rule_group_name` - (Required) The managed rule group to override.

* `exclusion` - (Optional) One or more `exclusion` blocks as defined below.

* `rule` - (Optional) One or more `rule` blocks as defined below. If none are specified, all of the rules in the group will be disabled.

---

The `rule` block supports the following:

* `rule_id` - (Required) Identifier for the managed rule.

* `action` - (Required) The action to be applied when the rule matches. Possible values are `Allow`, `Block`, `Log`, or `Redirect`.

* `enabled` - (Optional) Is the managed rule override enabled or disabled. Defaults to `false`

* `exclusion` - (Optional) One or more `exclusion` blocks as defined below.

---

The `exclusion` block supports the following:

* `match_variable` - (Required) The variable type to be excluded. Possible values are `QueryStringArgNames`, `RequestBodyPostArgNames`, `RequestCookieNames`, `RequestHeaderNames`.

* `operator` - (Required) Comparison operator to apply to the selector when specifying which elements in the collection this exclusion applies to. Possible values are: `Equals`, `Contains`, `StartsWith`, `EndsWith`, `EqualsAny`.

* `selector` - (Required) Selector for the value in the `match_variable` attribute this exclusion applies to.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the FrontDoor Firewall Policy.

* `location` - The Azure Region where this FrontDoor Firewall Policy exists.

* `frontend_endpoint_ids` - The Frontend Endpoints associated with this Front Door Web Application Firewall policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the FrontDoor Web Application Firewall Policy.
* `update` - (Defaults to 30 minutes) Used when updating the FrontDoor Web Application Firewall Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the FrontDoor Web Application Firewall Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the FrontDoor Web Application Firewall Policy.

## Import

FrontDoor Web Application Firewall Policy can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_frontdoor_firewall_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/example-fdwafpolicy
```
