---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_table"
description: |-
  Manages a Kusto / Data Explorer Table.
---

# azurerm_kusto_table

Manages a Kusto (also known as Azure Data Explorer) Table

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "East US"
}

resource "azurerm_kusto_cluster" "example" {
  name                = "example-cluster"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }
}

resource "azurerm_kusto_database" "example" {
  name                = "example-database"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.example.name

  hot_cache_period   = "P7D"
  soft_delete_period = "P31D"
}

resource "azurerm_kusto_table" "example" {
  name        = "exampletable"
  database_id = azurerm_kusto_database.example.id
  doc         = "Optional documentation"
  folder      = "Optional folder"

  column {
    name = "col_1"
    type = "string"
  }

  column {
    name = "col_2"
    type = "real"
  }

  column {
    name = "col_n"
    type = "dynamic"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Kusto Table. Changing this forces a new Kusto Table to be created.

* `database_id` - (Required) The ARM ID of the Kusto Database within which this table should be created. Changing this forces a new Kusto Table to be created.

* `column` - (Required) One or more `column` blocks as defined below.

~> **NOTE:** For table resource creation the Kusto Management command `.create table` is used. For resource update the Kusto Management commands `.alter table`, `.alter table T folder` as well as `.alter table T docstring` are used. Refer to the [command docs](https://docs.microsoft.com/de-de/azure/data-explorer/kusto/management/tables) to learn what effect changes to columns blocks will have.

~> **NOTE:** Changes to column data types are not supported so far (see `.alter table` docs). If this is attempted anyway, it will result in an Kusto specific error.

---

* `doc` - (Optional) An optional documentation string for the table.

* `folder` - (Optional) An optional folder name within this table should be created.

---

A `column` block supports the following:

* `name` - (Required) The name of the column.

* `type` - (Required) The type of the column. Possible values are: `bool`, `datetime`, `decimal`, `dynamic`, `guid`, `int`, `real`, `string` and `timespan`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Kusto Table.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Kusto Table.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Table.
* `update` - (Defaults to 1 hour) Used when updating the Kusto Table.
* `delete` - (Defaults to 1 hour) Used when deleting the Kusto Table.

## Import

Kusto Tables can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_table.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1/tables/table1
```
