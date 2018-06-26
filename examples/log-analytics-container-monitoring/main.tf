resource "azurerm_resource_group" "test" {
  name     = "k8s-log-analytics-test"
  location = "westeurope"
}

resource "random_id" "workspace" {
  keepers = {
    # Generate a new id each time we switch to a new resource group
    group_name = "${azurerm_resource_group.test.name}"
  }

  byte_length = 8
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "k8s-workspace-${random_id.workspace.hex}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Free"
}

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "Containers"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  workspace_resource_id = "${azurerm_log_analytics_workspace.test.id}"
  workspace_name        = "${azurerm_log_analytics_workspace.test.name}"

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/Containers"
  }
}

resource "azurerm_log_analytics_solution" "test2" {
  solution_name         = "Security"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  workspace_resource_id = "${azurerm_log_analytics_workspace.test.id}"
  workspace_name        = "${azurerm_log_analytics_workspace.test.name}"

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/Security"
  }
}
