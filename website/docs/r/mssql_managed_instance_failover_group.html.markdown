---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_instance_failover_group"
description: |-
  Manages an Azure SQL Managed Instance Failover Group.
---

# azurerm_mssql_managed_instance_failover_group

Manages an Azure SQL Managed Instance Failover Group.

## Example Usage

-> **Note:** For a more complete example, see the [`./examples/sql-azure/managed_instance_failover_group` directory](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/sql-azure/managed_instance_failover_group) within the GitHub Repository.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_security_group" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_route_table" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet_route_table_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  route_table_id = azurerm_route_table.example.id
}

resource "azurerm_mssql_managed_instance" "primary" {
  name                         = "example-primary"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
  license_type                 = "BasePrice"
  subnet_id                    = azurerm_subnet.example.id
  sku_name                     = "GP_Gen5"
  vcores                       = 4
  storage_size_in_gb           = 32

  depends_on = [
    azurerm_subnet_network_security_group_association.example,
    azurerm_subnet_route_table_association.example,
  ]

  tags = {
    environment = "prod"
  }
}

resource "azurerm_mssql_managed_instance" "secondary" {
  name                         = "example-secondary"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
  license_type                 = "BasePrice"
  subnet_id                    = azurerm_subnet.example.id
  sku_name                     = "GP_Gen5"
  vcores                       = 4
  storage_size_in_gb           = 32

  depends_on = [
    azurerm_subnet_network_security_group_association.example,
    azurerm_subnet_route_table_association.example,
  ]

  tags = {
    environment = "prod"
  }
}

resource "azurerm_mssql_managed_instance_failover_group" "example" {
  name                        = "example-failover-group"
  location                    = azurerm_mssql_managed_instance.primary.location
  managed_instance_id         = azurerm_mssql_managed_instance.primary.id
  partner_managed_instance_id = azurerm_mssql_managed_instance.secondary.id

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Managed Instance Failover Group. Changing this forces a new resource to be created.

* `location` - The Azure Region where the Managed Instance Failover Group should exist. Changing this forces a new resource to be created.

* `managed_instance_id` - (Required) The ID of the Azure SQL Managed Instance which will be replicated using a Managed Instance Failover Group. Changing this forces a new resource to be created.

* `partner_managed_instance_id` - (Required) The ID of the Azure SQL Managed Instance which will be replicated to. Changing this forces a new resource to be created.

* `read_write_endpoint_failover_policy` - (Required) A `read_write_endpoint_failover_policy` block as defined below.

* `readonly_endpoint_failover_policy_enabled` - (Optional) Failover policy for the read-only endpoint. Defaults to `false`.

---

A `read_write_endpoint_failover_policy` block supports the following:

* `mode` - (Required) The failover mode. Possible values are `Automatic` or `Manual`.

* `grace_minutes` - (Optional) Applies only if `mode` is `Automatic`. The grace period in minutes before failover with data loss is attempted.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Managed Instance Failover Group.

* `partner_region` - A `partner_region` block as defined below.

* `role` - The local replication role of the Managed Instance Failover Group.

---

A `partner_region` block exports the following:

* `location` - The Azure Region where the Managed Instance Failover Group partner exists.

* `role` - The partner replication role of the Managed Instance Failover Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Instance Failover Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Instance Failover Group.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Instance Failover Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Instance Failover Group.

## Import

SQL Instance Failover Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_managed_instance_failover_group.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/locations/Location/instanceFailoverGroups/failoverGroup1
```
