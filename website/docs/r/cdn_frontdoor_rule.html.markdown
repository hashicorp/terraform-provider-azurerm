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

* `name` - (Required) The name which should be used for this Front Door Rule. Possible values must be between 1 and 260 characters in length, begin with a letter and may contain only letters and numbers. Changing this forces a new Front Door Rule to be created.

* `cdn_frontdoor_rule_set_id` - (Required) The resource ID of the Front Door Rule Set for this Front Door Rule. Changing this forces a new Front Door Rule to be created.

* `order` - (Required) The order in which the rules will be applied for the Front Door Endpoint. The order value should be sequential and begin at `1`(e.g. `1`, `2`, `3`...). A Front Door Rule with a lesser order value will be applied before a rule with a greater order value.

-> **Note:** If the Front Door Rule has an order value of `0` they do not require any conditions and the actions will always be applied.

* `actions` - (Required) An `actions` block as defined below.

* `behavior_on_match` - (Optional) If this rule is a match should the rules engine continue processing the remaining rules or stop? Possible values are `Continue` and `Stop`. Defaults to `Continue`.

* `conditions` - (Optional) A `conditions` block as defined below.

---

An `actions` block supports the following:

-> **Note:** You may include up to 5 separate actions in the `actions` block.

Some actions support `Action Server Variables` which provide access to structured information about the request. For more information about `Action Server Variables` see the `Action Server Variables` as defined below.

* `url_rewrite_action` - (Optional) A `url_rewrite_action` block as defined below. You may **not** have a `url_rewrite_action` **and** a `url_redirect_action` defined in the same `actions` block.

* `url_redirect_action` - (Optional) A `url_redirect_action` block as defined below. You may **not** have a `url_redirect_action` **and** a `url_rewrite_action` defined in the same `actions` block.

* `route_configuration_override_action` - (Optional) A `route_configuration_override_action` block as defined below.

* `request_header_action` - (Optional) A `request_header_action` block as defined below.

* `response_header_action` - (Optional) A `response_header_action` block as defined below.

---

An `url_redirect_action` block supports the following:

* `redirect_type` - (Required) The response type to return to the requestor. Possible values include `Moved`, `Found` , `TemporaryRedirect` or `PermanentRedirect`.

* `destination_hostname` - (Required) The host name you want the request to be redirected to. The value must be a string between `0` and `2048` characters in length, leave blank to preserve the incoming host.

* `redirect_protocol` - (Optional) The protocol the request will be redirected as. Possible values include `MatchRequest`, `Http` or `Https`. Defaults to `MatchRequest`.

* `destination_path` - (Optional) The path to use in the redirect. The value must be a string and include the leading `/`, leave blank to preserve the incoming path. Defaults to `""`.

* `query_string` - (Optional) The query string used in the redirect URL. The value must be in the &lt;key>=&lt;value> or &lt;key>={`action_server_variable`} format and must not include the leading `?`, leave blank to preserve the incoming query string. Maximum allowed length for this field is `2048` characters. Defaults to `""`.

* `destination_fragment` - (Optional) The fragment to use in the redirect. The value must be a string between `0` and `1024` characters in length, leave blank to preserve the incoming fragment. Defaults to `""`.

---

A `route_configuration_override_action` block supports the following:

-> **Note:** In the v3.x of the provider the `cache_duration`, `cache_behavior` and `query_string_caching_behavior` will have default values. You can use Terraform's [ignore_changes](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#ignore_changes) functionality to ignore these default values. In v4.0 of the provider the `cache_duration`, `cache_behavior` and `query_string_caching_behavior` will **NOT** have default values and will need to be explicitly set in the configuration file.

* `cache_duration` - (Optional) When Cache behavior is set to `Override` or `SetIfMissing`, this field specifies the cache duration to use. The maximum duration is 366 days specified in the `d.HH:MM:SS` format(e.g. `365.23:59:59`). If the desired maximum cache duration is less than 1 day then the maximum cache duration should be specified in the `HH:MM:SS` format(e.g. `23:59:59`).

* `cdn_frontdoor_origin_group_id` - (Optional) The Front Door Origin Group resource ID that the request should be routed to. This overrides the configuration specified in the Front Door Endpoint route.

