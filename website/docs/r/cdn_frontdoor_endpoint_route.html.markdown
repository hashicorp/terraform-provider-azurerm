
---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_endpoint_route"
description: |-
  Manages an Azure Front Door (Standard/Premium) instance. (currently in public preview)
---

# azurerm_cdn_frontdoor_endpoint_route

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}


resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "afdpremv2"
  location            = "global"
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Premium_AzureFrontDoor"
}


resource "azurerm_cdn_frontdoor_endpoint" "example1" {
  name       = "afdendpoint1"
  profile_id = azurerm_cdn_frontdoor_profile.example.id
  enabled    = true
}

resource "azurerm_cdn_frontdoor_endpoint_route" "example" {
  name        = "example-route"
  endpoint_id = azurerm_cdn_frontdoor_endpoint.example1.id
  enabled     = true

  origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id

  query_string_caching_behavior = "NotSet2"

  origin_path = "/*"

  forwarding_protocol = "MatchRequest"

  supported_protocols = [
    "Https",
    "Http"
  ]

  link_to_default_domain = false

  custom_domains = [
    azurerm_cdn_frontdoor_custom_domain.example1.id
  ]
}
````

## Argument Reference

The following arguments are supported:

* `name` - (Required) Endpoint name.

* `enabled` - (Required) Can be set to `true` or `false`.

* `endpoint_id` - (Required) Refers to the Front Door endpoint.

* `custom_domains` - List of custom domains (resource id).

* `origin_group_id` - (Required) Refers to the Front Door origin group (resource id).

* `origin_path` - A directory path on the origin that AzureFrontDoor can use to retrieve content from, e.g. `contoso.cloudapp.net/originpath`.

* `rule_sets`

* `supported_protocols` - List of supported protocols. Can be set to `Http` and/or `Https`.

* `patterns_to_match`

* `forwarding_protocol` - Protocol this rule will use when forwarding traffic to backends. Can be set to `HttpsOnly`, `HttpOnly` or `MatchRequest`. Defaults to `MatchRequest`.

* `link_to_default_domain` - Whether this route will be linked to the default endpoint domain. Can be set to `true` or `false`.

* `https_redirect` - Whether to automatically redirect HTTP traffic to HTTPS traffic. Note that this is a easy way to set up this rule and it will be the first rule that gets executed. Can be set to `Enabled` or `Disabled`.

* `enable_caching` - Indicates whether content compression is enabled on Azure Front Door. Defaults to `false`. If compression is enabled, content will be served as compressed if user requests for a compressed version. Content won't be compressed on Azure Front Door when requested content is smaller than 1 byte or larger than 1 MB.

* `content_types_to_compress` - List of MIME content types to compress. I.e. `application/json`.

* `query_string_caching_behavior` - Can be set to `NotSet`, `UseQueryString` or `IgnoreQueryString`. Defaults to `NotSet`.
