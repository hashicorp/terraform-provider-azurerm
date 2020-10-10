provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
}

resource "azurerm_network_security_group" "example" {
  name                = "acceptanceTestSecurityGroup1"
  location            = data.azurerm_resource_group.example.location
  resource_group_name = data.azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-network"
  resource_group_name = data.azurerm_resource_group.example.name
  location            = data.azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = data.azurerm_resource_group.example.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "miDelegation"

    service_delegation {
      name    = "Microsoft.Sql/managedInstances"
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_route_table" "example" {
  name                = "example-routetable"
  location            = data.azurerm_resource_group.example.location
  resource_group_name = data.azurerm_resource_group.example.name

  route {
    name                   = "example"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_subnet_route_table_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  route_table_id = azurerm_route_table.example.id
}

resource "azurerm_mssql_managed_instance" "dns" {
  name                 = "${var.prefix}-dns-mi"
  resource_group_name = data.azurerm_resource_group.example.name
  location            = data.azurerm_resource_group.example.location
  administrator_login = "demoReadUser"
  administrator_login_password = "ReadUserDemo@12345"
  subnet_id = "${azurerm_subnet.example.id}"
  identity {
    type = "SystemAssigned"
  }
   sku {
        capacity = 8
        family = "Gen5"
        name = "GP_Gen5"
        tier = "GeneralPurpose"
      }
      license_type = "BasePrice"
      collation =  "SQL_Latin1_General_CP1_CI_AS"
      proxy_override = "Redirect"
      storage_size_gb = 64
      vcores = 8
      public_data_endpoint_enabled = true
      timezone_id = "UTC"
      minimal_tls_version = "1.1"
}


resource "azurerm_mssql_managed_instance" "example" {
  name                 = "${var.prefix}-mi"
  resource_group_name = data.azurerm_resource_group.example.name
  location            = data.azurerm_resource_group.example.location
  administrator_login = "demoReadUser"
  administrator_login_password = "ReadUserDemo@123456"
  subnet_id = "${azurerm_subnet.example.id}"
  identity {
    type = "SystemAssigned"
  }
   sku {
        capacity = 8
        family = "Gen5"
        name = "GP_Gen5"
        tier = "GeneralPurpose"
      }
      license_type = "LicenseIncluded"
      collation =  "SQL_Latin1_General_CP1_CI_AS"
      proxy_override = "Redirect"
      storage_size_gb = 64
      vcores = 8
      public_data_endpoint_enabled = false
      timezone_id = "UTC"
      minimal_tls_version = "1.2"
      dns_zone_partner = azurerm_mssql_managed_instance.dns.id
}