* `forwarding_protocol` - (Optional) The forwarding protocol the request will be redirected as. This overrides the configuration specified in the route to be associated with. Possible values include `MatchRequest`, `HttpOnly` or `HttpsOnly`.

-> **Note:** If the `cdn_frontdoor_origin_group_id` is not defined you cannot set the `forwarding_protocol`.

* `query_string_caching_behavior` - (Optional) `IncludeSpecifiedQueryStrings` query strings specified in the `query_string_parameters` field get included when the cache key gets generated. `UseQueryString` cache every unique URL, each unique URL will have its own cache key. `IgnoreSpecifiedQueryStrings` query strings specified in the `query_string_parameters` field get excluded when the cache key gets generated. `IgnoreQueryString` query strings aren't considered when the cache key gets generated. Possible values include `IgnoreQueryString`, `UseQueryString`, `IgnoreSpecifiedQueryStrings` or `IncludeSpecifiedQueryStrings`.

* `query_string_parameters` - (Optional) A list of query string parameter names.

-> **Note:** `query_string_parameters` is a required field when the `query_string_caching_behavior` is set to `IncludeSpecifiedQueryStrings` or `IgnoreSpecifiedQueryStrings`.

* `compression_enabled` - (Optional) Should the Front Door dynamically compress the content? Possible values include `true` or `false`.

-> **Note:** Content won't be compressed on AzureFrontDoor when requested content is smaller than `1 byte` or larger than `1 MB`.

* `cache_behavior` - (Optional) `HonorOrigin` the Front Door will always honor origin response header directive. If the origin directive is missing, Front Door will cache contents anywhere from `1` to `3` days. `OverrideAlways` the TTL value returned from your Front Door Origin is overwritten with the value specified in the action. This behavior will only be applied if the response is cacheable. `OverrideIfOriginMissing` if no TTL value gets returned from your Front Door Origin, the rule sets the TTL to the value specified in the action. This behavior will only be applied if the response is cacheable. `Disabled` the Front Door will not cache the response contents, irrespective of Front Door Origin response directives. Possible values include `HonorOrigin`, `OverrideAlways`, `OverrideIfOriginMissing` or `Disabled`.

---

An `url_rewrite_action` block supports the following:

* `source_pattern` - (Required) The source pattern in the URL path to replace. This uses prefix-based matching. For example, to match all URL paths use a forward slash `"/"` as the source pattern value.

* `destination` - (Required) The destination path to use in the rewrite. The destination path overwrites the source pattern.

* `preserve_unmatched_path` - (Optional) Append the remaining path after the source pattern to the new destination path? Possible values `true` or `false`. Defaults to `false`.

---

A `request_header_action` block supports the following:

* `header_action` - (Required) The action to be taken on the specified `header_name`. Possible values include `Append`, `Overwrite` or `Delete`.

-> **Note:** `Append` causes the specified header to be added to the request with the specified value. If the header is already present, the value is appended to the existing header value using string concatenation. No delimiters are added. `Overwrite` causes specified header to be added to the request with the specified value. If the header is already present, the specified value overwrites the existing value. `Delete` causes the header to be deleted from the request.

* `header_name` - (Required) The name of the header to modify.

* `value` - (Optional) The value to append or overwrite.

-> **Note:** `value` is required if the `header_action` is set to `Append` or `Overwrite`.

---

A `response_header_action` block supports the following:

* `header_action` - (Required) The action to be taken on the specified `header_name`. Possible values include `Append`, `Overwrite` or `Delete`.

-> **Note:** `Append` causes the specified header to be added to the request with the specified value. If the header is already present, the value is appended to the existing header value using string concatenation. No delimiters are added. `Overwrite` causes specified header to be added to the request with the specified value. If the header is already present, the specified value overwrites the existing value. `Delete` causes the header to be deleted from the request.

* `header_name` - (Required) The name of the header to modify.

* `value` - (Optional) The value to append or overwrite.

-> **Note:** `value` is required if the `header_action` is set to `Append` or `Overwrite`.

---

A `conditions` block supports the following:

-> **Note:** You may include up to 10 separate conditions in the `conditions` block.

* `remote_address_condition` - (Optional) A `remote_address_condition` block as defined below.

* `request_method_condition` - (Optional) A `request_method_condition` block as defined below.

* `query_string_condition` - (Optional) A `query_string_condition` block as defined below.

