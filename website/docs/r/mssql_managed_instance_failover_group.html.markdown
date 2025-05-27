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
locals {
  name              = "mymssqlmitest"
  primary_name      = "${local.name}-primary"
  primary_location  = "West Europe"
  failover_name     = "${local.name}-failover"
  failover_location = "North Europe"
}

resource "azurerm_mssql_managed_instance_failover_group" "example" {
  name                        = "example-failover-group"
  location                    = azurerm_mssql_managed_instance.primary.location
  managed_instance_id         = azurerm_mssql_managed_instance.primary.id
  partner_managed_instance_id = azurerm_mssql_managed_instance.failover.id
  secondary_type              = "Geo"
  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }
  depends_on = [
    azurerm_private_dns_zone_virtual_network_link.primary,
    azurerm_private_dns_zone_virtual_network_link.failover,
  ]
}

resource "azurerm_private_dns_zone" "example" {
  name                = "${local.name}.private"
  resource_group_name = azurerm_resource_group.primary.name
}

## Primary SQL Managed Instance
resource "azurerm_resource_group" "primary" {
  name     = local.primary_name
  location = local.primary_location
}

resource "azurerm_virtual_network" "primary" {
  name                = local.primary_name
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_private_dns_zone_virtual_network_link" "primary" {
  name                  = "primary-link"
  resource_group_name   = azurerm_resource_group.primary.name
  private_dns_zone_name = azurerm_private_dns_zone.example.name
  virtual_network_id    = azurerm_virtual_network.primary.id
}

resource "azurerm_subnet" "primary" {
  name                 = local.primary_name
  resource_group_name  = azurerm_resource_group.primary.name
  virtual_network_name = azurerm_virtual_network.primary.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
      name = "Microsoft.Sql/managedInstances"
    }
  }
}

resource "azurerm_network_security_group" "primary" {
  name                = local.primary_name
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name
}

resource "azurerm_subnet_network_security_group_association" "primary" {
  subnet_id                 = azurerm_subnet.primary.id
  network_security_group_id = azurerm_network_security_group.primary.id
}

resource "azurerm_route_table" "primary" {
  name                = local.primary_name
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name
}

resource "azurerm_subnet_route_table_association" "primary" {
  subnet_id      = azurerm_subnet.primary.id
  route_table_id = azurerm_route_table.primary.id
}

resource "azurerm_mssql_managed_instance" "primary" {
  name                         = local.primary_name
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
}

resource "azurerm_virtual_network_peering" "primary_to_failover" {
  name                      = "primary-to-failover"
  remote_virtual_network_id = azurerm_virtual_network.failover.id
  resource_group_name       = azurerm_resource_group.primary.name
  virtual_network_name      = azurerm_virtual_network.primary.name
}

## Secondary (Fail-over) SQL Managed Instance
resource "azurerm_resource_group" "failover" {
  name     = local.failover_name
  location = local.failover_location
}

resource "azurerm_virtual_network" "failover" {
  name                = local.failover_name
  location            = azurerm_resource_group.failover.location
  resource_group_name = azurerm_resource_group.failover.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_private_dns_zone_virtual_network_link" "failover" {
  name                  = "failover-link"
  resource_group_name   = azurerm_private_dns_zone.example.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.example.name
  virtual_network_id    = azurerm_virtual_network.failover.id
}

resource "azurerm_subnet" "default" {
  name                 = "default"
  resource_group_name  = azurerm_resource_group.failover.name
  virtual_network_name = azurerm_virtual_network.failover.name
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "failover" {
  name                 = "ManagedInstance"
  resource_group_name  = azurerm_resource_group.failover.name
  virtual_network_name = azurerm_virtual_network.failover.name
  address_prefixes     = ["10.1.1.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
      name = "Microsoft.Sql/managedInstances"
    }
  }
}

resource "azurerm_network_security_group" "failover" {
  name                = local.failover_name
  location            = azurerm_resource_group.failover.location
  resource_group_name = azurerm_resource_group.failover.name
}

resource "azurerm_subnet_network_security_group_association" "failover" {
  subnet_id                 = azurerm_subnet.failover.id
  network_security_group_id = azurerm_network_security_group.failover.id
}

resource "azurerm_route_table" "failover" {
  name                = local.failover_name
  location            = azurerm_resource_group.failover.location
  resource_group_name = azurerm_resource_group.failover.name
}

resource "azurerm_subnet_route_table_association" "failover" {
  subnet_id      = azurerm_subnet.failover.id
  route_table_id = azurerm_route_table.failover.id
}

resource "azurerm_mssql_managed_instance" "failover" {
  name                         = local.failover_name
  resource_group_name          = azurerm_resource_group.failover.name
  location                     = azurerm_resource_group.failover.location
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
  license_type                 = "BasePrice"
  subnet_id                    = azurerm_subnet.failover.id
  sku_name                     = "GP_Gen5"
  vcores                       = 4
  storage_size_in_gb           = 32
  dns_zone_partner_id          = azurerm_mssql_managed_instance.primary.id

  depends_on = [
    azurerm_subnet_network_security_group_association.failover,
    azurerm_subnet_route_table_association.failover,
  ]
}

resource "azurerm_virtual_network_peering" "failover_to_primary" {
  name                      = "failover-to-primary"
  remote_virtual_network_id = azurerm_virtual_network.primary.id
  resource_group_name       = azurerm_resource_group.failover.name
  virtual_network_name      = azurerm_virtual_network.failover.name
}
```

-> **Note:** There are many prerequisites that must be in place before creating the failover group. To see them all, refer to [Configure a failover group for Azure SQL Managed Instance](https://learn.microsoft.com/en-us/azure/azure-sql/managed-instance/failover-group-configure-sql-mi).

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Managed Instance Failover Group. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Managed Instance Failover Group should exist. Changing this forces a new resource to be created.

* `managed_instance_id` - (Required) The ID of the Azure SQL Managed Instance which will be replicated using a Managed Instance Failover Group. Changing this forces a new resource to be created.

* `partner_managed_instance_id` - (Required) The ID of the Azure SQL Managed Instance which will be replicated to. Changing this forces a new resource to be created.

* `read_write_endpoint_failover_policy` - (Required) A `read_write_endpoint_failover_policy` block as defined below.

* `readonly_endpoint_failover_policy_enabled` - (Optional) Failover policy for the read-only endpoint. Defaults to `true`.

* `secondary_type` - (Optional) The type of the secondary Managed Instance. Possible values are `Geo`, `Standby`. Defaults to `Geo`.

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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Instance Failover Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Instance Failover Group.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Instance Failover Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Instance Failover Group.

## Import

SQL Instance Failover Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_managed_instance_failover_group.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/locations/Location/instanceFailoverGroups/failoverGroup1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Sql`: 2023-08-01-preview
