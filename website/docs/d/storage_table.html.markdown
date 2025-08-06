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

* `storage_account_name` - (Required) The name of the Storage Account where the Table exists.

## Attributes Reference

* `acl` - A mapping of ACLs for this Table.

* `id` - The ID of the Storage Table.

* `resource_manager_id` - The Resource Manager ID of this Storage Table.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage.
