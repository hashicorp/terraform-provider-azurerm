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

resource "azurerm_storage_account_static_website" "test" {
  storage_account_id = azurerm_storage_account.test.id
  error_404_document = "custom_not_found.html"
  index_document     = "custom_index.html"
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the Storage Account to set Static Website on. Changing this forces a new resource to be created.

* `error_404_document` - (Optional) The absolute path to a custom webpage that should be used when a request is made which does not correspond to an existing file.

* `index_document` - (Optional) The webpage that Azure Storage serves for requests to the root of a website or any subfolder. For example, index.html.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Account Static Website.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Static Website.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Account Static Website.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Account Static Website.

## Import

Storage Account Static Websites can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account_static_website.mysite /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
