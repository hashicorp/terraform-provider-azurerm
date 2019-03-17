resource "random_integer" "ri" {
  min = 100
  max = 999
}

resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group_name}"
  location = "${var.resource_group_location}"
}

module "vm" {
  source = "modules/vm"

  resource_group_name = "${azurerm_resource_group.rg.name}"
  vm_size             = "Standard_F2"
  prefix              = "tfexrecove${random_integer.ri.result}"
  hostname            = "tfexrecove${random_integer.ri.result}"
  dns_name            = "tfexrecove${random_integer.ri.result}"
  admin_username      = "vmadmin"
  admin_password      = "Password123!@#"
}

resource "azurerm_recovery_services_vault" "example" {
  name                = "tfex-recovery-vault"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  sku                 = "Standard"
}

resource "azurerm_recovery_services_protection_policy_vm" "simple" {
  name                = "tfex-policy-simple"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.example.name}"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}

resource "azurerm_recovery_services_protection_policy_vm" "advanced" {
  name                = "tfex-policy-advanced"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.example.name}"

  timezone = "UTC"

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

resource "azurerm_recovery_services_protected_vm" "example" {
  resource_group_name = "${azurerm_resource_group.rg.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.example.name}"
  source_vm_id        = "${module.vm.vm-id}"
  backup_policy_id    = "${azurerm_recovery_services_protection_policy_vm.simple.id}"
}
