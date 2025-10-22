---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_front_door_cache_purge"
description: |-
  Purges the cache on an Azure Front Door Endpoint.
---

# Action: azurerm_cdn_front_door_cache_purge

~> **Note:** `azurerm_cdn_front_door_cache_purge` is in beta. Its interface and behaviour may change as the feature evolves, and breaking changes are possible. It is offered as a technical preview without compatibility guarantees until Terraform 1.14 is generally available.

Purges the cache on an Azure Front Door Endpoint.

## Example Usage

### Basic Usage

```terraform
# ... additional resource config

resource "azurerm_cdn_frontdoor_endpoint" "example" {
  name                     = "example"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
}

resource "terraform_data" "trigger" {
  input = "example-trigger"
  lifecycle {
    action_trigger {
      events  = [after_update]
      actions = [action.azurerm_cdn_front_door_cache_purge.example]
    }
  }
}

action "azurerm_cdn_front_door_cache_purge" "example" {
  config {
    front_door_id = azurerm_cdn_frontdoor_endpoint.test.id
    content_paths = [
      "/images/*"
    ]
  }
}
```

### Custom Domains Usage

```terraform
resource "azurerm_cdn_frontdoor_endpoint" "example" {
  name                     = "example"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
}

resource "terraform_data" "trigger" {
  input = "example-trigger"
  lifecycle {
    action_trigger {
      events  = [after_update]
      actions = [action.azurerm_cdn_front_door_cache_purge.example]
    }
  }
}

action "azurerm_cdn_front_door_cache_purge" "example" {
  config {
    front_door_id = azurerm_cdn_frontdoor_endpoint.test.id
    content_paths = [
      "/*"
    ]
    domains = [
      "examplehost.contoso.com"
    ]
  }
}
```


## Argument Reference

This action supports the following arguments:

* `front_door_endpoint_id` - (Required) The ID of the Front Door Endpoint to purge the cache of.

* `content_paths` - (Required) The paths to purge from the Front Door Endpoint.

* `domains` - (Optional) The Custom Domain names associated with and bound to the Front Door Endpoint to purge for the `content_paths`.

* `timeout` - (Optional) Timeout in seconds to wait for the Front Door Cache Purge action to complete. Defaults to 900 seconds (15 minutes). Purge operations typically take up to 10 minutes.