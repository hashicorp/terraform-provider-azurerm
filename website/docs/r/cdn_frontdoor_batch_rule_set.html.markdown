---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_batch_rule_set"
description: |-
  Manages a Front Door (standard/premium) Batch Rule Set.
---

# azurerm_cdn_frontdoor_batch_rule_set

Manages a Front Door (standard/premium) Batch Rule Set.

~> **Note:** This resource creates the Front Door Rule Set in batch mode and manages the full ordered batch rule collection for it. Any change to the configured `rule` blocks sends the desired final ordered rule list to the Resource Provider in a single request.

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

resource "azurerm_cdn_frontdoor_endpoint" "example" {
  name                     = "example-endpoint"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  tags = {
    endpoint = "contoso.com"
  }
}

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                     = "example-originGroup"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  session_affinity_enabled = true

  restore_traffic_time_to_healed_or_new_endpoint_in_minutes = 10

  health_probe {
    interval_in_seconds = 240
    path                = "/healthProbe"
    protocol            = "Https"
    request_type        = "GET"
  }

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}

resource "azurerm_cdn_frontdoor_origin" "example" {
  name                          = "example-origin"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id
  enabled                       = true

  certificate_name_check_enabled = false

  host_name          = azurerm_cdn_frontdoor_endpoint.example.host_name
  http_port          = 80
  https_port         = 443
  origin_host_header = "contoso.com"
  priority           = 1
  weight             = 500
}

resource "azurerm_cdn_frontdoor_batch_rule_set" "example" {
  name                     = "examplebatchruleset"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  rule {
    name               = "examplebatchrule"
    order              = 1
    behaviour_on_match = "Continue"

    actions {
      route_configuration_override {
        origin_group {
          cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id
          forwarding_protocol           = "HttpsOnly"
        }

        caching {
          behaviour               = "OverrideIfOriginMissing"
          duration                = "365.23:59:59"
          compression_enabled     = true
          query_string_behaviour  = "IncludeSpecifiedQueryStrings"
          query_string_parameters = ["foo", "clientIp={client_ip}"]
        }
      }
    }

    conditions {
      host_name {
        operator   = "Equal"
        values     = ["www.contoso.com", "images.contoso.com", "video.contoso.com"]
        transforms = ["Lowercase", "Trim"]
      }

      device_type {
        operator = "Equal"
        values   = ["Mobile"]
      }

      post_argument {
        name       = "customerName"
        operator   = "BeginsWith"
        values     = ["J", "K"]
        transforms = ["Uppercase"]
      }

      request_method {
        operator = "Equal"
        values   = ["DELETE"]
      }

      request_filename {
        operator   = "Equal"
        values     = ["media.mp4"]
        transforms = ["Lowercase", "RemoveNulls", "Trim"]
      }
    }
  }
}

