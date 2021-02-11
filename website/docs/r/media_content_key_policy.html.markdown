---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_content_key_policy"
description: |-
  Manages a Content Key Policy.
---

# azurerm_media_content_key_policy

Manages a Content Key Policy.

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

resource "azurerm_media_content_key_policy" "example" {
  name                        = "example"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  policy_option {
    name = "fairPlay"
    fairplay_configuration {
      ask                       = "myKeyForFairPlay"
      pfx                       = "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"
      pfx_password              = "password"
      rental_duration_seconds   = 2249
      rental_and_lease_key_type = "PersistentUnlimited"

    }
    open_restriction_enabled = true
  }
  policy_option {
    name = "playReady"
    playready_configuration_license {
      allow_test_devices = true
      begin_date         = "2017-10-16T18:22:53Z"
      play_right {
        scms_restriction                                         = 2
        digital_video_only_content_restriction                   = false
        image_constraint_for_analog_component_video_restriction  = false
        image_constraint_for_analog_computer_monitor_restriction = false
        allow_passing_video_content_to_unknown_output            = "NotAllowed"
        uncompressed_digital_video_opl                           = 100
        uncompressed_digital_audio_opl                           = 100
        analog_video_opl                                         = 150
        compressed_digital_audio_opl                             = 150
      }
      license_type                             = "Persistent"
      content_type                             = "UltraVioletDownload"
      content_key_location_from_header_enabled = true
    }
    open_restriction_enabled = true
  }
  policy_option {
    name                            = "clearKey"
    clear_key_configuration_enabled = true
    token_restriction {
      issuer                      = "urn:issuer"
      audience                    = "urn:audience"
      token_type                  = "Swt"
      primary_symmetric_token_key = "AAAAAAAAAAAAAAAAAAAAAA=="
    }
  }
  policy_option {
    name = "widevine"
    widevine_configuration_template = jsonencode({
      "allowed_track_types" : "SD_HD",
      "content_key_specs" : [{
        "track_type" : "SD",
        "security_level" : 1,
        "required_output_protection" : {
          "hdcp" : "HDCP_V2"
        },
      }],
      "policy_overrides" : {
        "can_play" : true,
        "can_persist" : true,
        "can_renew" : false,
      },
    })
    open_restriction_enabled = true
  }
}
```

## Arguments Reference

The following arguments are supported:

* `media_services_account_name` - (Required) The Media Services account name. Changing this forces a new Content Key Policy to be created.

* `name` - (Required) The name which should be used for this Content Key Policy. Changing this forces a new Content Key Policy to be created.

* `policy_option` - (Required) One or more `policy_option` blocks as defined below.

* `resource_group_name` - (Required) The name of the Resource Group where the Content Key Policy should exist. Changing this forces a new Content Key Policy to be created.

---

* `description` - (Optional) A description for the Policy.

---

A `fairplay_configuration` block supports the following:

* `ask` - (Optional) The key that must be used as FairPlay Application Secret key.

* `offline_rental_configuration` - (Optional) A `offline_rental_configuration` block as defined below.

* `pfx` - (Optional) The Base64 representation of FairPlay certificate in PKCS 12 (pfx) format (including private key).

* `pfx_password` - (Optional) The password encrypting FairPlay certificate in PKCS 12 (pfx) format.

* `rental_and_lease_key_type` - (Optional) The rental and lease key type. Supported values are `DualExpiry`, `PersistentLimited`, `PersistentUnlimited` or `Undefined`.

* `rental_duration_seconds` - (Optional) The rental duration. Must be greater than 0.

---

A `offline_rental_configuration` block supports the following:

* `playback_duration_seconds` - (Optional) Playback duration.

* `storage_duration_seconds` - (Optional) Storage duration.

---

A `play_right` block supports the following:

* `agc_and_color_stripe_restriction` - (Optional) Configures Automatic Gain Control (AGC) and Color Stripe in the license. Must be between 0 and 3 inclusive.

* `allow_passing_video_content_to_unknown_output` - (Optional) Configures Unknown output handling settings of the license. Supported values are `Allowed`, `AllowedWithVideoConstriction` or `NotAllowed`.

* `analog_video_opl` - (Optional) Specifies the output protection level for compressed digital audio. Supported values are 100, 150 or 200.

* `compressed_digital_audio_opl` - (Optional) Specifies the output protection level for compressed digital audio.Supported values are 100, 150 or 200.

* `digital_video_only_content_restriction` - (Optional) Enables the Image Constraint For Analog Component Video Restriction in the license.

* `first_play_expiration` - (Optional) The amount of time that the license is valid after the license is first used to play content.

* `image_constraint_for_analog_component_video_restriction` - (Optional) Enables the Image Constraint For Analog Component Video Restriction in the license.

* `image_constraint_for_analog_computer_monitor_restriction` - (Optional) Enables the Image Constraint For Analog Component Video Restriction in the license.

* `scms_restriction` - (Optional) Configures the Serial Copy Management System (SCMS) in the license. Must be between 0 and 3 inclusive.

* `uncompressed_digital_audio_opl` - (Optional) Specifies the output protection level for uncompressed digital audio. Supported values are 100, 150, 250 or 300.

* `uncompressed_digital_video_opl` - (Optional) Specifies the output protection level for uncompressed digital video. Supported values are 100, 150, 250 or 300.

---

A `playready_configuration_license` block supports the following:

* `allow_test_devices` - (Optional) A flag indicating whether test devices can use the license.

* `begin_date` - (Optional) The begin date of license.

* `content_key_location_from_header_enabled` - (Optional) Specifies that the content key ID is in the PlayReady header. 
* `content_key_location_from_key_id` - (Optional) The content key ID. Specifies that the content key ID is specified in the PlayReady configuration. 

-> **NOTE:** You can only specify one content key location. For example if you specify content_key_location_from_header_enabled in true, you shouldn't specify content_key_location_from_key_id and vice versa.


* `content_type` - (Optional) The PlayReady content type. Supported values are `UltraVioletDownload`, `UltraVioletStreaming` or `Unspecified`.

* `expiration_date` - (Optional) The expiration date of license.

* `grace_period` - (Optional) The grace period of license.

* `license_type` - (Optional) The license type. Supported values are `NonPersistent` or `Persistent`.

* `play_right` - (Optional) A `play_right` block as defined above.

* `relative_begin_date` - (Optional) The relative begin date of license.

* `relative_expiration_date` - (Optional) The relative expiration date of license.

---

A `policy_option` block supports the following:

* `name` - (Required) The name which should be used for this Policy Option.

* `clear_key_configuration_enabled` - (Optional) Enable a configuration for non-DRM keys.

* `fairplay_configuration` - (Optional) A `fairplay_configuration` block as defined above. Check license requirements here https://docs.microsoft.com/en-us/azure/media-services/latest/fairplay-license-overview.

* `open_restriction_enabled` - (Optional) Enable an open restriction. License or key will be delivered on every request.

* `playready_configuration_license` - (Optional) One or more `playready_configuration_license` blocks as defined above.

* `token_restriction` - (Optional) A `token_restriction` block as defined below.

* `widevine_configuration_template` - (Optional) The Widevine template.

-> **NOTE:** Each policy_option can only have one type of configuration: fairplay_configuration,clear_key_configuration_enabled, playready_configuration_license or widevine_configuration_template. And is possible to assign only one type of restriction: open_restriction_enabled or token_restriction.

---

A `required_claim` block supports the following:

* `type` - (Optional) Token claim type.

* `value` - (Optional) Token claim value.

---

A `token_restriction` block supports the following:

* `audience` - (Optional) The audience for the token.

* `issuer` - (Optional) The token issuer.

* `open_id_connect_discovery_document` - (Optional) The OpenID connect discovery document.

* `primary_rsa_token_key_exponent` - (Optional) The RSA Parameter exponent.

* `primary_rsa_token_key_modulus` - (Optional) The RSA Parameter modulus.

* `primary_symmetric_token_key` - (Optional) The key value of the key. Specifies a symmetric key for token validation.

* `primary_x509_token_key_raw` - (Optional) The raw data field of a certificate in PKCS 12 format (X509Certificate2 in .NET). Specifies a certificate for token validation.

* `required_claim` - (Optional) One or more `required_claim` blocks as defined above.

* `token_type` - (Optional) The type of token. Supported values are `Jwt` or `Swt`.

-> **NOTE:** Each token_restriction can only have one type of primary verification key: if you want use RSA you must provide primary_rsa_token_key_exponent and primary_rsa_token_key_modulus, if you want to use symmetric you need to provide primary_symmetric_token_key and for x509 you must provide primary_x509_token_key_raw. For more information about Token access please refer to https://docs.microsoft.com/en-us/azure/media-services/latest/content-protection-overview#controlling-content-access 

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Group.

## Import

Resource Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_content_key_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaservices/account1/contentkeypolicies/policy1
```
