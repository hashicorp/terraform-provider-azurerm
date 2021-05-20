---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_live_event"
description: |-
  Manages a Live Event.
---

# azurerm_media_live_event

Manages a Live Event.

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

resource "azurerm_media_live_event" "example" {
  name                        = "example"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  media_services_account_name = azurerm_media_services_account.example.name
  description                 = "My Event Description"

  input {
    streaming_protocol = "RTMP"
    ip_access_control_allow {
      name                 = "AllowAll"
      address              = "0.0.0.0"
      subnet_prefix_length = 0
    }
  }

  encoding {
    type               = "Standard"
    preset_name        = "Default720p"
    stretch_mode       = "AutoFit"
    key_frame_interval = "PT2S"
  }

  preview {
    ip_access_control_allow {
      name                 = "AllowAll"
      address              = "0.0.0.0"
      subnet_prefix_length = 0
    }
  }

  use_static_hostname     = true
  hostname_prefix         = "special-event"
  transcription_languages = ["en-US"]
}
```

## Arguments Reference

The following arguments are supported:

* `input` - (Required) A `input` block as defined below.

* `location` - (Required) The Azure Region where the Live Event should exist. Changing this forces a new Live Event to be created.

* `media_services_account_name` - (Required)  The Media Services account name. Changing this forces a new Live Event to be created.

* `name` - (Required) The name which should be used for this Live Event. Changing this forces a new Live Event to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Live Event should exist. Changing this forces a new Live Event to be created.

---

* `auto_start_enabled` - (Optional) The flag indicates if the resource should be automatically started on creation. Default is `false`.

* `cross_site_access_policy` - (Optional) A `cross_site_access_policy` block as defined below.

* `description` - (Optional) A description for the live event.

* `encoding` - (Optional) A `encoding` block as defined below.

* `hostname_prefix` - (Optional) When `use_static_hostname` is set to true, the `hostname_prefix` specifies the first part of the hostname assigned to the live event preview and ingest endpoints. The final hostname would be a combination of this prefix, the media service account name and a short code for the Azure Media Services data center.

* `preview` - (Optional) A `preview` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Live Event.

* `transcription_languages` - (Optional) Specifies a list of languages (locale) to be used for speech-to-text transcription – it should match the spoken language in the audio track. The value should be in `BCP-47` format (e.g: `en-US`). [See the Microsoft Documentation for more information about the live transcription feature and the list of supported languages](https://go.microsoft.com/fwlink/?linkid=2133742 ).

* `use_static_hostname` - (Optional) Specifies whether a static hostname would be assigned to the live event preview and ingest endpoints. Changing this forces a new Live Event to be created.
---

A `cross_site_access_policy` block supports the following:

* `client_access_policy` - (Optional) The content of clientaccesspolicy.xml used by Silverlight.

* `cross_domain_policy` - (Optional) The content of the Cross Domain Policy (`crossdomain.xml`).

---

A `encoding` block supports the following:

* `key_frame_interval` - (Optional) Use an `ISO 8601` time value between 0.5 to 20 seconds to specify the output fragment length for the video and audio tracks of an encoding live event. For example, use `PT2S` to indicate 2 seconds. For the video track it also defines the key frame interval, or the length of a GoP (group of pictures). If this value is not set for an encoding live event, the fragment duration defaults to 2 seconds. The value cannot be set for pass-through live events.

* `preset_name` - (Optional) The optional encoding preset name, used when `type` is not `None`. If the `type` is set to `Standard`, then the default preset name is `Default720p`. Else if the `type` is set to `Premium1080p`, the default preset is `Default1080p`. Changing this forces a new resource to be created.

* `stretch_mode` - (Optional) Specifies how the input video will be resized to fit the desired output resolution(s). Allowed values are `None`, `AutoFit` or `AutoSize`. Default is `None`.

* `type` - (Optional) Live event type. Allowed values are `None`, `Premium1080p` or `Standard`. When set to `None`, the service simply passes through the incoming video and audio layer(s) to the output. When `type` is set to `Standard` or `Premium1080p`, a live encoder transcodes the incoming stream into multiple bitrates or layers. Defaults to `None`. Changing this forces a new resource to be created.

-> [More information can be found in the Microsoft Documentation](https://go.microsoft.com/fwlink/?linkid=2095101).

---

A `input` block supports the following:

* `access_token` - (Optional) A UUID in string form to uniquely identify the stream. If omitted, the service will generate a unique value. Changing this forces a new value to be created.

* `ip_access_control_allow` - (Optional) One or more `ip_access_control_allow` blocks as defined below.

* `key_frame_interval_duration` - (Optional) ISO 8601 time duration of the key frame interval duration of the input. This value sets the `EXT-X-TARGETDURATION` property in the HLS output. For example, use PT2S to indicate 2 seconds. This field cannot be set when `type` is set to `Encoding`.

* `streaming_protocol` - (Optional) The input protocol for the live event. Allowed values are `FragmentedMP4` and `RTMP`. Changing this forces a new resource to be created.

---

A `ip_access_control_allow` block supports the following:

* `address` - (Optional) The IP address or CIDR range.

* `name` - (Optional) The friendly name for the IP address range.

* `subnet_prefix_length` - (Optional) The subnet mask prefix length (see CIDR notation).

---

A `preview` block supports the following:

* `alternative_media_id` - (Optional) An alternative media identifier associated with the streaming locator created for the preview. The identifier can be used in the `CustomLicenseAcquisitionUrlTemplate` or the `CustomKeyAcquisitionUrlTemplate` of the Streaming Policy specified in the `streaming_policy_name` field. Changing this forces a new resource to be created.

* `ip_access_control_allow` - (Optional) One or more `ip_access_control_allow` blocks as defined above.

* `preview_locator` - (Optional) The identifier of the preview locator in Guid format. Specifying this at creation time allows the caller to know the preview locator url before the event is created. If omitted, the service will generate a random identifier. Changing this forces a new resource to be created.

* `streaming_policy_name` - (Optional) The name of streaming policy used for the live event preview. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Live Event.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Live Event.
* `read` - (Defaults to 5 minutes) Used when retrieving the Live Event.
* `update` - (Defaults to 30 minutes) Used when updating the Live Event.
* `delete` - (Defaults to 30 minutes) Used when deleting the Live Event.

## Import

Live Events can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_live_event.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Media/mediaservices/account1/liveevents/event1
```
