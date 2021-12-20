---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_origin_group"
description: |-
  Manages an Azure Front Door (Standard/Premium) instance. (currently in public preview)
---

# azurerm_cdn_frontdoor_origin_group

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Origin Group.

* `profile_id` - (Required) Azure Front Door Profile ID.

* `session_affinity_state` - Whether to allow session affinity on this host. Can be set to `true` or `false`.

---

The `load_balancing` block supports the following:

* `sample_size` - The number of samples to consider for load balancing decisions.

* `successful_samples_required` - The number of samples within the sample period that must succeed.

* `additional_latency_in_ms` - The additional latency in milliseconds for probes to fall into the lowest latency bucket.

---

The `health_probe` block supports the following:

* `path` - The path relative to the origin that is used to determine the health of the origin.

* `request_type` - The type of health probe request that is made. Can be set to `GET`, `HEAD` or `NotSet`. Defaults to `NotSet`.

* `protocol` - Protocol to use for health probe. Can be set to `Http` or `Https`. Defaults to `Http`.

* `interval_in_seconds` - The number of seconds between health probes. Defaults to `240` seconds.
