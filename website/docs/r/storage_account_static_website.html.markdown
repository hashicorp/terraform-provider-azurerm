---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_static_website"
description: |-
  Manages the Static Website of an Azure Storage Account.
---

# azurerm_storage_account_static_website

Manages the Static Website of an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
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

resource "azurerm_storage_account_static_website" "example" {
  storage_account_id = azurerm_storage_account.example.id
  error_404_document = "custom_not_found.html"
  index_document     = "custom_index.html"
}
```

## Arguments Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the Storage Account to set Static Website on. Changing this forces a new resource to be created.

* `error_404_document` - (Optional) The absolute path to a custom webpage that should be used when a request is made which does not correspond to an existing file.

* `index_document` - (Optional) The webpage that Azure Storage serves for requests to the root of a website or any subfolder. For example, index.html.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Account Static Website.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Account Static Website.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Static Website.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Account Static Website.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Account Static Website.

## Import

`azurerm_storage_account_static_website` resources can be imported using one of the following methods:

* The `terraform import` CLI command with an `id` string:

  ```shell
  terraform import azurerm_storage_account_static_website.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{storageAccountName}
  ```

* An `import` block with an `id` argument:
  
  ```hcl
  import {
    to = azurerm_storage_account_static_website.example
    id = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{storageAccountName}"
  }
  ```

* An `import` block with an `identity` argument:

  ```hcl
  import {
    to       = azurerm_storage_account_static_website.example
    identity = {
      TODO Resource Identity Format
    }
  }
  ```
