---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_cdn_frontdoor_batch_rule_set"
description: |-
  Gets information about an existing Front Door (standard/premium) Batch Rule Set.
---

# Data Source: azurerm_cdn_frontdoor_batch_rule_set

Gets information about an existing Front Door (standard/premium) Batch Rule Set.

~> **Note:** This data source can only read Rule Sets that were provisioned in batch mode. Use the `azurerm_cdn_frontdoor_rule_set` data source for Rule Sets that were not provisioned in batch mode.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_batch_rule_set" "example" {
  name                = "existing"
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

* `cdn_frontdoor_profile_id` - The ID of the Front Door Profile associated with this Front Door Batch Rule Set.

* `rules` - A `rules` block as defined below.

---

A `rules` block exports the following:

* `actions` - An `actions` block as defined below.

* `behaviour_on_match` - Whether the rules engine continues processing after this rule matches.

* `conditions` - A `conditions` block as defined below.

* `name` - The name of the Front Door Batch Rule.

* `order` - The order in which this rule is applied.

---

An `actions` block exports the following:

* `modify_request_header` - One or more `modify_request_header` blocks as defined below.

* `modify_response_header` - One or more `modify_response_header` blocks as defined below.

* `route_configuration_override` - A `route_configuration_override` block as defined below.

* `url_redirect` - A `url_redirect` block as defined below.

* `url_rewrite` - A `url_rewrite` block as defined below.

---

A `modify_request_header` block exports the following:

* `header_name` - The name of the request header.

* `header_value` - The value associated with the request header action.

* `operator` - The action applied to the request header.

---

A `modify_response_header` block exports the following:

* `header_name` - The name of the response header.

* `header_value` - The value associated with the response header action.

* `operator` - The action applied to the response header.

---

A `route_configuration_override` block exports the following:

* `caching` - A `caching` block as defined below.

* `origin_group` - An `origin_group` block as defined below.

---

An `origin_group` block exports the following:

* `cdn_frontdoor_origin_group_id` - The ID of the Front Door Origin Group associated with this action.

* `forwarding_protocol` - The forwarding protocol applied to this action.

---

A `caching` block exports the following:

* `behaviour` - The cache behaviour applied to this action.

* `compression_enabled` - Whether compression is enabled.

* `duration` - The cache duration applied to this action.

* `query_string_behaviour` - The query string caching behaviour applied to this action.

* `query_string_parameters` - The query string parameters associated with this action.

---

A `url_redirect` block exports the following:

* `destination_fragment` - The destination fragment for the redirect action.

* `destination_host_name` - The destination host name for the redirect action.

* `destination_path` - The destination path for the redirect action.

* `query_string` - The query string for the redirect action.

* `redirect_protocol` - The redirect protocol for the redirect action.

* `redirect_type` - The redirect type for the redirect action.

---

A `url_rewrite` block exports the following:

* `destination_path` - The destination path for the rewrite action.

* `preserve_unmatched_path_enabled` - Whether to preserve the unmatched part of the path.

* `source_pattern` - The source pattern for the rewrite action.

---

A `conditions` block exports the following:

* `client_port` - One or more `client_port` blocks as defined below.

* `device_type` - One or more `device_type` blocks as defined below.

* `host_name` - One or more `host_name` blocks as defined below.

* `http_version` - One or more `http_version` blocks as defined below.

* `post_argument` - One or more `post_argument` blocks as defined below.

* `query_string` - One or more `query_string` blocks as defined below.

* `remote_address` - One or more `remote_address` blocks as defined below.

* `request_body` - One or more `request_body` blocks as defined below.

* `request_cookies` - One or more `request_cookies` blocks as defined below.

* `request_file_extension` - One or more `request_file_extension` blocks as defined below.

* `request_filename` - One or more `request_filename` blocks as defined below.

* `request_header` - One or more `request_header` blocks as defined below.

* `request_method` - One or more `request_method` blocks as defined below.

* `request_path` - One or more `request_path` blocks as defined below.

* `request_scheme` - One or more `request_scheme` blocks as defined below.

* `request_url` - One or more `request_url` blocks as defined below.

* `server_port` - One or more `server_port` blocks as defined below.

* `socket_address` - One or more `socket_address` blocks as defined below.

* `ssl_protocol` - One or more `ssl_protocol` blocks as defined below.


---

A `client_port` block exports the following:

* `operator` - The operator for this condition.

* `values` - The client port values associated with this condition.

---

A `device_type` block exports the following:

* `operator` - The operator for this condition.

* `values` - The device types associated with this condition.

---

A `host_name` block exports the following:

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

* `values` - The host names associated with this condition.

---

A `http_version` block exports the following:

* `operator` - The operator for this condition.

* `values` - The HTTP versions associated with this condition.

---

A `post_argument` block exports the following:

* `name` - The POST argument name associated with this condition.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

* `values` - The POST argument values associated with this condition.

---

A `query_string` block exports the following:

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

* `values` - The query string values associated with this condition.

---

A `remote_address` block exports the following:

* `operator` - The operator for this condition.

* `values` - The remote address values associated with this condition.

---

A `request_body` block exports the following:

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

* `values` - The request body values associated with this condition.

---

A `request_cookies` block exports the following:

* `name` - The cookie name associated with this condition.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

* `values` - The cookie values associated with this condition.

---

A `request_file_extension` block exports the following:

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

* `values` - The request file extension values associated with this condition.

---

A `request_filename` block exports the following:

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

* `values` - The request file name values associated with this condition.

---

A `request_header` block exports the following:

* `name` - The request header name associated with this condition.

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

* `values` - The request header values associated with this condition.

---

A `request_method` block exports the following:

* `operator` - The operator for this condition.

* `values` - The request methods associated with this condition.

---

A `request_path` block exports the following:

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

* `values` - The request path values associated with this condition.

---

A `request_scheme` block exports the following:

* `operator` - The operator for this condition.

* `values` - The request schemes associated with this condition.

---

A `request_url` block exports the following:

* `operator` - The operator for this condition.

* `transforms` - The transforms associated with this condition.

* `values` - The request URL values associated with this condition.

---

A `server_port` block exports the following:

* `operator` - The operator for this condition.

* `values` - The server port values associated with this condition.

---

A `socket_address` block exports the following:

* `operator` - The operator for this condition.

* `values` - The socket address values associated with this condition.

---

A `ssl_protocol` block exports the following:

* `operator` - The operator for this condition.

* `values` - The SSL protocol values associated with this condition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Batch Rule Set.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Cdn` - 2025-12-01