* `post_args_condition` - (Optional) A `post_args_condition` block as defined below.

* `request_uri_condition` - (Optional) A `request_uri_condition` block as defined below.

* `request_header_condition` - (Optional) A `request_header_condition` block as defined below.

* `request_body_condition` - (Optional) A `request_body_condition` block as defined below.

* `request_scheme_condition` - (Optional) A `request_scheme_condition` block as defined below.

* `url_path_condition` - (Optional) A `url_path_condition` block as defined below.

* `url_file_extension_condition` - (Optional) A `url_file_extension_condition` block as defined below.

* `url_filename_condition` - (Optional) A `url_filename_condition` block as defined below.

* `http_version_condition` - (Optional) A `http_version_condition` block as defined below.

* `cookies_condition` - (Optional) A `cookies_condition` block as defined below.

* `is_device_condition` - (Optional) A `is_device_condition` block as defined below.

* `socket_address_condition` - (Optional) A `socket_address_condition` block as defined below.

* `client_port_condition` - (Optional) A `client_port_condition` block as defined below.

* `server_port_condition` - (Optional) A `server_port_condition` block as defined below.

* `host_name_condition` - (Optional) A `host_name_condition` block as defined below.

* `ssl_protocol_condition` - (Optional) A `ssl_protocol_condition` block as defined below.

---

A `ssl_protocol_condition` block supports the following:

-> **Note:** The `ssl_protocol_condition` identifies requests based on the SSL protocol of an established TLS connection.

* `match_values` - (Required) A list of one or more HTTP methods. Possible values are `TLSv1`, `TLSv1.1` and `TLSv1.2` logic.

* `operator` - (Optional) Possible value `Equal`. Defaults to `Equal`.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

---

A `host_name_condition` block supports the following:

-> **Note:** The `host_name_condition` identifies requests based on the specified hostname in the request from client.

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) A list of one or more string values representing the value of the request hostname to match. If multiple values are specified, they're evaluated using `OR` logic.

* `transforms` - (Optional) A Conditional operator. Possible values include `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` or `UrlEncode`. Details can be found in the `Condition Transform List` below.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

---

A `server_port_condition` block supports the following:

-> **Note:** The `server_port_condition` identifies requests based on which port of the Front Door server accepted the request on.

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Required) A list of one or more integer values(e.g. "1") representing the value of the client port to match. Possible values include `80` or `443`. If multiple values are specified, they're evaluated using `OR` logic.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

---

A `client_port_condition` block supports the following:

-> **Note:** The `client_port_condition` identifies requests based on the port of the client which made the request.

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) One or more integer values(e.g. "1") representing the value of the client port to match. If multiple values are specified, they're evaluated using `OR` logic.

---

A `socket_address_condition` block supports the following:

-> **Note:** The `socket_address_condition` identifies requests based on the IP address of the direct connection to the Front Door Profiles edge. If the client used an HTTP proxy or a load balancer to send the request, the value of Socket address is the IP address of the proxy or load balancer.

-> **Note:** Remote Address represents the original client IP that is either from the network connection or typically the `X-Forwarded-For` request header if the user is behind a proxy.

* `operator` - (Optional) The type of match. The Possible values are `IpMatch` or `Any`. Defaults to `IPMatch`.

-> **Note:** If the value of the `operator` field is set to `IpMatch` then the `match_values` field is also required.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) Specify one or more IP address ranges. If multiple IP address ranges are specified, they're evaluated using `OR` logic.

-> **Note:** See the `Specifying IP Address Ranges` section below on how to correctly define the `match_values` field.

---

A `remote_address_condition` block supports the following:

-> **Note:** Remote Address represents the original client IP that is either from the network connection or typically the `X-Forwarded-For` request header if the user is behind a proxy.

* `operator` - (Optional) The type of the remote address to match. Possible values include `Any`, `GeoMatch` or `IPMatch`. Use the `negate_condition` to specify Not `GeoMatch` or Not `IPMatch`. Defaults to `IPMatch`.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) For the IP Match or IP Not Match operators: specify one or more IP address ranges. If multiple IP address ranges are specified, they're evaluated using `OR` logic. For the Geo Match or Geo Not Match operators: specify one or more locations using their country code.

