---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_cassandra_datacenter"
description: |-
  Manages a Cassandra Datacenter.
---

# azurerm_cosmosdb_cassandra_datacenter

Manages a Cassandra Datacenter.

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

resource "azurerm_cosmosdb_cassandra_datacenter" "example" {
  name                           = "example-datacenter"
  location                       = azurerm_cosmosdb_cassandra_cluster.example.location
  cassandra_cluster_id           = azurerm_cosmosdb_cassandra_cluster.example.id
  delegated_management_subnet_id = azurerm_subnet.example.id
  node_count                     = 3
  disk_count                     = 4
  sku_name                       = "Standard_DS14_v2"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cassandra Datacenter. Changing this forces a new Cassandra Datacenter to be created.

* `location` - (Required) The Azure Region where the Cassandra Datacenter should exist. Changing this forces a new Cassandra Datacenter to be created.

* `cassandra_cluster_id` - (Required) The ID of the Cassandra Cluster. Changing this forces a new Cassandra Datacenter to be created.

* `delegated_management_subnet_id` - (Required) The ID of the delegated management subnet for this Cassandra Datacenter. Changing this forces a new Cassandra Datacenter to be created.

* `node_count` - (Required) The number of nodes the Cassandra Datacenter should have. The number should be equal or greater than 3. Defaults to 3.

---

* `sku_name` - (Optional) Determines the selected sku. Defaults to Standard_DS14_v2. 

* `disk_count` - (Optional) Determines the number of p30 disks that are attached to each node. Defaults to 4.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cassandra Datacenter.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cassandra Datacenter.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cassandra Datacenter.
* `update` - (Defaults to 30 minutes) Used when updating the Cassandra Datacenter.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cassandra Datacenter.

## Import

Cassandra Datacenters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_cassandra_datacenter.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/cassandraClusters/cluster1/dataCenters/dc1
```
