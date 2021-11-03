---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_cassandra_cluster"
description: |-
  Manages a Cassandra Cluster.
---

# azurerm_cosmosdb_cassandra_cluster

Manages a Cassandra Cluster.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "accexample-rg"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_virtual_network.example.id
  role_definition_name = "Network Contributor"
  principal_id         = "e5007d2c-4b13-4a74-9b6a-605d99f03501"
}

resource "azurerm_cosmosdb_cassandra_cluster" "example" {
  name                           = "example-cluster"
  resource_group_name            = azurerm_resource_group.example.name
  location                       = azurerm_resource_group.example.location
  delegated_management_subnet_id = azurerm_subnet.example.id
  default_admin_password         = "Password1234"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cassandra Cluster. Changing this forces a new Cassandra Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Cassandra Cluster should exist. Changing this forces a new Cassandra Cluster to be created.

* `location` - (Required) The Azure Region where the Cassandra Cluster should exist. Changing this forces a new Cassandra Cluster to be created.

* `delegated_management_subnet_id` - (Required) The ID of the delegated management subnet for this Cassandra Cluster. Changing this forces a new Cassandra Cluster to be created.

* `default_admin_password` - (Required) The initial admin password for this Cassandra Cluster.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cassandra Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cassandra Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cassandra Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cassandra Cluster.

## Import

Cassandra Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_cassandra_cluster.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/cassandraClusters/cluster1
```
