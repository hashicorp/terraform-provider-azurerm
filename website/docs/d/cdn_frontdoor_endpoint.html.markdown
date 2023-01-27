---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_endpoint"
description: |-
  Gets information about an existing Front Door (standard/premium) Endpoint.
---

# Data Source: azurerm_cdn_frontdoor_endpoint

Use this data source to access information about an existing Front Door (standard/premium) Endpoint.

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

* `name` - (Required) Specifies the name of the Front Door Endpoint.

* `profile_name` - (Required) The name of the Front Door Profile within which Front Door Endpoint exists.

* `resource_group_name` - (Required) The name of the Resource Group where the Front Door Profile exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of this Front Door Endpoint.

* `enabled` - Specifies whether this Front Door Endpoint is enabled or not.

* `host_name` - Specifies the host name of the Front Door Endpoint, in the format `{endpointName}.{dnsZone}` (for example, `contoso.azureedge.net`).

* `tags` - Specifies a mapping of Tags assigned to this Front Door Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Endpoint.
