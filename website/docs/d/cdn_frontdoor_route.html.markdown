---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_route"
description: |-
  Gets information about an existing Front Door (standard/premium) Route.
---

# Data Source: azurerm_cdn_frontdoor_route

Use this data source to access information about an existing Front Door (standard/premium) Route.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_route" "example" {
  cdn_frontdoor_route_id = "existing-cdn-route-id"
}
```

## Argument Reference

The following arguments are supported:

* `cdn_frontdoor_route_id` - Specifies The Resource ID of the existing Front Door Route.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Front Door Route.

* `name` - Specifies The name which should be used for this Front Door Route. Valid values must begin with a letter or number, end with a letter or number and may only contain letters, numbers and hyphens with a maximum length of 90 characters. Changing this forces a new Front Door Route to be created.

* `cdn_frontdoor_endpoint_id` - Specifies The resource ID of the Front Door Endpoint where this Front Door Route should exist. Changing this forces a new Front Door Route to be created.

* `cdn_frontdoor_origin_group_id` - Specifies The resource ID of the Front Door Origin Group where this Front Door Route should be created.

* `forwarding_protocol` - Specifies The Protocol that will be use when forwarding traffic to backends. Possible values are `HttpOnly`, `HttpsOnly` or `MatchRequest`.

* `patterns_to_match` - Specifies The route patterns of the rule.

* `supported_protocols` - Specifies One or more Protocols supported by this Front Door Route. Possible values are `Http` or `Https`.

* `cache` - A `cache` block as defined below.

* `cdn_frontdoor_custom_domain_ids` - Specifies The IDs of the Front Door Custom Domains which are associated with this Front Door Route.

* `cdn_frontdoor_origin_path` - Specifies A directory path on the Front Door Origin that can be used to retrieve content (e.g. `contoso.cloudapp.net/originpath`).

* `cdn_frontdoor_rule_set_ids` - Specifies A list of the Front Door Rule Set IDs which should be assigned to this Front Door Route.

* `enabled` - Specifies Is this Front Door Route enabled? Possible values are `true` or `false`. Defaults to `true`.

* `https_redirect_enabled` - Specifies Automatically redirect HTTP traffic to HTTPS traffic? Possible values are `true` or `false`. Defaults to `true`.

* `link_to_default_domain` - Specifies Should this Front Door Route be linked to the default endpoint? Possible values include `true` or `false`. Defaults to `true`.

---

A `cache` block supports the following:

* `query_string_caching_behavior` - Specifies Defines how the Front Door Route will cache requests that include query strings. Possible values include `IgnoreQueryString`, `IgnoreSpecifiedQueryStrings`, `IncludeSpecifiedQueryStrings` or `UseQueryString`. Defaults it `IgnoreQueryString`.

* `query_strings` - Specifies Query strings to include or ignore.

* `compression_enabled` - Specifies Is content compression enabled? Possible values are `true` or `false`. Defaults to `false`.

* `content_types_to_compress` - Specifies A list of one or more `Content types` (formerly known as `MIME types`) to compress.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Route.
