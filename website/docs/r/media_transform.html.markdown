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
      preset_configuration {
        complexity                    = "Balanced"
        interleave_output             = "NonInterleavedOutput"
        key_frame_interval_in_seconds = 123122.5
        max_bitrate_bps               = 300000
        max_height                    = 480
        max_layers                    = 14
        min_bitrate_bps               = 200000
        min_height                    = 360
      }
    }
  }

  output {
    relative_priority = "Low"
    on_error_action   = "ContinueJob"
    audio_analyzer_preset {
      audio_language      = "en-US"
      audio_analysis_mode = "Basic"
      experimental_options = {
        env = "test"
      }
    }
  }

  output {
    relative_priority = "Low"
    on_error_action   = "StopProcessingJob"
    face_detector_preset {
      analysis_resolution = "StandardDefinition"
      blur_type           = "Med"
      face_redactor_mode  = "Combined"
      experimental_options = {
        env = "test"
      }
    }
  }

  output {
    relative_priority = "Normal"
    on_error_action   = "StopProcessingJob"
    video_analyzer_preset {
      audio_language      = "en-US"
      audio_analysis_mode = "Basic"
      insights_type       = "AllInsights"
      experimental_options = {
        env = "test"
      }
    }
  }

  output {
    relative_priority = "Low"
    on_error_action   = "ContinueJob"
    custom_preset {
      codec {
        aac_audio {
          bitrate       = 128000
          channels      = 2
          sampling_rate = 48000
          profile       = "AacLc"
        }
      }

      codec {
        copy_audio {
          label = "test"
        }
      }

      codec {
        copy_video {
          label = "test"
        }
      }

      codec {
        h264_video {
          key_frame_interval             = "PT1S"
          stretch_mode                   = "AutoSize"
          sync_mode                      = "Auto"
          scene_change_detection_enabled = false
          rate_control_mode              = "ABR"
          complexity                     = "Quality"
          layer {
            width                    = "64"
            height                   = "64"
            bitrate                  = 1045000
            max_bitrate              = 1045000
            b_frames                 = 3
            slices                   = 0
            adaptive_b_frame_enabled = true
            profile                  = "Auto"
            level                    = "auto"
            buffer_window            = "PT5S"
            reference_frames         = 4
            crf                      = 23
            entropy_mode             = "Cabac"
          }
          layer {
            width                    = "64"
            height                   = "64"
            bitrate                  = 1000
            max_bitrate              = 1000
            b_frames                 = 3
            frame_rate               = "32"
            slices                   = 1
            adaptive_b_frame_enabled = true
            profile                  = "High444"
            level                    = "auto"
            buffer_window            = "PT5S"
            reference_frames         = 4
            crf                      = 23
            entropy_mode             = "Cavlc"
          }
        }
      }

      codec {
        h265_video {
          key_frame_interval             = "PT2S"
          stretch_mode                   = "AutoSize"
          sync_mode                      = "Auto"
          scene_change_detection_enabled = false
          complexity                     = "Speed"
          layer {
            width                    = "64"
            height                   = "64"
            bitrate                  = 1045000
            max_bitrate              = 1045000
            b_frames                 = 3
            slices                   = 5
            adaptive_b_frame_enabled = true
            profile                  = "Auto"
            label                    = "test"
            level                    = "auto"
            buffer_window            = "PT5S"
            frame_rate               = "32"
            reference_frames         = 4
            crf                      = 23
          }
        }
      }

      codec {
        jpg_image {
          stretch_mode  = "AutoSize"
          sync_mode     = "Auto"
          start         = "10"
          range         = "100%%"
          sprite_column = 1
          step          = "10"
          layer {
            quality = 70
            height  = "180"
            label   = "test"
            width   = "120"
          }
        }
      }

      codec {
        png_image {
          stretch_mode = "AutoSize"
          sync_mode    = "Auto"
          start        = "{Best}"
          range        = "80"
          step         = "10"
          layer {
            height = "180"
            label  = "test"
            width  = "120"
          }
        }
      }

      format {
        jpg {
          filename_pattern = "test{Basename}"
        }
      }

      format {
        mp4 {
          filename_pattern = "test{Bitrate}"
          output_file {
            labels = ["test", "ppe"]
          }
        }
      }

      format {
        png {
          filename_pattern = "test{Basename}"
        }
      }

      format {
        transport_stream {
          filename_pattern = "test{Bitrate}"
          output_file {
            labels = ["prod"]
          }
        }
      }

      filter {
        crop_rectangle {
          height = "240"
          left   = "30"
          top    = "360"
          width  = "70"
        }
        deinterlace {
          parity = "TopFieldFirst"
          mode   = "AutoPixelAdaptive"
        }
        fade_in {
          duration   = "PT5S"
          fade_color = "0xFF0000"
          start      = "10"
        }
        fade_out {
          duration   = "90%%"
          fade_color = "#FF0C7B"
          start      = "10%%"
        }
        rotation = "Auto"
        overlay {
          audio {
            input_label       = "label.jpg"
            start             = "PT5S"
            end               = "PT30S"
            fade_in_duration  = "PT1S"
            fade_out_duration = "PT2S"
            audio_gain_level  = 1.0
          }
        }
        overlay {
          video {
            input_label       = "label.jpg"
            start             = "PT5S"
            end               = "PT30S"
            fade_in_duration  = "PT1S"
            fade_out_duration = "PT2S"
            audio_gain_level  = 1.0
            opacity           = 1.0
            position {
              height = "180"
              left   = "20"
              top    = "240"
              width  = "140"
            }
            crop_rectangle {
              height = "240"
              left   = "30"
              top    = "360"
              width  = "70"
            }
          }
        }
      }
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

A `aac_audio` block supports the following:

* `bitrate` - (Optional) The bitrate of the audio in bits per second. Default to `128000`.

* `channels` - (Optional) The number of audio channels. Default to `2`.

* `label` - (Optional) Specifies the label for the codec. The label can be used to control muxing behavior.

* `profile` - (Optional) The encoding profile to be used when encoding audio with AAC. Possible values are `AacLc`, `HeAacV1`,and `HeAacV2`. Default to `AacLc`.

* `sampling_rate` - (Optional) The sampling rate to use for encoding in Hertz. Default to `48000`.

---

A `audio_analyzer_preset` block supports the following:

* `audio_language` - (Optional) The language for the audio payload in the input using the BCP-47 format of 'language tag-region' (e.g: 'en-US'). If you know the language of your content, it is recommended that you specify it. The language must be specified explicitly for AudioAnalysisMode:Basic, since automatic language detection is not included in basic mode. If the language isn't specified, automatic language detection will choose the first language detected and process with the selected language for the duration of the file. It does not currently support dynamically switching between languages after the first language is detected. The automatic detection works best with audio recordings with clearly discernible speech. If automatic detection fails to find the language, transcription would fall back to `en-US`. The list of supported languages is available here: <https://go.microsoft.com/fwlink/?linkid=2109463>.

* `audio_analysis_mode` - (Optional) Possible values are `Basic` or `Standard`. Determines the set of audio analysis operations to be performed. Default to `Standard`.

* `experimental_options` - (Optional) Dictionary containing key value pairs for parameters not exposed in the preset itself.

---

An `audio` block supports the following:

* `input_label` - (Required) The label of the job input which is to be used as an overlay. The input must specify exact one file. You can specify an image file in JPG, PNG, GIF or BMP format, or an audio file (such as a WAV, MP3, WMA or M4A file), or a video file.

* `audio_gain_level` - (Optional) The gain level of audio in the overlay. The value should be in the range `0` to `1.0`. The default is `1.0`.

* `end` - (Optional) The end position, with reference to the input video, at which the overlay ends. The value should be in ISO 8601 format. For example, `PT30S` to end the overlay at 30 seconds into the input video. If not specified or the value is greater than the input video duration, the overlay will be applied until the end of the input video if the overlay media duration is greater than the input video duration, else the overlay will last as long as the overlay media duration.

* `fade_in_duration` - (Optional) The duration over which the overlay fades in onto the input video. The value should be in ISO 8601 duration format. If not specified the default behavior is to have no fade in (same as `PT0S`).

* `fade_out_duration` - (Optional) The duration over which the overlay fades out of the input video. The value should be in ISO 8601 duration format. If not specified the default behavior is to have no fade out (same as `PT0S`).

* `start` - (Optional) The start position, with reference to the input video, at which the overlay starts. The value should be in ISO 8601 format. For example, `PT05S` to start the overlay at 5 seconds into the input video. If not specified the overlay starts from the beginning of the input video.

---

A `builtin_preset` block supports the following:

* `preset_name` - (Required) The built-in preset to be used for encoding videos. The Possible values are `AACGoodQualityAudio`, `AdaptiveStreaming`, `ContentAwareEncoding`, `ContentAwareEncodingExperimental`, `CopyAllBitrateNonInterleaved`, `DDGoodQualityAudio`, `H265AdaptiveStreaming`, `H265ContentAwareEncoding`, `H265SingleBitrate4K`, `H265SingleBitrate1080p`, `H265SingleBitrate720p`, `H264MultipleBitrate1080p`, `H264MultipleBitrateSD`, `H264MultipleBitrate720p`, `H264SingleBitrate1080p`, `H264SingleBitrateSD` and `H264SingleBitrate720p`.

* `preset_configuration` - (Optional) A `present_configuration` block as defined below.

---

A `codec` block supports the following:

* `aac_audio` - (Optional) A `aac_audio` block as defined above.
 
* `copy_audio` - (Optional) A `copy_audio` block as defined below.

* `copy_video` - (Optional) A `copy_video` block as defined below.

* `dd_audio` - (Optional) A `dd_audio` block as defined below.

* `h264_video` - (Optional) A `h264_video` block as defined below.

* `h265_video` - (Optional) A `h265_video` block as defined below.

* `jpg_image` - (Optional) A `jpg_image` block as defined below.

* `png_image` - (Optional) A `png_image` block as defined below.

-> **NOTE:** Each codec can only have one type: `aac_audio`, `copy_audio`, `copy_video`, `dd_audio`, `h264_video`, `h265_video`, `jpg_image` or `png_image`. If you need to apply different codec you must create one codec for each one.

---

A `copy_audio` block supports the following:

* `label` - (Optional) Specifies the label for the codec. The label can be used to control muxing behavior. 

---

A `copy_video` block supports the following:

* `label` - (Optional) Specifies the label for the codec. The label can be used to control muxing behavior. 
 
---

A `crop_rectangle` block supports the following:

* `height` - (Optional) The height of the rectangular region in pixels. This can be absolute pixel value (e.g `100`), or relative to the size of the video (For example, `50%`).

* `left` - (Optional) The number of pixels from the left-margin. This can be absolute pixel value (e.g `100`), or relative to the size of the video (For example, `50%`).

* `top` - (Optional) 	
  The number of pixels from the top-margin. This can be absolute pixel value (e.g `100`), or relative to the size of the video (For example, `50%`).

* `width` - (Optional) The width of the rectangular region in pixels. This can be absolute pixel value (e.g` 100`), or relative to the size of the video (For example, `50%`).

---

A `custom_preset` block supports the following:

* `codec` - (Required) One or more `codec` blocks as defined above.

* `format` - (Required) One or more `format` blocks as defined below.

* `experimental_options` - (Optional) Dictionary containing key value pairs for parameters not exposed in the preset itself.

* `filter` - (Optional) A `filter` block as defined below.
 
---

A `dd_audio` block supports the following:

* `bitrate` - (Optional) The bitrate of the audio in bits per second. Default to `192000`.

* `channels` - (Optional) The number of audio channels. Default to `2`.

* `label` - (Optional) Specifies the label for the codec. The label can be used to control muxing behavior.

* `sampling_rate` - (Optional) The sampling rate to use for encoding in Hertz. Default to `48000`.

---

A `deinterlace` block supports the following:

* `parity` - (Optional) The field parity to use for deinterlacing. Possible values are `Auto`, `TopFieldFirst` or `BottomFieldFirst`. Default to `Auto`.

* `mode` - (Optional) The deinterlacing mode. Possible values are `AutoPixelAdaptive` or `Off`. Default to `AutoPixelAdaptive`.

---

A `face_detector_preset` block supports the following:

* `analysis_resolution` - (Optional) Possible values are `SourceResolution` or `StandardDefinition`. Specifies the maximum resolution at which your video is analyzed. which will keep the input video at its original resolution when analyzed. Using `StandardDefinition` will resize input videos to standard definition while preserving the appropriate aspect ratio. It will only resize if the video is of higher resolution. For example, a 1920x1080 input would be scaled to 640x360 before processing. Switching to `StandardDefinition` will reduce the time it takes to process high resolution video. It may also reduce the cost of using this component (see <https://azure.microsoft.com/en-us/pricing/details/media-services/#analytics> for details). However, faces that end up being too small in the resized video may not be detected. Default to `SourceResolution`.

* `blur_type` - (Optional) Specifies the type of blur to apply to faces in the output video. Possible values are `Black`, `Box`, `High`, `Low`,and `Med`.

* `experimental_options` - (Optional) Dictionary containing key value pairs for parameters not exposed in the preset itself.

* `face_redactor_mode` - (Optional) This mode provides the ability to choose between the following settings: 1) `Analyze` - For detection only. This mode generates a metadata JSON file marking appearances of faces throughout the video. Where possible, appearances of the same person are assigned the same ID. 2) `Combined` - Additionally redacts(blurs) detected faces. 3) `Redact` - This enables a 2-pass process, allowing for selective redaction of a subset of detected faces. It takes in the metadata file from a prior analyze pass, along with the source video, and a user-selected subset of IDs that require redaction. Default to `Analyze`.

