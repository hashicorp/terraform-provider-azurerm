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

resource "azurerm_media_streaming_policy" "example" {
  name                        = "Policy-1"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  no_encryption_enabled_protocols {
    download         = true
    dash             = true
    hls              = true
    smooth_streaming = true
  }
}
```

## Example Usage with Secure Streaming

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

resource "azurerm_media_streaming_policy" "example" {
  name                        = "Policy-1"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  common_encryption_cenc {
    enabled_protocols {
      download         = false
      dash             = true
      hls              = false
      smooth_streaming = false
    }
    drm_playready {
      custom_license_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/playready/{ContentKeyId}"
      custom_attributes                       = "PlayReady CustomAttributes"
    }
    drm_widevine_custom_license_acquisition_url_template = "https://contoso.com/{AssetAlternativeId}/widevine/{ContentKeyId"
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

* `no_encryption_enabled_protocols` - (Optional) A `no_encryption_enabled_protocols` block as defined below. Changing this forces a new Streaming Policy to be created.

---

A `common_encryption_cbcs` block supports the following:

* `default_content_key` - (Optional) A `default_content_key` block as defined below. Changing this forces a new Streaming Policy to be created.

* `drm_fairplay` - (Optional) A `drm_fairplay` block as defined below. Changing this forces a new Streaming Policy to be created.

* `enabled_protocols` - (Optional) A `enabled_protocols` block as defined below. Changing this forces a new Streaming Policy to be created.

---

A `common_encryption_cenc` block supports the following:

* `default_content_key` - (Optional) A `default_content_key` block as defined below. Changing this forces a new Streaming Policy to be created.

* `drm_playready` - (Optional) A `drm_playready` block as defined below. Changing this forces a new Streaming Policy to be created.

* `drm_widevine_custom_license_acquisition_url_template` - (Optional) TODO. Changing this forces a new Streaming Policy to be created.

* `enabled_protocols` - (Optional) A `enabled_protocols` block as defined below. Changing this forces a new Streaming Policy to be created.

---

A `default_content_key` block supports the following:

* `label` - (Optional) Label can be used to specify Content Key when creating a Streaming Locator. Changing this forces a new Streaming Policy to be created.

* `policy_name` - (Optional) Policy used by Default Key. Changing this forces a new Streaming Policy to be created.

---

A `drm_fairplay` block supports the following:

* `allow_persistent_license` - (Optional) All license to be persistent or not. Changing this forces a new Streaming Policy to be created.

* `custom_license_acquisition_url_template` - (Optional) Template for the URL of the custom service delivering licenses to end user players. Not required when using Azure Media Services for issuing licenses. The template supports replaceable tokens that the service will update at runtime with the value specific to the request. The currently supported token values are {AlternativeMediaId}, which is replaced with the value of StreamingLocatorId.AlternativeMediaId, and {ContentKeyId}, which is replaced with the value of identifier of the key being requested. Changing this forces a new Streaming Policy to be created.

---

A `drm_playready` block supports the following:

* `custom_attributes` - (Optional) Custom attributes for PlayReady. Changing this forces a new Streaming Policy to be created.

* `custom_license_acquisition_url_template` - (Optional) Template for the URL of the custom service delivering licenses to end user players. Not required when using Azure Media Services for issuing licenses. The template supports replaceable tokens that the service will update at runtime with the value specific to the request. The currently supported token values are {AlternativeMediaId}, which is replaced with the value of StreamingLocatorId.AlternativeMediaId, and {ContentKeyId}, which is replaced with the value of identifier of the key being requested. Changing this forces a new Streaming Policy to be created.

---

A `enabled_protocols` block supports the following:

* `dash` - (Optional) Enable DASH protocol or not. Changing this forces a new Streaming Policy to be created.

* `download` - (Optional) Enable Download protocol or not. Changing this forces a new Streaming Policy to be created.

* `hls` - (Optional) Enable HLS protocol or not. Changing this forces a new Streaming Policy to be created.

* `smooth_streaming` - (Optional) Enable SmoothStreaming protocol or not. Changing this forces a new Streaming Policy to be created.

---

A `no_encryption_enabled_protocols` block supports the following:

* `dash` - (Optional) Enable DASH protocol or not. Changing this forces a new Streaming Policy to be created.

* `download` - (Optional) Enable Download protocol or not. Changing this forces a new Streaming Policy to be created.

* `hls` - (Optional) Enable HLS protocol or not. Changing this forces a new Streaming Policy to be created.

* `smooth_streaming` - (Optional) Enable SmoothStreaming protocol or not. Changing this forces a new Streaming Policy to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Streaming Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Streaming Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Streaming Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Streaming Policy.

## Import

Streaming Policys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_streaming_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaservices/account1/streamingpolicies/policy1
```
