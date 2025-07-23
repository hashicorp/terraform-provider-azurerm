---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_cassandra_datacenter"
description: |-
  Manages a Cassandra Datacenter.
---

# azurerm_cosmosdb_cassandra_datacenter

Manages a Cassandra Datacenter.

~> **Note:** In order for the `Azure Managed Instances for Apache Cassandra` to work properly the product requires the `Azure Cosmos DB` Application ID to be present and working in your tenant. If the `Azure Cosmos DB` Application ID is missing in your environment you will need to have an administrator of your tenant run the following command to add the `Azure Cosmos DB` Application ID to your tenant:

```powershell
New-AzADServicePrincipal -ApplicationId a232010e-820c-4083-83bb-3ace5fc29d0b
```

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

data "azuread_service_principal" "example" {
  display_name = "Azure Cosmos DB"
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_virtual_network.example.id
  role_definition_name = "Network Contributor"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_cosmosdb_cassandra_cluster" "example" {
  name                           = "example-cluster"
  resource_group_name            = azurerm_resource_group.example.name
  location                       = azurerm_resource_group.example.location
  delegated_management_subnet_id = azurerm_subnet.example.id
  default_admin_password         = "Password1234"

  depends_on = [azurerm_role_assignment.example]
}

resource "azurerm_cosmosdb_cassandra_datacenter" "example" {
  name                           = "example-datacenter"
  location                       = azurerm_cosmosdb_cassandra_cluster.example.location
  cassandra_cluster_id           = azurerm_cosmosdb_cassandra_cluster.example.id
  delegated_management_subnet_id = azurerm_subnet.example.id
  node_count                     = 3
  disk_count                     = 4
  sku_name                       = "Standard_DS14_v2"
  availability_zones_enabled     = false
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cassandra Datacenter. Changing this forces a new Cassandra Datacenter to be created.

* `location` - (Required) The Azure Region where the Cassandra Datacenter should exist. Changing this forces a new Cassandra Datacenter to be created.

* `cassandra_cluster_id` - (Required) The ID of the Cassandra Cluster. Changing this forces a new Cassandra Datacenter to be created.

* `delegated_management_subnet_id` - (Required) The ID of the delegated management subnet for this Cassandra Datacenter. Changing this forces a new Cassandra Datacenter to be created.

* `node_count` - (Optional) The number of nodes the Cassandra Datacenter should have. The number should be equal or greater than `3`. Defaults to `3`.

---

* `backup_storage_customer_key_uri` - (Optional) The key URI of the customer key to use for the encryption of the backup Storage Account.

* `base64_encoded_yaml_fragment` - (Optional) The fragment of the cassandra.yaml configuration file to be included in the cassandra.yaml for all nodes in this Cassandra Datacenter. The fragment should be Base64 encoded and only a subset of keys is allowed.

* `disk_sku` - (Optional) The Disk SKU that is used for this Cassandra Datacenter. Defaults to `P30`.

* `managed_disk_customer_key_uri` - (Optional) The key URI of the customer key to use for the encryption of the Managed Disk.

* `sku_name` - (Optional) Determines the selected sku.

-> **Note:** In v4.0 of the provider the `sku_name` will have a default value of `Standard_E16s_v5`.

* `disk_count` - (Optional) Determines the number of p30 disks that are attached to each node.

* `availability_zones_enabled` - (Optional) Determines whether availability zones are enabled. Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cassandra Datacenter.

* `seed_node_ip_addresses` - A list of IP Address for the seed nodes in this Cassandra Datacenter. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Cassandra Datacenter.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cassandra Datacenter.
* `update` - (Defaults to 1 hour) Used when updating the Cassandra Datacenter.
* `delete` - (Defaults to 1 hour) Used when deleting the Cassandra Datacenter.

## Import

Cassandra Datacenters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_cassandra_datacenter.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/cassandraClusters/cluster1/dataCenters/dc1
```
