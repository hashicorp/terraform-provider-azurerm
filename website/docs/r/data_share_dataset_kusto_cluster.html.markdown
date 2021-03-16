---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_share_dataset_kusto_cluster"
description: |-
  Manages a Data Share Kusto Cluster Dataset.
---

# azurerm_data_share_dataset_kusto_cluster

Manages a Data Share Kusto Cluster Dataset.

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

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_kusto_cluster.example.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_data_share_account.example.identity.0.principal_id
}

resource "azurerm_data_share_dataset_kusto_cluster" "example" {
  name             = "example-dskc"
  share_id         = azurerm_data_share.example.id
  kusto_cluster_id = azurerm_kusto_cluster.example.id
  depends_on = [
    azurerm_role_assignment.example,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Data Share Kusto Cluster Dataset. Changing this forces a new Data Share Kusto Cluster Dataset to be created.

* `share_id` - (Required) The resource ID of the Data Share where this Data Share Kusto Cluster Dataset should be created. Changing this forces a new Data Share Kusto Cluster Dataset to be created.

* `kusto_cluster_id` - (Required) The resource ID of the Kusto Cluster to be shared with the receiver. Changing this forces a new Data Share Kusto Cluster Dataset to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The resource ID of the Data Share Kusto Cluster Dataset.

* `display_name` - The name of the Data Share Dataset.

* `kusto_cluster_location` - The location of the Kusto Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Share Kusto Cluster Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Kusto Cluster Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Share Kusto Cluster Dataset.

## Import

Data Share Kusto Cluster Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_share_dataset_kusto_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataShare/accounts/account1/shares/share1/dataSets/dataSet1
```
