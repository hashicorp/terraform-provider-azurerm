---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_firewall_policy"
description: |-
  Manages a Front Door (standard/premium) Firewall Policy instance.
---

# azurerm_cdn_frontdoor_firewall_policy

Manages a Front Door (standard/premium) Firewall Policy instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Premium_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_firewall_policy" "example" {
  name                              = "examplecdnfdwafpolicy"
  resource_group_name               = azurerm_resource_group.example.name
  sku_name                          = azurerm_cdn_frontdoor_profile.example.sku_name
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
      match_values       = ["10.0.1.0/24", "10.0.0.0/24"]
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
    action  = "Log"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `sku_name` - (Required) The sku's pricing tier for this Front Door Firewall Policy. Possible values include `Standard_AzureFrontDoor` or `Premium_AzureFrontDoor`. Changing this forces a new resource to be created.

-> **NOTE:** The `Standard_AzureFrontDoor` Front Door Firewall Policy sku may contain `custom` rules only. The `Premium_AzureFrontDoor` Front Door Firewall Policy skus may contain both `custom` and `managed` rules.

* `enabled` - (Optional) Is the Front Door Firewall Policy enabled? Defaults to `true`.

* `mode` - (Required) The Front Door Firewall Policy mode. Possible values are `Detection`, `Prevention`.

-> **NOTE:** When run in `Detection` mode, the Front Door Firewall Policy doesn't take any other actions other than monitoring and logging the request and its matched Front Door Rule to the Web Application Firewall logs.

* `redirect_url` - (Optional) If action type is redirect, this field represents redirect URL for the client.

* `custom_rule` - (Optional) One or more `custom_rule` blocks as defined below.

* `custom_block_response_status_code` - (Optional) If a `custom_rule` block's action type is `block`, this is the response status code. Possible values are `200`, `403`, `405`, `406`, or `429`.

* `custom_block_response_body` - (Optional) If a `custom_rule` block's action type is `block`, this is the response body. The body must be specified in base64 encoding.

* `managed_rule` - (Optional) One or more `managed_rule` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Front Door Firewall Policy.

---

A `custom_rule` block supports the following:

* `name` - (Required) Gets name of the resource that is unique within a policy. This name can be used to access the resource.

* `action` - (Required) The action to perform when the rule is matched. Possible values are `Allow`, `Block`, `Log`, or `Redirect`.

* `enabled` - (Optional) Is the rule is enabled or disabled? Defaults to `true`.

* `priority` - (Optional) The priority of the rule. Rules with a lower value will be evaluated before rules with a higher value. Defaults to `1`.

* `type` - (Required) The type of rule. Possible values are `MatchRule` or `RateLimitRule`.

* `match_condition` - (Optional) One or more `match_condition` block defined below. Can support up to `10` `match_condition` blocks.

* `rate_limit_duration_in_minutes` - (Optional) The rate limit duration in minutes. Defaults to `1`.

* `rate_limit_threshold` - (Optional) The rate limit threshold. Defaults to `10`.

---

A `match_condition` block supports the following:

* `match_variable` - (Required) The request variable to compare with. Possible values are `Cookies`, `PostArgs`, `QueryString`, `RemoteAddr`, `RequestBody`, `RequestHeader`, `RequestMethod`, `RequestUri`, or `SocketAddr`.

* `match_values` - (Required) Up to `600` possible values to match. Limit is in total across all `match_condition` blocks and `match_values` arguments. String value itself can be up to `256` characters in length.

* `operator` - (Required) Comparison type to use for matching with the variable value. Possible values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GeoMatch`, `GreaterThan`, `GreaterThanOrEqual`, `IPMatch`, `LessThan`, `LessThanOrEqual` or `RegEx`.

* `selector` - (Optional) Match against a specific key if the `match_variable` is `QueryString`, `PostArgs`, `RequestHeader` or `Cookies`.

* `negation_condition` - (Optional) Should the result of the condition be negated.

* `transforms` - (Optional) Up to `5` transforms to apply. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `URLDecode` or `URLEncode`.

---

A `managed_rule` block supports the following:

* `type` - (Required) The name of the managed rule to use with this resource. Possible values include `DefaultRuleSet`, `Microsoft_DefaultRuleSet`, `BotProtection` or `Microsoft_BotManagerRuleSet`.

* `version` - (Required) The version of the managed rule to use with this resource. Possible values depends on which DRS type you are using, for the `DefaultRuleSet` type the possible values include `1.0` or `preview-0.1`. For `Microsoft_DefaultRuleSet` the possible values include `1.1`, `2.0` or `2.1`. For `BotProtection` the value must be `preview-0.1` and for `Microsoft_BotManagerRuleSet` the value must be `1.0`.

* `action` - (Required) The action to perform for all DRS rules when the managed rule is matched or when the anomaly score is 5 or greater depending on which version of the DRS you are using. Possible values include `Allow`, `Log`, `Block`, and `Redirect`.

* `exclusion` - (Optional) One or more `exclusion` blocks as defined below.

* `override` - (Optional) One or more `override` blocks as defined below.

---

An `override` block supports the following:

* `rule_group_name` - (Required) The managed rule group to override.

* `exclusion` - (Optional) One or more `exclusion` blocks as defined below.

* `rule` - (Optional) One or more `rule` blocks as defined below. If none are specified, all of the rules in the group will be disabled.

---

A `rule` block supports the following:

* `rule_id` - (Required) Identifier for the managed rule.

* `action` - (Required) The action to be applied when the managed rule matches or when the anomaly score is 5 or greater. Possible values for DRS `1.1` and below are `Allow`, `Log`, `Block`, and `Redirect`. For DRS `2.0` and above the possible values are `Log` or `AnomalyScoring`.

->**NOTE:** Please see the DRS [product documentation](https://learn.microsoft.com/azure/web-application-firewall/afds/waf-front-door-drs?tabs=drs20#anomaly-scoring-mode) for more information.

* `enabled` - (Optional) Is the managed rule override enabled or disabled. Defaults to `false`

* `exclusion` - (Optional) One or more `exclusion` blocks as defined below.

---

An `exclusion` block supports the following:

* `match_variable` - (Required) The variable type to be excluded. Possible values are `QueryStringArgNames`, `RequestBodyPostArgNames`, `RequestCookieNames`, `RequestHeaderNames`, `RequestBodyJsonArgNames`

-> **NOTE:** `RequestBodyJsonArgNames` is only available on Default Rule Set (DRS) 2.0 or later

* `operator` - (Required) Comparison operator to apply to the selector when specifying which elements in the collection this exclusion applies to. Possible values are: `Equals`, `Contains`, `StartsWith`, `EndsWith`, `EqualsAny`.

* `selector` - (Required) Selector for the value in the `match_variable` attribute this exclusion applies to.

-> **NOTE:** `selector` must be set to `*` if `operator` is set to `EqualsAny`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Firewall Policy.

* `frontend_endpoint_ids` - The Front Door Profiles frontend endpoints associated with this Front Door Firewall Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Firewall Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Front Door Firewall Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Firewall Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Firewall Policy.

## Import

Front Door Firewall Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_firewall_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/firewallPolicy1
```
