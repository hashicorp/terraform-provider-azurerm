provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-anw-resources"
  location = "${var.location}"
}

resource "azurerm_route_table" "example" {
  name                = "${var.prefix}-routetable"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  route {
    name                   = "default"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = azurerm_resource_group.example.name
  address_prefix       = "10.1.0.0/22"
}

resource "azurerm_subnet_route_table_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  route_table_id = azurerm_route_table.example.id
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "${var.prefix}-anw"
  location            = "${azurerm_resource_group.example.location}"
  dns_prefix          = "${var.prefix}-anw"
  resource_group_name = "${azurerm_resource_group.example.name}"

  linux_profile {
    admin_username = "acctestuser1"

    ssh_key {
      key_data = "${file(var.public_ssh_key_path)}"
    }
  }

  agent_pool_profile {
    name            = "agentpool"
    count           = "2"
    vm_size         = "Standard_DS2_v2"
    os_type         = "Linux"
    os_disk_size_gb = 30

    # Required for advanced networking
    vnet_subnet_id = "${azurerm_subnet.example.id}"
  }

  service_principal {
    client_id     = "${var.kubernetes_client_id}"
    client_secret = "${var.kubernetes_client_secret}"
  }

  network_profile {
    network_plugin = "azure"
  }

  depends_on = [azurerm_subnet_route_table_association.example]
}
