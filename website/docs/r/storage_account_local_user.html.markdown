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
  storage_account_name = azurerm_storage_account.example.name
  ssh_key_enabled      = true
  ssh_password_enabled = true
  home_directory       = "example_path"
  ssh_authorized_key {
    description = "key1"
    key         = local.second_public_key
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

* `storage_account_name` - (Required) The name of the Storage Account that this Storage Account Local User resides in. Changing this forces a new Storage Account Local User to be created.

---

* `home_directory` - (Optional) The home directory of the Storage Account Local User.

* `permission_scope` - (Optional) One or more `permission_scope` blocks as defined below.

* `ssh_authorized_key` - (Optional) One or more `ssh_authorized_key` blocks as defined below. Changing this forces a new Storage Account Local User to be created.

* `ssh_key_enabled` - (Optional) Should the SSH key be enabled?

* `ssh_password_enabled` - (Optional) Should the SSh password be enabled?

---

A `permission_scope` block supports the following:

* `permissions` - (Required) A `permissions` block as defined below.

* `resource_name` - (Required) The container name (when `service` is set to `blob`) or the file share name (when `service` is set to `file`), used by the Storage Account Local User.

* `service` - (Required) The storage service used by this Storage Account Local User. Possible values are `blob` and `file`.

---

A `permissions` block supports the following:

* `create` - (Optional) The permission to create the resources defined in this scope.

* `delete` - (Optional) The permission to delete the resources defined in this scope.

* `list` - (Optional) The permission to list the resources defined in this scope.

* `read` - (Optional) The permission to read the resources defined in this scope.

* `write` - (Optional) The permission to modify the resources defined in this scope.

---

A `ssh_authorized_key` block supports the following:

* `description` - (Required) The description of this SSH authorized key.

* `key` - (Required) The public key value of this SSH authorized key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Account Local User.

* `password` - The value of the password, which is only available when `ssh_password_enabled` is set to `true`.

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
