---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_storage_defender"
description: |-
    Manages the Defender for Storage. 
---

# azurerm_security_center_storage_defender

Manages the Defender for Storage.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "westus2"
}

resource "azurerm_storage_account" "example" {
  name                = "exampleacc"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_security_center_storage_defender" "example" {
  storage_account_id = azurerm_storage_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the storage account the defender applied to. Changing this forces a new resource to be created.

* `override_subscription_settings_enabled` - (Optional) Whether the settings defined for this storage account should override the settings defined for the subscription. Defaults to `false`.

* `malware_scanning_on_upload_enabled` - (Optional) Whether On Upload malware scanning should be enabled. Defaults to `false`.

* `malware_scanning_on_upload_cap_gb_per_month` - (Optional) The max GB to be scanned per Month. Must be `-1` or above `0`. Omit this property or set to `-1` if no capping is needed. Defaults to `-1`.

* `malware_scanning_on_upload_filters` - (Optional) A `malware_scanning_on_upload_filters` block as documented below.

~> **Note:** `malware_scanning_on_upload_filters` cannot be set if `malware_scanning_on_upload_enabled` is `false`

* `malware_scanning_write_results_on_tags_enabled` - (Optional) Whether to write scan result on blobIndexTags. Defaults to `true`.

* `scan_results_event_grid_topic_id` - (Optional) The Event Grid Topic where every scan result will be sent to. When you set an Event Grid custom topic, you must set `override_subscription_settings_enabled` to `true` to override the subscription-level settings.

* `sensitive_data_discovery_enabled` - (Optional) Whether Sensitive Data Discovery should be enabled. Defaults to `false`.

---

A `malware_scanning_on_upload_filters` block supports the following:

* `exclude_blobs_larger_than_in_bytes` - (Optional) Specifies the maximum size in bytes for blobs to be scanned. This parameter accepts a single positive integer value larger than `0`. Blobs larger than this value will be excluded from scanning.

* `exclude_blobs_with_prefix` - (Optional) A list of prefixes to exclude from on-upload malware scanning. Format: `container-name/blob-name` (start with the container name; do not include the storage account name). Exclude entire containers: Use prefix of container names you want to exclude without a trailing /. Exclude a single container: Add a trailing slash / after the container name to avoid excluding other containers with similar prefixes.

* `exclude_blobs_with_suffix` - (Optional) A list of suffixes to exclude from on-upload malware scanning. Suffixes match only the end of blob names, and should be used for file extensions or blob name endings only.
 
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Defender for Storage id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Defender for Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Defender for Storage.
* `update` - (Defaults to 30 minutes) Used when updating the Defender for Storage.
* `delete` - (Defaults to 30 minutes) Used when deleting the Defender for Storage.

## Import

The setting can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_storage_defender.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/storageacc
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Security` - 2025-06-01
