---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: cassandra_managed_instance_datacenter"
description: |-
  Creates an Azure Managed Instance for Apache Cassandra Datacenter
---

# azurerm_cosmosdb_cassandra_kmanaged_instance_datacenter

Creates an [Azure Managed Instance for Apache Cassandra](https://docs.microsoft.com/azure/managed-instance-apache-cassandra/) Cluster and Datacenter.

## Example Usage

```hcl
provider "azurerm" {
    features {}
}

resource "azurerm_resource_group" "test" {
    name     = "myResourceGroup"
    location = "East US"
}
 
resource "azurerm_virtual_network" "test" {
    name                = "myVnet"
    location            = azurerm_resource_group.test.location
    resource_group_name = azurerm_resource_group.test.name
    address_space       = ["10.0.0.0/16"]

    tags = {
      environment = "Test"
    }
}

resource "azurerm_subnet" "test" {
    name                 = "mySubNet"
    resource_group_name  = azurerm_resource_group.test.name
    virtual_network_name = azurerm_virtual_network.test.name
    address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_role_assignment" "test" {
    scope                = azurerm_virtual_network.test.id
    role_definition_name = "Network Contributor"
    principal_id         = "e5007d2c-4b13-4a74-9b6a-605d99f03501"
}

resource "azurerm_cosmosdb_cassandra_managed_instance_cluster" "test" {
    name                             = "myCluster"
    resource_group_name              = azurerm_resource_group.test.name
    location                         = azurerm_resource_group.test.location
    delegated_management_subnet_id   = azurerm_subnet.test.id
    initial_cassandra_admin_password = "Password1234"  
}

resource "azurerm_cosmosdb_cassandra_managed_instance_datacenter" "test" {
    name                   = azurerm_cosmosdb_cassandra_managed_instance_cluster.test.name
    datacenter_name                = "myDatacenter"
    resource_group_name            = azurerm_resource_group.test.name
    location                       = azurerm_resource_group.test.location
    delegated_management_subnet_id = azurerm_subnet.test.id
    node_count                     = 3
    disk_capacity                  = 2
    sku                            = "Standard_D8s_v4"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required) Specifies the name of the Cassandra Cluster. Changing this forces a new resource to be created.

* `datacenter_name` - (Required) The name of the Managed Cassandra Datacenter to be greated in the cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cassandra Cluster is created. Changing this forces a new resource to be created.

* `location` - (Required) The region in which the datacenter will be created. 

* `delegated_management_subnet_id` - (Required) The resource id of the subnet in which the Datacenter nodes will be injected.

* `node_count` - (Required) The number of nodes that will be provisioned in the Datacenter. 

* `sku` - (Optional) Determines the selected sku. Defaults to Standard_DS14_v2. 

* `disk_capacity` - (Optional) Determines the number of p30 disks that are attached to each node. Defaults to 4. 


## Attributes Reference

The following attributes are exported:

* `id` - the ID of the Cassandra Datacenter.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CosmosDB Cassandra KeySpace.
* `update` - (Defaults to 30 minutes) Used when updating the CosmosDB Cassandra KeySpace.
* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB Cassandra KeySpace.
* `delete` - (Defaults to 30 minutes) Used when deleting the CosmosDB Cassandra KeySpace.

## Import

Managed Cassandra Datacenter can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cassandra_managed_instance_datacenter.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/cassandraClusters/clusterName/dataCenters/DatacenterName
```
