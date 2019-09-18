---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_lake_analytics_account"
sidebar_current: "docs-azurerm-resource-data-lake-analytics-account"
description: |-
  Manages an Azure Data Lake Analytics Account.
---

# azurerm_data_lake_analytics_account

Manages an Azure Data Lake Analytics Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-datalake-account"
  location = "northeurope"
}

resource "azurerm_data_lake_store" "example" {
  name                = "tfexdatalakestore"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
}

resource "azurerm_data_lake_analytics_account" "example" {
  name                = "tfexdatalakeaccount"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  default_store_account_name = "${azurerm_data_lake_store.example.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Lake Analytics Account. Changing this forces a new resource to be created. Has to be between 3 to 24 characters.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Lake Analytics Account.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `default_store_account_name` - (Required) Specifies the data lake store to use by default. Changing this forces a new resource to be created.

* `tier` - (Optional) The monthly commitment tier for Data Lake Analytics Account. Accepted values are `Consumption`, `Commitment_100000AUHours`, `Commitment_10000AUHours`, `Commitment_1000AUHours`, `Commitment_100AUHours`, `Commitment_500000AUHours`, `Commitment_50000AUHours`, `Commitment_5000AUHours`, or `Commitment_500AUHours`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Date Lake Store ID.

## Import

Date Lake Analytics Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_lake_analytics_account.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DataLakeAnalytics/accounts/mydatalakeaccount
```
