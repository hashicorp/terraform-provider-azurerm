provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

module "virtual-network" {
  source              = "./modules/virtual-network"
  resource_group_name = "${azurerm_resource_group.example.name}"
  prefix              = "${var.prefix}"
}

module "virtual-machine" {
  source              = "./modules/virtual-machine"
  resource_group_name = "${azurerm_resource_group.example.name}"
  prefix              = "${var.prefix}"
  subnet_id           = "${module.virtual-network.subnet_id}"
}

resource "azurerm_recovery_services_vault" "example" {
  name                = "${var.prefix}-vault"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "Standard"
}

resource "azurerm_backup_policy_vm" "example" {
  name                = "tfex-policy-simple"
  resource_group_name = "${azurerm_resource_group.example.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.example.name}"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}

resource "azurerm_recovery_services_protected_vm" "example" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.example.name}"
  source_vm_id        = "${module.virtual-machine.id}"
  backup_policy_id    = "${azurerm_backup_policy_vm.example.id}"
}
