---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_share_dataset_kusto_database"
description: |-
  Manages a Data Share Kusto Database Dataset.
---

# azurerm_data_share_dataset_kusto_database

Manages a Data Share Kusto Database Dataset.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_share_account" "example" {
  name                = "example-dsa"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_share" "example" {
  name       = "example_ds"
  account_id = azurerm_data_share_account.example.id
  kind       = "InPlace"
}

resource "azurerm_kusto_cluster" "example" {
  name                = "examplekc"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "example" {
  name                = "examplekd"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.example.name
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_kusto_cluster.example.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_data_share_account.example.identity[0].principal_id
}

resource "azurerm_data_share_dataset_kusto_database" "example" {
  name              = "example-dskd"
  share_id          = azurerm_data_share.example.id
  kusto_database_id = azurerm_kusto_database.example.id
  depends_on = [
    azurerm_role_assignment.example,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Data Share Kusto Database Dataset. Changing this forces a new Data Share Kusto Database Dataset to be created.

* `share_id` - (Required) The resource ID of the Data Share where this Data Share Kusto Database Dataset should be created. Changing this forces a new Data Share Kusto Database Dataset to be created.

* `kusto_database_id` - (Required) The resource ID of the Kusto Cluster Database to be shared with the receiver. Changing this forces a new Data Share Kusto Database Dataset to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The resource ID of the Data Share Kusto Database Dataset.

* `display_name` - The name of the Data Share Dataset.

* `kusto_cluster_location` - The location of the Kusto Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Share Kusto Database Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Kusto Database Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Share Kusto Database Dataset.

## Import

Data Share Kusto Database Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_share_dataset_kusto_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataShare/accounts/account1/shares/share1/dataSets/dataSet1
```
