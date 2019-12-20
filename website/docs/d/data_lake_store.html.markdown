---
subcategory: "Data Lake"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_lake_store"
sidebar_current: "docs-azurerm-datasource-data-lake-store"
description: |-
  Gets information about an existing Data Lake Store

---

# Data Source: azurerm_data_lake_store

Use this data source to access information about an existing Data Lake Store.

## Example Usage

```hcl
data "azurerm_data_lake_store" "example" {
  name                = "testdatalake"
  resource_group_name = "testdatalake"
}

output "data_lake_store_id" {
  value = "${data.azurerm_data_lake_store.example.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the Data Lake Store.

* `resource_group_name` - (Required) The Name of the Resource Group where the Data Lake Store exists.

## Attributes Reference

* `id` - The ID of the Data Lake Store.

* `encryption_state` - the Encryption State of this Data Lake Store Account, such as `Enabled` or `Disabled`.

* `encryption_type` - the Encryption Type used for this Data Lake Store Account.

* `firewall_allow_azure_ips` - are Azure Service IP's allowed through the firewall?

* `firewall_state` - the state of the firewall, such as `Enabled` or `Disabled`.

* `tier` - Current monthly commitment tier for the account.

* `tags` - A mapping of tags to assign to the Data Lake Store.
