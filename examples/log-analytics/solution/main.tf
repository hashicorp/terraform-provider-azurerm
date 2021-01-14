provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "${var.prefix}-laworkspace"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "ContainerInsights"
  location              = "${azurerm_resource_group.example.location}"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  workspace_resource_id = "${azurerm_log_analytics_workspace.example.id}"
  workspace_name        = "${azurerm_log_analytics_workspace.example.name}"

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/Containers"
  }
}

resource "azurerm_log_analytics_solution" "test2" {
  solution_name         = "Security"
  location              = "${azurerm_resource_group.example.location}"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  workspace_resource_id = "${azurerm_log_analytics_workspace.example.id}"
  workspace_name        = "${azurerm_log_analytics_workspace.example.name}"

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/Security"
  }
}
