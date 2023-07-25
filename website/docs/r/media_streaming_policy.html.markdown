---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_streaming_policy"
description: |-
  Manages a Streaming Policy.
---

# azurerm_media_streaming_policy

Manages a Streaming Policy.

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
      ask                       = "bb566284cc124a21c435a92cd3c108c4"
      pfx                       = "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"
      pfx_password              = "password"
      rental_duration_seconds   = 2249
      rental_and_lease_key_type = "PersistentUnlimited"

    }
    open_restriction_enabled = true
  }
}

resource "azurerm_media_streaming_policy" "example" {
  name                        = "Policy-1"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  common_encryption_cenc {
    clear_track {
      condition {
        property  = "FourCC"
        operation = "Equal"
        value     = "hev2"
      }
    }

    enabled_protocols {
      download         = false
      dash             = true
      hls              = false
      smooth_streaming = false
    }

    default_content_key {
      label       = "aesDefaultKey"
      policy_name = azurerm_media_content_key_policy.example.name
    }

    drm_playready {
      custom_license_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/playready/{ContentKeyId}"
      custom_attributes                       = "PlayReady CustomAttributes"
    }
    drm_widevine_custom_license_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/widevine/{ContentKeyId}"
  }

  common_encryption_cbcs {
    enabled_protocols {
      download         = false
      dash             = true
      hls              = false
      smooth_streaming = false
    }
    drm_fairplay {
      custom_license_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/fairplay/{ContentKeyId}"
      allow_persistent_license                = true
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `media_services_account_name` - (Required) The Media Services account name. Changing this forces a new Streaming Policy to be created.

* `name` - (Required) The name which should be used for this Streaming Policy. Changing this forces a new Streaming Policy to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Streaming Policy should exist. Changing this forces a new Streaming Policy to be created.

---

* `common_encryption_cbcs` - (Optional) A `common_encryption_cbcs` block as defined below. Changing this forces a new Streaming Policy to be created.

* `common_encryption_cenc` - (Optional) A `common_encryption_cenc` block as defined below. Changing this forces a new Streaming Policy to be created.

* `default_content_key_policy_name` - (Optional) Default Content Key used by current Streaming Policy. Changing this forces a new Streaming Policy to be created.

* `envelope_encryption` - (Optional) A `envelope_encryption` block as defined below. Changing this forces a new Streaming Policy to be created. 

* `no_encryption_enabled_protocols` - (Optional) A `no_encryption_enabled_protocols` block as defined below. Changing this forces a new Streaming Policy to be created.

---

A `common_encryption_cbcs` block supports the following:

* `clear_key_encryption` - (Optional) A `clear_key_encryption` block as defined below. Changing this forces a new Streaming Policy to be created.

* `default_content_key` - (Optional) A `default_content_key` block as defined below. Changing this forces a new Streaming Policy to be created.

* `drm_fairplay` - (Optional) A `drm_fairplay` block as defined below. Changing this forces a new Streaming Policy to be created.

* `enabled_protocols` - (Optional) A `enabled_protocols` block as defined below. Changing this forces a new Streaming Policy to be created.

---

A `common_encryption_cenc` block supports the following:

* `clear_key_encryption` - (Optional) A `clear_key_encryption` block as defined below. Changing this forces a new Streaming Policy to be created.

* `clear_track` - (Optional) One or more `clear_track` blocks as defined below. Changing this forces a new Streaming Policy to be created.

* `content_key_to_track_mapping` - (Optional) One or more `content_key_to_track_mapping` blocks as defined below. Changing this forces a new Streaming Policy to be created.

* `default_content_key` - (Optional) A `default_content_key` block as defined below. Changing this forces a new Streaming Policy to be created.

* `drm_playready` - (Optional) A `drm_playready` block as defined below. Changing this forces a new Streaming Policy to be created.

* `drm_widevine_custom_license_acquisition_url_template` - (Optional) The URL template for the custom service that delivers licenses to the end user. This is not required when using Azure Media Services for issuing licenses. Changing this forces a new Streaming Policy to be created.

* `enabled_protocols` - (Optional) A `enabled_protocols` block as defined below. Changing this forces a new Streaming Policy to be created.

---

A `clear_key_encryption` block supports the following:

* `custom_keys_acquisition_url_template` - (Required) The URL template for the custom service that delivers content keys to the end user. This is not required when using Azure Media Services for issuing keys. Changing this forces a new Streaming Policy to be created.

-> **Note** Either `clear_key_encryption` or `drm` must be specified.

---

A `clear_track` block supports the following:

* `condition` - (Required) One or more `condition` blocks as defined below. Changing this forces a new Streaming Policy to be created. 

---

A `condition` block supports the following:

* `operation` - (Required) The track property condition operation. Possible value is `Equal`. Changing this forces a new Streaming Policy to be created.

* `property` - (Required) The track property type. Possible value is `FourCC`. Changing this forces a new Streaming Policy to be created.

* `value` - (Required) The track property value. Changing this forces a new Streaming Policy to be created.

---

A `content_key_to_track_mapping` block supports the following:

* `label` - (Optional) Specifies the content key when creating a Streaming Locator. Changing this forces a new Streaming Policy to be created.

* `policy_name` - (Optional) The policy used by the default key. Changing this forces a new Streaming Policy to be created.

* `track` - (Optional) One or more `track` blocks as defined below. Changing this forces a new Streaming Policy to be created.

---

A `default_content_key` block supports the following:

* `label` - (Optional) Label can be used to specify Content Key when creating a Streaming Locator. Changing this forces a new Streaming Policy to be created.

* `policy_name` - (Optional) Policy used by Default Key. Changing this forces a new Streaming Policy to be created.

---

A `drm_fairplay` block supports the following:

* `allow_persistent_license` - (Optional) All license to be persistent or not. Changing this forces a new Streaming Policy to be created.

* `custom_license_acquisition_url_template` - (Optional) The URL template for the custom service that delivers licenses to the end user. This is not required when using Azure Media Services for issuing licenses. Changing this forces a new Streaming Policy to be created.

---

A `drm_playready` block supports the following:

* `custom_attributes` - (Optional) Custom attributes for PlayReady. Changing this forces a new Streaming Policy to be created.

* `custom_license_acquisition_url_template` - (Optional) The URL template for the custom service that delivers licenses to the end user. This is not required when using Azure Media Services for issuing licenses. Changing this forces a new Streaming Policy to be created.

---

A `enabled_protocols` block supports the following:

* `dash` - (Optional) Enable DASH protocol or not. Changing this forces a new Streaming Policy to be created.

* `download` - (Optional) Enable Download protocol or not. Changing this forces a new Streaming Policy to be created.

* `hls` - (Optional) Enable HLS protocol or not. Changing this forces a new Streaming Policy to be created.

* `smooth_streaming` - (Optional) Enable SmoothStreaming protocol or not. Changing this forces a new Streaming Policy to be created.

---

A `envelope_encryption` block supports the following:

* `custom_keys_acquisition_url_template` - (Optional) The URL template for the custom service that delivers content keys to the end user. This is not required when using Azure Media Services for issuing keys. Changing this forces a new Streaming Policy to be created.

* `default_content_key` - (Optional) A `default_content_key` block as defined above. Changing this forces a new Streaming Policy to be created.

* `enabled_protocols` - (Optional) A `enabled_protocols` block as defined above. Changing this forces a new Streaming Policy to be created.

---

A `no_encryption_enabled_protocols` block supports the following:

* `dash` - (Optional) Enable DASH protocol or not. Changing this forces a new Streaming Policy to be created.

* `download` - (Optional) Enable Download protocol or not. Changing this forces a new Streaming Policy to be created.

* `hls` - (Optional) Enable HLS protocol or not. Changing this forces a new Streaming Policy to be created.

* `smooth_streaming` - (Optional) Enable SmoothStreaming protocol or not. Changing this forces a new Streaming Policy to be created.

---

A `track` block supports the following:

* `condition` - (Required) One or more `condition` blocks as defined below. Changing this forces a new Streaming Policy to be created. 


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Streaming Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Streaming Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Streaming Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Streaming Policy.

## Import

Streaming Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_streaming_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaServices/account1/streamingPolicies/policy1
```
