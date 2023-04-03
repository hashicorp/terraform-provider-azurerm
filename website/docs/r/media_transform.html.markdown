---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_transform"
description: |-
  Manages a Transform.
---

# azurerm_media_transform

Manages a Transform.

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

resource "azurerm_media_transform" "example" {
  name                        = "transform1"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  description                 = "My transform description"
  output {
    relative_priority = "Normal"
    on_error_action   = "ContinueJob"
    builtin_preset {
      preset_name = "AACGoodQualityAudio"
    }
  }
}

```

## Example Usage with Multiple Outputs

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

resource "azurerm_media_transform" "example" {
  name                        = "transform1"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  description                 = "My transform description"
  output {
    relative_priority = "Normal"
    on_error_action   = "ContinueJob"
    builtin_preset {
      preset_name = "AACGoodQualityAudio"
    }
  }

  output {
    relative_priority = "Low"
    on_error_action   = "ContinueJob"
    audio_analyzer_preset {
      audio_language      = "en-US"
      audio_analysis_mode = "Basic"
    }
  }

  output {
    relative_priority = "Low"
    on_error_action   = "StopProcessingJob"
    face_detector_preset {
      analysis_resolution = "StandardDefinition"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `media_services_account_name` - (Required) The Media Services account name. Changing this forces a new Transform to be created.

* `name` - (Required) The name which should be used for this Transform. Changing this forces a new Transform to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Transform should exist. Changing this forces a new Transform to be created.

---

* `description` - (Optional) An optional verbose description of the Transform.

* `output` - (Optional) One or more `output` blocks as defined below. At least one `output` must be defined.

---

A `output` block supports the following:

* `audio_analyzer_preset` - (Optional) An `audio_analyzer_preset` block as defined below.

* `builtin_preset` - (Optional) A `builtin_preset` block as defined below.

* `custom_preset` - (Optional) A `custom_preset` block as defined below.

* `face_detector_preset` - (Optional) A `face_detector_preset` block as defined below.

* `on_error_action` - (Optional) A Transform can define more than one outputs. This property defines what the service should do when one output fails - either continue to produce other outputs, or, stop the other outputs. The overall Job state will not reflect failures of outputs that are specified with `ContinueJob`. Possible values are `StopProcessingJob` or `ContinueJob`. The default is `StopProcessingJob`.

* `relative_priority` - (Optional) Sets the relative priority of the TransformOutputs within a Transform. This sets the priority that the service uses for processing Transform Outputs. Possible values are `High`, `Normal` or `Low`.

* `video_analyzer_preset` - (Optional) A `video_analyzer_preset` block as defined below.

-> **NOTE:** Each output can only have one type of preset: builtin_preset,audio_analyzer_preset,face_detector_preset or video_analyzer_preset. If you need to apply different presets you must create one output for each one.

---

A `builtin_preset` block supports the following:

* `preset_name` - (Required) The built-in preset to be used for encoding videos. The Possible values are `AACGoodQualityAudio`, `AdaptiveStreaming`, `ContentAwareEncoding`, `ContentAwareEncodingExperimental`, `CopyAllBitrateNonInterleaved`, `H265AdaptiveStreaming`, `H265ContentAwareEncoding`, `H265SingleBitrate4K`, `H265SingleBitrate1080p`, `H265SingleBitrate720p`, `H264MultipleBitrate1080p`, `H264MultipleBitrateSD`, `H264MultipleBitrate720p`, `H264SingleBitrate1080p`, `H264SingleBitrateSD` and `H264SingleBitrate720p`.

* `preset_configuration` - (Optional) A `present_configuration` block as defined below.

---

A `audio_analyzer_preset` block supports the following:

* `audio_language` - (Optional) The language for the audio payload in the input using the BCP-47 format of 'language tag-region' (e.g: 'en-US'). If you know the language of your content, it is recommended that you specify it. The language must be specified explicitly for AudioAnalysisMode:Basic, since automatic language detection is not included in basic mode. If the language isn't specified, automatic language detection will choose the first language detected and process with the selected language for the duration of the file. It does not currently support dynamically switching between languages after the first language is detected. The automatic detection works best with audio recordings with clearly discernible speech. If automatic detection fails to find the language, transcription would fall back to `en-US`. The list of supported languages is available here: <https://go.microsoft.com/fwlink/?linkid=2109463>.

* `audio_analysis_mode` - (Optional) Possible values are `Basic` or `Standard`. Determines the set of audio analysis operations to be performed. Default to `Standard`.

* `experimental_options` - (Optional) Dictionary containing key value pairs for parameters not exposed in the preset itself.

---

A `preset_configuration` block supports the following:

* `compalexity` - (Optional) The complexity of the encoding. Possible values are `Balanced`, `Speed` or `Quality`.

* `interleave_output` - (Optional) Specifies the interleave mode of the output to control how audio are stored in the container format. Possible values are `InterleavedOutput` and `NonInterleavedOutput`. 

* `key_frame_interval_in_seconds` - (Optional) The key frame interval in seconds. Possible value is a positive float. For example, set as `2.0` to reduce the playback buffering for some players.

* `max_bitrate_bps` - (Optional) The maximum bitrate in bits per second (threshold for the top video layer). For example, set as `6000000` to avoid producing very high bitrate outputs for contents with high complexity.

* `max_height` - (Optional) The maximum height of output video layers. For example, set as `720` to produce output layers up to 720P even if the input is 4K.

* `max_layers` - (Optional) The maximum number of output video layers. For example, set as `4` to make sure at most 4 output layers are produced to control the overall cost of the encoding job.

* `min_bitrate_bps` - (Optional) The minimum bitrate in bits per second (threshold for the bottom video layer). For example, set as `200000` to have a bottom layer that covers users with low network bandwidth.

* `min_height` - (Optional) The minimum height of output video layers. For example, set as `360` to avoid output layers of smaller resolutions like 180P.

---

A `video_analyzer_preset` block supports the following:

* `audio_language` - (Optional) The language for the audio payload in the input using the BCP-47 format of 'language tag-region' (e.g: 'en-US'). If you know the language of your content, it is recommended that you specify it. The language must be specified explicitly for AudioAnalysisMode:Basic, since automatic language detection is not included in basic mode. If the language isn't specified, automatic language detection will choose the first language detected and process with the selected language for the duration of the file. It does not currently support dynamically switching between languages after the first language is detected. The automatic detection works best with audio recordings with clearly discernible speech. If automatic detection fails to find the language, transcription would fall back to `en-US`. The list of supported languages is available here: <https://go.microsoft.com/fwlink/?linkid=2109463>. 

* `audio_analysis_mode` - (Optional) Possible values are `Basic` or `Standard`. Determines the set of audio analysis operations to be performed. Default to `Standard`.

* `experimental_options` - (Optional) Dictionary containing key value pairs for parameters not exposed in the preset itself.
 
* `insights_type` - (Optional) Defines the type of insights that you want the service to generate. The allowed values are `AudioInsightsOnly`, `VideoInsightsOnly`, and `AllInsights`. If you set this to `AllInsights` and the input is audio only, then only audio insights are generated. Similarly, if the input is video only, then only video insights are generated. It is recommended that you not use `AudioInsightsOnly` if you expect some of your inputs to be video only; or use `VideoInsightsOnly` if you expect some of your inputs to be audio only. Your Jobs in such conditions would error out. Default to `AllInsights`.

---

A `face_detector_preset` block supports the following:

* `analysis_resolution` - (Optional) Possible values are `SourceResolution` or `StandardDefinition`. Specifies the maximum resolution at which your video is analyzed. The default behavior is `SourceResolution` which will keep the input video at its original resolution when analyzed. Using `StandardDefinition` will resize input videos to standard definition while preserving the appropriate aspect ratio. It will only resize if the video is of higher resolution. For example, a 1920x1080 input would be scaled to 640x360 before processing. Switching to `StandardDefinition` will reduce the time it takes to process high resolution video. It may also reduce the cost of using this component (see <https://azure.microsoft.com/en-us/pricing/details/media-services/#analytics> for details). However, faces that end up being too small in the resized video may not be detected. Default to `SourceResolution`.

* `blur_type` - (Optional) Specifies the type of blur to apply to faces in the output video. Possible values are `Black`, `Box`, `High`, `Low`,and `Med`.

* `experimental_options` - (Optional) Dictionary containing key value pairs for parameters not exposed in the preset itself.

* `mode` - (Optional) This mode provides the ability to choose between the following settings: 1) `Analyze` - For detection only. This mode generates a metadata JSON file marking appearances of faces throughout the video.Where possible, appearances of the same person are assigned the same ID. 2) `Combined` - Additionally redacts(blurs) detected faces. 3) `Redact` - This enables a 2-pass process, allowing for selective redaction of a subset of detected faces. It takes in the metadata file from a prior analyze pass, along with the source video, and a user-selected subset of IDs that require redaction. Default to `Analyze`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Transform.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Transform.
* `read` - (Defaults to 5 minutes) Used when retrieving the Transform.
* `update` - (Defaults to 30 minutes) Used when updating the Transform.
* `delete` - (Defaults to 30 minutes) Used when deleting the Transform.

## Import

Transforms can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_transform.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaServices/media1/transforms/transform1
```