-> **Note:** See the `Specifying IP Address Ranges` section below on how to correctly define the `match_values` field.

---

A `request_method_condition` block supports the following:

-> **Note:** The `request_method_condition` identifies requests that use the specified HTTP request method.

* `match_values` - (Required) A list of one or more HTTP methods. Possible values include `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `OPTIONS` or `TRACE`. If multiple values are specified, they're evaluated using `OR` logic.

* `operator` - (Optional) Possible value `Equal`. Defaults to `Equal`.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

---

A `query_string_condition` block supports the following:

-> **Note:** Use the `query_string_condition` to identify requests that contain a specific query string.

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) One or more string or integer values(e.g. "1") representing the value of the query string to match. If multiple values are specified, they're evaluated using `OR` logic.

* `transforms` - (Optional) A Conditional operator. Possible values include `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` or `UrlEncode`. Details can be found in the `Condition Transform List` below.

---

A `post_args_condition` block supports the following:

-> **Note:** Use the `post_args_condition` to identify requests based on the arguments provided within a `POST` request's body. A single match condition matches a single argument from the `POST` request's body.

* `post_args_name` - (Required) A string value representing the name of the `POST` argument.

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) One or more string or integer values(e.g. "1") representing the value of the `POST` argument to match. If multiple values are specified, they're evaluated using `OR` logic.

* `transforms` - (Optional) A Conditional operator. Possible values include `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` or `UrlEncode`. Details can be found in the `Condition Transform List` below.

---

A `request_uri_condition` block supports the following:

-> **Note:** The `request_uri_condition` identifies requests that match the specified URL. The entire URL is evaluated, including the protocol and query string, but not the fragment. When you use this rule condition, be sure to include the protocol(e.g. For example, use `https://www.contoso.com` instead of just `www.contoso.com`).

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) One or more string or integer values(e.g. "1") representing the value of the request URL to match. If multiple values are specified, they're evaluated using `OR` logic.

* `transforms` - (Optional) A Conditional operator. Possible values include `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` or `UrlEncode`. Details can be found in the `Condition Transform List` below.

---

A `request_header_condition` block supports the following:

-> **Note:** The `request_header_condition` identifies requests that include a specific header in the request. You can use this match condition to check if a header exists whatever its value, or to check if the header matches a specified value.

* `header_name` - (Required) A string value representing the name of the `POST` argument.

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) One or more string or integer values(e.g. "1") representing the value of the request header to match. If multiple values are specified, they're evaluated using `OR` logic.

* `transforms` - (Optional) A Conditional operator. Possible values include `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` or `UrlEncode`. Details can be found in the `Condition Transform List` below.

---

A `request_body_condition` block supports the following:

-> **Note:** The `request_body_condition` identifies requests based on specific text that appears in the body of the request.

-> **Note:** If a request body exceeds `64 KB` in size, only the first `64 KB` will be considered for the request body match condition.

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Required) A list of one or more string or integer values(e.g. "1") representing the value of the request body text to match. If multiple values are specified, they're evaluated using `OR` logic.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `transforms` - (Optional) A Conditional operator. Possible values include `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` or `UrlEncode`. Details can be found in the `Condition Transform List` below.

---

A `request_scheme_condition` block supports the following:

-> **Note:** The `request_scheme_condition` identifies requests that use the specified protocol.

* `operator` - (Optional) Possible value `Equal`. Defaults to `Equal`.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) The requests protocol to match. Possible values include `HTTP` or `HTTPS`.

---

An `url_path_condition` block supports the following:

-> **Note:** The `url_path_condition` identifies requests that include the specified path in the request URL. The path is the part of the URL after the hostname and a slash(e.g. in the URL `https://www.contoso.com/files/secure/file1.pdf`, the path is `files/secure/file1.pdf`).

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) One or more string or integer values(e.g. "1") representing the value of the request path to match. Don't include the leading slash (`/`). If multiple values are specified, they're evaluated using `OR` logic.

* `transforms` - (Optional) A Conditional operator. Possible values include `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` or `UrlEncode`. Details can be found in the `Condition Transform List` below.

---

An `url_file_extension_condition` block supports the following:

