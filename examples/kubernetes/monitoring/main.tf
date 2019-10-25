resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-mon-resources"
  location = "${var.location}"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "${var.prefix}-law"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "example" {
  solution_name         = "Containers"
  location              = "${azurerm_resource_group.example.location}"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  workspace_resource_id = "${azurerm_log_analytics_workspace.example.id}"
  workspace_name        = "${azurerm_log_analytics_workspace.example.name}"

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/Containers"
  }
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "${var.prefix}-mon"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  dns_prefix          = "${var.prefix}-mon"

  agent_pool_profile {
    name            = "default"
    count           = 1
    vm_size         = "Standard_D1_v2"
    os_type         = "Linux"
    os_disk_size_gb = 30
  }

  service_principal {
    client_id     = "${var.kubernetes_client_id}"
    client_secret = "${var.kubernetes_client_secret}"
  }

  addon_profile {
    oms_agent {
      enabled                    = true
      log_analytics_workspace_id = "${azurerm_log_analytics_workspace.example.id}"
    }
  }

  tags = {
    Environment = "Production"
  }
}
