---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_local_user"
description: |-
  Manages a Storage Account Local User.
---

# azurerm_storage_account_local_user

Manages a Storage Account Local User.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "WestEurope"
}

resource "azurerm_storage_account" "example" {
  name                     = "example-account"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}

resource "azurerm_storage_container" "example" {
  name                 = "example-container"
  storage_account_name = azurerm_storage_account.example.name
}

resource "azurerm_storage_account_local_user" "example" {
  name                 = "user1"
  storage_account_id   = azurerm_storage_account.example.id
  ssh_key_enabled      = true
  ssh_password_enabled = true
  home_directory       = "example_path"
  ssh_authorized_key {
    description = "key1"
    key         = local.first_public_key
  }
  ssh_authorized_key {
    description = "key2"
    key         = local.second_public_key
  }
  permission_scope {
    permissions {
      read   = true
      create = true
    }
    service       = "blob"
    resource_name = azurerm_storage_container.example.name
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Storage Account Local User. Changing this forces a new Storage Account Local User to be created.

* `storage_account_id` - (Required) The ID of the Storage Account that this Storage Account Local User resides in. Changing this forces a new Storage Account Local User to be created.

---

* `home_directory` - (Optional) The home directory of the Storage Account Local User.

* `permission_scope` - (Optional) One or more `permission_scope` blocks as defined below.

* `ssh_authorized_key` - (Optional) One or more `ssh_authorized_key` blocks as defined below.

* `ssh_key_enabled` - (Optional) Specifies whether SSH Key Authentication is enabled. Defaults to `false`.

* `ssh_password_enabled` - (Optional) Specifies whether SSH Password Authentication is enabled. Defaults to `false`.

---

A `permission_scope` block supports the following:

* `permissions` - (Required) A `permissions` block as defined below.

* `resource_name` - (Required) The container name (when `service` is set to `blob`) or the file share name (when `service` is set to `file`), used by the Storage Account Local User.

* `service` - (Required) The storage service used by this Storage Account Local User. Possible values are `blob` and `file`.

---

A `permissions` block supports the following:

* `create` - (Optional) Specifies if the Local User has the create permission for this scope. Defaults to `false`.

* `delete` - (Optional) Specifies if the Local User has the delete permission for this scope. Defaults to `false`.

* `list` - (Optional) Specifies if the Local User has the list permission for this scope. Defaults to `false`.

* `read` - (Optional) Specifies if the Local User has the read permission for this scope. Defaults to `false`.

* `write` - (Optional) Specifies if the Local User has the write permission for this scope. Defaults to `false`.

---

A `ssh_authorized_key` block supports the following:

* `key` - (Required) The public key value of this SSH authorized key.

* `description` - (Optional) The description of this SSH authorized key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Account Local User.

* `password` - The value of the password, which is only available when `ssh_password_enabled` is set to `true`.

~> **Note:** The `password` will be updated everytime when `ssh_password_enabled` got updated. If `ssh_password_enabled` is updated from `false` to `true`, the `password` is updated to be the value of the SSH password. If `ssh_password_enabled` is updated from `true` to `false`, the `password` is reset to empty string.

* `sid` - The unique Security Identifier of this Storage Account Local User.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Account Local User.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Local User.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Account Local User.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Account Local User.

## Import

Storage Account Local Users can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account_local_user.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/localUsers/user1
```
