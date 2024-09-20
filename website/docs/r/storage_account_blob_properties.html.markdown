---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_blob_properties"
description: |-
  Manages a Azure Storage Account Blob Properties.
---

# azurerm_storage_account_blob_properties

Manages an Azure Storage Accounts Blob Properties.

## Disclaimers

~> **Note on Storage Accounts and Blob Properties:** Terraform currently provides both a standalone [Blob Properties resource](storage_account_blob_properties.html), and allows for Blob Properties to be defined in-line within the [Storage Account resource](storage_account.html). At this time you cannot use a Storage Account with in-line Blob Properties in conjunction with any Blob Properties resource. Doing so will cause a conflict of Blob Properties configurations and will overwrite the in-line Blob Properties.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account_blob_properties" "example" {
  storage_account_id = azurerm_storage_account.example.id

  properties {
    delete_retention_policy {}

    restore_policy {
      days = 6
    }

    versioning_enabled  = true
    change_feed_enabled = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) Specifies the resource id of the storage account.

* `properties` - (Required) A `properties` block as defined below.

---

A `properties` block supports the following:

* `cors_rule` - (Optional) A `cors_rule` block as defined below.

* `delete_retention_policy` - (Optional) A `delete_retention_policy` block as defined below.

* `restore_policy` - (Optional) A `restore_policy` block as defined below. This must be used together with `delete_retention_policy` set, `versioning_enabled` and `change_feed_enabled` set to `true`.

-> **Note:** This field cannot be configured when `kind` is set to `Storage` (V1).

-> **Note:** `restore_policy` can not be configured when `dns_endpoint_type` is `AzureDnsZone`.

* `versioning_enabled` - (Optional) Is versioning enabled? Default to `false`.

-> **Note:** This field cannot be configured when `kind` is set to `Storage` (V1).

* `change_feed_enabled` - (Optional) Is the blob service properties for change feed events enabled? Default to `false`.

-> **Note:** This field cannot be configured when `kind` is set to `Storage` (V1).

* `change_feed_retention_in_days` - (Optional) The duration of change feed events retention in days. The possible values are between 1 and 146000 days (400 years). Setting this to null (or omit this in the configuration file) indicates an infinite retention of the change feed.

-> **Note:** This field cannot be configured when `kind` is set to `Storage` (V1).

* `default_service_version` - (Optional) The API Version which should be used by default for requests to the Data Plane API if an incoming request doesn't specify an API Version.

* `last_access_time_enabled` - (Optional) Is the last access time based tracking enabled? Default to `false`.

-> **Note:** This field cannot be configured when `kind` is set to `Storage` (V1).

* `container_delete_retention_policy` - (Optional) A `container_delete_retention_policy` block as defined below.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Account Blob Properties.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Storage Account Blob Properties.
* `update` - (Defaults to 60 minutes) Used when updating the Storage Account Blob Properties.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Blob Properties.
* `delete` - (Defaults to 60 minutes) Used when deleting the Storage Account Blob Properties.

## Import

Storage Account Blob Properties can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account_blob_properties.blob1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
