---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_static_website"
description: |-
  Gets information about an existing Static Website of an Azure Storage Account.
---

# Data Source: azurerm_storage_account_static_website

Use this data source to access information about an existing Static Website of an Azure Storage Account.

## Example Usage

```hcl
data "azurerm_storage_account_static_website" "example" {
  storage_account_id = "Existing"

}

output "id" {
  value = data.azurerm_storage_account_static_website.example.id

}
```

## Arguments Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the Storage Account to set Static Website on.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Account Static Website.

* `error_404_document` - The absolute path to a custom webpage that is used when a request is made which does not correspond to an existing file.

* `index_document` - The webpage that Azure Storage serves for requests to the root of a website or any subfolder.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Static Website.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Storage` - 2025-08-01
