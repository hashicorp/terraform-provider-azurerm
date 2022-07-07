---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_endpoint"
description: |-
  Gets information about an existing CDN FrontDoor Endpoint.
---

# Data Source: azurerm_cdn_frontdoor_endpoint

Use this data source to access information about an existing CDN FrontDoor Endpoint.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_endpoint" "example" {
  name                = "existing-endpoint"
  profile_name        = "existing-cdn-profile"
  resource_group_name = "existing-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the FrontDoor Endpoint.

* `profile_name` - (Required) The name of the FrontDoor Profile within which CDN FrontDoor Endpoint exists.

* `resource_group_name` - (Required) The name of the Resource Group where the CDN FrontDoor Profile exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of this CDN FrontDoor Endpoint.

* `enabled` - Specifies whether this CDN FrontDoor Endpoint is enabled or not.

* `host_name` - Specifies the host name of the CDN FrontDoor Endpoint, in the format `{endpointName}.{dnsZone}` (for example, `contoso.azureedge.net`).

* `tags` - Specifies a mapping of Tags assigned to this CDN FrontDoor Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the CDN FrontDoor Endpoint.