---

A `fade_in` block supports the following:

* `duration` - (Required) The duration of the fade effect in the video. The value can be in ISO 8601 format (For example, PT05S to fade In/Out a color during 5 seconds), or a frame count (For example, 10 to fade 10 frames from the start time), or a relative value to stream duration (For example, 10% to fade 10% of stream duration).

* `fade_color` - (Required) 	
  The color for the fade in/out. It can be on the [CSS Level1 colors](https://developer.mozilla.org/en-US/docs/Web/CSS/color_value/color_keywords) or an RGB/hex value: e.g: `rgb(255,0,0)`, `0xFF0000` or `#FF0000`.

* `start` - (Optional) The position in the input video from where to start fade. The value can be in ISO 8601 format (For example, `PT05S` to start at 5 seconds), or a frame count (For example, `10` to start at the 10th frame), or a relative value to stream duration (For example, `10%` to start at 10% of stream duration). Default to `0`.

---

A `fade_out` block supports the following:

* `duration` - (Required) The duration of the fade effect in the video. The value can be in ISO 8601 format (For example, PT05S to fade In/Out a color during 5 seconds), or a frame count (For example, 10 to fade 10 frames from the start time), or a relative value to stream duration (For example, 10% to fade 10% of stream duration).

* `fade_color` - (Required) 	
  The color for the fade in/out. It can be on the [CSS Level1 colors](https://developer.mozilla.org/en-US/docs/Web/CSS/color_value/color_keywords) or an RGB/hex value: e.g: `rgb(255,0,0)`, `0xFF0000` or `#FF0000`.

* `start` - (Optional) The position in the input video from where to start fade. The value can be in ISO 8601 format (For example, `PT05S` to start at 5 seconds), or a frame count (For example, `10` to start at the 10th frame), or a relative value to stream duration (For example, `10%` to start at 10% of stream duration). Default to `0`.

---

A `filter` block supports the following:

* `crop_rectangle` - (Optional) A `crop_rectangle` block as defined above.

* `deinterlace` - (Optional) A `deinterlace` block as defined below.

* `fade_in` - (Optional) A `fade_in` block as defined above.

* `fade_out` - (Optional) A `fade_out` block as defined above.

* `overlay` - (Optional) One or more `overlay` blocks as defined below.

* `rotation` - (Optional) The rotation to be applied to the input video before it is encoded. Possible values are `Auto`, `None`, `Rotate90`, `Rotate180`, `Rotate270`,or `Rotate0`. Default to `Auto`.

---

A `format` block supports the following:

* `jpg` - (Optional) A `jpg` block as defined below.

* `mp4` - (Optional) A `mp4` block as defined below.

* `png` - (Optional) A `png` block as defined below.

* `transport_stream` - (Optional) A `transport_stream` block as defined below.
 
-> **NOTE:** Each format can only have one type: `jpg`, `mp4`, `png` or `transport_stream`. If you need to apply different type you must create one format for each one.

---

A `h264_video` block supports the following:

* `complexity` - (Optional) The complexity of the encoding. Possible values are `Balanced`, `Speed` or `Quality`. Default to `Balanced`.

* `key_frame_interval` - (Optional) The distance between two key frames. The value should be non-zero in the range `0.5` to `20` seconds, specified in ISO 8601 format. The default is `2` seconds (`PT2S`). Note that this setting is ignored if `sync_mode` is set to `Passthrough`, where the KeyFrameInterval value will follow the input source setting.

* `label` - (Optional) Specifies the label for the codec. The label can be used to control muxing behavior.

* `layer` - (Optional) One or more `layer` blocks as defined below.

* `rate_control_mode` - (Optional) The rate control mode. Possible values are `ABR`, `CBR` or `CRF`. Default to `ABR`.

* `scene_change_detection_enabled` - (Optional) Whether the encoder should insert key frames at scene changes. This flag should be set to true only when the encoder is being configured to produce a single output video. Default to `false`.

* `stretch_mode` - (Optional) Specifies the resizing mode - how the input video will be resized to fit the desired output resolution(s). Possible values are `AutoFit`, `AutoSize` or `None`. Default to `AutoSize`.

* `sync_mode` - (Optional) Specifies the synchronization mode for the video. Possible values are `Auto`, `Cfr`, `Passthrough` or `Vfr`. Default to `Auto`.

---

A `h265_video` block supports the following:

* `complexity` - (Optional) The complexity of the encoding. Possible values are `Balanced`, `Speed` or `Quality`. Default to `Balanced`.

* `key_frame_interval` - (Optional) The distance between two key frames. The value should be non-zero in the range `0.5` to `20` seconds, specified in ISO 8601 format. The default is `2` seconds (`PT2S`). Note that this setting is ignored if `sync_mode` is set to `Passthrough`, where the KeyFrameInterval value will follow the input source setting.

* `label` - (Optional) Specifies the label for the codec. The label can be used to control muxing behavior.

* `layer` - (Optional) One or more `layer` blocks as defined below.

* `scene_change_detection_enabled` - (Optional) Whether the encoder should insert key frames at scene changes. This flag should be set to true only when the encoder is being configured to produce a single output video. Default to `false`.

* `stretch_mode` - (Optional) Specifies the resizing mode - how the input video will be resized to fit the desired output resolution(s). Possible values are `AutoFit`, `AutoSize` or `None`. Default to `AutoSize`.

* `sync_mode` - (Optional) Specifies the synchronization mode for the video. Possible values are `Auto`, `Cfr`, `Passthrough` or `Vfr`. Default to `Auto`.

---

A `jpg` block supports the following:

* `filename_pattern` - (Required) The file naming pattern used for the creation of output files. The following macros are supported in the file name: `{Basename}` - An expansion macro that will use the name of the input video file. If the base name(the file suffix is not included) of the input video file is less than 32 characters long, the base name of input video files will be used. If the length of base name of the input video file exceeds 32 characters, the base name is truncated to the first 32 characters in total length. `{Extension}` - The appropriate extension for this format. `{Label}` - The label assigned to the codec/layer. `{Index}` - A unique index for thumbnails. Only applicable to thumbnails. `{AudioStream}` - string "Audio" plus audio stream number(start from 1). `{Bitrate}` - The audio/video bitrate in kbps. Not applicable to thumbnails. `{Codec}` - The type of the audio/video codec. `{Resolution}` - The video resolution. Any unsubstituted macros will be collapsed and removed from the filename.

---

A `jpg_image` block supports the following:

* `start` - (Required) The position in the input video from where to start generating thumbnails. The value can be in ISO 8601 format (For example, `PT05S` to start at 5 seconds), or a frame count (For example, `10` to start at the 10th frame), or a relative value to stream duration (For example, `10%` to start at 10% of stream duration). Also supports a macro `{Best}`, which tells the encoder to select the best thumbnail from the first few seconds of the video and will only produce one thumbnail, no matter what other settings are for `step` and `range`.

* `key_frame_interval` - (Optional) The distance between two key frames. The value should be non-zero in the range `0.5` to `20` seconds, specified in ISO 8601 format. The default is `2` seconds (`PT2S`). Note that this setting is ignored if `sync_mode` is set to `Passthrough`, where the KeyFrameInterval value will follow the input source setting.

* `label` - (Optional) Specifies the label for the codec. The label can be used to control muxing behavior.

* `layer` - (Optional) One or more `layer` blocks as defined below.

* `range` - (Optional) The position relative to transform preset start time in the input video at which to stop generating thumbnails. The value can be in ISO 8601 format (For example, `PT5M30S` to stop at 5 minutes and 30 seconds from start time), or a frame count (For example, `300` to stop at the 300th frame from the frame at start time. If this value is `1`, it means only producing one thumbnail at start time), or a relative value to the stream duration (For example, `50%` to stop at half of stream duration from start time). The default value is `100%`, which means to stop at the end of the stream. 

* `sprite_column` - (Optional) Sets the number of columns used in thumbnail sprite image. The number of rows are automatically calculated and a VTT file is generated with the coordinate mappings for each thumbnail in the sprite. Note: this value should be a positive integer and a proper value is recommended so that the output image resolution will not go beyond JPEG maximum pixel resolution limit `65535x65535`.

* `step` - (Optional) The intervals at which thumbnails are generated. The value can be in ISO 8601 format (For example, `PT05S` for one image every 5 seconds), or a frame count (For example, `30` for one image every 30 frames), or a relative value to stream duration (For example, `10%` for one image every 10% of stream duration). Note: Step value will affect the first generated thumbnail, which may not be exactly the one specified at transform preset start time. This is due to the encoder, which tries to select the best thumbnail between start time and Step position from start time as the first output. As the default value is `10%`, it means if stream has long duration, the first generated thumbnail might be far away from the one specified at start time. Try to select reasonable value for Step if the first thumbnail is expected close to start time, or set Range value at `1` if only one thumbnail is needed at start time.

* `stretch_mode` - (Optional) The resizing mode, which indicates how the input video will be resized to fit the desired output resolution(s). Possible values are `AutoFit`, `AutoSize` or `None`. Default to `AutoSize`.

* `sync_mode` - (Optional) Specifies the synchronization mode for the video. Possible values are `Auto`, `Cfr`, `Passthrough` or `Vfr`. Default to `Auto`.

---

A `layer` block within `h264_video` block supports the following:

* `bitrate` - (Required) The average bitrate in bits per second at which to encode the input video when generating this layer.

* `adaptive_b_frame_enabled` - (Optional) Whether adaptive B-frames are used when encoding this layer. If not specified, the encoder will turn it on whenever the video profile permits its use. Default to `true`.

* `b_frames` - (Optional) The number of B-frames to use when encoding this layer. If not specified, the encoder chooses an appropriate number based on the video profile and level.

* `buffer_window` - (Optional) Specifies the maximum amount of time that the encoder should buffer frames before encoding. The value should be in ISO 8601 format. The value should be in the range `0.1` to `100` seconds. The default is `5` seconds (`PT5S`).

* `crf` - (Optional) The value of CRF to be used when encoding this layer. This setting takes effect when `rate_control_mode` is set `CRF`. The range of CRF value is between `0` and `51`, where lower values would result in better quality, at the expense of higher file sizes. Higher values mean more compression, but at some point quality degradation will be noticed. Default to `23`.

* `entropy_mode` - (Optional) The entropy mode to be used for this layer. Possible values are `Cabac` or `Cavlc`. If not specified, the encoder chooses the mode that is appropriate for the profile and level.

* `frame_rate` - (Optional) The frame rate (in frames per second) at which to encode this layer. The value can be in the form of `M/N` where `M` and `N` are integers (For example, `30000/1001`), or in the form of a number (For example, `30`, or `29.97`). The encoder enforces constraints on allowed frame rates based on the profile and level. If it is not specified, the encoder will use the same frame rate as the input video.

* `height` - (Optional) The height of the output video for this layer. The value can be absolute (in pixels) or relative (in percentage). For example `50%` means the output video has half as many pixels in height as the input.

* `label` - (Optional) The alphanumeric label for this layer, which can be used in multiplexing different video and audio layers, or in naming the output file.

* `level` - (Optional) The H.264 levels. Currently, the resource support Level up to `6.2`. The value can be `auto`, or a number that matches the H.264 profile. If not specified, the default is `auto`, which lets the encoder choose the Level that is appropriate for this layer.

* `max_bitrate` - (Optional) The maximum bitrate (in bits per second), at which the VBV buffer should be assumed to refill. If not specified, defaults to the same value as bitrate.

* `profile` - (Optional) The H.264 profile. Possible values are `Auto`, `Baseline`, `High`, `High422`, `High444`,or `Main`. Default to `Auto`.

* `reference_frames` - (Optional) The number of reference frames to be used when encoding this layer. If not specified, the encoder determines an appropriate number based on the encoder complexity setting.

* `slices` - (Optional) The number of slices to be used when encoding this layer. If not specified, default is `1`, which means that encoder will use a single slice for each frame.

* `width` - (Optional) The width of the output video for this layer. The value can be absolute (in pixels) or relative (in percentage). For example `50%` means the output video has half as many pixels in width as the input.

---

A `layer` block within `h265_video` block supports the following:

* `bitrate` - (Required) The average bitrate in bits per second at which to encode the input video when generating this layer.

* `adaptive_b_frame_enabled` - (Optional) Whether adaptive B-frames are used when encoding this layer. If not specified, the encoder will turn it on whenever the video profile permits its use. Default to `true`.

* `b_frames` - (Optional) The number of B-frames to use when encoding this layer. If not specified, the encoder chooses an appropriate number based on the video profile and level.

* `buffer_window` - (Optional) Specifies the maximum amount of time that the encoder should buffer frames before encoding. The value should be in ISO 8601 format. The value should be in the range `0.1` to `100` seconds. The default is `5` seconds (`PT5S`).

* `crf` - (Optional) The value of CRF to be used when encoding this layer. This setting takes effect when `rate_control_mode` is set `CRF`. The range of CRF value is between `0` and `51`, where lower values would result in better quality, at the expense of higher file sizes. Higher values mean more compression, but at some point quality degradation will be noticed. Default to `28`.

* `entropy_mode` - (Optional) The entropy mode to be used for this layer. Possible values are `Cabac` or `Cavlc`. If not specified, the encoder chooses the mode that is appropriate for the profile and level.

* `frame_rate` - (Optional) 	
  The frame rate (in frames per second) at which to encode this layer. The value can be in the form of `M/N` where `M` and `N` are integers (For example, `30000/1001`), or in the form of a number (For example, `30`, or `29.97`). The encoder enforces constraints on allowed frame rates based on the profile and level. If it is not specified, the encoder will use the same frame rate as the input video.

* `height` - (Optional) The height of the output video for this layer. The value can be absolute (in pixels) or relative (in percentage). For example `50%` means the output video has half as many pixels in height as the input.

* `label` - (Optional) The alphanumeric label for this layer, which can be used in multiplexing different video and audio layers, or in naming the output file.

* `level` - (Optional) The H.264 levels. Currently, the resource support Level up to `6.2`. The value can be `auto`, or a number that matches the H.264 profile. If not specified, the default is `auto`, which lets the encoder choose the Level that is appropriate for this layer.

* `max_bitrate` - (Optional) The maximum bitrate (in bits per second), at which the VBV buffer should be assumed to refill. If not specified, defaults to the same value as bitrate.

* `profile` - (Optional) The H.264 profile. Possible values are `Auto`, `Baseline`, `High`, `High422`, `High444`,or `Main`. Default to `Auto`.

* `reference_frames` - (Optional) The number of reference frames to be used when encoding this layer. If not specified, the encoder determines an appropriate number based on the encoder complexity setting.

* `slices` - (Optional) The number of slices to be used when encoding this layer. If not specified, default is `1`, which means that encoder will use a single slice for each frame.

* `width` - (Optional) The width of the output video for this layer. The value can be absolute (in pixels) or relative (in percentage). For example `50%` means the output video has half as many pixels in width as the input.

---

A `layer` block within `jpg_image` block supports the following:

* `height` - (Optional) The height of the output video for this layer. The value can be absolute (in pixels) or relative (in percentage). For example `50%` means the output video has half as many pixels in height as the input.

* `label` - (Optional) The alphanumeric label for this layer, which can be used in multiplexing different video and audio layers, or in naming the output file.

* `quality` - (Optional) The compression quality of the JPEG output. Range is from `0` to `100` and the default is `70`.

* `width` - (Optional) The width of the output video for this layer. The value can be absolute (in pixels) or relative (in percentage). For example `50%` means the output video has half as many pixels in width as the input.

---

A `layer` block within `png_image` block supports the following:

* `height` - (Optional) The height of the output video for this layer. The value can be absolute (in pixels) or relative (in percentage). For example `50%` means the output video has half as many pixels in height as the input.

* `label` - (Optional) The alphanumeric label for this layer, which can be used in multiplexing different video and audio layers, or in naming the output file.

* `width` - (Optional) The width of the output video for this layer. The value can be absolute (in pixels) or relative (in percentage). For example `50%` means the output video has half as many pixels in width as the input.

---

A `mp4` block supports the following:

* `filename_pattern` - (Required) The file naming pattern used for the creation of output files. The following macros are supported in the file name: `{Basename}` - An expansion macro that will use the name of the input video file. If the base name(the file suffix is not included) of the input video file is less than 32 characters long, the base name of input video files will be used. If the length of base name of the input video file exceeds 32 characters, the base name is truncated to the first 32 characters in total length. `{Extension}` - The appropriate extension for this format. `{Label}` - The label assigned to the codec/layer. `{Index}` - A unique index for thumbnails. Only applicable to thumbnails. `{AudioStream}` - string "Audio" plus audio stream number(start from 1). `{Bitrate}` - The audio/video bitrate in kbps. Not applicable to thumbnails. `{Codec}` - The type of the audio/video codec. `{Resolution}` - The video resolution. Any unsubstituted macros will be collapsed and removed from the filename.

* `output_file` - (Optional) One or more `output_file` blocks as defined below.

---

An `output` block supports the following:

* `audio_analyzer_preset` - (Optional) An `audio_analyzer_preset` block as defined above.

* `builtin_preset` - (Optional) A `builtin_preset` block as defined above.

* `custom_preset` - (Optional) A `custom_preset` block as defined above.

* `face_detector_preset` - (Optional) A `face_detector_preset` block as defined above.

* `on_error_action` - (Optional) A Transform can define more than one outputs. This property defines what the service should do when one output fails - either continue to produce other outputs, or, stop the other outputs. The overall Job state will not reflect failures of outputs that are specified with `ContinueJob`. Possible values are `StopProcessingJob` or `ContinueJob`. The default is `StopProcessingJob`.

* `relative_priority` - (Optional) Sets the relative priority of the TransformOutputs within a Transform. This sets the priority that the service uses for processing Transform Outputs. Possible values are `High`, `Normal` or `Low`. Defaults to `Normal`.

* `video_analyzer_preset` - (Optional) A `video_analyzer_preset` block as defined below.

-> **NOTE:** Each output can only have one type of preset: `builtin_preset`, `audio_analyzer_preset`, `custom_preset`, `face_detector_preset` or `video_analyzer_preset`. If you need to apply different presets you must create one output for each one.

---

An `output_file` block supports the following:

* `labels` - (Required) The list of labels that describe how the encoder should multiplex video and audio into an output file. For example, if the encoder is producing two video layers with labels `v1` and `v2`, and one audio layer with label `a1`, then an array like `["v1", "a1"]` tells the encoder to produce an output file with the video track represented by `v1` and the audio track represented by `a1`.

---

An `overlay` block supports the following:

* `audio` - (Optional) An `audio` block as defined above.

* `video` - (Optional) A `video` block as defined below.

-> **NOTE:** Each overlay can only have one type: `audio` or `video`. If you need to apply different type you must create one overlay for each one.

---

A `png` block supports the following:

* `filename_pattern` - (Required) The file naming pattern used for the creation of output files. The following macros are supported in the file name: `{Basename}` - An expansion macro that will use the name of the input video file. If the base name(the file suffix is not included) of the input video file is less than 32 characters long, the base name of input video files will be used. If the length of base name of the input video file exceeds 32 characters, the base name is truncated to the first 32 characters in total length. `{Extension}` - The appropriate extension for this format. `{Label}` - The label assigned to the codec/layer. `{Index}` - A unique index for thumbnails. Only applicable to thumbnails. `{AudioStream}` - string "Audio" plus audio stream number(start from 1). `{Bitrate}` - The audio/video bitrate in kbps. Not applicable to thumbnails. `{Codec}` - The type of the audio/video codec. `{Resolution}` - The video resolution. Any unsubstituted macros will be collapsed and removed from the filename.

---

A `png_image` block supports the following:

* `start` - (Required) The position in the input video from where to start generating thumbnails. The value can be in ISO 8601 format (For example, `PT05S` to start at 5 seconds), or a frame count (For example, `10` to start at the 10th frame), or a relative value to stream duration (For example, `10%` to start at 10% of stream duration). Also supports a macro `{Best}`, which tells the encoder to select the best thumbnail from the first few seconds of the video and will only produce one thumbnail, no matter what other settings are for `step` and `range`.

* `key_frame_interval` - (Optional) The distance between two key frames. The value should be non-zero in the range `0.5` to `20` seconds, specified in ISO 8601 format. The default is `2` seconds (`PT2S`). Note that this setting is ignored if `sync_mode` is set to `Passthrough`, where the KeyFrameInterval value will follow the input source setting.

* `label` - (Optional) Specifies the label for the codec. The label can be used to control muxing behavior.

* `layer` - (Optional) One or more `layer` blocks as defined below.

* `range` - (Optional) The position relative to transform preset start time in the input video at which to stop generating thumbnails. The value can be in ISO 8601 format (For example, `PT5M30S` to stop at `5` minutes and `30` seconds from start time), or a frame count (For example, `300` to stop at the 300th frame from the frame at start time. If this value is `1`, it means only producing one thumbnail at start time), or a relative value to the stream duration (For example, `50%` to stop at half of stream duration from start time). The default value is `100%`, which means to stop at the end of the stream.

* `step` - (Optional) The intervals at which thumbnails are generated. The value can be in ISO 8601 format (For example, `PT05S` for one image every 5 seconds), or a frame count (For example, `30` for one image every 30 frames), or a relative value to stream duration (For example, `10%` for one image every 10% of stream duration). Note: Step value will affect the first generated thumbnail, which may not be exactly the one specified at transform preset start time. This is due to the encoder, which tries to select the best thumbnail between start time and Step position from start time as the first output. As the default value is `10%`, it means if stream has long duration, the first generated thumbnail might be far away from the one specified at start time. Try to select reasonable value for Step if the first thumbnail is expected close to start time, or set Range value at `1` if only one thumbnail is needed at start time.

* `stretch_mode` - (Optional) The resizing mode, which indicates how the input video will be resized to fit the desired output resolution(s). Possible values are `AutoFit`, `AutoSize` or `None`. Default to `AutoSize`.

* `sync_mode` - (Optional) Specifies the synchronization mode for the video. Possible values are `Auto`, `Cfr`, `Passthrough` or `Vfr`. Default to `Auto`.

---

A `position` block supports the following:

* `height` - (Optional) The height of the rectangular region in pixels. This can be absolute pixel value (e.g `100`), or relative to the size of the video (For example, `50%`).

* `left` - (Optional) The number of pixels from the left-margin. This can be absolute pixel value (e.g `100`), or relative to the size of the video (For example, `50%`).

* `top` - (Optional) 	
  The number of pixels from the top-margin. This can be absolute pixel value (e.g `100`), or relative to the size of the video (For example, `50%`).

* `width` - (Optional) The width of the rectangular region in pixels. This can be absolute pixel value (e.g` 100`), or relative to the size of the video (For example, `50%`).

---

A `preset_configuration` block supports the following:

* `complexity` - (Optional) The complexity of the encoding. Possible values are `Balanced`, `Speed` or `Quality`.

* `interleave_output` - (Optional) Specifies the interleave mode of the output to control how audio are stored in the container format. Possible values are `InterleavedOutput` and `NonInterleavedOutput`. 

* `key_frame_interval_in_seconds` - (Optional) The key frame interval in seconds. Possible value is a positive float. For example, set as `2.0` to reduce the playback buffering for some players.

* `max_bitrate_bps` - (Optional) The maximum bitrate in bits per second (threshold for the top video layer). For example, set as `6000000` to avoid producing very high bitrate outputs for contents with high complexity.

* `max_height` - (Optional) The maximum height of output video layers. For example, set as `720` to produce output layers up to 720P even if the input is 4K.

* `max_layers` - (Optional) The maximum number of output video layers. For example, set as `4` to make sure at most 4 output layers are produced to control the overall cost of the encoding job.

* `min_bitrate_bps` - (Optional) The minimum bitrate in bits per second (threshold for the bottom video layer). For example, set as `200000` to have a bottom layer that covers users with low network bandwidth.

* `min_height` - (Optional) The minimum height of output video layers. For example, set as `360` to avoid output layers of smaller resolutions like 180P.

---

A `transport_stream` block supports the following:

* `filename_pattern` - (Required) The file naming pattern used for the creation of output files. The following macros are supported in the file name: `{Basename}` - An expansion macro that will use the name of the input video file. If the base name(the file suffix is not included) of the input video file is less than 32 characters long, the base name of input video files will be used. If the length of base name of the input video file exceeds 32 characters, the base name is truncated to the first 32 characters in total length. `{Extension}` - The appropriate extension for this format. `{Label}` - The label assigned to the codec/layer. `{Index}` - A unique index for thumbnails. Only applicable to thumbnails. `{AudioStream}` - string "Audio" plus audio stream number(start from 1). `{Bitrate}` - The audio/video bitrate in kbps. Not applicable to thumbnails. `{Codec}` - The type of the audio/video codec. `{Resolution}` - The video resolution. Any unsubstituted macros will be collapsed and removed from the filename.

* `output_file` - (Optional) One or more `output_file` blocks as defined above.

---

An `video` block supports the following:

* `input_label` - (Required) The label of the job input which is to be used as an overlay. The input must specify exact one file. You can specify an image file in JPG, PNG, GIF or BMP format, or an audio file (such as a WAV, MP3, WMA or M4A file), or a video file.

* `audio_gain_level` - (Optional) The gain level of audio in the overlay. The value should be in range between `0` to `1.0`. The default is `1.0`.

* `crop_rectangle` - (Optional) A `crop_rectangle` block as defined above.

* `end` - (Optional) The end position, with reference to the input video, at which the overlay ends. The value should be in ISO 8601 format. For example, `PT30S` to end the overlay at 30 seconds into the input video. If not specified or the value is greater than the input video duration, the overlay will be applied until the end of the input video if the overlay media duration is greater than the input video duration, else the overlay will last as long as the overlay media duration.

* `fade_in_duration` - (Optional) The duration over which the overlay fades in onto the input video. The value should be in ISO 8601 duration format. If not specified the default behavior is to have no fade in (same as `PT0S`).

* `fade_out_duration` - (Optional) The duration over which the overlay fades out of the input video. The value should be in ISO 8601 duration format. If not specified the default behavior is to have no fade out (same as `PT0S`).

* `opacity` - (Optional) The opacity of the overlay. The value should be in the range between `0` to `1.0`. Default to `1.0`, which means the overlay is opaque.

* `position` - (Optional) A `position` block as defined above.

* `start` - (Optional) The start position, with reference to the input video, at which the overlay starts. The value should be in ISO 8601 format. For example, `PT05S` to start the overlay at 5 seconds into the input video. If not specified the overlay starts from the beginning of the input video.

---

A `video_analyzer_preset` block supports the following:

* `audio_language` - (Optional) The language for the audio payload in the input using the BCP-47 format of 'language tag-region' (e.g: 'en-US'). If you know the language of your content, it is recommended that you specify it. The language must be specified explicitly for AudioAnalysisMode:Basic, since automatic language detection is not included in basic mode. If the language isn't specified, automatic language detection will choose the first language detected and process with the selected language for the duration of the file. It does not currently support dynamically switching between languages after the first language is detected. The automatic detection works best with audio recordings with clearly discernible speech. If automatic detection fails to find the language, transcription would fall back to `en-US`. The list of supported languages is available here: <https://go.microsoft.com/fwlink/?linkid=2109463>. 

* `audio_analysis_mode` - (Optional) Possible values are `Basic` or `Standard`. Determines the set of audio analysis operations to be performed. Default to `Standard`.

* `experimental_options` - (Optional) Dictionary containing key value pairs for parameters not exposed in the preset itself.
 
* `insights_type` - (Optional) Defines the type of insights that you want the service to generate. The allowed values are `AudioInsightsOnly`, `VideoInsightsOnly`, and `AllInsights`. If you set this to `AllInsights` and the input is audio only, then only audio insights are generated. Similarly, if the input is video only, then only video insights are generated. It is recommended that you not use `AudioInsightsOnly` if you expect some of your inputs to be video only; or use `VideoInsightsOnly` if you expect some of your inputs to be audio only. Your Jobs in such conditions would error out. Default to `AllInsights`.

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
