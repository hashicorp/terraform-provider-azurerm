provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  location = var.location
  name     = "${var.prefix}-rg"
}

module "network" {
  source = "./modules/network"

  location            = var.location
  prefix              = var.prefix
  resource_group_name = azurerm_resource_group.example.name
}

module "active-directory-domain" {
  source = "./modules/active-directory-domain"

  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  active_directory_domain_name  = "${var.prefix}.local"
  active_directory_netbios_name = var.prefix
  admin_username                = var.admin_username
  admin_password                = var.admin_password
  prefix                        = var.prefix
  subnet_id                     = module.network.domain_controllers_subnet_id
}

module "active-directory-member" {
  source = "./modules/domain-member"

  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  prefix              = var.prefix

  active_directory_domain_name = "${var.prefix}.local"
  active_directory_username    = var.admin_username
  active_directory_password    = var.admin_password
  admin_username               = var.admin_username
  admin_password               = var.admin_password
  subnet_id                    = module.network.domain_members_subnet_id
}
