---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_connection"
description: |-
  Manages a service connector for kubernetes cluster.
---

# azurerm_kubernetes_cluster_connection

Manages a service connector for kubernetes cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-cosmosdb-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "example" {
  name                = "cosmos-sql-db"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "example" {
  name                = "example-container"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  database_name       = azurerm_cosmosdb_sql_database.example.name
  partition_key_path  = "/definition"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks-cluster"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "exampleaks"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_connection" "example" {
  name                  = "example-serviceconnector"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  target_resource_id    = azurerm_cosmosdb_sql_database.example.id
  authentication {
    type   = "secret"
    name   = "foo"
    secret = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service connection. Changing this forces a new resource to be created.

* `kubernetes_cluster_id` - (Required) The ID of the kubernetes cluster. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the target resource. Changing this forces a new resource to be created. Possible target resources are `Postgres`, `PostgresFlexible`, `Mysql`, `Sql`, `Redis`, `RedisEnterprise`, `CosmosCassandra`, `CosmosGremlin`, `CosmosMongo`, `CosmosSql`, `CosmosTable`, `StorageBlob`, `StorageQueue`, `StorageFile`, `StorageTable`, `AppConfig`, `EventHub`, `ServiceBus`, `SignalR`, `WebPubSub`, `ConfluentKafka`. The integration guide can be found [here](https://learn.microsoft.com/en-us/azure/service-connector/how-to-integrate-postgres).

* `authentication` - (Required) The authentication info. An `authentication` block as defined below.

* `client_type` - (Optional) The application client type. Possible values are `none`, `dotnet`, `java`, `python`, `go`, `php`, `ruby`, `django`, `nodejs` and `springBoot`. Defaults to `none`.

* `vnet_solution` - (Optional) The type of VNet solution. Possible values are `serviceEndpoint`, `privateLink`.

* `secret_store` - (Optional) An option to store secret value in secure place. An `secret_store` block as defined below.

-> **Note:** If a Managed Identity is used, this will need to be configured on the Kubernetes Cluster.

---

An `authentication` block supports the following:

* `type` - (Required) The authentication type. Possible values are `userAssignedIdentity`, `servicePrincipalSecret`, `servicePrincipalCertificate`, `secret`. Changing this forces a new resource to be created.

* `certificate` - (Optional) The certificate for `servicePrincipalCertificate` authentication. Should be specified when `type` is set to `servicePrincipalCertificate`.

* `client_id` - (Optional) The client ID for `userAssignedIdentity` or `servicePrincipalSecret` authentication. Should be specified when `type` is set to `servicePrincipalSecret` or `userAssignedIdentity`.

* `name` - (Optional) The name of the secret for `secret` authentication. Should be specified when `type` is set to `secret`.

* `principal_id` - (Optional) The principal ID for `servicePrincipalSecret` authentication. Should be specified when `type` is set to `servicePrincipalSecret`.

* `secret` - (Optional) The secret for `secret` or `servicePrincipalSecret` authentication. Should be specified when `type` is set to `secret` or `servicePrincipalSecret`.

* `subscription_id` - (Optional) The subscription ID for `userAssignedIdentity` authentication. Should be specified when `type` is set to `userAssignedIdentity`.

---

An `secret_store` block supports the following:

* `key_vault_id` - (Required) The key vault id to store secret.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the service connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Connector.
* `update` - (Defaults to 30 minutes) Used when updating the Service Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Connector.

## Import

Service Connector for kubernetes cluster can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/providers/Microsoft.ServiceLinker/linkers/link1
```
