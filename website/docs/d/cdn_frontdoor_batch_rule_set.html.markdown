---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_cdn_frontdoor_batch_rule_set"
description: |-
  Gets information about an existing Front Door (standard/premium) Batch Rule Set.
---

# Data Source: azurerm_cdn_frontdoor_batch_rule_set

Gets information about an existing Front Door (standard/premium) Batch Rule Set.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_batch_rule_set" "example" {
  name                = "existingbatchruleset"
  profile_name        = "existing-profile"
  resource_group_name = "existing-resources"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Front Door Batch Rule Set.

* `profile_name` - (Required) The name of the Front Door Profile.

* `resource_group_name` - (Required) The name of the Resource Group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Batch Rule Set.

* `batch_mode_enabled` - Whether `batch mode` is enabled.

* `cdn_frontdoor_profile_id` - The ID of the Front Door Profile associated with this Front Door Batch Rule Set.

* `rules` - A `rules` block as defined below.

---

A `rules` block exports the following:

* `actions` - An `actions` block as defined below.

* `behavior_on_match` - Whether the rules engine continues processing after this rule matches.

* `conditions` - A `conditions` block as defined below.

* `name` - The name of the Front Door Batch Rule.

* `order` - The order in which this rule is applied.

---

An `actions` block exports the following:

* `request_header_action` - A `request_header_action` block as defined below.

* `response_header_action` - A `response_header_action` block as defined below.

* `route_configuration_override_action` - A `route_configuration_override_action` block as defined below.

* `url_redirect_action` - An `url_redirect_action` block as defined below.

* `url_rewrite_action` - An `url_rewrite_action` block as defined below.

---

A `request_header_action` block exports the following:

* `header_action` - The action applied to the request header.

* `header_name` - The name of the request header.

* `value` - The value associated with the request header action.

---

A `response_header_action` block exports the following:

* `header_action` - The action applied to the response header.

* `header_name` - The name of the response header.

* `value` - The value associated with the response header action.

---

A `route_configuration_override_action` block exports the following:

* `cache_behavior` - The cache behavior applied to this action.

* `cache_duration` - The cache duration applied to this action.

* `cdn_frontdoor_origin_group_id` - The ID of the Front Door Origin Group associated with this action.

* `compression_enabled` - Whether `compression` is enabled.

* `forwarding_protocol` - The forwarding protocol applied to this action.

* `query_string_caching_behavior` - The query string caching behavior applied to this action.

* `query_string_parameters` - The query string parameters associated with this action.

---

An `url_redirect_action` block exports the following:

* `destination_fragment` - The destination fragment for the redirect action.

* `destination_hostname` - The destination host name for the redirect action.

* `destination_path` - The destination path for the redirect action.

* `query_string` - The query string for the redirect action.

* `redirect_protocol` - The redirect protocol for the redirect action.

* `redirect_type` - The redirect type for the redirect action.

---

An `url_rewrite_action` block exports the following:

* `destination` - The destination path for the rewrite action.

* `preserve_unmatched_path` - Whether `preserve unmatched path` is enabled.

* `source_pattern` - The source pattern for the rewrite action.

---

A `conditions` block exports the following:

* `client_port_condition` - A `client_port_condition` block as defined below.
* `cookies_condition` - A `cookies_condition` block as defined below.
* `host_name_condition` - A `host_name_condition` block as defined below.
* `http_version_condition` - A `http_version_condition` block as defined below.
* `is_device_condition` - An `is_device_condition` block as defined below.
* `post_args_condition` - A `post_args_condition` block as defined below.
* `query_string_condition` - A `query_string_condition` block as defined below.
* `remote_address_condition` - A `remote_address_condition` block as defined below.
* `request_body_condition` - A `request_body_condition` block as defined below.
* `request_header_condition` - A `request_header_condition` block as defined below.
* `request_method_condition` - A `request_method_condition` block as defined below.
* `request_scheme_condition` - A `request_scheme_condition` block as defined below.
* `request_uri_condition` - A `request_uri_condition` block as defined below.
* `server_port_condition` - A `server_port_condition` block as defined below.
* `socket_address_condition` - A `socket_address_condition` block as defined below.
* `ssl_protocol_condition` - A `ssl_protocol_condition` block as defined below.
* `url_file_extension_condition` - An `url_file_extension_condition` block as defined below.
* `url_filename_condition` - An `url_filename_condition` block as defined below.
* `url_path_condition` - An `url_path_condition` block as defined below.

---

A `client_port_condition` block exports the following:

* `match_values` - The client port values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

---

A `cookies_condition` block exports the following:

* `cookie_name` - The cookie name associated with this condition.

* `match_values` - The cookie values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

---

A `host_name_condition` block exports the following:

* `match_values` - The host names associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

---

A `http_version_condition` block exports the following:

* `match_values` - The HTTP versions associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

---

An `is_device_condition` block exports the following:

* `match_values` - The device types associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

---

A `post_args_condition` block exports the following:

* `match_values` - The POST argument values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

* `post_args_name` - The POST argument name associated with this condition.

* `transforms` - The transforms associated with this condition.

---

A `query_string_condition` block exports the following:

* `match_values` - The query string values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

---

A `remote_address_condition` block exports the following:

* `match_values` - The remote address values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

---

A `request_body_condition` block exports the following:

* `match_values` - The request body values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

---

A `request_header_condition` block exports the following:

* `header_name` - The request header name associated with this condition.

* `match_values` - The request header values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

---

A `request_method_condition` block exports the following:

* `match_values` - The request methods associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

---

A `request_scheme_condition` block exports the following:

* `match_values` - The request schemes associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

---

A `request_uri_condition` block exports the following:

* `match_values` - The request URIs associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

---

A `server_port_condition` block exports the following:

* `match_values` - The server port values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

---

A `socket_address_condition` block exports the following:

* `match_values` - The socket address values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

---

A `ssl_protocol_condition` block exports the following:

* `match_values` - The SSL protocol values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

---

An `url_file_extension_condition` block exports the following:

* `match_values` - The URL file extension values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

---

An `url_filename_condition` block exports the following:

* `match_values` - The URL file name values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

---

An `url_path_condition` block exports the following:

* `match_values` - The URL path values associated with this condition.

* `negate_condition` - Whether `negate condition` is enabled.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Batch Rule Set.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Cdn` - 2025-12-01
