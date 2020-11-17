---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_endpoint"
description: |-
  Manages a CDN Endpoint.
---

# azurerm_cdn_endpoint

A CDN Endpoint is the entity within a CDN Profile containing configuration information regarding caching behaviours and origins. The CDN Endpoint is exposed using the URL format <endpointname>.azureedge.net.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_profile" "example" {
  name                = "example-cdn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "example" {
  name                = "example"
  profile_name        = azurerm_cdn_profile.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  origin {
    name      = "example"
    host_name = "www.contoso.com"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the CDN Endpoint. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the CDN Endpoint.

* `profile_name` - (Required) The CDN Profile to which to attach the CDN Endpoint.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `is_http_allowed` - (Optional) Defaults to `true`.

* `is_https_allowed` - (Optional) Defaults to `true`.

* `content_types_to_compress` - (Optional) An array of strings that indicates a content types on which compression will be applied. The value for the elements should be MIME types.

* `geo_filter` - (Optional) A set of Geo Filters for this CDN Endpoint. Each `geo_filter` block supports fields documented below.

* `is_compression_enabled` - (Optional) Indicates whether compression is to be enabled.

* `querystring_caching_behaviour` - (Optional) Sets query string caching behavior. Allowed values are `IgnoreQueryString`, `BypassCaching` and `UseQueryString`. `NotSet` value can be used for `Premium Verizon` CDN profile. Defaults to `IgnoreQueryString`.

* `optimization_type` - (Optional) What types of optimization should this CDN Endpoint optimize for? Possible values include `DynamicSiteAcceleration`, `GeneralMediaStreaming`, `GeneralWebDelivery`, `LargeFileDownload` and `VideoOnDemandMediaStreaming`.

* `origin` - (Required) The set of origins of the CDN endpoint. When multiple origins exist, the first origin will be used as primary and rest will be used as failover options. Each `origin` block supports fields documented below.

* `origin_host_header` - (Optional) The host header CDN provider will send along with content requests to origins. Defaults to the host name of the origin.

* `origin_path` - (Optional) The path used at for origin requests.

* `probe_path` - (Optional) the path to a file hosted on the origin which helps accelerate delivery of the dynamic content and calculate the most optimal routes for the CDN. This is relative to the `origin_path`.

-> **NOTE:** `global_delivery_rule` and `delivery_rule` are currently only available for `Microsoft_Standard` CDN profiles.

* `global_delivery_rule` - (Optional) Actions that are valid for all resources regardless of any conditions. A `global_delivery_rule` block as defined below.

* `delivery_rule` - (Optional) Rules for the rules engine. An endpoint can contain up until 4 of those rules that consist of conditions and actions. A `delivery_rule` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `origin` block supports:

* `name` - (Required) The name of the origin. This is an arbitrary value. However, this value needs to be unique under the endpoint. Changing this forces a new resource to be created.

* `host_name` - (Required) A string that determines the hostname/IP address of the origin server. This string can be a domain name, Storage Account endpoint, Web App endpoint, IPv4 address or IPv6 address. Changing this forces a new resource to be created.

* `http_port` - (Optional) The HTTP port of the origin. Defaults to `80`. Changing this forces a new resource to be created.

* `https_port` - (Optional) The HTTPS port of the origin. Defaults to `443`. Changing this forces a new resource to be created.

The `geo_filter` block supports:

* `relative_path` - (Required) The relative path applicable to geo filter.

* `action` - (Required) The Action of the Geo Filter. Possible values include `Allow` and `Block`.

* `country_codes` - (Required) A List of two letter country codes (e.g. `US`, `GB`) to be associated with this Geo Filter.

---

A `global_delivery_rule` block supports the following:

* `cache_expiration_action` - (Optional) A `cache_expiration_action` block as defined above.

* `cache_key_query_string_action` - (Optional) A `cache_key_query_string_action` block as defined above.

* `modify_request_header_action` - (Optional) A `modify_request_header_action` block as defined below.

* `modify_response_header_action` - (Optional) A `modify_response_header_action` block as defined below.

* `url_redirect_action` - (Optional) A `url_redirect_action` block as defined below.

* `url_rewrite_action` - (Optional) A `url_rewrite_action` block as defined below.

---

A `delivery_rule` block supports the following:

* `name` - (Required) The Name which should be used for this Delivery Rule.

* `order` - (Required) The order used for this rule, which must be larger than 1.

* `cache_expiration_action` - (Optional) A `cache_expiration_action` block as defined above.

* `cache_key_query_string_action` - (Optional) A `cache_key_query_string_action` block as defined above.

* `cookies_condition` - (Optional) A `cookies_condition` block as defined above.

* `device_condition` - (Optional) A `device_condition` block as defined below.

* `http_version_condition` - (Optional) A `http_version_condition` block as defined below.

* `modify_request_header_action` - (Optional) A `modify_request_header_action` block as defined below.

* `modify_response_header_action` - (Optional) A `modify_response_header_action` block as defined below.

* `post_arg_condition` - (Optional) A `post_arg_condition` block as defined below.

* `query_string_condition` - (Optional) A `query_string_condition` block as defined below.

* `remote_address_condition` - (Optional) A `remote_address_condition` block as defined below.

* `request_body_condition` - (Optional) A `request_body_condition` block as defined below.

* `request_header_condition` - (Optional) A `request_header_condition` block as defined below.

* `request_method_condition` - (Optional) A `request_method_condition` block as defined below.

* `request_scheme_condition` - (Optional) A `request_scheme_condition` block as defined below.

* `request_uri_condition` - (Optional) A `request_uri_condition` block as defined below.

* `url_file_extension_condition` - (Optional) A `url_file_extension_condition` block as defined below.

* `url_file_name_condition` - (Optional) A `url_file_name_condition` block as defined below.

* `url_path_condition` - (Optional) A `url_path_condition` block as defined below.

* `url_redirect_action` - (Optional) A `url_redirect_action` block as defined below.

* `url_rewrite_action` - (Optional) A `url_rewrite_action` block as defined below.

---

A `cache_expiration_action` block supports the following:

* `behavior` - (Required) The behavior of the cache. Valid values are `BypassCache`, `Override` and `SetIfMissing`.

* `duration` - (Optional) Duration of the cache. Only allowed when `behavior` is set to `Override` or `SetIfMissing`. Format: `[d.]hh:mm:ss`

---

A `cache_key_query_string_action` block supports the following:

* `behavior` - (Required) The behavior of the cache key for query strings. Valid values are `Exclude`, `ExcludeAll`, `Include` and `IncludeAll`.

* `parameters` - (Optional) Comma separated list of parameter values.

---

A `modify_request_header_action` block supports the following:

* `action` - (Required) Action to be executed on a header value. Valid values are `Append`, `Delete` and `Overwrite`.

* `name` - (Required) The header name.

* `value` - (Optional) The value of the header. Only needed when `action` is set to `Append` or `overwrite`.

---

A `modify_response_header_action` block supports the following:

* `action` - (Required) Action to be executed on a header value. Valid values are `Append`, `Delete` and `Overwrite`.

* `name` - (Required) The header name.

* `value` - (Optional) The value of the header. Only needed when `action` is set to `Append` or `overwrite`.

---

A `url_redirect_action` block supports the following:

* `redirect_type` - (Required) Type of the redirect. Valid values are `Found`, `Moved`, `PermanentRedirect` and `TemporaryRedirect`.

* `protocol` - (Optional) Specifies the protocol part of the URL. Valid values are `Http` and `Https`.

* `hostname` - (Optional) Specifies the hostname part of the URL.

* `path` - (Optional) Specifies the path part of the URL. This value must begin with a `/`.

* `fragment` - (Optional) Specifies the fragment part of the URL. This value must not start with a `#`.

* `query_string` - (Optional) Specifies the query string part of the URL. This value must not start with a `?` or `&` and must be in `<key>=<value>` format separated by `&`.

---

A `url_rewrite_action` block supports the following:

* `source_pattern` - (Required) This value must start with a `/` and can't be longer than 260 characters.

* `destination` - (Required) This value must start with a `/` and can't be longer than 260 characters.

* `preserve_unmatched_path` - (Optional) Defaults to `true`.

---

A `cookies_condition` block supports the following:

* `selector` - (Required) Name of the cookie.

* `operator` - (Required) Valid values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) List of values for the cookie.

