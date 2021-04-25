---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app_cosmosdb_association"
description: |-
  Associates a [Spring Cloud Application](spring_cloud_app.html) with a [CosmosDB Account](cosmosdb_account.html).
---

# azurerm_spring_cloud_app_cosmosdb_association

Associates a [Spring Cloud Application](spring_cloud_app.html) with a [CosmosDB Account](cosmosdb_account.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example-springcloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_spring_cloud_app" "example" {
  name                = "example-springcloudapp"
  resource_group_name = azurerm_resource_group.example.name
  service_name        = azurerm_spring_cloud_service.example.name
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-cosmosdb-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_spring_cloud_app_cosmosdb_association" "example" {
  name                = "example-bind"
  spring_cloud_app_id = azurerm_spring_cloud_app.example.id
  cosmosdb_account_id = azurerm_cosmosdb_account.example.id
  api_type            = "table"
  cosmosdb_access_key = azurerm_cosmosdb_account.example.primary_key
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Application Association. Changing this forces a new resource to be created.

* `spring_cloud_app_id` - (Required) Specifies the ID of the Spring Cloud Application where this Association is created. Changing this forces a new resource to be created.

* `cosmosdb_account_id` - (Required) Specifies the ID of the CosmosDB Account. Changing this forces a new resource to be created.

* `api_type` - (Required) Specifies the api type which should be used when connecting to the CosmosDB Account. Possible values are `cassandra`, `gremlin`, `mongo`, `sql` or `table`. Changing this forces a new resource to be created.
  
* `cosmosdb_access_key` - (Required) Specifies the CosmosDB Account access key.

* `cosmosdb_cassandra_keyspace_name` - (Optional) Specifies the name of the Cassandra Keyspace which the Spring Cloud App should be associated with. Should only be set when `api_type` is `cassandra`.

* `cosmosdb_gremlin_database_name` - (Optional) Specifies the name of the Gremlin Database which the Spring Cloud App should be associated with. Should only be set when `api_type` is `gremlin`.

* `cosmosdb_gremlin_graph_name` - (Optional) Specifies the name of the Gremlin Graph which the Spring Cloud App should be associated with. Should only be set when `api_type` is `gremlin`.

* `cosmosdb_mongo_database_name` - (Optional) Specifies the name of the Mongo Database which the Spring Cloud App should be associated with. Should only be set when `api_type` is `mongo`.

* `cosmosdb_sql_database_name` - (Optional) Specifies the name of the Sql Database which the Spring Cloud App should be associated with. Should only be set when `api_type` is `sql`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spring Cloud Application CosmosDB Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application CosmosDB Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application CosmosDB Association.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application CosmosDB Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application CosmosDB Association.

## Import

Spring Cloud Application CosmosDB Association can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_app_cosmosdb_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.AppPlatform/Spring/service1/apps/app1/bindings/bind1
```
