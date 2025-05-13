---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_cluster_managed_private_endpoint"
description: |-
  Manages a Managed Private Endpoint for a Kusto Cluster.
---

# azurerm_kusto_cluster_managed_private_endpoint

Manages a Managed Private Endpoint for a Kusto Cluster.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
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

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_kusto_cluster_managed_private_endpoint" "example" {
  name                         = "examplempe"
  resource_group_name          = azurerm_resource_group.example.name
  cluster_name                 = azurerm_kusto_cluster.example.name
  private_link_resource_id     = azurerm_storage_account.example.id
  private_link_resource_region = azurerm_storage_account.example.location
  group_id                     = "blob"
  request_message              = "Please Approve"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Private Endpoints to create. Changing this forces a new resource to be created.

* `cluster_name` - (Required) The name of the Kusto Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Kusto Cluster should exist. Changing this forces a new resource to be created.

* `private_link_resource_id` - (Required) The ARM resource ID of the resource for which the managed private endpoint is created. Changing this forces a new resource to be created.

* `group_id` - (Required) The group id in which the managed private endpoint is created. Changing this forces a new resource to be created.

* `private_link_resource_region` - (Optional) The region of the resource to which the managed private endpoint is created. Changing this forces a new resource to be created.

* `request_message` - (Optional) The user request message.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Kusto Cluster Customer Managed Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Cluster Customer Managed Key.
* `update` - (Defaults to 1 hour) Used when updating the Kusto Cluster Customer Managed Key.
* `delete` - (Defaults to 1 hour) Used when deleting the Kusto Cluster Customer Managed Key.

## Import

Managed Private Endpoint for a Kusto Cluster can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_cluster_managed_private_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1/managedPrivateEndpoints/managedPrivateEndpoint1
```