* `transforms` - (Optional) Valid values are `Lowercase` and `Uppercase`.

---

A `device_condition` block supports the following:

* `operator` - (Optional) Valid values are `Equal`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) Valid values are `Desktop` and `Mobile`.

---

A `http_version_condition` block supports the following:

* `operator` - (Optional) Valid values are `Equal`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) Valid values are `0.9`, `1.0`, `1.1` and `2.0`.

---

A `post_arg_condition` block supports the following:

* `selector` - (Required) Name of the post arg.

* `operator` - (Required) Valid values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) List of string values.

* `transforms` - (Optional) Valid values are `Lowercase` and `Uppercase`.

---

A `query_string_condition` block supports the following:

* `operator` - (Required) Valid values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) List of string values.

* `transforms` - (Optional) Valid values are `Lowercase` and `Uppercase`.

---

A `remote_address_condition` block supports the following:

* `operator` - (Required) Valid values are `Any`, `GeoMatch` and `IPMatch`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) List of string values. For `GeoMatch` `operator` this should be a list of country codes (e.g. `US` or `DE`). List of IP address if `operator` equals to `IPMatch`.

---

A `request_body_condition` block supports the following:

* `operator` - (Required) Valid values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) List of string values.

* `transforms` - (Optional) Valid values are `Lowercase` and `Uppercase`.

---

A `request_header_condition` block supports the following:

* `selector` - (Required) Header name.

* `operator` - (Required) Valid values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) List of header values.

* `transforms` - (Optional) Valid values are `Lowercase` and `Uppercase`.

---

A `request_method_condition` block supports the following:

* `operator` - (Optional) Valid values are `Equal`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) Valid values are `DELETE`, `GET`, `HEAD`, `OPTIONS`, `POST` and `PUT`.

---

A `request_scheme_condition` block supports the following:

* `operator` - (Optional) Valid values are `Equal`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) Valid values are `HTTP` and `HTTPS`.

---

A `request_uri_condition` block supports the following:

* `operator` - (Required) Valid values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) List of string values.

* `transforms` - (Optional) Valid values are `Lowercase` and `Uppercase`.

---

A `url_file_extension_condition` block supports the following:

* `operator` - (Required) Valid values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) List of string values.

* `transforms` - (Optional) Valid values are `Lowercase` and `Uppercase`.

---

A `url_file_name_condition` block supports the following:

* `operator` - (Required) Valid values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) List of string values.

* `transforms` - (Optional) Valid values are `Lowercase` and `Uppercase`.

---

A `url_path_condition` block supports the following:

* `operator` - (Required) Valid values are `Any`, `BeginsWith`, `Contains`, `EndsWith`, `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.

* `negate_condition` - (Optional) Defaults to `false`.

* `match_values` - (Required) List of string values.

* `transforms` - (Optional) Valid values are `Lowercase` and `Uppercase`.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CDN Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CDN Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the CDN Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the CDN Endpoint.

## Import

CDN Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/myprofile1/endpoints/myendpoint1
```
