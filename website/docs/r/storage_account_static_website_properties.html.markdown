---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_static_website_properties"
description: |-
  Manages an Azure Storage Accounts Static Website Properties.
---

# azurerm_storage_account_static_website_properties

Manages an Azure Storage Accounts Static Website Properties.

## Disclaimers

~> **Note on Storage Accounts and Static Web Site Properties:** Terraform currently provides both a standalone [Static Website Properties resource](storage_account_static_website_properties.html), and allows for Static Website Properties to be defined in-line within the [Storage Account resource](storage_account.html). At this time you cannot use a Storage Account with in-line Static Website Properties in conjunction with any Static Website Properties resource. Doing so will cause a conflict of Static Website Properties configurations and will overwrite the in-line Static Website Properties.

~> **Note:** An `azurerm_storage_account_static_website_properties` resource can only be defined when the referenced storage accounts `account_kind` is set to `StorageV2` or `BlockBlobStorage`.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_storage_account" "example" {
  name                = "storageaccountname"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account_static_website_properties" "example" {
  storage_account_id = azurerm_storage_account.example.id

  index_document     = "index.html"
  error_404_document = "error.html"
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) Specifies the resource id of the storage account.

* `index_document` - (Optional) The webpage that Azure Storage serves for requests to the root of a website or any subfolder (e.g., `index.html`). The value is case-sensitive.

* `error_404_document` - (Optional) The absolute path to a custom webpage that should be used when a request is made which does not correspond to an existing file.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Account Static Website Properties.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Storage Account Static Website Properties.
* `update` - (Defaults to 60 minutes) Used when updating the Storage Account Static Website Properties.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Static Website Properties.
* `delete` - (Defaults to 60 minutes) Used when deleting the Storage Account Static Website Properties.

## Import

Storage Accounts Static Website Properties can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account_static_website_properties.webSite1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
