---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_endpoint_cosmosdb_account"
description: |-
  Manages an IotHub Cosmos DB Account Endpoint
---

# azurerm_iothub_endpoint_cosmosdb_account

Manages an IotHub Cosmos DB Account Endpoint

~> **Note:** Endpoints can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_endpoint_*` resources - but the two ways of defining the endpoints cannot be used together. If both are used against the same IoTHub, spurious changes will occur. Also, defining a `azurerm_iothub_endpoint_*` resource and another endpoint of a different type directly on the `azurerm_iothub` resource is not supported.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_iothub" "example" {
  name                = "exampleIothub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "example"
  }
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "cosmosdb-account"
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

resource "azurerm_cosmosdb_sql_database" "example" {
  name                = "cosmos-sql-db"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
}

resource "azurerm_cosmosdb_sql_container" "example" {
  name                = "example-container"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  database_name       = azurerm_cosmosdb_sql_database.example.name
  partition_key_path  = "/definition/id"
}

resource "azurerm_iothub_endpoint_cosmosdb_account" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  iothub_id           = azurerm_iothub.example.id
  container_name      = azurerm_cosmosdb_sql_container.example.name
  database_name       = azurerm_cosmosdb_sql_database.example.name
  endpoint_uri        = azurerm_cosmosdb_account.example.endpoint
  primary_key         = azurerm_cosmosdb_account.example.primary_key
  secondary_key       = azurerm_cosmosdb_account.example.secondary_key
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the endpoint. The name must be unique across endpoint types. The following names are reserved: `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the Cosmos DB Account has been created. Changing this forces a new resource to be created.

* `iothub_id` - (Required) The ID of the IoT Hub to create the endpoint. Changing this forces a new resource to be created.

* `container_name` - (Required) The name of the Cosmos DB Container in the Cosmos DB Database. Changing this forces a new resource to be created.

* `database_name` - (Required) The name of the Cosmos DB Database in the Cosmos DB Account. Changing this forces a new resource to be created.

* `endpoint_uri` - (Required) The URI of the Cosmos DB Account. Changing this forces a new resource to be created.

* `authentication_type` - (Optional) The type used to authenticate against the Cosmos DB Account endpoint. Possible values are `keyBased` and `identityBased`. Defaults to `keyBased`.

* `identity_id` - (Optional) The ID of the User Managed Identity used to authenticate against the Cosmos DB Account endpoint.

~> **Note:** `identity_id` can only be specified when `authentication_type` is `identityBased`. It must be one of the `identity_ids` of the Iot Hub. If not specified when `authentication_type` is `identityBased`, System Assigned Managed Identity of the Iot Hub will be used.

* `partition_key_name` - (Optional) The name of the partition key associated with the Cosmos DB Container.

* `partition_key_template` - (Optional) The template for generating a synthetic partition key value for use within the Cosmos DB Container.

* `primary_key` - (Optional) The primary key of the Cosmos DB Account.

~> **Note:** `primary_key` must and can only be specified when `authentication_type` is `keyBased`.

* `secondary_key` - (Optional) The secondary key of the Cosmos DB Account.

~> **Note:** `secondary_key` must and can only be specified when `authentication_type` is `keyBased`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoTHub Cosmos DB Account Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub Cosmos DB Account Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub Cosmos DB Account Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub Cosmos DB Account Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub Cosmos DB Account Endpoint.

## Import

IoTHub Cosmos DB Account Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_endpoint_cosmosdb_account.endpoint1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/iotHubs/hub1/endpoints/cosmosDBAccountEndpoint1
```
