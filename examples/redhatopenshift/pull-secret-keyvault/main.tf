provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_key_vault" "example" {
  name                = var.key_vault_name
  resource_group_name	= var.key_vault_resource_group_name
}

data "azurerm_key_vault_secret" "pull_secret" {
  name			   = "pull-secret"
  key_vault_id = data.azurerm_key_vault.example.id
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-aro-resources"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-aro-vnet"
  address_space       = ["10.0.0.0/22"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "master_subnet" {
  name                 = "master-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/23"]
  service_endpoints    = ["Microsoft.ContainerRegistry"]
  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "worker_subnet" {
  name                 = "worker-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/23"]
  service_endpoints    = ["Microsoft.ContainerRegistry"]
}

resource "azurerm_redhatopenshift_cluster" "example" {
  name                = "${var.prefix}-aro"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  
  service_principal {
    client_id     = var.client_id
    client_secret = var.client_secret
  }

  cluster_profile {
    pull_secret = data.azurerm_key_vault_secret.pull_secret.value
  }

  master_profile {
    vm_size   = "Standard_D8s_v3"
    subnet_id = azurerm_subnet.master_subnet.id
  }
  
  worker_profile {
    vm_size   = "Standard_D4s_v3"
    subnet_id = azurerm_subnet.worker_subnet.id
  }
}
