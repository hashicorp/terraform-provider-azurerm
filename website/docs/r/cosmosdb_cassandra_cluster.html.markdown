---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_cassandra_cluster"
description: |-
  Manages a Cassandra Cluster.
---

# azurerm_cosmosdb_cassandra_cluster

Manages a Cassandra Cluster.

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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cassandra Cluster. Changing this forces a new Cassandra Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Cassandra Cluster should exist. Changing this forces a new Cassandra Cluster to be created.

* `location` - (Required) The Azure Region where the Cassandra Cluster should exist. Changing this forces a new Cassandra Cluster to be created.

* `delegated_management_subnet_id` - (Required) The ID of the delegated management subnet for this Cassandra Cluster. Changing this forces a new Cassandra Cluster to be created.

* `default_admin_password` - (Required) The initial admin password for this Cassandra Cluster. Changing this forces a new resource to be created.

* `authentication_method` - (Optional) The authentication method that is used to authenticate clients. Possible values are `None` and `Cassandra`. Defaults to `Cassandra`.

* `client_certificate_pems` - (Optional) A list of TLS certificates that is used to authorize client connecting to the Cassandra Cluster.

* `external_gossip_certificate_pems` - (Optional) A list of TLS certificates that is used to authorize gossip from unmanaged Cassandra Data Center.

* `external_seed_node_ip_addresses` - (Optional) A list of IP Addresses of the seed nodes in unmanaged the Cassandra Data Center which will be added to the seed node lists of all managed nodes.

* `hours_between_backups` - (Optional) The number of hours to wait between taking a backup of the Cassandra Cluster. Defaults to `24`.

~> **Note:** To disable this feature, set this property to `0`.

* `identity` - (Optional) An `identity` block as defined below.

* `repair_enabled` - (Optional) Is the automatic repair enabled on the Cassandra Cluster? Defaults to `true`.

* `version` - (Optional) The version of Cassandra what the Cluster converges to run. Possible values are `3.11` and `4.0`. Defaults to `3.11`. Changing this forces a new Cassandra Cluster to be created.

* `tags` - (Optional) A mapping of tags assigned to the resource.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Cassandra Cluster. The only possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cassandra Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cassandra Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cassandra Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Cassandra Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cassandra Cluster.

## Import

Cassandra Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_cassandra_cluster.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/cassandraClusters/cluster1
```