-> **Note:** The `url_file_extension_condition` identifies requests that include the specified file extension in the file name in the request URL. Don't include a leading period(e.g. use `html` instead of `.html`).

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Required) A list of one or more string or integer values(e.g. "1") representing the value of the request file extension to match. If multiple values are specified, they're evaluated using `OR` logic.

* `transforms` - (Optional) A Conditional operator. Possible values include `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` or `UrlEncode`. Details can be found in the `Condition Transform List` below.

---

An `url_filename_condition` block supports the following:

-> **Note:** The `url_filename_condition` identifies requests that include the specified file name in the request URL.

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) A list of one or more string or integer values(e.g. "1") representing the value of the request file name to match. If multiple values are specified, they're evaluated using `OR` logic.

-> **Note:** The `match_values` field is only optional if the `operator` is set to `Any`.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `transforms` - (Optional) A Conditional operator. Possible values include `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` or `UrlEncode`. Details can be found in the `Condition Transform List` below.

---

A `http_version_condition` block supports the following:

-> **Note:** Use the HTTP version match condition to identify requests that have been made by using a specific version of the HTTP protocol.

* `match_values` - (Required) What HTTP version should this condition match? Possible values `2.0`, `1.1`, `1.0` or `0.9`.

* `operator` - (Optional) Possible value `Equal`. Defaults to `Equal`.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

---

A `cookies_condition` block supports the following:

-> **Note:** Use the `cookies_condition` to identify requests that have include a specific cookie.

* `cookie_name` - (Required) A string value representing the name of the cookie.

* `operator` - (Required) A Conditional operator. Possible values include `Any`, `Equal`, `Contains`, `BeginsWith`, `EndsWith`, `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual` or `RegEx`. Details can be found in the `Condition Operator List` below.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) One or more string or integer values(e.g. "1") representing the value of the request header to match. If multiple values are specified, they're evaluated using `OR` logic.

* `transforms` - (Optional) A Conditional operator. Possible values include `Lowercase`, `RemoveNulls`, `Trim`, `Uppercase`, `UrlDecode` or `UrlEncode`. Details can be found in the `Condition Transform List` below.

---

An `is_device_condition` block supports the following:

-> **Note:** Use the `is_device_condition` to identify requests that have been made from a `mobile` or `desktop` device.

* `operator` - (Optional) Possible value `Equal`. Defaults to `Equal`.

* `negate_condition` - (Optional) If `true` operator becomes the opposite of its value. Possible values `true` or `false`. Defaults to `false`. Details can be found in the `Condition Operator List` below.

* `match_values` - (Optional) Which device should this rule match on? Possible values `Mobile` or `Desktop`.

---

## Specifying IP Address Ranges

When specifying IP address ranges in the `socket_address_condition` and the `remote_address_condition` `match_values` use the following format:

Use `CIDR` notation when specifying IP address blocks. This means that the syntax for an IP address block is the base IP address followed by a forward slash and the prefix size For example:

* `IPv4` example: `5.5.5.64/26` matches any requests that arrive from addresses `5.5.5.64` through `5.5.5.127`.
* `IPv6` example: `1:2:3:/48` matches any requests that arrive from addresses `1:2:3:0:0:0:0:0` through `1:2:3:ffff:ffff:ffff:ffff:ffff`.

When you specify multiple IP addresses and IP address blocks, `OR` logic is applied.

* `IPv4` example: if you add two IP addresses `1.2.3.4` and `10.20.30.40`, the condition is matched for any requests that arrive from either address `1.2.3.4` or `10.20.30.40`.
* `IPv6` example: if you add two IP addresses `1:2:3:4:5:6:7:8` and `10:20:30:40:50:60:70:80`, the condition is matched for any requests that arrive from either address `1:2:3:4:5:6:7:8` or `10:20:30:40:50:60:70:80`.

---

## Action Server Variables

Rule Set server variables provide access to structured information about the request. You can use server variables to dynamically change the request/response headers or URL rewrite paths/query strings, for example, when a new page load or when a form is posted.

### Supported Action Server Variables

