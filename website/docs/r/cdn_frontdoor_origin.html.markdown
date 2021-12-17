---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_origin"
description: |-
  Manages an Azure Front Door (Standard/Premium) instance. (currently in public preview)
---

# azurerm_cdn_frontdoor_origin

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

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name       = "afdorigingroup1"
  profile_id = azurerm_cdn_frontdoor_profile.example.id

  health_probe {
    protocol = "Http"
  }

  load_balancing {
    sample_size                 = 4
    successful_samples_required = 2
  }
}

resource "azurerm_cdn_frontdoor_origin" "example" {
  name            = "afdorigin2"
  origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id
  host_name       = "lmo.com"
  weight          = 14
  priority        = 2
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the origin.

* `origin_group_id` - (Required) Resource ID of the origin group.

* `enabled` - Can be set to `true` or `false`.

* `origin_host_header` - The host header value sent to the origin with each request.

* `host_name` - (Required) The address of the origin. Domain names, IPv4 addresses, and IPv6 addresses are supported.This should be unique across all origins in an endpoint.

* `priority` - Priority of origin in given origin group for load balancing. Higher priorities will not be used for load balancing if any lower priority origin is healthy.Must be between 1 and 5

* `weight` - Weight of the origin in given origin group for load balancing. Must be between 1 and 1000.

* `http_port` - The value of the HTTP port. Must be between 1 and 65535. Defaults to `80`.

* `https_port` - The value of the HTTPS port. Must be between 1 and 65535. Defaults to `443`.

---

The `private_link` block supports the following:

* `alias` -  The Alias of the Private Link resource. Populating this optional field indicates that this origin is 'Private'

* `approval_message` - A custom message to be included in the approval request to connect to the Private Link.

* `location` - The location of the Private Link resource. Required only if `resource_id` is populated.

* `resource_id` - The Resource Id of the Private Link resource. Populating this optional field indicates that this backend is 'Private'.
