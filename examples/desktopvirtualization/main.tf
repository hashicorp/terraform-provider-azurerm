#terraform
terraform {
    required_version        = ">=0.13.4"
}

#Azure provider
provider "azurerm" {
    version                 = "2.31.1"
    features {}
    subscription_id         = var.subscriptionID
}

#Create WVD workspace
resource "azurerm_virtual_desktop_workspace" "TFWVDWorkspace" {
    name                    = "TFWorkspace"
    location                = var.location
    resource_group_name     = var.resourcegroupname
}

# Create WVD host pool
resource "azurerm_virtual_desktop_host_pool" "TFWVDHP" {
    location                 = var.location
    name                     = "TF-WVD-HP"
    type                     = "Pooled"
    load_balancer_type       = "BreadthFirst"  ##[BreadthFirst DepthFirst Persistent]
    maximum_sessions_allowed = 16
    validate_environment     = false
    resource_group_name      = var.resourcegroupname
}

# Create WVD DAG
resource "azurerm_virtual_desktop_application_group" "TFWVDDAG" {
    location                = var.location
    type                    = "Desktop"
    name                    = "TF-WVD-DAG"
    host_pool_id            = azurerm_virtual_desktop_host_pool.TFWVDHP.id
    resource_group_name     = var.resourcegroupname
    depends_on              =[azurerm_virtual_desktop_host_pool.TFWVDHP]
}

# Associate Workspace and DAG
resource "azurerm_virtual_desktop_workspace_application_group_association" "TFWVDAS" {
    application_group_id    = azurerm_virtual_desktop_application_group.TFWVDDAG.id
    workspace_id            = azurerm_virtual_desktop_workspace.TFWVDWorkspace.id
}
