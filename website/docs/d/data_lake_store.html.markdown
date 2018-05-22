---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_lake_store"
sidebar_current: "docs-azurerm-datasource-data-lake-store"
description: |-
  Get information about a Data Lake Store

---

# Data Source: azurerm_data_lake_store

Use this data source to obtain information about a Data Lake Store.

## Example Usage

```hcl
data "azurerm_data_lake_store" "test" {
  name                = "testdatalake"
  resource_group_name = "testdatalake"
}

output "data_lake_store_id" {
  value = "${data.azurerm_data_lake_store.test.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the Data Lake Store.
* `resource_group_name` - (Required) The Name of the Resource Group where the Data Lake Store exists.

## Attributes Reference

* `id` - The ID of the Data Lake Store.
* `tier` - Current monthly commitment tier for the account.
* `tags` - A mapping of tags to assign to the Data Lake Store.
