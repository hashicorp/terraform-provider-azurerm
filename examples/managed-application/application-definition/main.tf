provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_role_definition" "builtin" {
  name = "Contributor"
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_managed_application_definition" "example" {
  name                 = "${var.prefix}managedappdef"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  lock_level           = "ReadOnly"
  package_file_uri     = "https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip"
  display_name         = "TestManagedApplicationDefinition"
  description          = "Test Managed Application Definition"

  authorization {
    service_principal_id = data.azurerm_client_config.current.object_id
    role_definition_id   = split("/", data.azurerm_role_definition.builtin.id)[length(split("/", data.azurerm_role_definition.builtin.id)) - 1]
  }
}
