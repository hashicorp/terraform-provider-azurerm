---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_streaming_locator"
description: |-
  Manages a Media Streaming Locator.
---

# azurerm_media_streaming_locator

Manages a Media Streaming Locator.

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

resource "azurerm_media_streaming_locator" "example" {
  name                        = "example"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  asset_name                  = azurerm_media_asset.example.name
  streaming_policy_name       = "Predefined_ClearStreamingOnly"
}
```

## Arguments Reference

The following arguments are supported:

* `asset_name` - (Required) Asset Name. Changing this forces a new Streaming Locator to be created.

* `media_services_account_name` - (Required) The Media Services account name. Changing this forces a new Streaming Locator to be created.

* `name` - (Required) The name which should be used for this Streaming Locator. Changing this forces a new Streaming Locator to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Streaming Locator should exist. Changing this forces a new Streaming Locator to be created.

* `streaming_policy_name` - (Required) Name of the Streaming Policy used by this Streaming Locator. Either specify the name of Streaming Policy you created or use one of the predefined Streaming Policies. The predefined Streaming Policies available are: `Predefined_DownloadOnly`, `Predefined_ClearStreamingOnly`, `Predefined_DownloadAndClearStreaming`, `Predefined_ClearKey`, `Predefined_MultiDrmCencStreaming` and `Predefined_MultiDrmStreaming`. Changing this forces a new Streaming Locator to be created.

---

* `alternative_media_id` - (Optional) Alternative Media ID of this Streaming Locator. Changing this forces a new Streaming Locator to be created.

* `content_key` - (Optional) One or more `content_key` blocks as defined below. Changing this forces a new Streaming Locator to be created.

* `default_content_key_policy_name` - (Optional) Name of the default Content Key Policy used by this Streaming Locator.Changing this forces a new Streaming Locator to be created.

* `end_time` - (Optional) The end time of the Streaming Locator. Changing this forces a new Streaming Locator to be created.

* `start_time` - (Optional) The start time of the Streaming Locator. Changing this forces a new Streaming Locator to be created.

* `streaming_locator_id` - (Optional) The ID of the Streaming Locator. Changing this forces a new Streaming Locator to be created.

---

A `content_key` block supports the following:

* `content_key_id` - (Optional) ID of Content Key. Changing this forces a new Streaming Locator to be created.

* `label_reference_in_streaming_policy` - (Optional) Label of Content Key as specified in the Streaming Policy. Changing this forces a new Streaming Locator to be created.

* `policy_name` - (Optional) Content Key Policy used by Content Key. Changing this forces a new Streaming Locator to be created.

* `type` - (Optional) Encryption type of Content Key. Supported values are `CommonEncryptionCbcs`, `CommonEncryptionCenc` or `EnvelopeEncryption`. Changing this forces a new Streaming Locator to be created.

* `value` - (Optional) Value of Content Key. Changing this forces a new Streaming Locator to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Streaming Locator.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Streaming Locator.
* `read` - (Defaults to 5 minutes) Used when retrieving the Streaming Locator.
* `delete` - (Defaults to 30 minutes) Used when deleting the Streaming Locator.

## Import

Streaming Locators can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_streaming_locator.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Media/mediaservices/account1/streaminglocators/locator1
```
