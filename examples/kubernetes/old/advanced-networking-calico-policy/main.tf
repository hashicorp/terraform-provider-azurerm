provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "${var.prefix}-anw-resources"
  location = "${var.location}"
}

resource "azurerm_route_table" "test" {
  name                = "${var.prefix}-routetable"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name                   = "default"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "${var.prefix}-network"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  address_prefixes     = ["10.1.0.0/24"]
  virtual_network_name = "${azurerm_virtual_network.test.name}"
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = "${azurerm_subnet.test.id}"
  route_table_id = "${azurerm_route_table.test.id}"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "${var.prefix}-anw"
  location            = "${azurerm_resource_group.test.location}"
  dns_prefix          = "${var.prefix}-anw"
  resource_group_name = "${azurerm_resource_group.test.name}"

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
    vnet_subnet_id = "${azurerm_subnet.test.id}"
  }

  service_principal {
    client_id     = "${var.kubernetes_client_id}"
    client_secret = "${var.kubernetes_client_secret}"
  }

  network_profile {
    network_plugin = "azure"
    network_policy = "calico"
  }
}
