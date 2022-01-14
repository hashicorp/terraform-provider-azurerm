---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_rules_engine"
description: |-
  Manages an Azure Front Door Rules Engine configuration and rules.
---

# azurerm_frontdoor_rules_engine

Manages an Azure Front Door Rules Engine configuration and rules.

## Example Usage

```hcl
resource "azurerm_frontdoor_rules_engine" "example_rules_engine" {
  name                = "exampleRulesEngineConfig1"
  frontdoor_name      = azurerm_frontdoor.example.name
  resource_group_name = azurerm_frontdoor.example.resource_group_name

  rule {
    name     = "debuggingoutput"
    priority = 1

    action {
      response_header {
        header_action_type = "Append"
        header_name        = "X-TEST-HEADER"
        value              = "Append Header Rule"
      }
    }
  }

  rule {
    name     = "overwriteorigin"
    priority = 2

    match_condition {
      variable = "RequestMethod"
      operator = "Equal"
      value    = ["GET", "POST"]
    }

    action {

      response_header {
        header_action_type = "Overwrite"
        header_name        = "Access-Control-Allow-Origin"
        value              = "*"
      }

      response_header {
        header_action_type = "Overwrite"
        header_name        = "Access-Control-Allow-Credentials"
        value              = "true"
      }

    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Rules engine configuration. Changing this forces a new resource to be created.

* `frontdoor_name` - (Required) The name of the Front Door instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `rule` - (Required) A `rule` block as defined below.

---

The `rule` block supports the following:

* `name` - (Required) The name of the rule.

* `priority` - (Required) Priority of the rule, must be unique per rules engine definition.

* `action` - (Required) A `rule_action` block as defined below.

* `match_condition` - One or more `match_condition` block as defined below.

---

The `action` block supports the following:

* `request_header` - (Optional) A `request_header` block as defined below.

* `response_header` - (Optional) A `response_header` block as defined below.

* `forwarding_configuration` - (Optional) A `forwarding_configuration` block as defined below, must be unique per action.

* `redirect_configuration`   - (Optional) A `redirect_configuration` block as defined below, must be unique per action.

---

The `request_header` block supports the following:

* `header_action_type` can be set to `Overwrite`, `Append` or `Delete`.

* `header_name` header name (string).

* `value` value name (string).

---

The `response_header` block supports the following:

* `header_action_type` can be set to `Overwrite`, `Append` or `Delete`.

* `header_name` header name (string).

* `value` value name (string).

---

The `forwarding_configuration` block supports the following:

* `backend_pool_name` - (Required) Specifies the name of the Backend Pool to forward the incoming traffic to.

* `cache_enabled` - (Optional) Specifies whether to Enable caching or not. Valid options are `true` or `false`. Defaults to `false`.

* `cache_use_dynamic_compression` - (Optional) Whether to use dynamic compression when caching. Valid options are `true` or `false`. Defaults to `false`.

* `cache_query_parameter_strip_directive` - (Optional) Defines cache behaviour in relation to query string parameters. Valid options are `StripAll`, `StripAllExcept`, `StripOnly` or `StripNone`. Defaults to `StripAll`.

* `cache_query_parameters` - (Optional) Specify query parameters (array). Works only in combination with `cache_query_parameter_strip_directive` set to `StripAllExcept` or `StripOnly`.

* `cache_duration` - (Optional) Specify the caching duration (in ISO8601 notation e.g. `P1DT2H` for 1 day and 2 hours). Needs to be greater than 0 and smaller than 365 days. `cache_duration` works only in combination with `cache_enabled` set to `true`.

* `custom_forwarding_path` - (Optional) Path to use when constructing the request to forward to the backend. This functions as a URL Rewrite. Default behaviour preserves the URL path.

* `forwarding_protocol` - (Optional) Protocol to use when redirecting. Valid options are `HttpOnly`, `HttpsOnly`, or `MatchRequest`. Defaults to `HttpsOnly`.

---

The `redirect_configuration` block supports the following:

* `custom_host` - (Optional)  Set this to change the URL for the redirection.

* `redirect_protocol` - (Optional) Protocol to use when redirecting. Valid options are `HttpOnly`, `HttpsOnly`, or `MatchRequest`. Defaults to `MatchRequest`

* `redirect_type` - (Required) Status code for the redirect. Valida options are `Moved`, `Found`, `TemporaryRedirect`, `PermanentRedirect`.

* `custom_fragment` - (Optional) The destination fragment in the portion of URL after '#'. Set this to add a fragment to the redirect URL.

* `custom_path` - (Optional) The path to retain as per the incoming request, or update in the URL for the redirection.

* `custom_query_string` - (Optional) Replace any existing query string from the incoming request URL.

---

The `match_condition` block supports the following:

* `variable` can be set to `IsMobile`, `RemoteAddr`, `RequestMethod`, `QueryString`, `PostArgs`, `RequestURI`, `RequestPath`, `RequestFilename`, `RequestFilenameExtension`,`RequestHeader`,`RequestBody` or `RequestScheme`.

* `selector` match against a specific key when `variable` is set to `PostArgs` or `RequestHeader`. It cannot be used with `QueryString` and `RequestMethod`. Defaults to `null`.

* `operator` can be set to `Any`, `IPMatch`, `GeoMatch`, `Equal`, `Contains`, `LessThan`, `GreaterThan`, `LessThanOrEqual`, `GreaterThanOrEqual`, `BeginsWith` or `EndsWith`

* `transform` can be set to one or more values out of `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` and `UrlEncode`

* `negate_condition` can be set to `true` or `false` to negate the given condition. Defaults to `true`.

* `value` (array) can contain one or more strings.

## Import

Azure Front Door Rules Engine's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_frontdoor_rules_engine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/frontdoors/frontdoor1/rulesengines/rule1
```