| Variable name | Description |
|---------------|-------------|
| `socket_ip`      | The IP address of the direct connection to Front Door Profiles edge. If the client used an HTTP proxy or a load balancer to send the request, the value of `socket_ip` is the IP address of the proxy or load balancer. |
| `client_ip`      | The IP address of the client that made the original request. If there was an `X-Forwarded-For` header in the request, then the client IP address is picked from the header. |
| `client_port`    | The IP port of the client that made the request. |
| `hostname`       | The host name in the request from the client. |
| `geo_country`    | Indicates the requester's country/region of origin through its country/region code. |
| `http_method`    | The method used to make the URL request, such as `GET` or `POST`. |
| `http_version`   | The request protocol. Usually `HTTP/1.0`, `HTTP/1.1`, or `HTTP/2.0`. |
| `query_string`   | The list of variable/value pairs that follows the "?" in the requested URL. For example, in the request `http://contoso.com:8080/article.aspx?id=123&title=fabrikam`, the `query_string` value will be `id=123&title=fabrikam`. |
| `request_scheme` | The request scheme: `http` or `https`. |
| `request_uri`    | The full original request URI (with arguments). For example, in the request `http://contoso.com:8080/article.aspx?id=123&title=fabrikam`, the `request_uri` value will be `/article.aspx?id=123&title=fabrikam`. |
| `ssl_protocol`   | The protocol of an established TLS connection. |
| `server_port`    | The port of the server that accepted a request. |
| `url_path`       | Identifies the specific resource in the host that the web client wants to access. This is the part of the request URI without the arguments. For example, in the request `http://contoso.com:8080/article.aspx?id=123&title=fabrikam`, the `uri_path` value will be `/article.aspx`. |

### Action Server Variable Format

Server variables can be specified using the following formats:

* `{variable}` - Include the entire server variable. For example, if the client IP address is `111.222.333.444` then the `{client_ip}` token would evaluate to `111.222.333.444`.

* `{variable:offset}` - Include the server variable after a specific offset, until the end of the variable. The offset is zero-based. For example, if the client IP address is `111.222.333.444` then the `{client_ip:3}` token would evaluate to `.222.333.444`.

* `{variable:offset:length}` - Include the server variable after a specific offset, up to the specified length. The offset is zero-based. For example, if the client IP address is `111.222.333.444` then the `{client_ip:4:3}` token would evaluate to `222`.

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

| Operator                   | Description | Condition Value |
|----------------------------|-------------|-----------------|
| Any                        |Matches when there is any value, regardless of what it is. | Any |
| Equal                      | Matches when the value exactly matches the specified string. | Equal |
| Contains                   | Matches when the value contains the specified string. | Contains |
| Less Than                  | Matches when the length of the value is less than the specified integer. | LessThan |
| Greater Than               | Matches when the length of the value is greater than the specified integer. | GreaterThan |
| Less Than or Equal         | Matches when the length of the value is less than or equal to the specified integer. | LessThanOrEqual |
| Greater Than or Equal      | Matches when the length of the value is greater than or equal to the specified integer. | GreaterThanOrEqual |
| Begins With                | Matches when the value begins with the specified string. | BeginsWith |
| Ends With                  | Matches when the value ends with the specified string. | EndsWith |
| RegEx                      | Matches when the value matches the specified regular expression. See below for further details. | RegEx |
| Not Any                    | Matches when there is no value. | Any and negateCondition = true |
| Not Equal                  | Matches when the value does not match the specified string. | Equal and negateCondition : true |
| Not Contains               | Matches when the value does not contain the specified string. | Contains and negateCondition = true |
| Not Less Than              | Matches when the length of the value is not less than the specified integer. | LessThan and negateCondition = true |
| Not Greater Than           | Matches when the length of the value is not greater than the specified integer. | GreaterThan and negateCondition = true |
| Not Less Than or Equal     | Matches when the length of the value is not less than or equal to the specified integer. | LessThanOrEqual and negateCondition = true |
| Not Greater Than or Equals | Matches when the length of the value is not greater than or equal to the specified integer. | GreaterThanOrEqual and negateCondition = true |
| Not Begins With            | Matches when the value does not begin with the specified string. | BeginsWith and negateCondition = true |
| Not Ends With              | Matches when the value does not end with the specified string. | EndsWith and negateCondition = true |
| Not RegEx                  | Matches when the value does not match the specified regular expression. See `Condition Regular Expressions` for further details. | RegEx and negateCondition = true |

---

## Condition Regular Expressions

Regular expressions **don't** support the following operations:

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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Front Door Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Rule.

## Import

Front Door Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1/rules/rule1
```
