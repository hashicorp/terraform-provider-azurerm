---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_table"
description: |-
  Gets information about an existing Storage Table.
---

# Data Source: azurerm_storage_table

Use this data source to access information about an existing Storage Table.

## Example Usage

```hcl
data "azurerm_storage_table" "example" {
  name                 = "example-table-name"
  storage_account_name = "example-storage-account-name"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Table.

* `storage_account_name` - (Optional) The name of the Storage Account where the Table exists. This property is deprecated in favour of `storage_account_id`.

* `storage_account_id` - (Optional) The name of the Storage Account where the Table exists. This property will become Required in version 5.0 of the Provider.

~> **Note:** One of `storage_account_name` or `storage_account_id` must be specified. When specifying `storage_account_id` the resource will use the Resource Manager API, rather than the Data Plane API.

## Attributes Reference

* `acl` - A mapping of ACLs for this Table.

* `id` - The ID of the Storage Table.

* `resource_manager_id` - The Resource Manager ID of this Storage Table.

~> **Note:** `resource_manager_id` will be removed in version 5.0 of the Provider.

* `url` - The data plane URL of the Storage Table in the format of `<storage table endpoint>/Tables('<table name>')`. E.g. `https://example.table.core.windows.net/Tables('mytable')"`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Storage` - 2023-05-01
