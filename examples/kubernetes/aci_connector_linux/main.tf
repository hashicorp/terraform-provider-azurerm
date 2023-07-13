# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-k8s-rg"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.10.0.0/16"]
}

resource "azurerm_subnet" "example-nodepool" {
  name                 = "default"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = azurerm_resource_group.example.name
  address_prefixes     = ["10.10.1.0/24"]
}

resource "azurerm_subnet" "example-aci" {
  name                 = "aci"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = azurerm_resource_group.example.name
  address_prefixes     = ["10.10.3.0/24"]

  delegation {
    name = "aciDelegation"
    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "${var.prefix}-k8s"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "${var.prefix}-k8s"

  default_node_pool {
    name           = "default"
    node_count     = 1
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.example-nodepool.id
  }

  network_profile {
    network_plugin    = "azure"
    network_policy    = "azure"
    load_balancer_sku = "standard"
  }

  identity {
    type = "SystemAssigned"
  }

  aci_connector_linux {
    subnet_name = azurerm_subnet.example-aci.name
  }

  azure_policy_enabled             = false
  http_application_routing_enabled = false
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_subnet.example-aci.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_kubernetes_cluster.example.aci_connector_linux[0].connector_identity[0].object_id
}
