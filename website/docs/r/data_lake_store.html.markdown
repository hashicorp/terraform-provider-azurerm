---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_lake_store"
sidebar_current: "docs-azurerm-resource-data-lake-store"
description: |-
  Manage an Azure Data Lake Store.
---

# azurerm_data_lake_store

Manage an Azure Data Lake Store.

## Example Usage

```hcl
# Pay As You Go
resource "azurerm_resource_group" "test" {
  name     = "test"
  location = "northeurope"
}

resource "azurerm_data_lake_store" "consumption" {
  name                = "consumptiondatalake"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "northeurope"
}

# Monthly Commitment Tier
resource "azurerm_resource_group" "test" {
  name     = "test"
  location = "westeurope"
}

resource "azurerm_data_lake_store" "monthly" {
  name                = "monthlydatalake"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "westeurope"
  tier                = "Commitment_1TB"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Lake Store. Changing this forces a new resource to be created. Has to be between 3 to 24 characters.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Lake Store.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tier` - (Optional) The monthly commitment tier for Data Lake Store. Accepted values are `Consumption`, `Commitment_1TB`, `Commitment_10TB`, `Commitment_100TB`, `Commitment_500TB`, `Commitment_1PB` or `Commitment_5PB`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Date Lake Store ID.

## Import

Date Lake Store can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_lake_store.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DataLakeStore/accounts/mydatalakeaccount
```
