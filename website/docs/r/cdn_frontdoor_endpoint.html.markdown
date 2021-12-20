---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_endpoint"
description: |-
  Manages an Azure Front Door (Standard/Premium) instance. (currently in public preview)
---

# azurerm_cdn_frontdoor_endpoint

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

resource "azurerm_cdn_frontdoor_endpoint" "example" {
  name       = "afdendpoint1"
  profile_id = azurerm_cdn_frontdoor_profile.example.id
  enabled    = true

  tags = {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Custom Domain. I.e. `mydomain-com`.

* `profile_id` - (Required) Resource ID of the Azure Front Door Profile.

* `enabled` - (Required) Enabled state can be set to `true` or `false`.

* `origin_response_timeout_in_seconds` - Send and receive timeout on forwarding request to the origin. When timeout is reached, the request fails and returns. Must between `16` and `240`. Defaults to `60`.

* `tags` - (Optional) A mapping of tags to assign to the resource.
---

The following attributes are exported:

* `id` - The ID of the FrontDoor.

* `fqdn` - Returns the Front Door Endpoint FQDN i.e. `afdendpoint1.z01.azurefd.net`.