resource "azurerm_cdn_frontdoor_route" "example" {
  name                          = "example-cdn-frontdoor-route"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.example.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.example.id]
  cdn_frontdoor_rule_set_ids    = [azurerm_cdn_frontdoor_batch_rule_set.example.id]
  patterns_to_match             = ["/*"]
  supported_protocols           = ["Http", "Https"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Front Door Batch Rule Set. Changing this forces a new resource to be created.

* `cdn_frontdoor_profile_id` - (Required) The resource ID of the Front Door Profile where this Front Door Batch Rule Set should be created. Changing this forces a new resource to be created.

* `rule` - (Required) One or more `rule` blocks as defined below. The configured blocks represent the complete set of rules managed for this Front Door Batch Rule Set. The final rule ordering is determined by each block's `order` value. A maximum of `100` `rule` blocks may be defined.

~> **Note:** The `rule` blocks must be declared in ascending `order`, gaps between different rules are allowed. To insert, remove, or move a rule, update the full `rule` collection in the same ascending order that you want Terraform to store.

~> **Note:** Each `rule` block must use a unique `name` value and a unique `order` value.

~> **Note:** Each `rule` that enables caching (using the `route_configuration_override.caching` block with a `behaviour` other than `Disabled`) consumes two of the `100` available rule slots. The plan fails if the effective number of rule slots exceeds this service-side quota.

---

A `rule` block supports the following:

* `actions` - (Required) An `actions` block as defined below.

* `name` - (Required) The name which should be used for this Front Door Batch Rule.

-> **Note:** `name` must be between `1` and `260` characters in length, begin with a letter, and may contain only letters and numbers.

* `order` - (Required) The order in which this rule will be applied for the Front Door Endpoint. Rules with a lesser `order` value are applied before rules with a greater `order` value. Possible values are `0` or greater.

* `behaviour_on_match` - (Optional) The behaviour on a condition match. Possible values are `Continue` and `Stop`. Defaults to `Continue`.

* `conditions` - (Optional) A `conditions` block as defined below.

---

An `actions` block supports the following:

~> **Note:** The `actions` block must define at least one action, and may include up to `5` separate actions in total.

-> **Note:** Specific actions support Rule Set server variables, for more information reference [Azure Front Door Rule Set Server Variables](https://learn.microsoft.com/azure/frontdoor/rule-set-server-variables)

* `modify_request_header` - (Optional) One or more `modify_request_header` blocks as defined below.

* `modify_response_header` - (Optional) One or more `modify_response_header` blocks as defined below.

* `route_configuration_override` - (Optional) A `route_configuration_override` block as defined below.

~> **Note:** `route_configuration_override` conflicts with `url_redirect`.

* `url_redirect` - (Optional) A `url_redirect` block as defined below.

* `url_rewrite` - (Optional) A `url_rewrite` block as defined below.

~> **Note:** `url_rewrite` conflicts with `url_redirect` and vice-versa.

---

A `modify_request_header` block supports the following:

* `header_name` - (Required) The name of the header to modify.

* `operator` - (Required) The action to take on `header_name`. Possible values are `Append`, `Overwrite`, and `Delete`.

* `header_value` - (Optional) The value to append or overwrite.

~> **Note:** `header_value` is required when `operator` is set to `Append` or `Overwrite`, and must not be set when `operator` is set to `Delete`.

---

A `modify_response_header` block supports the following:

* `header_name` - (Required) The name of the header to modify.

* `operator` - (Required) The action to take on `header_name`. Possible values are `Append`, `Overwrite`, and `Delete`.

* `header_value` - (Optional) The value to append or overwrite.

~> **Note:** `header_value` is required when `operator` is set to `Append` or `Overwrite`, and must not be set when `operator` is set to `Delete`.

---

A `route_configuration_override` block supports the following:

* `caching` - (Required) A `caching` block as defined below.

* `origin_group` - (Optional) An `origin_group` block as defined below.

---

An `origin_group` block supports the following:

* `cdn_frontdoor_origin_group_id` - (Required) The Front Door Origin Group resource ID that the request should be routed to.

~> **Note:** If you remove the `origin_group` block from a rule that currently points at the only enabled origin in an Origin Group, apply the Batch Rule Set update first and then remove or disable the last origin in a separate apply. The service rejects deleting or disabling the last origin while the Origin Group is still associated with a route or a rule.

* `forwarding_protocol` - (Required) The forwarding protocol the request is redirected as. Possible values are `MatchRequest`, `HttpOnly`, and `HttpsOnly`.

---

A `caching` block supports the following:

* `behaviour` - (Required) Controls how Front Door handles cache behaviour for the response. Possible values are `HonorOrigin`, `OverrideAlways`, `OverrideIfOriginMissing`, and `Disabled`.

~> **Note:** If `behaviour` is set to `Disabled`, you cannot set `compression_enabled`, `duration`, `query_string_behaviour`, or `query_string_parameters`.

-> **Note:** Enabling caching in a `route_configuration_override` block affects the service-side quota used for rule operations. Each rule that enables caching consumes two of the `100` available rule slots during an update.

* `compression_enabled` - (Optional) Whether compression is enabled. Defaults to `false`.

* `duration` - (Optional) When `behaviour` is set to `OverrideAlways` or `OverrideIfOriginMissing`, this field specifies the cache duration to use and is required. The maximum allowed value is `365.23:59:59`. If the desired maximum cache duration is less than `1` day, specify it in the `HH:MM:SS` format, for example `23:59:59`.

~> **Note:** `duration` must not be set when `behaviour` is set to `HonorOrigin`.

* `query_string_behaviour` - (Optional) Controls how query strings contribute to the cache key. Possible values are `IgnoreQueryString`, `UseQueryString`, `IgnoreSpecifiedQueryStrings`, and `IncludeSpecifiedQueryStrings`.

~> **Note:** `query_string_behaviour` is required when `behaviour` is not set to `Disabled`.

* `query_string_parameters` - (Optional) A list of query string parameter names. A maximum of `100` parameters may be defined.

~> **Note:** `query_string_parameters` is required when `query_string_behaviour` is set to `IncludeSpecifiedQueryStrings` or `IgnoreSpecifiedQueryStrings`, and must not be set when `query_string_behaviour` is set to `UseQueryString` or `IgnoreQueryString`.

---

A `url_redirect` block supports the following:

* `redirect_type` - (Required) The response type to return to the requestor. Possible values are `Moved`, `Found`, `TemporaryRedirect`, and `PermanentRedirect`.

* `destination_fragment` - (Optional) The fragment to use in the redirect. The value must be a string between `1` and `1024` characters in length and must not start with `#`. Leave this unset to preserve the incoming fragment.

* `destination_host_name` - (Optional) The host name you want the request to be redirected to. The value must be a string between `1` and `2048` characters in length. Leave this unset to preserve the incoming host.

* `destination_path` - (Optional) The path to use in the redirect. The value must be a string and include the leading `/`. Leave this unset to preserve the incoming path.

* `query_string` - (Optional) The query string used in the redirect URL. The value must be in the `<key>=<value>` or `<key>={<action_server_variable>}` format and must not include the leading `?`. Leave this unset to preserve the incoming query string. The maximum allowed length for this field is `2048` characters.

* `redirect_protocol` - (Optional) The protocol the request is redirected as. Possible values are `MatchRequest`, `Http`, and `Https`. Defaults to `MatchRequest`.

---

A `url_rewrite` block supports the following:

* `destination_path` - (Required) The destination path to use in the rewrite.

* `source_pattern` - (Required) The source pattern in the URL path to replace.

* `preserve_unmatched_path_enabled` - (Optional) Whether to append the remaining path after the source pattern to the new destination path. Defaults to `false`.

---

A `conditions` block supports the following:

-> **Note:** You may include up to `10` separate conditions in the `conditions` block.

* `client_port` - (Optional) One or more `client_port` blocks as defined below.

* `device_type` - (Optional) One or more `device_type` blocks as defined below.

* `host_name` - (Optional) One or more `host_name` blocks as defined below.

* `http_version` - (Optional) One or more `http_version` blocks as defined below.

* `post_argument` - (Optional) One or more `post_argument` blocks as defined below.

* `query_string` - (Optional) One or more `query_string` blocks as defined below.

* `remote_address` - (Optional) One or more `remote_address` blocks as defined below.

* `request_body` - (Optional) One or more `request_body` blocks as defined below.

* `request_cookies` - (Optional) One or more `request_cookies` blocks as defined below.

* `request_file_extension` - (Optional) One or more `request_file_extension` blocks as defined below.

* `request_filename` - (Optional) One or more `request_filename` blocks as defined below.

* `request_header` - (Optional) One or more `request_header` blocks as defined below.

* `request_method` - (Optional) One or more `request_method` blocks as defined below.

* `request_path` - (Optional) One or more `request_path` blocks as defined below.

* `request_scheme` - (Optional) One or more `request_scheme` blocks as defined below.

* `request_url` - (Optional) One or more `request_url` blocks as defined below.

* `server_port` - (Optional) One or more `server_port` blocks as defined below.

* `socket_address` - (Optional) One or more `socket_address` blocks as defined below.

* `ssl_protocol` - (Optional) One or more `ssl_protocol` blocks as defined below.

---

A `client_port` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `values` - (Optional) One or more values representing the client port to match. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `device_type` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Equal` and `NotEqual`.

* `values` - (Required) The device type to match. Possible values are `Mobile` and `Desktop`.

~> **Note:** Currently, only a single value may be specified.

---

A `host_name` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `transforms` - (Optional) A list of condition transforms. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`. A maximum of `4` transforms may be defined.

* `values` - (Optional) A list of one or more values representing the request hostname to match. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `http_version` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Equal` and `NotEqual`.

* `values` - (Required) A list of one or more HTTP versions to match. Possible values are `2.0`, `1.1`, `1.0`, and `0.9`.

---

A `post_argument` block supports the following:

* `name` - (Required) A string value representing the name of the `POST` argument.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `transforms` - (Optional) A list of condition transforms. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`. A maximum of `4` transforms may be defined.

* `values` - (Optional) One or more values representing the `POST` argument value to match. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `query_string` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `transforms` - (Optional) A list of condition transforms. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`. A maximum of `4` transforms may be defined.

* `values` - (Optional) One or more values representing the query string value to match. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `remote_address` block supports the following:

* `operator` - (Required) The type of remote address to match. Possible values are `GeoMatch`, `IPMatch`, `NotGeoMatch`, and `NotIPMatch`.

* `values` - (Required) A list of CIDR ranges or country codes. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** When `operator` is set to `GeoMatch` or `NotGeoMatch`, each value in `values` must be a two-letter uppercase country code.

~> **Note:** When `operator` is set to `IPMatch` or `NotIPMatch`, each value in `values` must be a valid CIDR range.

---

A `request_body` block supports the following:

-> **Note:** If a request body exceeds 64 KB, only the first 64 KB is considered for this condition.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `transforms` - (Optional) A list of condition transforms. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`. A maximum of `4` transforms may be defined.

* `values` - (Optional) One or more values representing the request body text to match. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `request_cookies` block supports the following:

* `name` - (Required) The name of the cookie.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `transforms` - (Optional) A list of condition transforms. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`. A maximum of `4` transforms may be defined.

* `values` - (Optional) One or more values representing the cookie value to match. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `request_file_extension` block supports the following:

-> **Note:** `request_file_extension` identifies requests that include the specified file extension. Do not include a leading period.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `transforms` - (Optional) A list of condition transforms. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`. A maximum of `4` transforms may be defined.

* `values` - (Optional) One or more values representing the request file extension to match. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `request_filename` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `transforms` - (Optional) A list of condition transforms. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`. A maximum of `4` transforms may be defined.

* `values` - (Optional) One or more values representing the request file name to match. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `request_header` block supports the following:

* `name` - (Required) The name of the request header.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `transforms` - (Optional) A list of condition transforms. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`. A maximum of `4` transforms may be defined.

* `values` - (Optional) One or more values representing the request header value to match. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `request_method` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Equal` and `NotEqual`.

* `values` - (Required) A list of one or more HTTP methods. Possible values are `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `OPTIONS`, and `TRACE`. A maximum of `7` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

---

A `request_path` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `Wildcard`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, `NotRegEx`, and `NotWildcard`.

* `transforms` - (Optional) A list of condition transforms. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`. A maximum of `4` transforms may be defined.

* `values` - (Optional) One or more values representing the request path to match. Do not include the leading slash (`/`). A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `request_scheme` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Equal` and `NotEqual`.

* `values` - (Required) The request protocol to match. Possible values are `HTTP` and `HTTPS`.

~> **Note:** Currently, only a single value may be specified

---

A `request_url` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `transforms` - (Optional) A list of condition transforms. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`. A maximum of `4` transforms may be defined.

* `values` - (Optional) One or more values representing the request URL to match. A maximum of `25` values may be defined. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** `values` must not be set when `operator` is set to `Any` or `NotAny`, and is required for all other operators.

---

A `server_port` block supports the following:

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, `NotAny`, `NotEqual`, `NotContains`, `NotBeginsWith`, `NotEndsWith`, `NotLessThan`, `NotLessThanOrEqual`, `NotGreaterThan`, `NotGreaterThanOrEqual`, and `NotRegEx`.

* `values` - (Optional) A list of one or more values representing the server port to match. Possible values are `80` and `443`. If multiple values are specified, they are evaluated using `OR` logic.

---

A `socket_address` block supports the following:

* `operator` - (Required) The type of match. Possible values are `IPMatch` and `NotIPMatch`.

* `values` - (Required) One or more IP address ranges. A maximum of `25` values may be defined. If multiple IP address ranges are specified, they are evaluated using `OR` logic.

---

A `ssl_protocol` block supports the following:

-> **Note:** `ssl_protocol` identifies requests based on the SSL protocol of an established TLS connection.

* `operator` - (Required) A condition operator. Possible values are `Equal` and `NotEqual`.

* `values` - (Required) A list of one or more SSL protocol values. Possible values are `TLSv1`, `TLSv1.1`, and `TLSv1.2`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Batch Rule Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 4 hours) Used when creating the Front Door Batch Rule Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Batch Rule Set.
* `update` - (Defaults to 4 hours) Used when updating the Front Door Batch Rule Set.
* `delete` - (Defaults to 6 hours) Used when deleting the Front Door Batch Rule Set.

## Import

A Front Door Batch Rule Set can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_batch_rule_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1
```

~> **Note:** Only Rule Sets that were provisioned in batch mode can be managed by this resource. Importing a Rule Set that was not provisioned in batch mode returns an error - use `azurerm_cdn_frontdoor_rule_set` instead.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cdn` - 2025-12-01
