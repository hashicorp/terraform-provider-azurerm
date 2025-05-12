---
subcategory: "Video Indexer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_video_indexer_account"
description: |-
  Manages a Video Indexer Account
---

# azurerm_video_indexer_account

Manages a Video Indexer Account

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_video_indexer_account" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"

  storage {
    storage_account_id = azurerm_storage_account.example.id
  }

  identity {
    type = "SystemAssigned"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Video Indexer Account. Changing the name forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group that the Video Indexer Account will be associated with. Changing the name forces a new resource to be created.

* `location` - (Required) The Azure location where the Video Indexer Account exists. Changing this forces a new resource to be created.

* `storage` - (Required) A `storage` block as defined below.

* `identity` - (Required) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `storage` block supports the following: 

* `storage_account_id` - (Required) The ID of the storage account to be associated with the Video Indexer Account. Changing this forces a new Video Indexer Account to be created.

* `user_assigned_identity_id` - (Optional) The reference to the user assigned identity to use to access the Storage Account.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Video Indexer Account. Possible values are `SystemAssigned` (where Azure will generate a Service Principal for you), `UserAssigned` where you can specify the Service Principal IDs in the `identity_ids` field, and `SystemAssigned, UserAssigned` which assigns both a system managed identity as well as the specified user assigned identities.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned. Required if `type` is `UserAssigned`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Video Indexer Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Video Indexer Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Video Indexer Account.
* `update` - (Defaults to 1 hour) Used when updating the Video Indexer Account.
* `delete` - (Defaults to 1 hour) Used when deleting the Video Indexer Account.

## Import

Video Indexer Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_video_indexer_account.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.VideoIndexer/accounts/example-account-name
```
