---
subcategory: "Dev Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_test_lab"
sidebar_current: "docs-azurerm-datasource-dev-test-lab"
description: |-
  Gets information about an existing Dev Test Lab.
---

# Data Source: azurerm_dev_test_lab

Use this data source to access information about an existing Dev Test Lab.

## Example Usage

```hcl
data "azurerm_dev_test_lab" "example" {
  name                = "example-lab"
  resource_group_name = "example-resources"
}

output "unique_identifier" {
  value = "${data.azurerm_dev_test_lab.example.unique_identifier}"
}
```

## Argument Reference

* `name` - (Required) The name of the Dev Test Lab.

* `resource_group_name` - (Required) The Name of the Resource Group where the Dev Test Lab exists.

## Attributes Reference

* `id` - The ID of the Dev Test Lab.

* `artifacts_storage_account_id` - The ID of the Storage Account used for Artifact Storage.

* `default_storage_account_id` - The ID of the Default Storage Account for this Dev Test Lab.

* `default_premium_storage_account_id` - The ID of the Default Premium Storage Account for this Dev Test Lab.

* `key_vault_id` - The ID of the Key used for this Dev Test Lab.

* `location` - The Azure location where the Dev Test Lab exists.

* `premium_data_disk_storage_account_id` - The ID of the Storage Account used for Storage of Premium Data Disk.

* `storage_type` - The type of storage used by the Dev Test Lab.

* `tags` - A mapping of tags to assign to the resource.

* `unique_identifier` - The unique immutable identifier of the Dev Test Lab.
