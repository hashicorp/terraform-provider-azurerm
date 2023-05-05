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
  name                      = "existing-cdn-route-name"
  cdn_frontdoor_endpoint_id = "existing-cdn-endpoint-id"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Front Door Route.

* `cdn_frontdoor_endpoint_id` - (Required) Specifies The Resource ID of the Front Door Endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Front Door Route.

* `name` - Specifies The name of the Front Door Route.

* `cdn_frontdoor_endpoint_id` - Specifies the resource ID of the Front Door Endpoint.

* `cdn_frontdoor_origin_group_id` - Specifies The resource ID of the Front Door Origin Group.

* `forwarding_protocol` - Specifies the Protocol that will be use when forwarding traffic to backends.

* `patterns_to_match` - Specifies the route patterns of the rule.

* `supported_protocols` - Specifies the supported Protocols by this Front Door Route.

* `cache` - A `cache` block as defined below.

* `cdn_frontdoor_custom_domain_ids` - Specifies the IDs of the Front Door Custom Domains which are associated with this Front Door Route.

* `cdn_frontdoor_origin_path` - Specifies a directory path on the Front Door Origin that can be used to retrieve content.

* `cdn_frontdoor_rule_set_ids` - Specifies a list of the Front Door Rule Set IDs which are assigned to this Front Door Route.

* `enabled` - Specifies if this Front Door Route enabled or not.

* `https_redirect_enabled` - Specifies if automatically redirect HTTP traffic to HTTPS traffic are enabled or not.

* `link_to_default_domain` - Specifies if this Front Door Route is linked to the default endpoint or not.

---

A `cache` block supports the following:

* `query_string_caching_behavior` - Specifies the Front Door Routes query string behavior.

* `query_strings` - Specifies the Front Door Routes Query strings.

* `compression_enabled` - Specifies if content compression is enabled or not?

* `content_types_to_compress` - Specifies a list of `Content types` that are being compressed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Route.
