provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "rg-lucia"
  location = "australiaeast"
}

/*
resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-lucia"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  monitoring {
    alerts_for_all_job_failures_enabled            = true
    alerts_for_critical_operation_failures_enabled = false
  }

  soft_delete_enabled = false
}
*/

/*
resource "azurerm_recovery_services_vault" "test2" {
  name                = "acctest-Vault-lucia2"
  location            = "West Central US"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  monitoring {
    alerts_for_all_job_failures_enabled            = true
    alerts_for_critical_operation_failures_enabled = false
  }

  soft_delete_enabled = true
}
*/

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-lucia2-ae"
  location            = "australiaeast"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  monitoring {
    alerts_for_all_job_failures_enabled            = true
    alerts_for_critical_operation_failures_enabled = true
  }

  soft_delete_enabled = false
}