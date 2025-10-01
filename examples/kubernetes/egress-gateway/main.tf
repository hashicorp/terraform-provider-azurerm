# Example demonstrating AKS Egress Gateway feature
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-k8s-egress-gateway"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "exampleaks"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"
  }

  # Enable static egress gateway for the cluster
  network_profile {
    network_plugin = "azure"
    static_egress_gateway_profile {
      enabled = true
    }
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Demo"
    Feature     = "EgressGateway"
  }
}

# Create a node pool specifically for hosting egress gateways
resource "azurerm_kubernetes_cluster_node_pool" "egress_gateway" {
  name                  = "egressgw"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  vm_size               = "Standard_D2_v2"
  node_count            = 2
  mode                  = "Gateway"

  # Configure the gateway profile with public IP prefix size
  gateway_profile {
    public_ip_prefix_size = 30  # Supports 28-31, defaults to 31
  }

  tags = {
    Environment = "Demo"
    Purpose     = "EgressGateway"
  }
}

# Optional: Additional user node pool for regular workloads
resource "azurerm_kubernetes_cluster_node_pool" "user" {
  name                  = "user"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  vm_size               = "Standard_D2_v2"
  node_count            = 2
  mode                  = "User"

  tags = {
    Environment = "Demo"
    Purpose     = "UserWorkloads"
  }
}

output "kube_config" {
  description = "Kubernetes configuration for connecting to the cluster"
  value       = azurerm_kubernetes_cluster.example.kube_config_raw
  sensitive   = true
}

output "cluster_id" {
  description = "The Kubernetes Cluster ID"
  value       = azurerm_kubernetes_cluster.example.id
}