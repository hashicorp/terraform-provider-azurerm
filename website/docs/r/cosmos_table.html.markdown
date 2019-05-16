---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmos_table"
sidebar_current: "docs-azurerm-resource-cosmos-table"
description: |-
  Manages a Cosmos Table.
---

# azurerm_cosmos_table

Manages a Cosmos Table.

## Example Usage

```hcl
resource "azurerm_cosmos_table" "table" {
  name                = "tfex-cosmos-table"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  account_name        = "${azurerm_cosmosdb_account.account.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos Table. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos Table is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos Table to create the table within. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The Cosmos Table ID.

## Import

Cosmos Tables can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmos_table.t1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/apis/table/tables/t1
```
