---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_managed_instance_failover_group"
description: |-
  Manages a SQL Instance Failover Group.
---

# azurerm_sql_managed_instance_failover_group

Manages a SQL Instance Failover Group.

## Example Usage

-> **Note:** The `azurerm_sql_managed_instance_failover_group` resource is deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use the [`azurerm_mssql_managed_instance_failover_group`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_managed_instance_failover_group) resource instead.

~> **Note:** For a more complete example, see the [the `examples/sql-azure/managed_instance_failover_group` directory](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/sql-azure/managed_instance_failover_group) within the GitHub Repository.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

resource "azurerm_sql_managed_instance" "primary" {
  name                         = "example-primary"
  resource_group_name          = azurerm_resource_group.primary.name
  location                     = azurerm_resource_group.primary.location
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
  license_type                 = "BasePrice"
  subnet_id                    = azurerm_subnet.primary.id
  sku_name                     = "GP_Gen5"
  vcores                       = 4
  storage_size_in_gb           = 32

  depends_on = [
    azurerm_subnet_network_security_group_association.primary,
    azurerm_subnet_route_table_association.primary,
  ]

  tags = {
    environment = "prod"
  }
}

resource "azurerm_sql_managed_instance" "secondary" {
  name                         = "example-secondary"
  resource_group_name          = azurerm_resource_group.secondary.name
  location                     = azurerm_resource_group.secondary.location
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
  license_type                 = "BasePrice"
  subnet_id                    = azurerm_subnet.secondary.id
  sku_name                     = "GP_Gen5"
  vcores                       = 4
  storage_size_in_gb           = 32

  depends_on = [
    azurerm_subnet_network_security_group_association.secondary,
    azurerm_subnet_route_table_association.secondary,
  ]

  tags = {
    environment = "prod"
  }
}

resource "azurerm_sql_managed_instance_failover_group" "example" {
  name                        = "example-failover-group"
  resource_group_name         = azurerm_resource_group.primary.name
  location                    = azurerm_sql_managed_instance.primary.location
  managed_instance_name       = azurerm_sql_managed_instance.primary.name
  partner_managed_instance_id = azurerm_sql_managed_instance.secondary.id

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SQL Instance Failover Group. Changing this forces a new SQL Instance Failover Group to be created.

* `managed_instance_name` - (Required) The name of the SQL Managed Instance which will be replicated using a SQL Instance Failover Group. Changing this forces a new SQL Instance Failover Group to be created.

* `location` - (Required) The Azure Region where the SQL Instance Failover Group exists. Changing this forces a new resource to be created.

* `partner_managed_instance_id` - (Required) ID of the SQL Managed Instance which will be replicated to. Changing this forces a new resource to be created.

* `read_write_endpoint_failover_policy` - (Required) A `read_write_endpoint_failover_policy` block as defined below.

* `resource_group_name` - (Required) The name of the Resource Group where the SQL Instance Failover Group should exist. Changing this forces a new SQL Instance Failover Group to be created.

* `readonly_endpoint_failover_policy_enabled` - (Optional) Failover policy for the read-only endpoint. Defaults to `true`.

---

A `read_write_endpoint_failover_policy` block supports the following:

* `mode` - (Required) The failover mode. Possible values are `Manual`, `Automatic`

* `grace_minutes` - (Optional) Applies only if `mode` is `Automatic`. The grace period in minutes before failover with data loss is attempted.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SQL Instance Failover Group.

* `partner_region` - A `partner_region` block as defined below.

* `role` - The local replication role of the SQL Instance Failover Group.

---

A `partner_region` block exports the following:

* `location` - The Azure Region where the SQL Instance Failover Group partner exists.

* `role` - The partner replication role of the SQL Instance Failover Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SQL Instance Failover Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Instance Failover Group.
* `update` - (Defaults to 30 minutes) Used when updating the SQL Instance Failover Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the SQL Instance Failover Group.

## Import

SQL Instance Failover Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_managed_instance_failover_group.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/locations/Location/instanceFailoverGroups/failoverGroup1
```
