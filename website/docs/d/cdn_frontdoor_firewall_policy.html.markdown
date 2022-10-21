---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_firewall_policy"
description: |-
  Gets information about an existing CDN Front Door Firewall Policy.
---

# Data Source: azurerm_cdn_frontdoor_firewall_policy

Gets information about an existing CDN Front Door Firewall Policy.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_firewall_policy" "example" {
  name                = "examplecdnfdwafpolicy"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the CDN FrontDoor Firewall Policy.

* `resource_group_name` - (Required) The name of the resource group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CDN FrontDoor Firewall Policy.

* `enabled` - The enabled state of the CDN FrontDoor Firewall Policy.

* `frontend_endpoint_ids` - The CDN Frontend Endpoints associated with this CDN FrontDoor Firewall Policy.

* `mode` - The CDN FrontDoor Firewall Policy mode.

* `redirect_url` - The redirect URL for the client.

* `sku_name` - The sku's pricing tier for this CDN FrontDoor Firewall Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the CDN FrontDoor Firewall Policy.
