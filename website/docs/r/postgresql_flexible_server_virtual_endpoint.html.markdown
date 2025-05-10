---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_flexible_server_virtual_endpoint"
description: |-
  Manages a Virtual Endpoint on a PostgreSQL Flexible Server
---

# azurerm_postgresql_flexible_server_virtual_endpoint

Allows you to create a Virtual Endpoint associated with a Postgres Flexible Replica.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_postgresql_flexible_server" "example" {
  name                          = "example"
  resource_group_name           = azurerm_resource_group.example.name
  location                      = azurerm_resource_group.example.location
  version                       = "16"
  public_network_access_enabled = false
  administrator_login           = "psqladmin"
  administrator_password        = "H@Sh1CoR3!"
  zone                          = "1"

  storage_mb   = 32768
  storage_tier = "P30"

  sku_name = "GP_Standard_D2ads_v5"
}

resource "azurerm_postgresql_flexible_server" "example_replica" {
  name                          = "example-replica"
  resource_group_name           = azurerm_postgresql_flexible_server.example.resource_group_name
  location                      = azurerm_postgresql_flexible_server.example.location
  create_mode                   = "Replica"
  source_server_id              = azurerm_postgresql_flexible_server.example.id
  version                       = "16"
  public_network_access_enabled = false
  zone                          = "1"
  storage_mb                    = 32768
  storage_tier                  = "P30"

  sku_name = "GP_Standard_D2ads_v5"
}

resource "azurerm_postgresql_flexible_server_virtual_endpoint" "example" {
  name              = "example-endpoint-1"
  source_server_id  = azurerm_postgresql_flexible_server.example.id
  replica_server_id = azurerm_postgresql_flexible_server.example_replica.id
  type              = "ReadWrite"
}
```

-> **Note:** If creating multiple replicas, an error can occur if virtual endpoints are created before all replicas have been completed. To avoid this error, use a `depends_on` property on `azurerm_postgresql_flexible_server_virtual_endpoint` that references all Postgres Flexible Server Replicas.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Endpoint

* `source_server_id` - (Required) The Resource ID of the *Source* Postgres Flexible Server this should be associated with.

* `replica_server_id` - (Required) The Resource ID of the *Replica* Postgres Flexible Server this should be associated with

~> **Note:** If a fail-over has occurred, you will be unable to update `replica_server_id`. You can remove the resource from state and reimport it back in with `source_server_id` and `replica_server_id` flipped and then update `replica_server_id`.

* `type` - (Required) The type of Virtual Endpoint. Currently only `ReadWrite` is supported.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL Flexible Virtual Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the PostgreSQL Flexible Virtual Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Flexible Virtual Endpoint.
* `update` - (Defaults to 10 minutes) Used when updating the PostgreSQL Flexible Virtual Endpoint.
* `delete` - (Defaults to 5 minutes) Used when deleting the PostgreSQL Flexible Virtual Endpoint.

## Import

A PostgreSQL Flexible Virtual Endpoint can be imported using the `resource id`, e.g.
```shell
terraform import azurerm_postgresql_flexible_server_virtual_endpoint.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforPostgreSQL/flexibleServers/sourceServerName/virtualEndpoints/endpointName|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforPostgreSQL/flexibleServers/replicaServerName/virtualEndpoints/endpointName"
```
