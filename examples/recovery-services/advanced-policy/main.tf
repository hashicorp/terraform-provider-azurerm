provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_recovery_services_vault" "example" {
  name                = "${var.prefix}-vault"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "Standard"
}

resource "azurerm_backup_policy_vm" "example" {
  name                = "${var.prefix}-policy"
  resource_group_name = "${azurerm_resource_group.example.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.example.name}"
  timezone            = "UTC"

  backup {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Monday", "Wednesday"]
  }

  retention_weekly {
    weekdays = ["Monday", "Wednesday"]
    count    = 52
  }

  retention_monthly {
    weeks    = ["First", "Second"]
    weekdays = ["Monday", "Wednesday"]
    count    = 100
  }

  retention_yearly {
    months   = ["July"]
    weeks    = ["First", "Second"]
    weekdays = ["Monday", "Wednesday"]
    count    = 100
  }
}
