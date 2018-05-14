---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_endpoint"
sidebar_current: "docs-azurerm-resource-cdn-endpoint"
description: |-
  Manages a CDN Endpoint.

---

# azurerm_cdn_endpoint

A CDN Endpoint is the entity within a CDN Profile containing configuration information regarding caching behaviors and origins. The CDN Endpoint is exposed using the URL format <endpointname>.azureedge.net by default, but custom domains can also be created.

## Example Usage

```hcl
resource "random_id" "server" {
  keepers = {
    azi_id = 1
  }

  byte_length = 8
}

resource "azurerm_resource_group" "test" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_cdn_profile" "test" {
  name                = "exampleCdnProfile"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "${random_id.server.hex}"
  profile_name        = "${azurerm_cdn_profile.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  origin {
    name      = "exampleCdnOrigin"
    host_name = "www.example.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the CDN Endpoint. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the CDN Endpoint.

* `profile_name` - (Required) The CDN Profile to which to attach the CDN Endpoint.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `is_http_allowed` - (Optional) Defaults to `true`.

* `is_https_allowed` - (Optional) Defaults to `true`.

* `content_types_to_compress` - (Optional) An array of strings that indicates a content types on which compression will be applied. The value for the elements should be MIME types.

* `geo_filter` - (Optional) A set of Geo Filters for this CDN Endpoint. Each `geo_filter` block supports fields documented below.

* `is_compression_enabled` - (Optional) Indicates whether compression is to be enabled. Defaults to false.

* `querystring_caching_behaviour` - (Optional) Sets query string caching behavior. Allowed values are `IgnoreQueryString`, `BypassCaching` and `UseQueryString`. Defaults to `IgnoreQueryString`.

* `optimization_type` - (Optional) What types of optimization should this CDN Endpoint optimize for? Possible values include `DynamicSiteAcceleration`, `GeneralMediaStreaming`, `GeneralWebDelivery`, `LargeFileDownload` and `VideoOnDemandMediaStreaming`.

* `origin` - (Optional) The set of origins of the CDN endpoint. When multiple origins exist, the first origin will be used as primary and rest will be used as failover options. Each `origin` block supports fields documented below.

* `origin_host_header` - (Optional) The host header CDN provider will send along with content requests to origins. Defaults to the host name of the origin.

* `origin_path` - (Optional) The path used at for origin requests.

* `probe_path` - (Optional) the path to a file hosted on the origin which helps accelerate delivery of the dynamic content and calculate the most optimal routes for the CDN. This is relative to the `origin_path`.

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

## Attributes Reference

The following attributes are exported:

* `id` - The CDN Endpoint ID.

## Import

CDN Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_endpoint.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/myprofile1/endpoints/myendpoint1
```
