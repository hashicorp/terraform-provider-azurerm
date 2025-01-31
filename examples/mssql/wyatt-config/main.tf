terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "3.108.0"
    }
  }
}

provider "azurerm" {
  features {}
}

locals {
  name                     = "acctestwyatt5qlmi"
  primary_name             = "${local.name}-primary"
  primary_location         = "West Europe"
  failover_name            = "${local.name}-failover"
  failover_paired_location = "North Europe"
}

resource "azurerm_resource_group" "primary" {
  name     = local.primary_name
  location = local.primary_location
}

## Primary Instance

resource "azurerm_virtual_network" "primary" {
  name                = local.primary_name
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name
  address_space       = ["10.0.0.0/16"]
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

  tags = {
    environment = "prod"
  }
}

## Secondary (Fail-over) Instance

resource "azurerm_resource_group" "failover" {
  name     = local.failover_name
  location = local.failover_paired_location
}

resource "azurerm_virtual_network" "failover" {
  name                = local.failover_name
  location            = azurerm_resource_group.failover.location
  resource_group_name = azurerm_resource_group.failover.name
  address_space       = ["10.1.0.0/16"]
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

resource "azurerm_mssql_managed_instance_failover_group" "example" {
  name                        = "example-failover-group"
  location                    = azurerm_mssql_managed_instance.primary.location
  managed_instance_id         = azurerm_mssql_managed_instance.primary.id
  partner_managed_instance_id = azurerm_mssql_managed_instance.failover.id
  secondary_type = "Standby"

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }

  depends_on = [
    azurerm_private_dns_zone_virtual_network_link.primary,
    azurerm_private_dns_zone_virtual_network_link.failover,
  ]
}

resource "azurerm_virtual_network_peering" "failover_to_primary" {
  name                      = "failover-to-primary"
  remote_virtual_network_id = azurerm_virtual_network.primary.id
  resource_group_name       = azurerm_resource_group.failover.name
  virtual_network_name      = azurerm_virtual_network.failover.name
}

resource "azurerm_virtual_network_peering" "primary_to_failover" {
  name                      = "primary-to-failover"
  remote_virtual_network_id = azurerm_virtual_network.failover.id
  resource_group_name       = azurerm_resource_group.primary.name
  virtual_network_name      = azurerm_virtual_network.primary.name
}

resource "azurerm_private_dns_zone" "example" {
  name                = "${local.name}.private"
  resource_group_name = azurerm_resource_group.primary.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "primary" {
  name                  = "primary-link"
  resource_group_name   = azurerm_resource_group.primary.name
  private_dns_zone_name = azurerm_private_dns_zone.example.name
  virtual_network_id    = azurerm_virtual_network.primary.id
}

resource "azurerm_private_dns_zone_virtual_network_link" "failover" {
  name                  = "failover-link"
  resource_group_name   = azurerm_private_dns_zone.example.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.example.name
  virtual_network_id    = azurerm_virtual_network.failover.id
}

output "debug" {
  value =  azurerm_private_dns_zone.example.id
}