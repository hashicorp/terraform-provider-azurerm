---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_asset_filter"
description: |-
  Manages an Asset Filter.
---

# azurerm_media_asset_filter

Manages an Asset Filter.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "media-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "example" {
  name                = "examplemediaacc"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  storage_account {
    id         = azurerm_storage_account.example.id
    is_primary = true
  }
}

resource "azurerm_media_asset" "example" {
  name                        = "Asset1"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  description                 = "Asset description"
}

resource "azurerm_media_asset_filter" "example" {
  name                        = "Filter1"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  asset_name                  = azurerm_media_asset.example.name
  first_quality_bitrate       = 128000

  presentation_time_range {
    start_timestamp              = 0
    end_timestamp                = 170000000
    presentation_window_duration = 9223372036854775000
    live_backoff_duration        = 0
    timescale                    = 10000000
    force_end_timestamp          = false
  }

  track {
    selection {
      property  = "Type"
      operation = "Equal"
      value     = "Audio"
    }

    selection {
      property  = "Language"
      operation = "NotEqual"
      value     = "en"
    }

    selection {
      property  = "FourCC"
      operation = "NotEqual"
      value     = "EC-3"
    }
  }


  track {
    selection {
      property  = "Type"
      operation = "Equal"
      value     = "Video"
    }

    selection {
      property  = "Bitrate"
      operation = "Equal"
      value     = "3000000-5000000"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `asset_name` - (Required) The Asset name. Changing this forces a new Asset Filter to be created.

* `media_services_account_name` - (Required) Specifies the name of the Media Services Account. Changing this forces a new Asset Filter to be created.

* `name` - (Required) The name which should be used for this Asset Filter. Changing this forces a new Asset Filter to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Asset Filter should exist. Changing this forces a new Asset Filter to be created.

---

* `first_quality_bitrate` - (Optional) The first quality bitrate.

* `presentation_time_range` - (Optional) A `presentation_time_range` block as defined below.

* `track` - (Optional) One or more `track` blocks as defined below.

---

A `presentation_time_range` block supports the following:

* `end_timestamp` - (Optional) The absolute end time boundary.

* `force_end_timestamp` - (Optional) The indicator of forcing existing of end time stamp.

* `live_backoff_duration` - (Optional) The relative to end right edge.

* `presentation_window_duration` - (Optional) The relative to end sliding window.

* `start_timestamp` - (Optional) The absolute start time boundary.

* `timescale` - (Optional) The time scale of time stamps.

---

A `selection` block supports the following:

* `operation` - (Optional) The track property condition operation. Supported values are `Equal` and `NotEqual`.

* `property` - (Optional) The track property type. Supported values are `Bitrate`, `FourCC`, `Language`, `Name` and `Type`.

* `value` - (Optional) The track property value.

---

A `track` block supports the following:

* `selection` - (Optional) One or more `selection` blocks as defined above.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Asset Filter.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Asset Filter.
* `read` - (Defaults to 5 minutes) Used when retrieving the Asset Filter.
* `update` - (Defaults to 30 minutes) Used when updating the Asset Filter.
* `delete` - (Defaults to 30 minutes) Used when deleting the Asset Filter.

## Import

Asset Filters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_asset_filter.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaservices/account1/assets/asset1/assetFilters/filter1
```
