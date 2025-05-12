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
  location = "West Europe"
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
      version = "3.2"
      rule_group_override {
        rule_group_name = "REQUEST-920-PROTOCOL-ENFORCEMENT"
        rule {
          id      = "920300"
          enabled = true
          action  = "Log"
        }

        rule {
          id      = "920440"
          enabled = true
          action  = "Block"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `location` - (Required) Resource location. Changing this forces a new resource to be created.

* `custom_rules` - (Optional) One or more `custom_rules` blocks as defined below.

* `policy_settings` - (Optional) A `policy_settings` block as defined below.

* `managed_rules` - (Required) A `managed_rules` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Web Application Firewall Policy.

---

The `custom_rules` block supports the following:

* `enabled` - (Optional) Describes if the policy is in enabled state or disabled state. Defaults to `true`.

* `name` - (Optional) Gets name of the resource that is unique within a policy. This name can be used to access the resource.

* `priority` - (Required) Describes priority of the rule. Rules with a lower value will be evaluated before rules with a higher value.

* `rule_type` - (Required) Describes the type of rule. Possible values are `MatchRule`, `RateLimitRule` and `Invalid`.

* `match_conditions` - (Required) One or more `match_conditions` blocks as defined below.

* `action` - (Required) Type of action. Possible values are `Allow`, `Block` and `Log`.

~> **Note:** If the `rule_type` is specified as `RateLimitRule`, the `Allow` is not supported.

* `rate_limit_duration` - (Optional) Specifies the duration at which the rate limit policy will be applied. Should be used with `RateLimitRule` rule type. Possible values are `FiveMins` and `OneMin`.

* `rate_limit_threshold` - (Optional) Specifies the threshold value for the rate limit policy. Must be greater than or equal to 1 if provided.

* `group_rate_limit_by` - (Optional) Specifies what grouping the rate limit will count requests by. Possible values are `GeoLocation`, `ClientAddr` and `None`.

---

The `match_conditions` block supports the following:

* `match_variables` - (Required) One or more `match_variables` blocks as defined below.

* `match_values` - (Optional) A list of match values. This is **Required** when the `operator` is not `Any`.

* `operator` - (Required) Describes operator to be matched. Possible values are `Any`, `IPMatch`, `GeoMatch`, `Equal`, `Contains`, `LessThan`, `GreaterThan`, `LessThanOrEqual`, `GreaterThanOrEqual`, `BeginsWith`, `EndsWith` and `Regex`.

* `negation_condition` - (Optional) Describes if this is negate condition or not

* `transforms` - (Optional) A list of transformations to do before the match is attempted. Possible values are `HtmlEntityDecode`, `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` and `UrlEncode`.

---

The `match_variables` block supports the following:

* `variable_name` - (Required) The name of the Match Variable. Possible values are `RemoteAddr`, `RequestMethod`, `QueryString`, `PostArgs`, `RequestUri`, `RequestHeaders`, `RequestBody` and `RequestCookies`.

* `selector` - (Optional) Describes field of the matchVariable collection

---

The `policy_settings` block supports the following:

* `enabled` - (Optional) Describes if the policy is in enabled state or disabled state. Defaults to `true`.

* `mode` - (Optional) Describes if it is in detection mode or prevention mode at the policy level. Valid values are `Detection` and `Prevention`. Defaults to `Prevention`.

* `file_upload_limit_in_mb` - (Optional) The File Upload Limit in MB. Accepted values are in the range `1` to `4000`. Defaults to `100`.

* `request_body_check` - (Optional) Is Request Body Inspection enabled? Defaults to `true`.

* `max_request_body_size_in_kb` - (Optional) The Maximum Request Body Size in KB. Accepted values are in the range `8` to `2000`. Defaults to `128`.

* `log_scrubbing` - (Optional) One `log_scrubbing` block as defined below.

* `request_body_enforcement` - (Optional) Whether the firewall should block a request with body size greater then `max_request_body_size_in_kb`. Defaults to `true`.

* `request_body_inspect_limit_in_kb` - (Optional) Specifies the maximum request body inspection limit in KB for the Web Application Firewall. Defaults to `128`.

* `js_challenge_cookie_expiration_in_minutes` - (Optional) Specifies the JavaScript challenge cookie validity lifetime in minutes. The user is challenged after the lifetime expires. Accepted values are in the range `5` to `1440`. Defaults to `30`.

* `file_upload_enforcement` - (Optional)  - Whether the firewall should block a request with upload size greater then `file_upload_limit_in_mb`.

---

The `managed_rules` block supports the following:

* `exclusion` - (Optional) One or more `exclusion` block defined below.

* `managed_rule_set` - (Required) One or more `managed_rule_set` block defined below.

---

The `exclusion` block supports the following:

* `match_variable` - (Required) The name of the Match Variable. Possible values: `RequestArgKeys`, `RequestArgNames`, `RequestArgValues`, `RequestCookieKeys`, `RequestCookieNames`, `RequestCookieValues`, `RequestHeaderKeys`, `RequestHeaderNames`, `RequestHeaderValues`.

* `selector` - (Required) Describes field of the matchVariable collection.

* `selector_match_operator` - (Required) Describes operator to be matched. Possible values: `Contains`, `EndsWith`, `Equals`, `EqualsAny`, `StartsWith`.

* `excluded_rule_set` - (Optional) One or more `excluded_rule_set` block defined below.

---

The `excluded_rule_set` block supports the following:

* `type` - (Optional) The rule set type. Possible values are `Microsoft_DefaultRuleSet`, `Microsoft_BotManagerRuleSet` and `OWASP`. Defaults to `OWASP`.

* `version` - (Optional) The rule set version. Possible values are `1.0`, `1.1` (for rule set type `Microsoft_BotManagerRuleSet`), `2.1` (for rule set type `Microsoft_DefaultRuleSet`) and `3.2` (for rule set type `OWASP`). Defaults to `3.2`.

* `rule_group` - (Optional) One or more `rule_group` block defined below.

---

The `rule_group` block supports the following:

* `rule_group_name` - (Required) The name of rule group for exclusion. Possible values are `BadBots`, `crs_20_protocol_violations`, `crs_21_protocol_anomalies`, `crs_23_request_limits`, `crs_30_http_policy`, `crs_35_bad_robots`, `crs_40_generic_attacks`, `crs_41_sql_injection_attacks`, `crs_41_xss_attacks`, `crs_42_tight_security`, `crs_45_trojans`, `crs_49_inbound_blocking`, `General`, `GoodBots`, `KnownBadBots`, `Known-CVEs`, `REQUEST-911-METHOD-ENFORCEMENT`, `REQUEST-913-SCANNER-DETECTION`, `REQUEST-920-PROTOCOL-ENFORCEMENT`, `REQUEST-921-PROTOCOL-ATTACK`, `REQUEST-930-APPLICATION-ATTACK-LFI`, `REQUEST-931-APPLICATION-ATTACK-RFI`, `REQUEST-932-APPLICATION-ATTACK-RCE`, `REQUEST-933-APPLICATION-ATTACK-PHP`, `REQUEST-941-APPLICATION-ATTACK-XSS`, `REQUEST-942-APPLICATION-ATTACK-SQLI`, `REQUEST-943-APPLICATION-ATTACK-SESSION-FIXATION`, `REQUEST-944-APPLICATION-ATTACK-JAVA`, `UnknownBots`, `METHOD-ENFORCEMENT`, `PROTOCOL-ENFORCEMENT`, `PROTOCOL-ATTACK`, `LFI`, `RFI`, `RCE`, `PHP`, `NODEJS`, `XSS`, `SQLI`, `FIX`, `JAVA`, `MS-ThreatIntel-WebShells`, `MS-ThreatIntel-AppSec`, `MS-ThreatIntel-SQLI` and `MS-ThreatIntel-CVEs`.
	`MS-ThreatIntel-AppSec`, `MS-ThreatIntel-SQLI` and `MS-ThreatIntel-CVEs`.

* `excluded_rules` - (Optional) One or more Rule IDs for exclusion.

---

The `managed_rule_set` block supports the following:

* `type` - (Optional) The rule set type. Possible values: `Microsoft_BotManagerRuleSet`, `Microsoft_DefaultRuleSet` and `OWASP`. Defaults to `OWASP`.

* `version` - (Required) The rule set version. Possible values: `0.1`, `1.0`, `1.1`, `2.1`, `2.2.9`, `3.0`, `3.1` and `3.2`.

* `rule_group_override` - (Optional) One or more `rule_group_override` block defined below.

---

The `rule_group_override` block supports the following:

* `rule_group_name` - (Required) The name of the Rule Group. Possible values are `BadBots`, `crs_20_protocol_violations`, `crs_21_protocol_anomalies`, `crs_23_request_limits`, `crs_30_http_policy`, `crs_35_bad_robots`, `crs_40_generic_attacks`, `crs_41_sql_injection_attacks`, `crs_41_xss_attacks`, `crs_42_tight_security`, `crs_45_trojans`, `crs_49_inbound_blocking`, `General`, `GoodBots`, `KnownBadBots`, `Known-CVEs`, `REQUEST-911-METHOD-ENFORCEMENT`, `REQUEST-913-SCANNER-DETECTION`, `REQUEST-920-PROTOCOL-ENFORCEMENT`, `REQUEST-921-PROTOCOL-ATTACK`, `REQUEST-930-APPLICATION-ATTACK-LFI`, `REQUEST-931-APPLICATION-ATTACK-RFI`, `REQUEST-932-APPLICATION-ATTACK-RCE`, `REQUEST-933-APPLICATION-ATTACK-PHP`, `REQUEST-941-APPLICATION-ATTACK-XSS`, `REQUEST-942-APPLICATION-ATTACK-SQLI`, `REQUEST-943-APPLICATION-ATTACK-SESSION-FIXATION`, `REQUEST-944-APPLICATION-ATTACK-JAVA`, `UnknownBots`, `METHOD-ENFORCEMENT`, `PROTOCOL-ENFORCEMENT`, `PROTOCOL-ATTACK`, `LFI`, `RFI`, `RCE`, `PHP`, `NODEJS`, `XSS`, `SQLI`, `FIX`, `JAVA`, `MS-ThreatIntel-WebShells`, `MS-ThreatIntel-AppSec`, `MS-ThreatIntel-SQLI` and `MS-ThreatIntel-CVEs`MS-ThreatIntel-WebShells`,.

* `rule` - (Optional) One or more `rule` block defined below.

---

The `rule` block supports the following:

* `id` - (Required) Identifier for the managed rule.

* `enabled` - (Optional) Describes if the managed rule is in enabled state or disabled state. Defaults to `false`.

* `action` - (Optional) Describes the override action to be applied when rule matches. Possible values are `Allow`, `AnomalyScoring`, `Block`, `JSChallenge` and `Log`. `JSChallenge` is only valid for rulesets of type `Microsoft_BotManagerRuleSet`.

---

The `log_scrubbing` block supports the following:

* `enabled` - (Optional) Whether the log scrubbing is enabled or disabled. Defaults to `true`.

* `rule` - (Optional) One or more `scrubbing_rule` blocks as define below.

---

The `scrubbing_rule` block supports the following:

* `enabled` - (Optional) Whether this rule is enabled. Defaults to `true`.

* `match_variable` - (Required) Specifies the variable to be scrubbed from the logs. Possible values are `RequestHeaderNames`, `RequestCookieNames`, `RequestArgNames`, `RequestPostArgNames`, `RequestJSONArgNames` and `RequestIPAddress`.

* `selector_match_operator` - (Optional) Specifies the operating on the `selector`. Possible values are `Equals` and `EqualsAny`. Defaults to `Equals`.

* `selector` - (Optional) Specifies which elements in the collection this rule applies to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Web Application Firewall Policy.

* `http_listener_ids` - A list of HTTP Listener IDs from an `azurerm_application_gateway`.

* `path_based_rule_ids` - A list of URL Path Map Path Rule IDs from an `azurerm_application_gateway`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web Application Firewall Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web Application Firewall Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Web Application Firewall Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Web Application Firewall Policy.

## Import

Web Application Firewall Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_application_firewall_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/applicationGatewayWebApplicationFirewallPolicies/example-wafpolicy
```
