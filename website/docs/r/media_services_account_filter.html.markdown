---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_services_account_filter"
description: |-
  Manages a Media Services Account Filter.
---

# azurerm_media_services_account_filter

Manages a Media Services Account Filter.

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

resource "azurerm_media_services_account_filter" "example" {
  name                        = "Filter1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  first_quality_bitrate       = 128000

  presentation_time_range {
    start_in_units                 = 0
    end_in_units                   = 15
    presentation_window_in_units   = 90
    live_backoff_in_units          = 0
    unit_timescale_in_milliseconds = 1000
    force_end                      = false
  }

  track_selection {
    condition {
      property  = "Type"
      operation = "Equal"
      value     = "Audio"
    }

    condition {
      property  = "Language"
      operation = "NotEqual"
      value     = "en"
    }

    condition {
      property  = "FourCC"
      operation = "NotEqual"
      value     = "EC-3"
    }
  }


  track_selection {
    condition {
      property  = "Type"
      operation = "Equal"
      value     = "Video"
    }

    condition {
      property  = "Bitrate"
      operation = "Equal"
      value     = "3000000-5000000"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Account Filter. Changing this forces a new Account Filter to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Account Filter should exist. Changing this forces a new Account Filter to be created.

* `media_services_account_name` - (Required) The Media Services account name. Changing this forces a new Account Filter to be created.

---

* `first_quality_bitrate` - (Optional) The first quality bitrate. Sets the first video track to appear in the Live Streaming playlist to allow HLS native players to start downloading from this quality level at the beginning.

* `presentation_time_range` - (Optional) A `presentation_time_range` block as defined below.

* `track_selection` - (Optional) One or more `track_selection` blocks as defined below.

---

A `presentation_time_range` block supports the following:

* `unit_timescale_in_milliseconds` - (Required) Specified as the number of milliseconds in one unit timescale. For example, if you want to set a `start_in_units` at 30 seconds, you would use a value of 30 when using the `unit_timescale_in_milliseconds` in 1000. Or if you want to set `start_in_units` in 30 milliseconds, you would use a value of 30 when using the `unit_timescale_in_milliseconds` in 1. Applies timescale to `start_in_units`, `start_timescale` and `presentation_window_in_timescale` and `live_backoff_in_timescale`.
 
* `end_in_units` - (Optional) The absolute end time boundary. Applies to Video on Demand (VoD).
For the Live Streaming presentation, it is silently ignored and applied when the presentation ends and the stream becomes VoD. This is a long value that represents an absolute end point of the presentation, rounded to the closest next GOP start. The unit is defined by `unit_timescale_in_milliseconds`, so an `end_in_units` of 180 would be for 3 minutes. Use `start_in_units` and `end_in_units` to trim the fragments that will be in the playlist (manifest). For example, `start_in_units` set to 20 and `end_in_units` set to 60 using `unit_timescale_in_milliseconds` in 1000 will generate a playlist that contains fragments from between 20 seconds and 60 seconds of the VoD presentation. If a fragment straddles the boundary, the entire fragment will be included in the manifest.

* `force_end` - (Optional) Indicates whether the `end_in_units` property must be present. If true, `end_in_units` must be specified or a bad request code is returned. Applies to Live Streaming only. Allowed values: `false`, `true`.

* `live_backoff_in_units` - (Optional) The relative to end right edge. Applies to Live Streaming only.
This value defines the latest live position that a client can seek to. Using this property, you can delay live playback position and create a server-side buffer for players. The unit is defined by `unit_timescale_in_milliseconds`. The maximum live back off duration is 300 seconds. For example, a value of 20 means that the latest available content is 20 seconds delayed from the real live edge.

* `presentation_window_in_units` - (Optional) The relative to end sliding window. Applies to Live Streaming only. Use `presentation_window_in_units` to apply a sliding window of fragments to include in a playlist. The unit is defined by `unit_timescale_in_milliseconds`. For example, set `presentation_window_in_units` to 120 to apply a two-minute sliding window. Media within 2 minutes of the live edge will be included in the playlist. If a fragment straddles the boundary, the entire fragment will be included in the playlist. The minimum presentation window duration is 60 seconds.

* `start_in_units` - (Optional) The absolute start time boundary. Applies to Video on Demand (VoD) or Live Streaming. This is a long value that represents an absolute start point of the stream. The value gets rounded to the closest next GOP start. The unit is defined by `unit_timescale_in_milliseconds`, so a `start_in_units` of 15 would be for 15 seconds. Use `start_in_units` and `end_in_units` to trim the fragments that will be in the playlist (manifest). For example, `start_in_units` set to 20 and `end_in_units` set to 60 using `unit_timescale_in_milliseconds` in 1000 will generate a playlist that contains fragments from between 20 seconds and 60 seconds of the VoD presentation. If a fragment straddles the boundary, the entire fragment will be included in the manifest.

---

A `selection` block supports the following:

* `operation` - (Required) The condition operation to test a track property against. Supported values are `Equal` and `NotEqual`.

* `property` - (Required) The track property to compare. Supported values are `Bitrate`, `FourCC`, `Language`, `Name` and `Type`. Check [documentation](https://docs.microsoft.com/azure/media-services/latest/filters-concept) for more details.

* `value` - (Required) The track property value to match or not match.

---

A `track_selection` block supports the following:

* `condition` - (Required) One or more `selection` blocks as defined above.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Account Filter.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Account Filter.
* `read` - (Defaults to 5 minutes) Used when retrieving the Account Filter.
* `update` - (Defaults to 30 minutes) Used when updating the Account Filter.
* `delete` - (Defaults to 30 minutes) Used when deleting the Account Filter.

## Import

Account Filters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_services_account_filter.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaServices/account1/accountFilters/filter1
```
