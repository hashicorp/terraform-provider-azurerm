---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_firewall_policy"
description: |-
  Gets information about an existing Front Door (standard/premium) Firewall Policy.
---

# Data Source: azurerm_cdn_frontdoor_firewall_policy

Use this data source to access information about an existing Front Door (standard/premium) Firewall Policy.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_firewall_policy" "example" {
  name                = "examplecdnfdwafpolicy"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Front Door Firewall Policy.

* `resource_group_name` - (Required) The name of the resource group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Front Door Firewall Policy.

* `enabled` - The enabled state of the Front Door Firewall Policy.

* `frontend_endpoint_ids` - The Front Door Profiles frontend endpoints associated with this Front Door Firewall Policy.

* `mode` - The Front Door Firewall Policy mode.

* `redirect_url` - The redirect URL for the client.

* `sku_name` - The sku's pricing tier for this Front Door Firewall Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Firewall Policy.
