terraform {
  # required_providers {
  #   azurerm = {
  #     source  = "hashicorp/azurerm"
  #     version = "3.108.0"
  #   }
  # }
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "3.6.2"
    }
  }
}

provider "azurerm" {
  features {}
  subscription_id = "1a6092a6-137e-4025-9a7c-ef77f76f2c02"
}

resource "azurerm_resource_group" "example" {
  name     = "acctest-mssql-managed-instance" # TODO add "Wyatt"
  location = "eastus"
}

resource "azurerm_network_security_group" "example" {
  name                = "acctest-mssql-managed-instance"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}


resource "azurerm_network_security_rule" "allow_management_inbound" {
  name                        = "allow_management_inbound"
  priority                    = 106
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges = ["9000", "9003", "1438", "1440", "1452"]
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "allow_misubnet_inbound" {
  name                        = "allow_misubnet_inbound"
  priority                    = 200
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "10.0.0.0/24"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "allow_health_probe_inbound" {
  name                        = "allow_health_probe_inbound"
  priority                    = 300
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "AzureLoadBalancer"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "allow_tds_inbound" {
  name                        = "allow_tds_inbound"
  priority                    = 1000
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "1433"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "deny_all_inbound" {
  name                        = "deny_all_inbound"
  priority                    = 4096
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "allow_management_outbound" {
  name                        = "allow_management_outbound"
  priority                    = 106
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges = ["80", "443", "12000"]
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "allow_misubnet_outbound" {
  name                        = "allow_misubnet_outbound"
  priority                    = 200
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "10.0.0.0/24"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "deny_all_outbound" {
  name                        = "deny_all_outbound"
  priority                    = 4096
  direction                   = "Outbound"
  access                      = "Deny"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "vnet-mi"
  resource_group_name = azurerm_resource_group.example.name
  address_space = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
}

resource "azurerm_subnet" "example" {
  name                              = "subnet-mi"
  resource_group_name               = azurerm_resource_group.example.name
  virtual_network_name              = azurerm_virtual_network.example.name
  private_endpoint_network_policies = "Enabled"
  address_prefixes = ["10.0.0.0/24"]

  delegation {
    name = "managedinstancedelegation"

    service_delegation {
      name = "Microsoft.Sql/managedInstances"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"
      ]
    }
  }
}

resource "azurerm_subnet" "another_example" {
  name                              = "${random_pet.example.id}-subnet-mi"
  resource_group_name               = azurerm_resource_group.example.name
  virtual_network_name              = azurerm_virtual_network.example.name
  private_endpoint_network_policies = "Enabled"
  address_prefixes = ["10.0.1.0/24"]

  delegation {
    name = "managedinstancedelegation"

    service_delegation {
      name = "Microsoft.Sql/managedInstances"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_subnet_network_security_group_association" "another_example" {
  subnet_id                 = azurerm_subnet.another_example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_route_table" "example" {
  name     = "routetable-mi"
  location = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  # disable_bgp_route_propagation = false
  depends_on = [
    azurerm_subnet.example,
  ]
}

resource "azurerm_subnet_route_table_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  route_table_id = azurerm_route_table.example.id
}

resource "azurerm_subnet_route_table_association" "another_example" {
  subnet_id      = azurerm_subnet.another_example.id
  route_table_id = azurerm_route_table.example.id
}

resource "random_pet" "example" {}

resource "azurerm_mssql_managed_instance" "example" {
  name                = "${random_pet.example.id}-managedsqlinstance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.another_example.id
  vcores             = 4

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id
    ]
  }

  administrator_login          = "mradministrator"
  administrator_login_password = "UPPERlower123!SixteenCharacters"

  depends_on = [
    azurerm_subnet_network_security_group_association.example,
    azurerm_subnet_route_table_association.example,
    azurerm_subnet_network_security_group_association.another_example,
    azurerm_subnet_route_table_association.another_example,
  ]
}

resource "azurerm_user_assigned_identity" "example" {
  location            = azurerm_resource_group.example.location
  name                = "UAIexample"
  resource_group_name = azurerm_resource_group.example.name
}

import {
  id = "/subscriptions/1a6092a6-137e-4025-9a7c-ef77f76f2c02/resourceGroups/acctest-mssql-managed-instance/providers/Microsoft.Sql/managedInstances/${random_pet.example.id}-managedsqlinstance"
  to = azurerm_mssql_managed_instance.example
}