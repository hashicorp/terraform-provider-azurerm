resource "azurerm_resource_group" "akc-rg" {
  name     = "${var.resource_group_name}"
  location = "${var.resource_group_location}"
}

#an attempt to keep the aci container group name (and dns label) somewhat unique
resource "random_integer" "random_int" {
  min = 100
  max = 999
}

resource azurerm_network_security_group "aks_advanced_network" {
  name                = "akc-${random_integer.random_int.result}-nsg"
  location            = "${var.resource_group_location}"
  resource_group_name = "${azurerm_resource_group.akc-rg.name}"
}

resource "azurerm_virtual_network" "aks_advanced_network" {
  name                = "akc-${random_integer.random_int.result}-vnet"
  location            = "${var.resource_group_location}"
  resource_group_name = "${azurerm_resource_group.akc-rg.name}"
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "aks_subnet" {
  name                      = "akc-${random_integer.random_int.result}-subnet"
  resource_group_name       = "${azurerm_resource_group.akc-rg.name}"
  network_security_group_id = "${azurerm_network_security_group.aks_advanced_network.id}"
  address_prefix            = "10.1.0.0/24"
  virtual_network_name      = "${azurerm_virtual_network.aks_advanced_network.name}"
}

resource "azurerm_kubernetes_cluster" "aks_container" {
  name       = "akc-${random_integer.random_int.result}"
  location   = "${var.resource_group_location}"
  dns_prefix = "akc-${random_integer.random_int.result}"

  resource_group_name = "${azurerm_resource_group.akc-rg.name}"

  linux_profile {
    admin_username = "${var.linux_admin_username}"

    ssh_key {
      key_data = "${var.linux_admin_ssh_publickey}"
    }
  }

  agent_pool_profile {
    name    = "agentpool"
    count   = "2"
    vm_size = "Standard_DS2_v2"
    os_type = "Linux"

    # Required for advanced networking
    vnet_subnet_id = "${azurerm_subnet.aks_subnet.id}"
  }

  service_principal {
    client_id     = "${var.client_id}"
    client_secret = "${var.client_secret}"
  }

  network_profile {
    network_plugin     = "azure"
    dns_service_ip     = "10.0.0.10"
    docker_bridge_cidr = "172.17.0.1/16"
    service_cidr       = "10.0.0.0/16"
  }
}
