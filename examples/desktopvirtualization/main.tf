provider "azurerm" {
       features {}
}

resource "azurerm_resource_group" "example" {
    name       ="${var.prefix}-resources"
    location   = "${var.location}"
}

resource "azurerm_virtual_desktop_workspace" "example" {
    name                     = "${var.prefix}workspace"
    resource_group_name      = azurerm_resource_group.example.name
    location                 = azurerm_resource_group.example.location
}

resource "azurerm_virtual_desktop_host_pool" "example" {
    resource_group_name      = azurerm_resource_group.example.name
    name                     = "${var.prefix}hostpool"
    location                 = azurerm_resource_group.example.location
    
    validate_environment     = false
    type                     = "Pooled"
    maximum_sessions_allowed = 16
    load_balancer_type       = "BreadthFirst"
}

resource "azurerm_virtual_desktop_application_group" "example" {
    resource_group_name      = azurerm_resource_group.example.name
    host_pool_id             = azurerm_virtual_desktop_host_pool.example.id
    location                 = azurerm_resource_group.example.location
    type                     = "Desktop"
    name                     = "${var.prefix}dag"
    depends_on               = [azurerm_virtual_desktop_host_pool.example]
}

resource "azurerm_virtual_desktop_workspace_application_group_association" "example" {
    application_group_id     = azurerm_virtual_desktop_application_group.example.id
    workspace_id             = azurerm_virtual_desktop_workspace.example.id
}
