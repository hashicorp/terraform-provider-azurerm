---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_cosmosdb_data_connection"
description: |-
  Manages a Data Explorer.
---

# azurerm_kusto_cosmosdb_data_connection

Manages a Data Explorer.

## Example Usage

```hcl
resource "azurerm_kusto_cosmosdb_data_connection" "example" {
  name = "example"
  resource_group_name = "example"
  location = "West Europe"
  cosmosdb_account_id = "TODO"
  cosmosdb_container = "TODO"
  managed_identity_id = "TODO"
  cluster_name = "example"
  database_name = "example"
  cosmosdb_database = "TODO"
  table_name = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `cluster_name` - (Required) TODO. Changing this forces a new Data Explorer to be created.

* `cosmosdb_account_id` - (Required) The ID of the TODO. Changing this forces a new Data Explorer to be created.

* `cosmosdb_container` - (Required) TODO. Changing this forces a new Data Explorer to be created.

* `cosmosdb_database` - (Required) TODO. Changing this forces a new Data Explorer to be created.

* `database_name` - (Required) TODO. Changing this forces a new Data Explorer to be created.

* `location` - (Required) The Azure Region where the Data Explorer should exist. Changing this forces a new Data Explorer to be created.

* `managed_identity_id` - (Required) The ID of the TODO. Changing this forces a new Data Explorer to be created.

* `name` - (Required) The name which should be used for this Data Explorer. Changing this forces a new Data Explorer to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Data Explorer should exist. Changing this forces a new Data Explorer to be created.

* `table_name` - (Required) TODO. Changing this forces a new Data Explorer to be created.

---

* `mapping_rule_name` - (Optional) TODO.

* `retrieval_start_date` - (Optional) TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Explorer.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Explorer.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Explorer.
* `update` - (Defaults to 30 minutes) Used when updating the Data Explorer.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Explorer.

## Import

Data Explorers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_cosmosdb_data_connection.example C:/Program Files/Git/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1/dataConnections/dataConnection1
```