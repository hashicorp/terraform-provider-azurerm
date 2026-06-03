---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_rule"
description: |-
  Manages a Front Door (standard/premium) Rule.
---

# azurerm_cdn_frontdoor_rule

Manages a Front Door (standard/premium) Rule.

!> **Note:** The Rules resource **must** include a `depends_on` meta-argument which references the `azurerm_cdn_frontdoor_origin` and the `azurerm_cdn_frontdoor_origin_group`.

~> **Note:** Azure Front Door Rule operations are currently affected by a service-side regression where unattached rules or rule sets can fail with `400 Bad Request` until they are associated with a Front Door Route. As a result, unattached and attached scenarios can currently behave differently while the service-side fix is pending.

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

resource "azurerm_cdn_frontdoor_rule_set" "example" {
  name                     = "exampleruleset"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
}

resource "azurerm_cdn_frontdoor_rule" "example" {
  depends_on = [azurerm_cdn_frontdoor_origin_group.example, azurerm_cdn_frontdoor_origin.example]

  name                      = "examplerule"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.example.id
  order                     = 1
  behavior_on_match         = "Continue"

  actions {
    route_configuration_override_action {
      cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id
      forwarding_protocol           = "HttpsOnly"
      query_string_caching_behavior = "IncludeSpecifiedQueryStrings"
      query_string_parameters       = ["foo", "clientIp={client_ip}"]
      compression_enabled           = true
      cache_behavior                = "OverrideIfOriginMissing"
      cache_duration                = "365.23:59:59"
    }

    url_redirect_action {
      redirect_type        = "PermanentRedirect"
      redirect_protocol    = "MatchRequest"
      query_string         = "clientIp={client_ip}"
      destination_path     = "/exampleredirection"
      destination_hostname = "contoso.com"
      destination_fragment = "UrlRedirect"
    }
  }

  conditions {
    host_name_condition {
      operator         = "Equal"
      negate_condition = false
      match_values     = ["www.contoso.com", "images.contoso.com", "video.contoso.com"]
      transforms       = ["Lowercase", "Trim"]
    }

    is_device_condition {
      operator         = "Equal"
      negate_condition = false
      match_values     = ["Mobile"]
    }

    post_args_condition {
      post_args_name = "customerName"
      operator       = "BeginsWith"
      match_values   = ["J", "K"]
      transforms     = ["Uppercase"]
    }

    request_method_condition {
      operator         = "Equal"
      negate_condition = false
      match_values     = ["DELETE"]
    }

    url_filename_condition {
      operator         = "Equal"
      negate_condition = false
      match_values     = ["media.mp4"]
      transforms       = ["Lowercase", "RemoveNulls", "Trim"]
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Front Door Rule. Changing this forces a new resource to be created.

-> **Note:** Possible values must be between 1 and 260 characters in length, begin with a letter, and may contain only letters and numbers.

* `cdn_frontdoor_rule_set_id` - (Required) The resource ID of the Front Door Rule Set for this Front Door Rule. Changing this forces a new resource to be created.

~> **Note:** The parent Rule Set must use the existing Front Door Standard/Premium per-rule update mode. Batch-mode Rule Sets are managed by `azurerm_cdn_frontdoor_batch_rule_set` instead.

* `actions` - (Required) An `actions` block as defined below.

* `order` - (Required) The order in which the rules will be applied for the Front Door Endpoint. Rules with a lesser `order` value are applied before rules with a greater `order` value. Possible values are `0` or greater.

* `behavior_on_match` - (Optional) If this rule is a match, should the rules engine continue processing the remaining rules or stop? Possible values are `Continue` and `Stop`. Defaults to `Continue`.

* `conditions` - (Optional) A `conditions` block as defined below.

---

An `actions` block supports the following:

-> **Note:** You may include up to 5 separate actions in the `actions` block.

Some actions support `Action Server Variables` which provide access to structured information about the request. For more information about `Action Server Variables` see the `Action Server Variables` as defined below.

* `request_header_action` - (Optional) A `request_header_action` block as defined below.

* `response_header_action` - (Optional) A `response_header_action` block as defined below.

* `route_configuration_override_action` - (Optional) A `route_configuration_override_action` block as defined below.

* `url_redirect_action` - (Optional) A `url_redirect_action` block as defined below. You may **not** have a `url_redirect_action` **and** a `url_rewrite_action` defined in the same `actions` block.

* `url_rewrite_action` - (Optional) A `url_rewrite_action` block as defined below. You may **not** have a `url_rewrite_action` **and** a `url_redirect_action` defined in the same `actions` block.

---

An `url_redirect_action` block supports the following:

* `destination_hostname` - (Required) The host name you want the request to be redirected to. The value must be a string between `0` and `2048` characters in length. Leave this blank to preserve the incoming host.

* `redirect_type` - (Required) The response type to return to the requestor. Possible values are `Moved`, `Found`, `TemporaryRedirect`, and `PermanentRedirect`.

* `destination_fragment` - (Optional) The fragment to use in the redirect. The value must be a string between `0` and `1024` characters in length. Leave this blank to preserve the incoming fragment. Defaults to `""`.

* `destination_path` - (Optional) The path to use in the redirect. The value must be a string and include the leading `/`. Leave this blank to preserve the incoming path. Defaults to `""`.

* `query_string` - (Optional) The query string used in the redirect URL. The value must be in the `<key>=<value>` or `<key>={`action_server_variable`}` format and must not include the leading `?`. Leave this blank to preserve the incoming query string. The maximum allowed length for this field is `2048` characters. Defaults to `""`.

* `redirect_protocol` - (Optional) The protocol the request is redirected as. Possible values are `MatchRequest`, `Http`, and `Https`. Defaults to `MatchRequest`.

---

A `route_configuration_override_action` block supports the following:

* `cache_behavior` - (Optional) Controls how Front Door handles cache behavior for the response. Possible values are `Disabled`, `HonorOrigin`, `OverrideAlways`, and `OverrideIfOriginMissing`.

* `cache_duration` - (Optional) When `cache_behavior` is set to `OverrideAlways` or `OverrideIfOriginMissing`, this field specifies the cache duration to use. The maximum duration is 366 days specified in the `d.HH:MM:SS` format (for example `365.23:59:59`). If the desired maximum cache duration is less than 1 day, the maximum cache duration should be specified in the `HH:MM:SS` format (for example `23:59:59`).

* `cdn_frontdoor_origin_group_id` - (Optional) The Front Door Origin Group resource ID that the request should be routed to. This overrides the configuration specified in the Front Door Endpoint route.

* `compression_enabled` - (Optional) Whether dynamic compression is enabled. Possible values are `true` and `false`.

-> **Note:** Content is not compressed when the requested content is smaller than `1 byte` or larger than `1 MB`.

* `forwarding_protocol` - (Optional) The forwarding protocol the request is redirected as. This overrides the configuration specified in the associated route. Possible values are `MatchRequest`, `HttpOnly`, and `HttpsOnly`.

~> **Note:** If `cdn_frontdoor_origin_group_id` is not defined, you cannot set `forwarding_protocol`.

* `query_string_caching_behavior` - (Optional) Controls how query strings contribute to the cache key. Possible values are `IgnoreQueryString`, `UseQueryString`, `IgnoreSpecifiedQueryStrings`, and `IncludeSpecifiedQueryStrings`.

* `query_string_parameters` - (Optional) A list of query string parameter names.

~> **Note:** `query_string_parameters` is required when `query_string_caching_behavior` is set to `IncludeSpecifiedQueryStrings` or `IgnoreSpecifiedQueryStrings`.

---

An `url_rewrite_action` block supports the following:

* `destination` - (Required) The destination path to use in the rewrite. The destination path overwrites the source pattern.

* `source_pattern` - (Required) The source pattern in the URL path to replace. This uses prefix-based matching. For example, to match all URL paths, use `/` as the source pattern value.

* `preserve_unmatched_path` - (Optional) Whether to append the remaining path after the source pattern to the new destination path. Possible values are `true` and `false`. Defaults to `false`.

---

A `request_header_action` block supports the following:

* `header_action` - (Required) The action to take on `header_name`. Possible values are `Append`, `Overwrite`, and `Delete`.

-> **Note:** `Append` adds the specified header to the request with the specified value. If the header is already present, the value is appended to the existing header value using string concatenation. `Overwrite` adds the specified header to the request with the specified value, replacing any existing value. `Delete` removes the header from the request.

* `header_name` - (Required) The name of the header to modify.

* `value` - (Optional) The value to append or overwrite.

~> **Note:** `value` is required when `header_action` is set to `Append` or `Overwrite`.

---

A `response_header_action` block supports the following:

* `header_action` - (Required) The action to take on `header_name`. Possible values are `Append`, `Overwrite`, and `Delete`.

-> **Note:** `Append` adds the specified header to the response with the specified value. If the header is already present, the value is appended to the existing header value using string concatenation. `Overwrite` adds the specified header to the response with the specified value, replacing any existing value. `Delete` removes the header from the response.

* `header_name` - (Required) The name of the header to modify.

* `value` - (Optional) The value to append or overwrite.

~> **Note:** `value` is required when `header_action` is set to `Append` or `Overwrite`.

---

A `conditions` block supports the following:

-> **Note:** You may include up to 10 separate conditions in the `conditions` block.

* `client_port_condition` - (Optional) A `client_port_condition` block as defined below.
* `cookies_condition` - (Optional) A `cookies_condition` block as defined below.
* `host_name_condition` - (Optional) A `host_name_condition` block as defined below.
* `http_version_condition` - (Optional) A `http_version_condition` block as defined below.
* `is_device_condition` - (Optional) A `is_device_condition` block as defined below.
* `post_args_condition` - (Optional) A `post_args_condition` block as defined below.
* `query_string_condition` - (Optional) A `query_string_condition` block as defined below.
* `remote_address_condition` - (Optional) A `remote_address_condition` block as defined below.
* `request_body_condition` - (Optional) A `request_body_condition` block as defined below.
* `request_header_condition` - (Optional) A `request_header_condition` block as defined below.
* `request_method_condition` - (Optional) A `request_method_condition` block as defined below.
* `request_scheme_condition` - (Optional) A `request_scheme_condition` block as defined below.
* `request_uri_condition` - (Optional) A `request_uri_condition` block as defined below.
* `server_port_condition` - (Optional) A `server_port_condition` block as defined below.
* `socket_address_condition` - (Optional) A `socket_address_condition` block as defined below.
* `ssl_protocol_condition` - (Optional) A `ssl_protocol_condition` block as defined below.
* `url_file_extension_condition` - (Optional) A `url_file_extension_condition` block as defined below.
* `url_filename_condition` - (Optional) A `url_filename_condition` block as defined below.
* `url_path_condition` - (Optional) A `url_path_condition` block as defined below.

---

A `ssl_protocol_condition` block supports the following:

-> **Note:** The `ssl_protocol_condition` identifies requests based on the SSL protocol of an established TLS connection.

* `match_values` - (Required) A list of one or more SSL protocol values. Possible values are `TLSv1`, `TLSv1.1`, and `TLSv1.2`.
* `operator` - (Optional) The only possible value is `Equal`. Defaults to `Equal`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.

---

A `host_name_condition` block supports the following:

-> **Note:** The `host_name_condition` identifies requests based on the specified hostname in the request from client.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `match_values` - (Optional) A list of one or more string values representing the request hostname to match. If multiple values are specified, they are evaluated using `OR` logic.
* `transforms` - (Optional) A condition transform. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.

---

A `server_port_condition` block supports the following:

-> **Note:** The `server_port_condition` identifies requests based on which port of the Front Door server accepted the request on.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `match_values` - (Required) A list of one or more integer values representing the server port to match. Common values are `80` and `443`. If multiple values are specified, they are evaluated using `OR` logic.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.

---

A `client_port_condition` block supports the following:

-> **Note:** The `client_port_condition` identifies requests based on the port of the client which made the request.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `match_values` - (Optional) One or more integer values representing the client port to match. If multiple values are specified, they are evaluated using `OR` logic.

---

A `socket_address_condition` block supports the following:

-> **Note:** The `socket_address_condition` identifies requests based on the IP address of the direct connection to the Front Door Profiles edge. If the client used an HTTP proxy or a load balancer to send the request, the value of Socket address is the IP address of the proxy or load balancer.

-> **Note:** Remote Address represents the original client IP that is either from the network connection or typically the `X-Forwarded-For` request header if the user is behind a proxy.

* `operator` - (Optional) The type of match. Possible values are `IpMatch` and `Any`. Defaults to `IPMatch`.

* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `match_values` - (Optional) One or more IP address ranges. If multiple IP address ranges are specified, they are evaluated using `OR` logic.

~> **Note:** If `operator` is set to `IpMatch`, `match_values` is also required.

-> **Note:** See the `Specifying IP Address Ranges` section below for how to define `match_values`.

---

A `remote_address_condition` block supports the following:

-> **Note:** Remote Address represents the original client IP that is either from the network connection or typically the `X-Forwarded-For` request header if the user is behind a proxy.

* `operator` - (Optional) The type of remote address to match. Possible values are `Any`, `GeoMatch`, and `IPMatch`. Defaults to `IPMatch`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `match_values` - (Optional) For `IPMatch`, specify one or more IP address ranges. For `GeoMatch`, specify one or more country codes. If multiple values are specified, they are evaluated using `OR` logic.

~> **Note:** When `operator` is set to `GeoMatch`, each value in `match_values` must be a two-letter uppercase country code.

-> **Note:** See the `Specifying IP Address Ranges` section below for how to define `match_values`.

---

A `request_method_condition` block supports the following:

-> **Note:** The `request_method_condition` identifies requests that use the specified HTTP request method.

* `match_values` - (Required) A list of one or more HTTP methods. Possible values are `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `OPTIONS`, and `TRACE`. If multiple values are specified, they are evaluated using `OR` logic.
* `operator` - (Optional) The only possible value is `Equal`. Defaults to `Equal`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.

---

A `query_string_condition` block supports the following:

-> **Note:** Use the `query_string_condition` to identify requests that contain a specific query string.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `match_values` - (Optional) One or more string or integer values representing the query string value to match. If multiple values are specified, they are evaluated using `OR` logic.
* `transforms` - (Optional) A condition transform. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`.

---

A `post_args_condition` block supports the following:

-> **Note:** Use the `post_args_condition` to identify requests based on the arguments provided within a `POST` request's body. A single match condition matches a single argument from the `POST` request's body.

* `post_args_name` - (Required) A string value representing the name of the `POST` argument.
* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `match_values` - (Optional) One or more string or integer values representing the `POST` argument value to match. If multiple values are specified, they are evaluated using `OR` logic.
* `transforms` - (Optional) A condition transform. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`.

---

A `request_uri_condition` block supports the following:

-> **Note:** The `request_uri_condition` identifies requests that match the specified URL. The entire URL is evaluated, including the protocol and query string, but not the fragment. When you use this rule condition, be sure to include the protocol(e.g. For example, use `https://www.contoso.com` instead of just `www.contoso.com`).

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `match_values` - (Optional) One or more string or integer values representing the request URL to match. If multiple values are specified, they are evaluated using `OR` logic.
* `transforms` - (Optional) A condition transform. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`.

---

A `request_header_condition` block supports the following:

-> **Note:** The `request_header_condition` identifies requests that include a specific header in the request. You can use this match condition to check if a header exists whatever its value, or to check if the header matches a specified value.

* `header_name` - (Required) The name of the request header.
* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `match_values` - (Optional) One or more string or integer values representing the request header value to match. If multiple values are specified, they are evaluated using `OR` logic.
* `transforms` - (Optional) A condition transform. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`.

---

A `request_body_condition` block supports the following:

-> **Note:** The `request_body_condition` identifies requests based on specific text that appears in the body of the request.

-> **Note:** If a request body exceeds `64 KB` in size, only the first `64 KB` will be considered for the request body match condition.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `match_values` - (Required) A list of one or more string or integer values representing the request body text to match. If multiple values are specified, they are evaluated using `OR` logic.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `transforms` - (Optional) A condition transform. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`.

---

A `request_scheme_condition` block supports the following:

-> **Note:** The `request_scheme_condition` identifies requests that use the specified protocol.

* `match_values` - (Optional) The request protocol to match. Possible values are `HTTP` and `HTTPS`.

~> **Note:** `match_values` must be set when `request_scheme_condition` is used.

* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `operator` - (Optional) The only possible value is `Equal`. Defaults to `Equal`.

---

An `url_path_condition` block supports the following:

-> **Note:** The `url_path_condition` identifies requests that include the specified path in the request URL. The path is the part of the URL after the hostname and a slash(e.g. in the URL `https://www.contoso.com/files/secure/file1.pdf`, the path is `files/secure/file1.pdf`).

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `RegEx`, and `Wildcard`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `match_values` - (Optional) One or more values representing the request path to match. Do not include the leading slash (`/`). If multiple values are specified, they are evaluated using `OR` logic.
* `transforms` - (Optional) A condition transform. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`.

---

An `url_file_extension_condition` block supports the following:

-> **Note:** The `url_file_extension_condition` identifies requests that include the specified file extension in the file name in the request URL. Don't include a leading period(e.g. use `html` instead of `.html`).

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `match_values` - (Required) A list of one or more values representing the request file extension to match. If multiple values are specified, they are evaluated using `OR` logic.
* `transforms` - (Optional) A condition transform. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`.

---

An `url_filename_condition` block supports the following:

-> **Note:** The `url_filename_condition` identifies requests that include the specified file name in the request URL.

* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `match_values` - (Optional) A list of one or more values representing the request file name to match. If multiple values are specified, they are evaluated using `OR` logic.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `transforms` - (Optional) A condition transform. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`.

~> **Note:** `match_values` is only optional when `operator` is set to `Any`.

---

A `http_version_condition` block supports the following:

-> **Note:** Use the HTTP version match condition to identify requests that have been made by using a specific version of the HTTP protocol.

* `match_values` - (Required) The HTTP version to match. Possible values are `2.0`, `1.1`, `1.0`, and `0.9`.
* `operator` - (Optional) The only possible value is `Equal`. Defaults to `Equal`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.

---

A `cookies_condition` block supports the following:

-> **Note:** `cookies_condition` identifies requests that include a specific cookie.

* `cookie_name` - (Required) The name of the cookie.
* `operator` - (Required) A condition operator. Possible values are `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, and `RegEx`.
* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `match_values` - (Optional) One or more values representing the cookie value to match. If multiple values are specified, they are evaluated using `OR` logic.
* `transforms` - (Optional) A condition transform. Possible values are `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode`, and `UrlEncode`.

---

An `is_device_condition` block supports the following:

-> **Note:** `is_device_condition` identifies requests made from a `Mobile` or `Desktop` device.

* `match_values` - (Optional) The device type to match. Possible values are `Mobile` and `Desktop`.

~> **Note:** `match_values` must be set when `is_device_condition` is used.

* `negate_condition` - (Optional) Whether to negate the condition. Possible values are `true` and `false`. Defaults to `false`.
* `operator` - (Optional) The only possible value is `Equal`. Defaults to `Equal`.

---

## Specifying IP Address Ranges

When specifying IP address ranges in `socket_address_condition` and `remote_address_condition`, use `CIDR` notation. This means the syntax for an IP address block is the base IP address followed by a forward slash and the prefix size.

* `IPv4` example: `5.5.5.64/26` matches any requests that arrive from addresses `5.5.5.64` through `5.5.5.127`.
* `IPv6` example: `1:2:3:/48` matches any requests that arrive from addresses `1:2:3:0:0:0:0:0` through `1:2:3:ffff:ffff:ffff:ffff:ffff`.

When you specify multiple IP addresses and IP address blocks, `OR` logic is applied.

* `IPv4` example: if you add `1.2.3.4` and `10.20.30.40`, the condition matches requests from either address.
* `IPv6` example: if you add `1:2:3:4:5:6:7:8` and `10:20:30:40:50:60:70:80`, the condition matches requests from either address.

---

## Action Server Variables

Rule Set server variables provide access to structured information about the request. You can use server variables to dynamically change request and response headers or URL rewrite paths and query strings.

### Supported Action Server Variables

| Variable name | Description |
|---------------|-------------|
| `socket_ip` | The IP address of the direct connection to the Front Door edge. |
| `client_ip` | The IP address of the client that made the original request. |
| `client_port` | The client port used for the request. |
| `hostname` | The host name in the client request. |
| `geo_country` | The requester country or region code. |
| `http_method` | The HTTP method used to make the request. |
| `http_version` | The request protocol version. |
| `query_string` | The query string portion of the request URL. |
| `request_scheme` | The request scheme, either `http` or `https`. |
| `request_uri` | The full original request URI including arguments. |
| `ssl_protocol` | The protocol of an established TLS connection. |
| `server_port` | The server port that accepted the request. |
| `url_path` | The request URI path without the query string. |

### Action Server Variable Format

Server variables can be specified using the following formats:

* `{variable}` - Includes the entire server variable.
* `{variable:offset}` - Includes the server variable after a specific zero-based offset.
* `{variable:offset:length}` - Includes the server variable after a specific zero-based offset, up to the specified length.

### Action Server Variables Support

Action Server variables are supported on the following actions:

* `route_configuration_override_action`
* `request_header_action`
* `response_header_action`
* `url_redirect_action`
* `url_rewrite_action`

---

## Condition Operator list

For rules that accept values from the standard operator list, the following operators are valid:

| Operator | Description | Condition Value |
|----------|-------------|-----------------|
| Any | Matches when there is any value, regardless of what it is. | Any |
| Equal | Matches when the value exactly matches the specified string. | Equal |
| Contains | Matches when the value contains the specified string. | Contains |
| Less Than | Matches when the length of the value is less than the specified integer. | LessThan |
| Greater Than | Matches when the length of the value is greater than the specified integer. | GreaterThan |
| Less Than or Equal | Matches when the length of the value is less than or equal to the specified integer. | LessThanOrEqual |
| Greater Than or Equal | Matches when the length of the value is greater than or equal to the specified integer. | GreaterThanOrEqual |
| Begins With | Matches when the value begins with the specified string. | BeginsWith |
| Ends With | Matches when the value ends with the specified string. | EndsWith |
| RegEx | Matches when the value matches the specified regular expression. See `Condition Regular Expressions` below for more details. | RegEx |
| Wildcard | Matches when the request path matches a wildcard expression. See `Condition Wildcard Expression` below for more details. | Wildcard |
| Not Any | Matches when there is no value. | Any and negateCondition = true |
| Not Equal | Matches when the value does not match the specified string. | Equal and negateCondition = true |
| Not Contains | Matches when the value does not contain the specified string. | Contains and negateCondition = true |
| Not Less Than | Matches when the length of the value is not less than the specified integer. | LessThan and negateCondition = true |
| Not Greater Than | Matches when the length of the value is not greater than the specified integer. | GreaterThan and negateCondition = true |
| Not Less Than or Equal | Matches when the length of the value is not less than or equal to the specified integer. | LessThanOrEqual and negateCondition = true |
| Not Greater Than or Equals | Matches when the length of the value is not greater than or equal to the specified integer. | GreaterThanOrEqual and negateCondition = true |
| Not Begins With | Matches when the value does not begin with the specified string. | BeginsWith and negateCondition = true |
| Not Ends With | Matches when the value does not end with the specified string. | EndsWith and negateCondition = true |
| Not RegEx | Matches when the value does not match the specified regular expression. See `Condition Regular Expressions` below for more details. | RegEx and negateCondition = true |
| Not Wildcard | Matches when the request path does not match a wildcard expression. See `Condition Wildcard Expression` below for more details. | Wildcard and negateCondition = true |

---

## Condition Regular Expressions

Regular expressions do **not** support the following operations:

* Backreferences and capturing subexpressions.
* Arbitrary zero-width assertions.
* Subroutine references and recursive patterns.
* Conditional patterns.
* Backtracking control verbs.
* The `\C` single-byte directive.
* The `\R` newline match directive.
* The `\K` start of match reset directive.
* Callouts and embedded code.
* Atomic grouping and possessive quantifiers.

## Condition Wildcard Expression

A wildcard expression can include the `*` character to match zero or more characters within the path. For example, `files/customer*/file.pdf` matches `files/customer1/file.pdf`, `files/customer109/file.pdf`, and `files/customer/file.pdf`, but does not match `files/customer2/anotherfile.pdf`.

---

## Condition Transform List

For rules that can transform strings, the following transforms are valid:

| Transform   | Description |
|-------------|-------------|
| Lowercase   | Converts the string to the lowercase representation. |
| Uppercase   | Converts the string to the uppercase representation. |
| Trim        | Trims leading and trailing whitespace from the string. |
| RemoveNulls | Removes null values from the string. |
| URLEncode   | URL-encodes the string. |
| URLDecode   | URL-decodes the string. |

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Rule.

* `cdn_frontdoor_rule_set_name` - The name of the Front Door Rule Set containing this Front Door Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 4 hours) Used when creating the Front Door Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Rule.
* `update` - (Defaults to 4 hours) Used when updating the Front Door Rule.
* `delete` - (Defaults to 6 hours) Used when deleting the Front Door Rule.

## Import

A Front Door Rule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1/rules/rule1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cdn` - 2024-09-01
