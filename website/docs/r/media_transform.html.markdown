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
    on_error_type     = "ContinueJob"
    preset {
      type        = "BuiltInStandardEncoderPreset"
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
    on_error_type     = "ContinueJob"
    preset {
      type        = "BuiltInStandardEncoderPreset"
      preset_name = "AACGoodQualityAudio"
    }
  }

  output {
    relative_priority = "Low"
    on_error_type     = "ContinueJob"
    preset {
      type        = "AudioAnalyzerPreset"
      audio_language= "en-US"
      audio_analysis_mode = "Basic"
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

* `output` - (Required) One or more `output` blocks as defined below. At least one `output` must be defined.

---

A `output` block supports the following:

* `on_error_type` - (Optional) A Transform can define more than one outputs. This property defines what the service should do when one output fails - either continue to produce other outputs, or, stop the other outputs. The overall Job state will not reflect failures of outputs that are specified with `ContinueJob`. Possibles value are `StopProcessingJob` or `ContinueJob`. The default is `StopProcessingJob`.

* `preset` - (Optional) A `preset` block as defined below.

* `relative_priority` - (Optional) Sets the relative priority of the TransformOutputs within a Transform. This sets the priority that the service uses for processing TransformOutputs. Possibles value are `High`, `Normal` or `Low`. The default priority is `Normal`.

---

A `preset` block supports the following:

* `analysis_resolution` - (Optional) Possibles value are `SourceResolution` or `StandardDefinition`. Specifies the maximum resolution at which your video is analyzed. The default behavior is `SourceResolution` which will keep the input video at its original resolution when analyzed. Using `StandardDefinition` will resize input videos to standard definition while preserving the appropriate aspect ratio. It will only resize if the video is of higher resolution. For example, a 1920x1080 input would be scaled to 640x360 before processing. Switching to `StandardDefinition` will reduce the time it takes to process high resolution video. It may also reduce the cost of using this component (see https://azure.microsoft.com/en-us/pricing/details/media-services/#analytics for details). However, faces that end up being too small in the resized video may not be detected. This property is only used if the `type` of preset is `FaceDetectorPreset` otherwise will be ignored.

* `audio_analysis_mode` - (Optional) Possibles value are `Basic` or `Standard`. Determines the set of audio analysis operations to be performed. If unspecified, the `Standard` AudioAnalysisMode would be chosen. This property is only used if the `type` of preset is `AudioAnalyzerPreset` or `VideoAnalyzerPreset`  otherwise will be ignored.

* `audio_language` - (Optional) The language for the audio payload in the input using the BCP-47 format of 'language tag-region' (e.g: 'en-US'). If you know the language of your content, it is recommended that you specify it. The language must be specified explicitly for AudioAnalysisMode:Basic, since automatic language detection is not included in basic mode. If the language isn't specified, automatic language detection will choose the first language detected and process with the selected language for the duration of the file. It does not currently support dynamically switching between languages after the first language is detected. The automatic detection works best with audio recordings with clearly discernable speech. If automatic detection fails to find the language, transcription would fallback to 'en-US'." The list of supported languages is available here: https://go.microsoft.com/fwlink/?linkid=2109463. This property is only used if the `type` of preset is `AudioAnalyzerPreset` or `VideoAnalyzerPreset`  otherwise will be ignored.

* `insights_type` - (Optional) Defines the type of insights that you want the service to generate. The allowed values are `AudioInsightsOnly`, `VideoInsightsOnly`, and `AllInsights`. The default is `AllInsights`. If you set this to `AllInsights` and the input is audio only, then only audio insights are generated. Similarly if the input is video only, then only video insights are generated. It is recommended that you not use `AudioInsightsOnly` if you expect some of your inputs to be video only; or use `VideoInsightsOnly` if you expect some of your inputs to be audio only. Your Jobs in such conditions would error out. This attribute is only used if the `type` of preset is `VideoAnalyzerPreset`  otherwise will be ignored.

* `preset_name` - (Optional) The built-in preset to be used for encoding videos. This property is only used if the `type` of preset is `BuiltInStandardEncoderPreset`  otherwise will be ignored.

* `type` - (Optional) The type of preset. Possibles values are `BuiltInStandardEncoderPreset`,`AudioAnalyzerPreset`,`VideoAnalyzerPreset` or `FaceDetectorPreset`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Transform.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Transform.
* `read` - (Defaults to 5 minutes) Used when retrieving the Transform.
* `update` - (Defaults to 30 minutes) Used when updating the Transform.
* `delete` - (Defaults to 30 minutes) Used when deleting the Transform.

## Import

Transforms can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_transform.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaservices/media1/transforms/transform1
```