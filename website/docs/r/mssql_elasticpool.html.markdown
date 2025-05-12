---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_elasticpool"
description: |-
  Manages an Azure SQL Elastic Pool.
---

# azurerm_mssql_elasticpool

Allows you to manage an Azure SQL Elastic Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "my-resource-group"
  location = "West Europe"
}

resource "azurerm_mssql_server" "example" {
  name                         = "my-sql-server"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "example" {
  name                = "test-epool"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  server_name         = azurerm_mssql_server.example.name
  license_type        = "LicenseIncluded"
  max_size_gb         = 756

  sku {
    name     = "BasicPool"
    tier     = "Basic"
    family   = "Gen4"
    capacity = 4
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the elastic pool. This needs to be globally unique. Changing this forces a new resource to be created. 

* `resource_group_name` - (Required) The name of the resource group in which to create the elastic pool. This must be the same as the resource group of the underlying SQL server. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `server_name` - (Required) The name of the SQL Server on which to create the elastic pool. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `per_database_settings` - (Required) A `per_database_settings` block as defined below.

* `maintenance_configuration_name` - (Optional) The name of the Public Maintenance Configuration window to apply to the elastic pool. Valid values include `SQL_Default`, `SQL_EastUS_DB_1`, `SQL_EastUS2_DB_1`, `SQL_SoutheastAsia_DB_1`, `SQL_AustraliaEast_DB_1`, `SQL_NorthEurope_DB_1`, `SQL_SouthCentralUS_DB_1`, `SQL_WestUS2_DB_1`, `SQL_UKSouth_DB_1`, `SQL_WestEurope_DB_1`, `SQL_EastUS_DB_2`, `SQL_EastUS2_DB_2`, `SQL_WestUS2_DB_2`, `SQL_SoutheastAsia_DB_2`, `SQL_AustraliaEast_DB_2`, `SQL_NorthEurope_DB_2`, `SQL_SouthCentralUS_DB_2`, `SQL_UKSouth_DB_2`, `SQL_WestEurope_DB_2`, `SQL_AustraliaSoutheast_DB_1`, `SQL_BrazilSouth_DB_1`, `SQL_CanadaCentral_DB_1`, `SQL_CanadaEast_DB_1`, `SQL_CentralUS_DB_1`, `SQL_EastAsia_DB_1`, `SQL_FranceCentral_DB_1`, `SQL_GermanyWestCentral_DB_1`, `SQL_CentralIndia_DB_1`, `SQL_SouthIndia_DB_1`, `SQL_JapanEast_DB_1`, `SQL_JapanWest_DB_1`, `SQL_NorthCentralUS_DB_1`, `SQL_UKWest_DB_1`, `SQL_WestUS_DB_1`, `SQL_AustraliaSoutheast_DB_2`, `SQL_BrazilSouth_DB_2`, `SQL_CanadaCentral_DB_2`, `SQL_CanadaEast_DB_2`, `SQL_CentralUS_DB_2`, `SQL_EastAsia_DB_2`, `SQL_FranceCentral_DB_2`, `SQL_GermanyWestCentral_DB_2`, `SQL_CentralIndia_DB_2`, `SQL_SouthIndia_DB_2`, `SQL_JapanEast_DB_2`, `SQL_JapanWest_DB_2`, `SQL_NorthCentralUS_DB_2`, `SQL_UKWest_DB_2`, `SQL_WestUS_DB_2`, `SQL_WestCentralUS_DB_1`, `SQL_FranceSouth_DB_1`, `SQL_WestCentralUS_DB_2`, `SQL_FranceSouth_DB_2`, `SQL_SwitzerlandNorth_DB_1`, `SQL_SwitzerlandNorth_DB_2`, `SQL_BrazilSoutheast_DB_1`, `SQL_UAENorth_DB_1`, `SQL_BrazilSoutheast_DB_2`, `SQL_UAENorth_DB_2`, `SQL_SouthAfricaNorth_DB_1`, `SQL_SouthAfricaNorth_DB_2`, `SQL_WestUS3_DB_1`, `SQL_WestUS3_DB_2`, `SQL_SwedenCentral_DB_1`, `SQL_SwedenCentral_DB_2`. Defaults to `SQL_Default`.

* `max_size_gb` - (Optional) The max data size of the elastic pool in gigabytes. Conflicts with `max_size_bytes`.

* `max_size_bytes` - (Optional) The max data size of the elastic pool in bytes. Conflicts with `max_size_gb`.

-> **Note:** One of either `max_size_gb` or `max_size_bytes` must be specified.

* `enclave_type` - (Optional) Specifies the type of enclave to be used by the elastic pool. When `enclave_type` is not specified (e.g., the default) enclaves are not enabled on the elastic pool. Once enabled (e.g., by specifying `Default` or `VBS`) removing the `enclave_type` field from the configuration file will force the creation of a new resource. Possible values are `Default` or `VBS`.

-> **Note:** All databases that are added to the elastic pool must have the same `enclave_type` as the elastic pool.

-> **Note:** `enclave_type` is not supported for DC-series SKUs.

~> **Note:** The default value for `enclave_type` field is unset not `Default`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `zone_redundant` - (Optional) Whether or not this elastic pool is zone redundant. `tier` needs to be `Premium` for `DTU` based or `BusinessCritical` for `vCore` based `sku`.

* `license_type` - (Optional) Specifies the license type applied to this database. Possible values are `LicenseIncluded` and `BasePrice`.

---

The `sku` block supports the following:

* `name` - (Required) Specifies the SKU Name for this Elasticpool. The name of the SKU, will be either `vCore` based or `DTU` based. Possible `DTU` based values are `BasicPool`, `StandardPool`, `PremiumPool` while possible `vCore` based values are `GP_Gen4`, `GP_Gen5`, `GP_Fsv2`, `GP_DC`, `BC_Gen4`, `BC_Gen5`, `BC_DC`, `HS_PRMS`, `HS_MOPRMS`, or `HS_Gen5`.

* `capacity` - (Required) The scale up/out capacity, representing server's compute units. For more information see the documentation for your Elasticpool configuration: [vCore-based](https://docs.microsoft.com/azure/sql-database/sql-database-vcore-resource-limits-elastic-pools) or [DTU-based](https://docs.microsoft.com/azure/sql-database/sql-database-dtu-resource-limits-elastic-pools).

* `tier` - (Required) The tier of the particular SKU. Possible values are `GeneralPurpose`, `BusinessCritical`, `Basic`, `Standard`, `Premium`, or `HyperScale`. For more information see the documentation for your Elasticpool configuration: [vCore-based](https://docs.microsoft.com/azure/sql-database/sql-database-vcore-resource-limits-elastic-pools) or [DTU-based](https://docs.microsoft.com/azure/sql-database/sql-database-dtu-resource-limits-elastic-pools).

* `family` - (Optional) The `family` of hardware `Gen4`, `Gen5`, `Fsv2`, `MOPRMS`, `PRMS`, or `DC`.

---

The `per_database_settings` block supports the following:

* `min_capacity` - (Required) The minimum capacity all databases are guaranteed.

* `max_capacity` - (Required) The maximum capacity any one database can consume.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MS SQL Elastic Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MS SQL Elastic Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the MS SQL Elastic Pool.
* `update` - (Defaults to 30 minutes) Used when updating the MS SQL Elastic Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the MS SQL Elastic Pool.

## Import

SQL Elastic Pool can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_elasticpool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/elasticPools/myelasticpoolname
```
