provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_network_security_group" "example" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-network"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  resource_group_name  = "${azurerm_resource_group.example.name}"
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
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

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

data "azurerm_mssql_managed_instance" "example" {
  name = "${var.prefix}-mi"
  resource_group_name =  "${azurerm_resource_group.example.name}"
}

resource "azurerm_mssql_managed_instance" "example" {
  name                 = "${var.prefix}-mi"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  administrator_login = "priyanka"
  administrator_login_password = "priya@123"
  subnet_id = "${azurerm_subnet.example.id}"
  identity {
    type = "SystemAssigned"
  }
      license_type = "LicenseIncluded"
      collation =  "SQL_Latin1_General_CP1_CI_AS"
      proxy_override = "Redirect"
      storage_size_gb = 32
        vcores = 16
        minimal_tls_version = "1.2"
}
