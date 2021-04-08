---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_live_event_output"
description: |-
  Manages an Azure Media Live Event Output.
---

# azurerm_media_live_event_output

Manages a Azure Media Live Event Output.

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
  name                        = "inputAsset"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
}

resource "azurerm_media_live_event" "example" {
  name                        = "exampleevent"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  media_services_account_name = azurerm_media_services_account.example.name
  description                 = "My Event Description"

  input {
    streaming_protocol          = "RTMP"
    key_frame_interval_duration = "PT6S"
    ip_access_control_allow {
      name                 = "AllowAll"
      address              = "0.0.0.0"
      subnet_prefix_length = 0
    }
  }
}

resource "azurerm_media_live_event_output" "example" {
  name                         = "exampleoutput"
  live_event_id                = azurerm_media_live_event.example.id
  archive_window_length        = "PT5M"
  asset_name                   = azurerm_media_asset.example.name
  description                  = "Test live output 1"
  manifest_name                = "testmanifest"
  output_snap_time_in_seconds  = 0
  hls_fragments_per_ts_segment = 5
}
```

## Arguments Reference

The following arguments are supported:

* `archive_window_duration` - (Required) `ISO 8601` time between 1 minute to 25 hours to indicate the maximum content length that can be archived in the asset for this live output. This also sets the maximum content length for the rewind window. For example, use `PT1H30M` to indicate 1 hour and 30 minutes of archive window. Changing this forces a new Live Output to be created.

* `asset_name` - (Required) The asset that the live output will write to. Changing this forces a new Live Output to be created.

* `live_event_id` - (Required) The id of the live event. Changing this forces a new Live Output to be created.

* `name` - (Required) The name which should be used for this Live Event Output. Changing this forces a new Live Output to be created.

---

* `description` - (Optional) The description of the live output. Changing this forces a new Live Output to be created.

* `hls_fragments_per_ts_segment` - (Optional) The number of fragments in an HTTP Live Streaming (HLS) TS segment in the output of the live event. This value does not affect the packing ratio for HLS CMAF output. Changing this forces a new Live Output to be created.

* `manifest_name` - (Optional) The manifest file name. If not provided, the service will generate one automatically. Changing this forces a new Live Output to be created.

* `output_snap_timestamp_in_seconds` - (Optional) The initial timestamp that the live output will start at, any content before this value will not be archived. Changing this forces a new Live Output to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Live Output.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Live Output.
* `read` - (Defaults to 5 minutes) Used when retrieving the Live Output.
* `update` - (Defaults to 30 minutes) Used when updating the Live Output.
* `delete` - (Defaults to 30 minutes) Used when deleting the Live Output.

## Import

Live Outputs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_live_output.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaservices/account1/liveevents/event1/liveoutputs/output1
```
