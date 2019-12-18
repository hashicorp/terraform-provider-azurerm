---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_table"
sidebar_current: "docs-azurerm-resource-cosmosdb-table"
description: |-
  Manages a Table within a Cosmos DB Account.
---

# azurerm_cosmosdb_table

Manages a Table within a Cosmos DB Account.

## Example Usage

```hcl
data "azurerm_cosmosdb_account" "example" {
  name                = "tfex-cosmosdb-account"
  resource_group_name = "tfex-cosmosdb-account-rg"
}

resource "azurerm_cosmosdb_table" "example" {
  name                = "tfex-cosmos-table"
  resource_group_name = "${data.azurerm_cosmosdb_account.example.resource_group_name}"
  account_name        = "${data.azurerm_cosmosdb_account.example.name}"
  throughput          = 400
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB Table. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB Table is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB Table to create the table within. Changing this forces a new resource to be created.

* `throughput` - (Optional) The throughput of Table (RU/s). Must be set in increments of `100`. The minimum value is `400`. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.


## Attributes Reference

The following attributes are exported:

* `id` - the Cosmos DB Table ID.

## Import

Cosmos Tables can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_table.t1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/apis/table/tables/t1
```
