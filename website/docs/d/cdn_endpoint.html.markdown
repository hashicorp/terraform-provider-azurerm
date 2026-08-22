---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_endpoint"
description: |-
  Gets information about an existing CDN Endpoint.
---

# Data Source: azurerm_cdn_endpoint

Use this data source to access information about an existing CDN Endpoint.

## Example Usage

```hcl
data "azurerm_cdn_endpoint" "example" {
  name                = "myfirstcdnendpoint"
  profile_name        = "myfirstcdnprofile"
  resource_group_name = "example-resources"
}

output "cdn_endpoint_fqdn" {
  value = data.azurerm_cdn_endpoint.example.fqdn
}
```

### Using with a DNS Record

```hcl
data "azurerm_cdn_endpoint" "app" {
  name                = "app-cdn-endpoint"
  profile_name        = "app-cdn-profile"
  resource_group_name = "cdn-resources"
}

resource "azurerm_dns_cname_record" "cdn" {
  name                = "cdn"
  zone_name           = azurerm_dns_zone.example.name
  resource_group_name = azurerm_dns_zone.example.resource_group_name
  ttl                 = 300
  record              = data.azurerm_cdn_endpoint.app.fqdn
}
```

## Arguments Reference

* `name` - The name of the CDN Endpoint.

* `profile_name` - The name of the CDN Profile containing this CDN Endpoint.

* `resource_group_name` - The name of the resource group in which the CDN Endpoint exists.

## Attributes Reference

* `id` - The ID of the CDN Endpoint.

* `location` - The Azure Region where the CDN Endpoint exists.

* `fqdn` - The Fully Qualified Domain Name of the CDN Endpoint.

* `is_http_allowed` - Indicates whether HTTP traffic is allowed on the endpoint.

* `is_https_allowed` - Indicates whether HTTPS traffic is allowed on the endpoint.

* `origin_host_header` - The host header the CDN provider sends along with content requests to origins.

* `origin_path` - The path used for origin requests.

* `probe_path` - The path used for health probe requests.

* `querystring_caching_behaviour` - The query string caching behaviour for requests.

* `optimization_type` - The type of optimization used for this CDN Endpoint.

* `is_compression_enabled` - Indicates whether content compression is enabled on the endpoint.

* `content_types_to_compress` - A set of content types that are compressed when served by the CDN Endpoint.

* `geo_filter` - A `geo_filter` block as defined below.

* `origin` - An `origin` block as defined below.

* `tags` - A mapping of tags assigned to the resource.

---

A `geo_filter` block exports the following:

* `relative_path` - The relative path applicable to the geo filter.

* `action` - The action of the geo filter, either `Allow` or `Block`.

* `country_codes` - A list of two letter country codes that the filter applies to.

---

An `origin` block exports the following:

* `name` - The name of the origin.

* `host_name` - The hostname of the origin.

* `http_port` - The HTTP port of the origin.

* `https_port` - The HTTPS port of the origin.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the CDN Endpoint.
