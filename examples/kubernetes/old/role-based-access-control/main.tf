# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-rbac-resources"
  location = var.location
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "${var.prefix}-rbac"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "${var.prefix}-rbac"

  agent_pool_profile {
    name            = "default"
    count           = 1
    vm_size         = "Standard_D1_v2"
    os_type         = "Linux"
    os_disk_size_gb = 30
  }

  service_principal {
    client_id     = var.kubernetes_client_id
    client_secret = var.kubernetes_client_secret
  }

  role_based_access_control {
    enabled = true
  }

  tags = {
    Environment = "Production"
  }
}
