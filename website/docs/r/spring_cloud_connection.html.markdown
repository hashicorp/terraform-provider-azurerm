---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_connection"
description: |-
  Manages a service connector for spring cloud app.
---

# azurerm_spring_cloud_connection

Manages a service connector for spring cloud app.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_connection` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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

resource "azurerm_spring_cloud_service" "example" {
  name                = "examplespringcloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_spring_cloud_app" "example" {
  name                = "examplespringcloudapp"
  resource_group_name = azurerm_resource_group.example.name
  service_name        = azurerm_spring_cloud_service.example.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_spring_cloud_java_deployment" "example" {
  name                = "exampledeployment"
  spring_cloud_app_id = azurerm_spring_cloud_app.example.id
}

resource "azurerm_spring_cloud_connection" "example" {
  name               = "example-serviceconnector"
  spring_cloud_id    = azurerm_spring_cloud_java_deployment.example.id
  target_resource_id = azurerm_cosmosdb_sql_database.example.id
  authentication {
    type = "systemAssignedIdentity"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service connection. Changing this forces a new resource to be created.

* `spring_cloud_id` - (Required) The ID of the data source spring cloud. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the target resource. Changing this forces a new resource to be created. Possible target resources are `Postgres`, `PostgresFlexible`, `Mysql`, `Sql`, `Redis`, `RedisEnterprise`, `CosmosCassandra`, `CosmosGremlin`, `CosmosMongo`, `CosmosSql`, `CosmosTable`, `StorageBlob`, `StorageQueue`, `StorageFile`, `StorageTable`, `AppConfig`, `EventHub`, `ServiceBus`, `SignalR`, `WebPubSub`, `ConfluentKafka`. The integration guide can be found [here](https://learn.microsoft.com/en-us/azure/service-connector/how-to-integrate-postgres).

* `authentication` - (Required) The authentication info. An `authentication` block as defined below.

---

An `authentication` block supports the following:

* `type` - (Required) The authentication type. Possible values are `systemAssignedIdentity`, `userAssignedIdentity`, `servicePrincipalSecret`, `servicePrincipalCertificate`, `secret`. Changing this forces a new resource to be created.

* `name` - (Optional) Username or account name for secret auth. `name` and `secret` should be either both specified or both not specified when `type` is set to `secret`.

* `secret` - (Optional) Password or account key for secret auth. `secret` and `name` should be either both specified or both not specified when `type` is set to `secret`.

* `client_id` - (Optional) Client ID for `userAssignedIdentity` or `servicePrincipal` auth. Should be specified when `type` is set to `servicePrincipalSecret` or `servicePrincipalCertificate`. When `type` is set to `userAssignedIdentity`, `client_id` and `subscription_id` should be either both specified or both not specified.

* `subscription_id` - (Optional) Subscription ID for `userAssignedIdentity`. `subscription_id` and `client_id` should be either both specified or both not specified.

* `principal_id` - (Optional) Principal ID for `servicePrincipal` auth. Should be specified when `type` is set to `servicePrincipalSecret` or `servicePrincipalCertificate`.

* `certificate` - (Optional) Service principal certificate for `servicePrincipal` auth. Should be specified when `type` is set to `servicePrincipalCertificate`.

---

* `client_type` - (Optional) The application client type. Possible values are `none`, `dotnet`, `java`, `python`, `go`, `php`, `ruby`, `django`, `nodejs` and `springBoot`. Defaults to `none`.

* `vnet_solution` - (Optional) The type of the VNet solution. Possible values are `serviceEndpoint`, `privateLink`.

* `secret_store` - (Optional) An option to store secret value in secure place. An `secret_store` block as defined below.

---

An `secret_store` block supports the following:

* `key_vault_id` - (Required) The key vault id to store secret.


## Attribute Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the service connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Connector for spring cloud.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Connector for spring cloud.
* `update` - (Defaults to 30 minutes) Used when updating the Service Connector for spring cloud.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Connector for spring cloud.

## Import

Service Connector for spring cloud can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AppPlatform/Spring/springcloud/apps/springcloudapp/deployments/deployment/providers/Microsoft.ServiceLinker/linkers/serviceconnector1
```
